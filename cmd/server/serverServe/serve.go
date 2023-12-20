package serverServe

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/core"
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// WebServer 服务器
type WebServer struct {
	app               *iris.Application
	modules           []core.WebModule
	idleConnClosed    chan struct{}
	addr              string
	timeFormat        string
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
	staticPath        string
	webPath           string
}

// InitRouter 初始化模块路由
func (webServer *WebServer) InitRouter() error {
	webServer.app.UseRouter()

	{
		webServer.initModule()

		err := webServer.app.Build()
		if err != nil {
			return fmt.Errorf("build router %w", err)
		}

		return nil
	}
}

func (webServer *WebServer) initModule() {
	if len(webServer.modules) > 0 {
		for _, mod := range webServer.modules {
			mod := mod
			webServer.wg.Add(1)
			func(mod core.WebModule) {
				sub := webServer.app.PartyFunc(mod.RelativePath, mod.Handler)
				if len(mod.Modules) > 0 {
					for _, subModule := range mod.Modules {
						sub.PartyFunc(subModule.RelativePath, subModule.Handler)
					}
				}
				webServer.wg.Done()
			}(mod)
		}
		webServer.wg.Wait()
	}
}
