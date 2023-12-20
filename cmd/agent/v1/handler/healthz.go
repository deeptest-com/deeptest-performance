package handler

import (
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	"github.com/kataras/iris/v12"
)

type HealthCtrl struct {
}

func (c *HealthCtrl) Get(ctx iris.Context) {
	ctx.JSON(_domain.Response{Code: 200, Msg: "health"})
}
