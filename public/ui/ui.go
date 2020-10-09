package views

import (
	"fmt"
	"github.com/kataras/iris"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"resk/core/users"
	"resk/infra"
	"resk/infra/base"
	"resk/services"
	"runtime"
	"strconv"
	"time"
)

func init() {
	infra.RegisterApi(&MobileView{})
}

type MobileView struct {
	UserService *users.UserService
	groupRouter iris.Party
}

func (v *MobileView) Init() {
	v.UserService = new(users.UserService)
	_, f, _, _ := runtime.Caller(0)
	dir := filepath.Dir(f)
	app := base.Iris()
	views := iris.HTML(dir, ".html")
	views.Layout("layouts/layout.html")
	views.Reload(true) // reload templates on each request (development mode)
	app.RegisterView(views)
	app.Favicon(filepath.Join(dir, "/favicon.ico"))
	//contextPath := ""
	app.StaticWeb("/public/static", filepath.Join(dir, "../static"))
	app.StaticWeb("/public/ui", dir)
	v.groupRouter = app.Party("/envelope")
	v.groupRouter.Use(func(ctx iris.Context) {
		userId := ctx.GetCookie("userId")
		log.Info(userId)
		if userId == "" {
			ctx.Redirect("/login")
		} else {
			ctx.Next()
		}
	})
	app.Any("/", v.indexHandler)
	app.Any("", v.indexHandler)
	//登入登出
	app.Get("/login", v.loginHandler)
	app.Post("/login", v.loginSubmitHandler)
	app.Get("/logout", v.logoutHandler)
	//我发的红包列表
	v.groupRouter.Get("/home", v.homeHandler)
	//我抢到的红包列表
	v.groupRouter.Get("/recvd/list", v.receivedListHandler)
	//红包记录
	v.groupRouter.Get("/list", v.listHandler)
	//红包详情
	v.groupRouter.Get("/details", v.detailsHandler)
	//可抢红包
	v.groupRouter.Get("/rev/home", v.receiveHomeHandler)
	//抢红包
	v.groupRouter.Get("/recd", v.receiveSubmitHandler)
	//发红包
	v.groupRouter.Get("/sending", v.sendingHandler)
	v.groupRouter.Post("/sending", v.sendingSubmitHandler)
}
func (v *MobileView) loginSubmitHandler(ctx iris.Context) {
	form := UserForm{}
	err := ctx.ReadForm(&form)
	if err != nil {
		log.Error(err)
	}
	if form.Mobile == "" {
		ctx.ViewData("msg", "用户名称不能为空！")
		ctx.View("index.html")
		return
	}
	if form.Username == "" {
		ctx.ViewData("msg", "用户密码不能为空！")
		ctx.View("index.html")
		return
	}
	user := v.UserService.Login(form.Mobile, form.Username)
	if user == nil {
		ctx.ViewData("msg", "系统出错了！")
		ctx.View("index.html")
		return
	}

	ctx.SetCookieKV("userId", user.UserId, iris.CookieExpires(1*time.Hour))
	ctx.SetCookieKV("username", user.Username, iris.CookieExpires(1*time.Hour))

	ctx.Redirect("/envelope/home")
}
func (v *MobileView) indexHandler(ctx iris.Context) {
	ctx.View("index.html")
}
func (v *MobileView) homeHandler(ctx iris.Context) {

	userId := ctx.GetCookie("userId")
	es := services.GetRedEnvelopeServiceInterface()
	orders := es.ListSended(userId, 0, 200)
	ctx.ViewData("orders", orders)
	ctx.ViewData("format", services.DefaultTimeFormat)
	_ = ctx.View("home.html")

}

func (v *MobileView) logoutHandler(ctx iris.Context) {
	ctx.RemoveCookie("userId")
	ctx.RemoveCookie("username")

	ctx.View("index.html")
}
func (v *MobileView) loginHandler(ctx iris.Context) {
	ctx.View("index.html")

}

//我抢到的红包列表：recvd_list.html /recvd/list
func (v *MobileView) receivedListHandler(ctx iris.Context) {
	userId := ctx.GetCookie("userId")
	es := services.GetRedEnvelopeServiceInterface()
	items := es.ListReceived(userId, 0, 100)
	ctx.ViewData("items", items)
	ctx.ViewData("format", services.DefaultTimeFormat)

	ctx.View("recvd_list.html")
}

