package views

import (
	"github.com/kataras/iris"
	"imooc.com/resk/core/users"
	"imooc.com/resk/infra"
	"imooc.com/resk/infra/base"
	"path/filepath"
	"runtime"
)

func init() {
	infra.RegisterApi(&View{})
}

type View struct {
	UserService *users.UserService
	groupRouter iris.Party
}

func (v *View) Init() {
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	app := base.Iris()
	//views := iris.HTML(dir, ".html")
	//views.Layout("layouts/layout.html")
	//views.Reload(true) // reload templates on each request (development mode)
	//app.RegisterView(views)
	app.Favicon(filepath.Join(dir, "/favicon.ico"))
	//contextPath := ""
	app.StaticWeb("/public/static", filepath.Join(dir, "../static"))
	app.StaticWeb("/public/views", dir)
	v.groupRouter = app.Party("/v1/envelope")
	v.groupRouter.Layout("views/layouts/layout.html")
	v.index()
	v.SendingRedEnvelopeIndex()

}
func (v *View) index() {
	base.Iris().Get("/index", func(ctx iris.Context) {
		ctx.View("views/index.html")
	})
	base.Iris().Get("/home", func(ctx iris.Context) {
		ctx.View("views/index.html")
	})
}

func (v *View) SendingRedEnvelopeIndex() {
	v.groupRouter.Get("/Sending", func(ctx iris.Context) {
		ctx.View("app.html")
	})
}
