package resource

type CommType string
type CommStatusType string

const (
	COMM_REQUEST  CommType = "req"
	COMM_RESPONSE CommType = "resp"
)

const (
	COMM_STATUS_NONE    CommStatusType = "none"
	COMM_STATUS_SUCCESS CommStatusType = "success"
	COMM_STATUS_FAILURE CommStatusType = "failure"
)

type CommJSON struct {
	Type    CommType       `json:"type"`
	Status  CommStatusType `json:"status"`
	Message string         `json:"message"`
	Data    []byte         `json:"data"`
}
