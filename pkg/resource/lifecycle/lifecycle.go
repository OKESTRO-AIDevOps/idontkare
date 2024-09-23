package lifecycle

import "time"

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

	Log *string `yaml:"log"`
}
