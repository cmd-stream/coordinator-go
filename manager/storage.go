package manager

import (
	"fmt"

	"github.com/cmd-stream/core-go"
	"github.com/cmd-stream/transport-go"
	com "github.com/mus-format/common-go"
)

type (
	PartID  int
	PartSeq int64
)

type LoadCallback func(partID PartID, partSeq PartSeq, bs []byte) error

type Storage interface {
	Save(bs []byte) (partID PartID, partSeq PartSeq, err error)
	SetCompleted(partID PartID, partSeq PartSeq) (err error)
	LoadUncompleted(callback LoadCallback) error
}

func loadFromStorage[T any](manager *Manager[T],
	codec transport.Codec[core.Result, core.Cmd[T]],
	outcomeSer OutcomeSerializer,
) (lastPartSeq PartSeq, err error) {
	idsMap := make(map[WorkflowID]CmdRecord[T])
	lastPartSeq = -1

	callback := func(partID PartID, partSeq PartSeq, bs []byte) (err error) {
		lastPartSeq = partSeq
		if bs == nil {
			panic(fmt.Sprintf("found nil item on PartID %v, PartSeq %v",
				partID, partSeq))
		}
		dtm := com.DTM(bs[0])
		switch dtm {
		case CmdEnvelopeDTM:
			var cmdRecord CmdRecord[T]
			cmdRecord, err = unmarshalCmdRecord(codec, bs)
			if err != nil {
				return
			}
			var (
				id       = WorkflowID{PartID: partID, PartSeq: partSeq}
				workflow = NewWorkflow(id, manager)
			)
			cmdRecord.Cmd = cmdRecord.Cmd.(Cmd[T]).SetWorkflow(workflow)
			select {
			case manager.queue <- cmdRecord:
				idsMap[id] = cmdRecord
			default:
				// TODO
				panic("queue is full")
			}
		case OutcomeEnvelopeDTM:
			var (
				id      WorkflowID
				outcome Outcome
			)
			id, outcome, err = unmarshalOutcome(outcomeSer, bs)
			if err != nil {
				return
			}
			// TODO Casting may panic.
			idsMap[id].Cmd.(Cmd[T]).Workflow().appendOutcome(outcome)
		default:
			panic(NewUnexpectedDTMError(dtm, partID, partSeq))
		}
		return
	}
	err = manager.options.Storage.LoadUncompleted(callback)
	manager.suspended = len(manager.queue) > 0
	return
}
