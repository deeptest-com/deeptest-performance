package _domain

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/kataras/iris/v12"
)

type WsResp struct {
	Uuid     string               `json:"uuid"`
	Category consts.WsMsgCategory `json:"category"`

	Msg  string      `json:"msg"`
	Info iris.Map    `json:"info,omitempty"`
	Data interface{} `json:"data,omitempty"`
}

type MqMsg struct {
	Namespace string `json:"namespace"`
	Room      string `json:"room"`
	Event     string `json:"event"`
	Content   string `json:"content"`
}
