package manager

type Outcome any

type OutcomeSerializer interface {
	Marshal(o Outcome) ([]byte, error)
	Unmarshal(bs []byte) (Outcome, error)
}

func unmarshalOutcome(outcomeSer OutcomeSerializer, bs []byte) (id WorkflowID,
	outcome Outcome, err error,
) {
	var envelope OutcomeEnvelope
	envelope, _, err = OutcomeEnvelopeDTS.Unmarshal(bs)
	if err != nil {
		return
	}
	outcome, err = outcomeSer.Unmarshal(envelope.Bs)
	if err != nil {
		return
	}
	return envelope.WorkflowID, outcome, nil
}
