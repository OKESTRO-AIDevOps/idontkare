package resource

import "github.com/gorilla/websocket"

const (
	AGENT_STATUS_CHALLENGE AgentStatusType = 0
	AGENT_STATUS_SUCCESS   AgentStatusType = 1
)

type AgentStatusType int

type AgentRegister map[string]AgentData

type AgentAddressRegister map[*websocket.Conn]string

type AgentData struct {
	C      *websocket.Conn
	Status AgentStatusType
	Key    string
}

type ChallengeData struct {
	Pass string `json:"pass"`
	Key  string `json:"key"`
}
