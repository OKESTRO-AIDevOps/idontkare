package apix

const (
	V1KindClientRequest     string = "client.request"
	V1KindClientRequestPriv string = "client.request.priv"
	V1KindAgentRequest      string = "agent.request"
	V1KindAgentRequestPriv  string = "agent.request.priv"
	V1KindAgentPush         string = "agent.push"
	V1KindServerWrite       string = "server.write"
	V1KindServerRead        string = "server.read"
)

const (
	V1HeadFromFile string = "from-file"
	V1HeadHelp     string = "help"
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
	FromFile V1HeadFromFileType `yaml:"from-file,omitempty"`
	Help     V1HeadHelpType     `yaml:"help,omitempty"`
}

type V1HeadFromFileType *string

type V1HeadHelpType map[string]string

type V1Body map[string]string

type V1ResultData struct {
	Output string `yaml:"output"`
}
