package manager

type WorkflowID struct {
	PartID  PartID
	PartSeq PartSeq
}

func NewWorkflow[T any](id WorkflowID, manager *Manager[T]) *Workflow[T] {
	return &Workflow[T]{id: id, manager: manager, outcomes: []Outcome{}}
}

type Workflow[T any] struct {
	id       WorkflowID
	manager  *Manager[T]
	outcomes []Outcome
}

func (w *Workflow[T]) ID() WorkflowID {
	return w.id
}

func (w *Workflow[T]) AppendOutcome(serID ServiceID, outcome Outcome) (err error) {
	err = w.manager.SaveOutcome(serID, w.id, outcome)
	if err != nil {
		return
	}
	w.appendOutcome(outcome)
	return
}

func (w *Workflow[T]) Outcomes() []Outcome {
	return w.outcomes
}

func (w *Workflow[T]) appendOutcome(outcome Outcome) {
	w.outcomes = append(w.outcomes, outcome)
}