//红包记录：re_one.html /list
func (v *MobileView) listHandler(ctx iris.Context) {
	envelopeNo := ctx.URLParamTrim("id")
	es := services.GetRedEnvelopeServiceInterface()
	order := es.Get(envelopeNo)
	if order != nil {

		items := es.ListItems(envelopeNo)
		totalAmount := decimal.NewFromFloat(0)
		t1, t2 := time.Unix(int64(0), int64(0)), time.Unix(int64(0), int64(0))

		for i, v := range items {
			if i == 0 {
				t1 = v.CreatedAt
				t2 = v.CreatedAt
			} else {
				if t1.After(v.CreatedAt) {
					t1 = v.CreatedAt
				}
				if t2.Before(v.CreatedAt) {
					t2 = v.CreatedAt
				}
			}

			totalAmount = totalAmount.Add(v.Amount)
			if order.RemainQuantity > 0 {
				v.IsLuckiest = false
			}
		}
		ctx.ViewData("items", items)
		ctx.ViewData("size", len(items))
		ctx.ViewData("totalAmount", totalAmount)
		seconds := t2.UnixNano() - t1.UnixNano()
		h := seconds / int64(time.Hour)

		msg := ""
		if h > 0 {
			msg += strconv.Itoa(int(h)) + "小时"
			seconds -= h * int64(time.Hour)
		}
		m := seconds / int64(time.Minute)

		if m > 0 {
			msg += strconv.Itoa(int(m)) + "分钟"
			seconds -= m * int64(time.Minute)
		}
		s := seconds / int64(time.Second)
		if s > 0 {
			msg += strconv.Itoa(int(s)) + "秒"
		}
		if msg == "" {
			msg = "0秒"
		}
		fmt.Println(t1, t2, seconds)

		ctx.ViewData("timeMsg", msg)
		ctx.ViewData("isReceived", len(items) == order.Quantity)
		ctx.ViewData("remainQuantity", order.Quantity-len(items))
	}

	ctx.ViewData("order", order)
	ctx.ViewData("isLucky", order.EnvelopeType == 2)
	ctx.ViewData("hasOrder", order != nil)

	ctx.ViewData("format", services.DefaultTimeFormat)
	ctx.View("re_one.html")
}

//红包详情：re_details.html /details
func (v *MobileView) detailsHandler(ctx iris.Context) {
	envelopeNo := ctx.URLParamTrim("id")
	es := services.GetRedEnvelopeServiceInterface()
	order := es.Get(envelopeNo)
	ctx.ViewData("order", order)
	ctx.ViewData("hasOrder", order != nil)
	ctx.ViewData("format", services.DefaultTimeFormat)
	ctx.View("re_details.html")
}

//可抢红包：rev_home.html /rev/home
func (v *MobileView) receiveHomeHandler(ctx iris.Context) {
	es := services.GetRedEnvelopeServiceInterface()
	orders := es.ListReceivable(0, 200)
	ctx.ViewData("orders", orders)
	ctx.ViewData("hasOrders", len(orders) > 0)
	ctx.ViewData("format", services.DefaultTimeFormat)
	ctx.View("rev_home.html")
}

//抢红包
func (v *MobileView) receiveSubmitHandler(ctx iris.Context) {
	envelopeNo := ctx.URLParamTrim("id")
	userId := ctx.GetCookie("userId")
	username := ctx.GetCookie("username")

	es := services.GetRedEnvelopeServiceInterface()
	dto := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   envelopeNo,
		RecvUserId:   userId,
		RecvUserName: username,
	}
	item, err := es.Receive(dto)
	msg := ""
	if err == nil {
		ctx.ViewData("hasReceived", true)
	} else {
		ctx.ViewData("hasReceived", false)
		msg = err.Error()
	}
	order := es.Get(envelopeNo)
	ctx.ViewData("order", order)
	ctx.ViewData("item", item)
	ctx.ViewData("hasOrder", order != nil)
	ctx.ViewData("format", services.DefaultTimeFormat)
	ctx.ViewData("msg", msg)
	ctx.View("recd.html")
}

//发红包
func (v *MobileView) sendingHandler(ctx iris.Context) {
	ctx.View("sending.html")
}

//发红包
func (v *MobileView) sendingSubmitHandler(ctx iris.Context) {
	form := RedEnvelopeSendingForm{}
	err := ctx.ReadForm(&form)
	if err != nil {
		log.Error(err)
		ctx.ViewData("msg", "读取数据出错")
		ctx.View("sending.html")
		return
	}
	userId := ctx.GetCookie("userId")
	username := ctx.GetCookie("username")
	dto := services.RedEnvelopeSendOutDTO{
		UserId:       userId,
		UserName:     username,
		Amount:       form.Amount,
		Blessing:     form.Blessing,
		EnvelopeType: form.EnvelopeType,
		Quantity:     form.Quantity,
	}

	service := services.GetRedEnvelopeServiceInterface()
	activity, err := service.SendOut(dto)
	if err != nil {
		log.Error(err)
		ctx.ViewData("msg", "发红包失败，系统出错")
		ctx.View("sending.html")
		return
	}
	ctx.ViewData("activity", activity)
	ctx.Redirect("/envelope/home")
}
