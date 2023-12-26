package serverServe

import (
	stdContext "context"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	middleware "github.com/aaronchen2k/deeptest/internal/pkg/core"
	commUtils "github.com/aaronchen2k/deeptest/internal/pkg/utils"
	_i118Utils "github.com/aaronchen2k/deeptest/pkg/lib/i118"
	"github.com/facebookgo/inject"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kataras/iris/v12/context"
	"github.com/snowlyg/helper/tests"

	_ "github.com/aaronchen2k/deeptest/cmd/server/docs"
	"github.com/kataras/iris/v12"
)

var client *tests.Client

// Start 初始化web服务
func Start() {
	inits()

	idleConnClosed := make(chan struct{})
	irisApp := createIrisApp(&idleConnClosed)

	server := &WebServer{
		app:               irisApp,
		addr:              "0.0.0.0:8087",
		timeFormat:        "2006-01-02 15:04:05",
		idleConnClosed:    idleConnClosed,
		globalMiddlewares: []context.Handler{middleware.Error()},
	}

	server.InjectModule()
	server.Start()
}

func inits() {
	consts.WorkDir = commUtils.GetWorkDir()

	_i118Utils.Init(consts.Language, "")
}

func createIrisApp(idleConnClosed *chan struct{}) (irisApp *iris.Application) {
	irisApp = iris.New()
	irisApp.Validator = validator.New()
	irisApp.Logger().SetLevel("debug")

	iris.RegisterOnInterrupt(func() {
		timeout := 10 * time.Second
		ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
		defer cancel()
		irisApp.Shutdown(ctx) // close all hosts

		close(*idleConnClosed)
	})

	return
}

func (webServer *WebServer) Start() {
	webServer.app.UseGlobal(webServer.globalMiddlewares...)
	err := webServer.InitRouter()
	if err != nil {
		fmt.Printf("初始化路由错误： %v\n", err)
		panic(err)
	}
	// 添加上传文件路径
	webServer.app.Listen(
		webServer.addr,
		iris.WithoutInterruptHandler,
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
		iris.WithTimeFormat(webServer.timeFormat),
	)
	<-webServer.idleConnClosed
}

// GetAddr 获取web服务地址
func (webServer *WebServer) GetAddr() string {
	return webServer.addr
}

// AddModule 添加模块
func (webServer *WebServer) AddModule(module ...middleware.WebModule) {
	webServer.modules = append(webServer.modules, module...)
}

// GetModules 获取模块
func (webServer *WebServer) GetModules() []middleware.WebModule {
	return webServer.modules
}

// Init 加载模块
func (webServer *WebServer) InjectModule() {
	var g inject.Graph
	g.Logger = logrus.StandardLogger()

	indexModule := v1.NewIndexModule()

	// inject objects
	if err := g.Provide(
		&inject.Object{Value: indexModule},
	); err != nil {
		logrus.Fatalf("provide usecase objects to the Graph: %v", err)
	}
	err := g.Populate()
	if err != nil {
		logrus.Fatalf("populate the incomplete Objects: %v", err)
	}

	webServer.AddModule(indexModule.ApiParty())
}
