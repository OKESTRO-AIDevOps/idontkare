package lifecycle

import (
	pkgresourcecd "github.com/OKESTRO-AIDevOps/idontkare/pkg/resource/cd"
)

type LifecycleManifest struct {
	Request *LifecycleManifest_Request `yaml:"request"`

	Process struct {
		ProjectId    int    `yaml:"project_id"`
		ProjectName  string `yaml:"project_name"`
		UserId       int    `yaml:"user_id"`
		UserName     string `yaml:"user_name"`
		ClusterId    int    `yaml:"cluster_id"`
		ClusterName  string `yaml:"cluster_name"`
		LifecylcleId int    `yaml:"lifecycle_id"`
	} `yaml:"process"`

	Service *pkgresourcecd.NodePort `yaml:"service,omitempty"`

	Deployment *pkgresourcecd.Deployment `yaml:"deployment,omitempty"`
}

type LifecycleManifest_Request struct {
}
