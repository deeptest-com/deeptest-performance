package serverDomain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/proto"
)

type PlanExecReq struct {
	PlanId int    `json:"planId"`
	Uuid   string `json:"uuid"`
	Title  string `json:"title"`

	GenerateType consts.GeneratorType `json:"generateType"`
	Stages       []*proto.Stage       `json:"stages"`

	Scenarios []*proto.Scenario `json:"scenarios"`

	NsqServerAddress string `json:"nsqServerAddress,omitempty"`
	NsqLookupAddress string `json:"nsqLookupAddress,omitempty"`
}
