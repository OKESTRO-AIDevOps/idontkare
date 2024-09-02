package cd

const (
	STATUS_READY     CdStatusType = "ready"
	STATUS_RUNNING   CdStatusType = "running"
	STATUS_ERROR     CdStatusType = "error"
	STATUS_COMPLETED CdStatusType = "completed"
)

type CdStatusType string
