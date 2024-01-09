package v1

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/router"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/core"
	"github.com/kataras/iris/v12"
)

type IndexModule struct {
	PerformanceTestModule *router.PerformanceTestModule `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

func (m *IndexModule) ApiParty() core.WebModule {
	handler := func(v1 iris.Party) {}
	modules := []core.WebModule{
		m.PerformanceTestModule.Party(),
	}

	return core.NewModule(consts.ApiPathServer, handler, modules...)
}
