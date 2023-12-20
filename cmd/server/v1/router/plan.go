package router

import (
	"github.com/aaronchen2k/deeptest/cmd/server/v1/handler"
	core "github.com/aaronchen2k/deeptest/internal/pkg/core"
	"github.com/kataras/iris/v12"
)

type PlanModule struct {
	PlanCtrl *handler.PlanCtrl `inject:""`
}

// Party 计划
func (m *PlanModule) Party() core.WebModule {
	handler := func(index iris.Party) {
		index.Get("/{id:uint}", m.PlanCtrl.Get).Name = "详情"
	}

	return core.NewModule("/plans", handler)
}
