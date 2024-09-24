package lifecycle

import "time"

type LifecycleStatusType string

const (
	LIFECYCLE_STATUS_OKAY  LifecycleStatusType = "okay"
	LIFECYCLE_STATUS_ERROR LifecycleStatusType = "error"
)

type LifecycleReport struct {
	SentTimestamp     time.Time              `yaml:"sent_timestamp"`
	ReceivedTimestamp time.Time              `yaml:"received_timestamp"`
	Obsolete          bool                   `yaml:"obsolete"`
	Process           LifecycleProcessInfo   `yaml:"process"`
	Detail            *LifecycleReportDetail `yaml:"detail"`
}

type LifecycleReportDetail struct {
	Service *string `yaml:"service"`

	Deployment *string `yaml:"deployment"`

	AppInfo *string `yaml:"app_info"`

	AppErrInfo *string `yaml:"app_err_info"`

	AppStatus LifecycleStatusType `yaml:"status"`
}
