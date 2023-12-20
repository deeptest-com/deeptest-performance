package router

import (
	"github.com/aaronchen2k/deeptest/cmd/agent/v1/handler"
	middleware "github.com/aaronchen2k/deeptest/internal/pkg/core"
	"github.com/kataras/iris/v12"
)

type HealthModule struct {
	HealthCtrl *handler.HealthCtrl `inject:""`
}

// Party
func (m *HealthModule) Party() middleware.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", m.HealthCtrl.Get).Name = "健康检查"
	}
	return middleware.NewModule("/health", handler)
}
