package manager

type ServiceID string

type OutcomeEnvelope struct {
	ServiceID  ServiceID
	WorkflowID WorkflowID
	Bs         []byte
}
