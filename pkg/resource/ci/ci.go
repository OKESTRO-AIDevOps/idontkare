package ci

const (
	STATUS_READY     CiStatusType = "ready"
	STATUS_RUNNING   CiStatusType = "running"
	STATUS_ERROR     CiStatusType = "error"
	STATUS_COMPLETED CiStatusType = "completed"
)

var BUILD_PROTOCOL string = "git://"

var BUILD_EXECUTOR string = "gcr.io/kaniko-project/executor:latest"

var BUILD_FILE_DEAFULT string = "Dockerfile"

var BUILD_LOG_SEP string = "\n===========BUILD LOG==========\n"

type CiStatusType string

type Builder struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Containers    []Builder_Container `yaml:"containers"`
		RestartPolicy string              `yaml:"restartPolicy"`
		Volumes       []Builder_Volume    `yaml:"volumes"`
	} `yaml:"spec"`
}

type Builder_Container struct {
	Name         string                          `yaml:"name"`
	Image        string                          `yaml:"image"`
	Args         []string                        `yaml:"args"`
	VolumeMounts []Builder_Container_VolumeMount `yaml:"volumeMounts"`
	Env          []Builder_Container_Env         `yaml:"env"`
}

type Builder_Container_VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type Builder_Container_Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Builder_Volume struct {
	Name   string `yaml:"name"`
	Secret struct {
		SecretName string                `yaml:"secretName"`
		Items      []Builder_Volume_Item `yaml:"items"`
	} `yaml:"secret"`
}

type Builder_Volume_Item struct {
	Key  string `yaml:"key"`
	Path string `yaml:"path"`
}
