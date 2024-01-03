package generator

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/nsqio/go-nsq"
)

type VuGenerator interface {
	Generate(param VuGeneratorParam) error
}

type ConstantVuGenerator struct {
	param VuGeneratorParam
}

func (s ConstantVuGenerator) Generate(param VuGeneratorParam) (err error) {

	return
}

type RampVuGenerator struct {
	NsqServerAddress string
	NsqLookupAddress string

	Topic    string
	producer *nsq.Producer
}

func (s RampVuGenerator) Generate(param VuGeneratorParam) (err error) {

	return
}

type VuGeneratorParam struct {
	Type   consts.GeneratorType `json:"type,omitempty"`
	Target int                  `json:"target,omitempty"`
	Stages []VuGeneratorStage   `json:"stages,omitempty"`
}
type VuGeneratorStage struct {
	Duration int `json:"duration"`
	Target   int `json:"target"`
}
