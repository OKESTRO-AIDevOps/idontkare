package ci

import (
	"time"
)

type CiOption struct {
	Request *CiOption_Request `yaml:"request,omitempty"`

	Process *struct {
		StoredRequest CiOption_Request `yaml:"stored_request"`
		ProjectIndex  int              `yaml:"project_index"`
		ProjectId     int              `yaml:"project_id"`
		UserId        int              `yaml:"user_id"`
		ProjectName   string           `yaml:"project_name"`
		LinkToCd      int              `yaml:"link_to_cd"`
		Error         bool             `yaml:"error"`
		Log           string           `yaml:"log"`
	} `yaml:"process,omitempty"`

	Response *struct {
		ProcessedTimestamp time.Time `yaml:"processed_timestamp"`
		Error              bool      `yaml:"error"`
		Log                string    `yaml:"log"`
	} `yaml:"response,omitempty"`
}

type CiOption_Request struct {
	Immediate bool `yaml:"immediate"`
}
