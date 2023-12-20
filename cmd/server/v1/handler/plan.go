package handler

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type PlanCtrl struct {
	PlanService *service.PlanService `inject:""`
}

func (c *PlanCtrl) Get(ctx iris.Context) {
	data := iris.Map{}
	ctx.JSON(_domain.Response{Code: _domain.NoErr.Code, Data: data, Msg: _domain.NoErr.Msg})
}
