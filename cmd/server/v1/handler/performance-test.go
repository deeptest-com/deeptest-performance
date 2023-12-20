package handler

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type PerformanceTestCtrl struct {
	PerformanceTestServices *service.PerformanceTestServices `inject:""`
}

func (c *PerformanceTestCtrl) Exec(ctx iris.Context) {
	if c.PerformanceTestServices.PerformanceServiceClient == nil {
		c.PerformanceTestServices.Connect(ctx)
	}

	c.PerformanceTestServices.Exec(ctx)

	data := iris.Map{}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data, Msg: _domain.NoErr.Msg})
}
