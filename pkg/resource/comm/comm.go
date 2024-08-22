package resource

import (
	"github.com/gorilla/websocket"
)

type CommStatusType string

type CommConnectionMap map[*websocket.Conn]string

const (
	COMM_STATUS_REQUEST CommStatusType = "request"
	COMM_STATUS_SUCCESS CommStatusType = "success"
	COMM_STATUS_FAILURE CommStatusType = "failure"
)

type CommJSON struct {
	Status  CommStatusType `json:"status"`
	Message string         `json:"message"`
	Data    []byte         `json:"data"`
}
