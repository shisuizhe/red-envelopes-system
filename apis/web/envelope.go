package web

import (
	"github.com/kataras/iris"
	"resk/infra"
	"resk/infra/base"
	"resk/services"
)

func init() {
	infra.RegisterApi(new(EnvelopeApi))
}

type EnvelopeApi struct {
	service services.RedEnvelopeService
}

func (e *EnvelopeApi) Init() {
	e.service = services.GetRedEnvelopeServiceInterface()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sending", e.sendOutHandler)
	groupRouter.Post("/receive", e.receiveHandler)
}

/*
{
	"envelope_type": 2,
	"user_name": "测试账户1",
	"user_id": "1fSJJYtPJBQFItcEV0NmvWgOjYi",
	"blessing": "",
	"amount": "10",
	"quantity": 5
}
*/

func (e *EnvelopeApi) sendOutHandler(ctx iris.Context) {
	dto := services.RedEnvelopeSendOutDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Response{
		Code: base.ResponseOK,
	}
	if err != nil {
		r.Code = base.ResponseRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	activity, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResponseInnerServerError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	r.Data = activity
	_, _ = ctx.JSON(r)
}

/*
{
	"envelope_no": "1fSLxfWalocjJXqWDddFfYyO4pm",
	"recv_user_name": "测试用户101",
	"recv_user_id": "1fSJJVWU3cAOgxpZvCkhuAAR8hx"
}
*/
func (e *EnvelopeApi) receiveHandler(ctx iris.Context) {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Response{
		Code: base.ResponseOK,
	}
	if err != nil {
		r.Code = base.ResponseRequestParamsError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	item, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResponseInnerServerError
		r.Message = err.Error()
		_, _ = ctx.JSON(&r)
		return
	}
	r.Data = item
	_, _ = ctx.JSON(r)
}
