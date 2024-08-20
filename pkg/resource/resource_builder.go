package resource

type KanikoBuilder struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Containers    []KanikoBuilder_Container `yaml:"containers"`
		RestartPolicy string                    `yaml:"restartPolicy"`
		Volumes       []KanikoBuilder_Volume    `yaml:"volumes"`
	} `yaml:"spec"`
}

type KanikoBuilder_Container struct {
	Name         string                                `yaml:"name"`
	Image        string                                `yaml:"image"`
	Args         []string                              `yaml:"args"`
	VolumeMounts []KanikoBuilder_Container_VolumeMount `yaml:"volumeMounts"`
	Env          []KanikoBuilder_Container_Env         `yaml:"env"`
}

type KanikoBuilder_Container_VolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
}

type KanikoBuilder_Container_Env struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type KanikoBuilder_Volume struct {
	Name   string `yaml:"name"`
	Secret struct {
		SecretName string                      `yaml:"secretName"`
		Items      []KanikoBuilder_Volume_Item `yaml:"items"`
	} `yaml:"secret"`
}

type KanikoBuilder_Volume_Item struct {
	Key  string `yaml:"key"`
	Path string `yaml:"path"`
}
