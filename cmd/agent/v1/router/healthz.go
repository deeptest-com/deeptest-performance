package router

import (
	"github.com/aaronchen2k/deeptest/cmd/agent/v1/handler"
	middleware "github.com/aaronchen2k/deeptest/internal/pkg/core"
	"github.com/kataras/iris/v12"
)

type HealthzModule struct {
	HealthzCtrl *handler.HealthCtrl `inject:""`
}

func NewHealthzModule() *HealthzModule {
	return &HealthzModule{}
}

// Party
func (m *HealthzModule) Party() middleware.WebModule {
	handler := func(index iris.Party) {
		index.Get("/", m.HealthzCtrl.Get).Name = "健康检查"
	}
	return middleware.NewModule("/healthz", handler)
}
