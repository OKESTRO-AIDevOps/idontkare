package apix

const (
	V1KindClientRequest     string = "client.request"
	V1KindClientRequestPriv string = "client.request.priv"
	V1KindAgentRequest      string = "agent.request"
	V1KindAgentRequestPriv  string = "agent.request.priv"
	V1KindAgentPush         string = "agent.push"
	V1KindServerWrite       string = "server.write"
)

type V1Manifest struct {
	Main []V1Main `yaml:"main"`
}

type V1Main struct {
	Kind string `yaml:"kind"`
	Path string `yaml:"path"`
	Head V1Head `yaml:"head,omitempty"`
	Body V1Body `yaml:"body"`
}

type V1Head struct {
	Help V1Help `yaml:"help,omitempty"`
}

type V1Help map[string]string

type V1Body map[string]string

type V1ResultData struct {
	Output string `yaml:"output"`
}
