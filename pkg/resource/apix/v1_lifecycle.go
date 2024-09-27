package apix

import (
	"time"
)

type V1LifecycleReportExport struct {
	ProjectName     string    `yaml:"project_name"`
	LifecycleReport string    `yaml:"lifecycle_report"`
	LifecycleStart  time.Time `yaml:"lifecycle_start"`
}
