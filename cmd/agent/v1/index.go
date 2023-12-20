package v1

import (
	"github.com/aaronchen2k/deeptest/cmd/agent/v1/router"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core"
	"github.com/kataras/iris/v12"
)

type IndexModule struct {
	HealthModule *router.HealthModule `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

// Party v1 模块
func (m *IndexModule) Party() core.WebModule {
	handler := func(v1 iris.Party) {}
	modules := []core.WebModule{
		m.HealthModule.Party(),
	}
	return core.NewModule(consts.ApiPathAgent, handler, modules...)
}
