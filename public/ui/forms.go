package views

import "github.com/shopspring/decimal"

type UserForm struct {
	UserId   string `form:"user_id"`
	Mobile   string `form:"mobile"`
	Username string `form:"username"`
}

type RedEnvelopeSendingForm struct {
	EnvelopeType int             `form:"envelopeType"` //红包类型：普通红包，碰运气红包
	Blessing     string          `form:"blessing"`     //祝福语
	Amount       decimal.Decimal `form:"amount"`       //红包金额:普通红包指单个红包金额，碰运气红包指总金额
	Quantity     int             `form:"quantity"`     //红包总数量
}
