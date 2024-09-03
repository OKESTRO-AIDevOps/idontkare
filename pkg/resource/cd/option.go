package cd

import (
	"time"
)

type CdOption struct {
	Request *CdOption_Request `yaml:"request,omitempty"`

	Process *struct {
		StoredRequest CdOption_Request `yaml:"stored_request"`
		ProjectIndex  int              `yaml:"project_index"`
		ProjectId     int              `yaml:"project_id"`
		UserId        int              `yaml:"user_id"`
		ProjectName   string           `yaml:"project_name"`
		LifecycleId   int              `yaml:"lifecycle_id"`
		Error         bool             `yaml:"error"`
		Log           string           `yaml:"log"`
	} `yaml:"process,omitempty"`

	Service *NodePort `yaml:"service,omitempty"`

	Deployment *Deployment `yaml:"deployment,omitempty"`

	Response *struct {
		ProcessedTimestamp time.Time `yaml:"processed_timestamp"`
		Error              bool      `yaml:"error"`
		Log                string    `yaml:"log"`
	} `yaml:"response,omitempty"`
}

type CdOption_Request struct {
	DependOnCI bool        `yaml:"depend_on_ci"`
	Expose     map[int]int `yaml:"expose"`
}

type NodePort struct {
	Kind       string `yaml:"kind"`
	APIVersion string `yaml:"apiVersion"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Type     string `yaml:"type"`
		Selector struct {
			App string `yaml:"app"`
		} `yaml:"selector"`
		Ports []NodePort_Ports `yaml:"ports"`
	} `yaml:"spec"`
}

type NodePort_Ports struct {
	NodePort   int `yaml:"nodePort"`
	Port       int `yaml:"port"`
	TargetPort int `yaml:"targetPort"`
}

type Deployment struct {
	APIVersion string `yaml:"apiVersion"`
	Kind       string `yaml:"kind"`
	Metadata   struct {
		Name string `yaml:"name"`
	} `yaml:"metadata"`
	Spec struct {
		Selector struct {
			MatchLabels struct {
				App string `yaml:"app"`
			} `yaml:"matchLabels"`
		} `yaml:"selector"`
		Replicas int `yaml:"replicas"`
		Template struct {
			Metadata struct {
				Labels struct {
					App string `yaml:"app"`
				} `yaml:"labels"`
			} `yaml:"metadata"`
			Spec struct {
				ImagePullSecrets []Deployment_ImagePullSecrets `yaml:"imagePullSecrets"`
				Containers       []Deployment_Containers       `yaml:"containers"`
			} `yaml:"spec"`
		} `yaml:"template"`
	} `yaml:"spec"`
}

type Deployment_ImagePullSecrets struct {
	Name string `yaml:"name"`
}

type Deployment_Containers struct {
	Name            string                        `yaml:"name"`
	Image           string                        `yaml:"image"`
	ImagePullPolicy string                        `yaml:"imagePullPolicy"`
	Ports           []Deployment_Containers_Ports `yaml:"ports"`
}

type Deployment_Containers_Ports struct {
	ContainerPort int `yaml:"containerPort"`
}
