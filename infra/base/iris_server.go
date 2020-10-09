package base

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	irisrecover "github.com/kataras/iris/middleware/recover"
	"github.com/sirupsen/logrus"
	"resk/infra"
	"time"
)

var irisApplication *iris.Application

func Iris() *iris.Application {
	Check(irisApplication)
	return irisApplication
}

type IrisServerStarter struct {
	infra.BaseStarter
}

func (s *IrisServerStarter) Init(infra.StarterContext) {
	// 创建iris application实例
	irisApplication = initIris()
	// 日志组件配置和扩展
	irisLoger := irisApplication.Logger()
	irisLoger.Install(logrus.StandardLogger())
}

func (s *IrisServerStarter) Start(ctx infra.StarterContext) {
	// 和logrus日志级别保持一致
	Iris().Logger().SetLevel(ctx.Props().GetDefault("log.level", "info"))

	// 把路由信息打印到控制台
	routers := irisApplication.GetRoutes()
	for _, r := range routers {
		logrus.Info(r.Trace())
	}

	// 启动iris
	port, _ := ctx.Props().Get("app.server.port")
	_ = irisApplication.Run(iris.Addr(":" + port))
}

func (s *IrisServerStarter) StartBlocking() bool {
	return true
}

func initIris() *iris.Application {
	app := iris.New()
	app.Use(irisrecover.New())
	// 主要中间件的配置 :recover，日志输出中间件的自定义
	cfg := logger.Config{
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
		Query:   true,
		Columns: false,
		LogFunc: func(
			endTime time.Time,
			latency time.Duration,
			status, ip,
			method, path string,
			message interface{},
			headerMessage interface{}) {
			app.Logger().Infof("| %s | %s | %s | %s | %s | %s | %s | %s",
				time.Now().Format("2006-01-02 15:04:05.000000"),
				latency.String(), status, ip, method, path, headerMessage, message)
		},
	}
	app.Use(logger.New(cfg))
	return app
}
