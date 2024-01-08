package domain

import (
	"github.com/aaronchen2k/deeptest/proto"
)

type Task struct {
	Uuid   string         `json:"uuid,omitempty"`
	Stages []*proto.Stage `json:"stages"`
	Dur    int            `json:"dur,omitempty"`

	VuNo     int               `json:"vuNo,omitempty"`
	Scenario []*proto.Scenario `json:"scenario,omitempty"`

	NsqServerAddress string `json:"nsqServerAddress,omitempty"`
	NsqLookupAddress string `json:"nsqLookupAddress,omitempty"`
}

type Scenario struct {
	Name       string           `json:"name"`
	Processors []proto.Scenario `json:"processors"`
	Dur        int              `json:"dur,omitempty"`
}

type Metrics struct {
	Name      string `json:"name"`
	Value     string `gorm:"type:text" json:"value"`
	Timestamp string `json:"timestamp"`
}
