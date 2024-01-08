package handler

import (
	"encoding/json"
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	websocketHelper "github.com/aaronchen2k/deeptest/internal/pkg/websocket"
	"github.com/aaronchen2k/deeptest/internal/server/modules/service"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	_i118Utils "github.com/aaronchen2k/deeptest/pkg/lib/i118"
	_logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/kataras/iris/v12/websocket"
)

const (
	result = "result"
	outPut = "output"
)

var (
	ch chan int
)

type WebSocketCtrl struct {
	Namespace         string
	*websocket.NSConn `stateless:"true"`

	PerformanceTestServices *service.PerformanceTestServices `inject:""`
}

func NewWebsocketCtrl() *WebSocketCtrl {
	inst := &WebSocketCtrl{Namespace: consts.WsDefaultNamespace}
	return inst
}

func (c *WebSocketCtrl) OnNamespaceConnected(wsMsg websocket.Message) error {
	websocketHelper.SetConn(c.Conn)

	_logUtils.Infof(_i118Utils.Sprintf("ws_namespace_connected", c.Conn.ID(), wsMsg.Room))

	resp := _domain.WsResp{Msg: "from server: connected to websocket"}
	bytes, _ := json.Marshal(resp)
	mqData := _domain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	websocketHelper.PubWsMsg(mqData)

	return nil
}

// OnNamespaceDisconnect
// This will call the "OnVisit" event on all clients, except the current one,
// it can't because it's left but for any case use this type of design.
func (c *WebSocketCtrl) OnNamespaceDisconnect(wsMsg websocket.Message) error {
	_logUtils.Infof(_i118Utils.Sprintf("ws_namespace_disconnected", c.Conn.ID()))

	resp := _domain.WsResp{Msg: "from server: disconnected to websocket"}
	bytes, _ := json.Marshal(resp)
	mqData := _domain.MqMsg{Namespace: wsMsg.Namespace, Room: wsMsg.Room, Event: wsMsg.Event, Content: string(bytes)}
	websocketHelper.PubWsMsg(mqData)

	return nil
}

// OnChat This will call the "OnVisit" event on all clients, including the current one, with the 'newCount' variable.
func (c *WebSocketCtrl) OnChat(wsMsg websocket.Message) (err error) {
	ctx := websocket.GetContext(c.Conn)

	_logUtils.Infof("WebSocket OnChat: remote address=%s, room=%s, msg=%s", ctx.RemoteAddr(), wsMsg.Room, string(wsMsg.Body))

	req := serverDomain.WsReq{}
	err = json.Unmarshal(wsMsg.Body, &req)
	if err != nil {
		websocketHelper.SendExecMsg("exec failed", err, "performance testing", "NO_UUID", &wsMsg)
		return
	}

	if req.Act == consts.Init {
		return
	}

	if req.Act == consts.ExecPlan {
		err = c.PerformanceTestServices.ExecStart(req.PlanExecReq, &wsMsg)

	} else if req.Act == consts.ExecStop {
		err = c.PerformanceTestServices.ExecStop(req.PlanExecReq, &wsMsg)

	}

	return
}
