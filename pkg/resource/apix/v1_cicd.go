package apix

import (
	"time"
)

type V1ProjectCiExport struct {
	ProjectName     string    `yaml:"project_name"`
	ClusterName     string    `yaml:"cluster_name"`
	ProjectCiStatus string    `yaml:"project_ci_status"`
	ProjectCiLog    string    `yaml:"project_ci_log"`
	ProjectCiStart  time.Time `yaml:"project_ci_start"`
	ProjectCiEnd    time.Time `yaml:"project_ci_end"`
}

type V1ProjectCdExport struct {
	ProjectName     string    `yaml:"project_name"`
	ClusterName     string    `yaml:"cluster_name"`
	ProjectCdStatus string    `yaml:"project_cd_status"`
	ProjectCdLog    string    `yaml:"project_cd_log"`
	ProjectCdStart  time.Time `yaml:"project_cd_start"`
	ProjectCdEnd    time.Time `yaml:"project_cd_end"`
}
