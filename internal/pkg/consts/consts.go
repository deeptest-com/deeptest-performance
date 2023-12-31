package consts

import "time"

const (
	App = "deeptest"

	ApiPathServer = "/api/v1"
	WsPathServer  = ApiPathServer + "/ws"
	ApiPathAgent  = "/api/v1"

	WsDefaultNamespace = "default"
	WsDefaultRoom      = "default_room"
	WsChatEvent        = "OnChat"

	SupportEmail = "chenqi@deeptest.com"

	EmailSmtpAddress = "smtp.exmail.qq.com"
	EmailSmtpPort    = 465
	EmailAccount     = "chenqi@deeptest.com"
	EmailPassword    = ""

	ExecTimeout = 12 * time.Hour
)

var (
	DirUpload     = "upload"
	HeaderOptions = []string{"Accept", "Accept-Encoding", "Accept-Language", "Connection", "Host", "Referer", "User-Agent", "Cache-Control", "Cookie", "Range"}
)
