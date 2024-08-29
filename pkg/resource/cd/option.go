package cd

import (
	"time"
)

type CdOption struct {
	Request *CdOption_Request `yaml:"request,omitempty"`

	Process *struct {
		StoredRequest CdOption_Request `yaml:"stored_request"`
		ProjectIndex  int              `yaml:"project_index"`
		UserId        int              `yaml:"user_id"`
		ProjectName   string           `yaml:"project_name"`
		Error         error            `yaml:"error"`
	} `yaml:"process,omitempty"`

	Response *struct {
		ProcessedTimestamp time.Time `yaml:"processed_timestamp"`
		Error              error     `yaml:"error"`
	} `yaml:"response,omitempty"`
}

type CdOption_Request struct {
	DependOnCI bool `yaml:"depend_on_ci"`
}
