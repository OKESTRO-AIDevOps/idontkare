package resource

type AgentBuilder struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Containers    []AgentBuilder_Container `yaml:"containers"`
		RestartPolicy string                   `yaml:"restartPolicy"`
		Volumes       []AgentBuilder_Volume    `yaml:"volumes"`
	} `yaml:"spec"`
}

type AgentBuilder_Container struct {
	Name         string                               `yaml:"name"`
	Image        string                               `yaml:"image"`
	Args         []string                             `yaml:"args"`
	VolumeMounts []AgentBuilder_Container_VolumeMount `yaml:"volumeMounts"`
	Env          []AgentBuilder_Container_Env         `yaml:"env"`
}

type AgentBuilder_Container_VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type AgentBuilder_Container_Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type AgentBuilder_Volume struct {
	Name   string `yaml:"name"`
	Secret struct {
		SecretName string                     `yaml:"secretName"`
		Items      []AgentBuilder_Volume_Item `yaml:"items"`
	} `yaml:"secret"`
}

type AgentBuilder_Volume_Item struct {
	Key  string `yaml:"key"`
	Path string `yaml:"path"`
}
