package apix

const (
	V1KindClientRequest     string = "client.request"
	V1KindClientRequestPriv string = "client.request.priv"
	V1KindAgentRequest      string = "agent.request"
	V1KindAgentRequestPriv  string = "agent.request.priv"
	V1KindAgentPush         string = "agent.push"
	V1KindServerPush        string = "server.push"
	V1KindServerRead        string = "server.read"
)

const (
	V1HeadFromFile string = "from-file"
	V1HeadHelp     string = "help"
)

const (
	V1RESULT_STATUS_SUCCESS V1ResultStatusType = "success"
	V1RESULT_STATUS_FAILURE V1ResultStatusType = "failure"
)

const (
	V1RESULT_FORMAT_STRING V1ResultFormatType = "string"
	V1RESULT_FORMAT_JSON   V1ResultFormatType = "json"
	V1RESULT_FORMAT_YAML   V1ResultFormatType = "yaml"
)

const (
	V1EXPORT_STRING_NULL    = "null"
	V1EXPORT_STRING_UNKNOWN = "__unknown__"
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

type V1ResultStatusType string

type V1ResultFormatType string

type V1ResultData struct {
	Status V1ResultStatusType `yaml:"status,omitempty"`
	Format V1ResultFormatType `yaml:"format,omitempty"`
	Output string             `yaml:"output"`
}
