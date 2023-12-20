package agentServe

import (
	"fmt"
	middleware "github.com/aaronchen2k/deeptest/internal/pkg/core"
	"sync"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

// AgentServer 服务器
type AgentServer struct {
	app               *iris.Application
	modules           []middleware.WebModule
	idleConnClosed    chan struct{}
	addr              string
	timeFormat        string
	globalMiddlewares []context.Handler
	wg                sync.WaitGroup
	staticPath        string
	webPath           string
}

// InitRouter 初始化模块路由
func (s *AgentServer) InitRouter() error {
	s.app.UseRouter()

	s.initModule()
	err := s.app.Build()
	if err != nil {
		return fmt.Errorf("build router %w", err)
	}

	return nil
}

func (s *AgentServer) initModule() {
	if len(s.modules) > 0 {
		for _, mod := range s.modules {
			mod := mod
			s.wg.Add(1)
			go func(mod middleware.WebModule) {
				sub := s.app.PartyFunc(mod.RelativePath, mod.Handler)
				if len(mod.Modules) > 0 {
					for _, subModule := range mod.Modules {
						sub.PartyFunc(subModule.RelativePath, subModule.Handler)
					}
				}
				s.wg.Done()
			}(mod)
		}
		s.wg.Wait()
	}
}
