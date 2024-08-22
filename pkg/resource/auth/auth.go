package resource

const (
	AGENT_STATUS_CHALLENGE AgentStatusType = 0
	AGENT_STATUS_SUCCESS   AgentStatusType = 1
)

type AgentStatusType int

type AgentRegister map[string]AgentData

type AgentData struct {
	Status AgentStatusType
	Key    string
}

type ChallengeData struct {
	Pass string `json:"pass"`
	Key  string `json:"key"`
}
