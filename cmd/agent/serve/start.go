package agentServe

import (
	stdContext "context"
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	middleware "github.com/aaronchen2k/deeptest/internal/pkg/core"
	commUtils "github.com/aaronchen2k/deeptest/internal/pkg/utils"
	_i118Utils "github.com/aaronchen2k/deeptest/pkg/lib/i118"
	"sync"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/str"
	"github.com/snowlyg/helper/tests"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
)

var client *tests.Client

// Init 初始化web服务
func Init() *AgentServer {
	consts.RunFrom = consts.FromAgent
	consts.WorkDir = commUtils.GetWorkDir()

	_i118Utils.Init(consts.Language, "")

	app := iris.New()
	app.Validator = validator.New()
	app.Logger().SetLevel("debug")
	idleConnClosed := make(chan struct{})
	iris.RegisterOnInterrupt(func() { //优雅退出
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()

		app.Shutdown(ctx)
		close(idleConnClosed)
	})

	// init grpc
	mvc.New(app)

	return &AgentServer{
		app:               app,
		addr:              "0.0.0.0:8088",
		timeFormat:        "2006-01-02 15:04:05",
		idleConnClosed:    idleConnClosed,
		globalMiddlewares: []context.Handler{middleware.Error()},
	}
}

// AddModule 添加模块
func (s *AgentServer) AddModule(module ...middleware.WebModule) {
	s.modules = append(s.modules, module...)
}

// GetModules 获取模块
func (s *AgentServer) GetModules() []middleware.WebModule {
	return s.modules
}

// GetTestAuth 获取测试验证客户端
func (s *AgentServer) GetTestAuth(t *testing.T) *tests.Client {
	var once sync.Once
	once.Do(
		func() {
			client = tests.New(str.Join("http://", s.addr), t, s.app)
			if client == nil {
				t.Fatalf("client is nil")
			}
		},
	)

	return client
}

// Init 启动web服务
func (s *AgentServer) Start() {
	s.app.UseGlobal(s.globalMiddlewares...)
	err := s.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}
	// 添加上传文件路径
	s.app.Listen(
		s.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(s.timeFormat),
	)
	<-s.idleConnClosed
}
