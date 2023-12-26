package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	core "github.com/aaronchen2k/deeptest/internal/pkg/core"
	"github.com/kataras/iris/v12"
)

type PerformanceTestModule struct {
	PerformanceTestCtrl *handler.PerformanceTestCtrl `inject:""`
}

func (m *PerformanceTestModule) Party() core.WebModule {
	handler := func(index iris.Party) {
		index.Post("/exec", m.PerformanceTestCtrl.Exec).Name = "详情"
	}

	return core.NewModule("/performanceTest", handler)
}
