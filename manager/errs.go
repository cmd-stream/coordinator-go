package manager

import (
	"errors"
	"fmt"

	com "github.com/mus-format/common-go"
)

var ErrAnotherOrchestratorDetected = errors.New("another orchestrator detected")

// TODO Rename.
var ErrServiceSuspended = errors.New("service suspended")

func NewUnexpectedDTMError(dtm com.DTM, partID PartID, partSeq PartSeq) error {
	return fmt.Errorf("unexpected DTM %v on PartID %v, PartSeq %v", dtm, partID,
		partSeq)
}
