package services

import (
	"github.com/shopspring/decimal"
	"resk/infra/base"
	"time"
)

var RedEnvelopeServiceInterface RedEnvelopeService

func GetRedEnvelopeServiceInterface() RedEnvelopeService {
	base.Check(RedEnvelopeServiceInterface)
	return RedEnvelopeServiceInterface
}

type RedEnvelopeService interface {
	// 发红包
	SendOut(dto RedEnvelopeSendOutDTO) (activity *RedEnvelopeActivity, err error)
	// 收红包
	Receive(dto RedEnvelopeReceiveDTO) (item *RedEnvelopeItemDTO, err error)
	// 退款
	Refund(envelopeNo string) (order *RedEnvelopeGoodsDTO)
	// 查询红包订单
	Get(envelopeNo string) (order *RedEnvelopeGoodsDTO)

	ListSended(userId string, page, size int) (orders []*RedEnvelopeGoodsDTO)
	ListReceived(userId string, page, size int) (items []*RedEnvelopeItemDTO)
	ListReceivable(page, size int) (orders []*RedEnvelopeGoodsDTO)
	ListItems(envelopeNo string) (items []*RedEnvelopeItemDTO)
}

type RedEnvelopeSendOutDTO struct {
	EnvelopeType int             `json:"envelope_type" varlidate:"required"`    // 红包类型：普通红包，碰运气红包
	UserName     string          `json:"user_name" varlidate:"required"`        // 用户名称
	UserId       string          `json:"user_id" varlidate:"required"`          // 用户编号, 红包所属用户
	Blessing     string          `json:"blessing"`                              // 祝福语
	Amount       decimal.Decimal `json:"amount" varlidate:"required,numeric"`   // 红包金额:普通红包指单个红包金额，碰运气红包指总金额
	Quantity     int             `json:"quantity" varlidate:"required,numeric"` // 红包总数量
}

func (r *RedEnvelopeSendOutDTO) ToGoods() *RedEnvelopeGoodsDTO {
	goods := &RedEnvelopeGoodsDTO{
		EnvelopeType: r.EnvelopeType,
		UserName:     r.UserName,
		UserId:       r.UserId,
		Blessing:     r.Blessing,
		Amount:       r.Amount,
		Quantity:     r.Quantity,
	}
	return goods
}

type RedEnvelopeReceiveDTO struct {
	EnvelopeNo   string `json:"envelope_no" varlidate:"required"`    // 红包编号,红包唯一标识
	RecvUserName string `json:"recv_user_name" varlidate:"required"` // 红包接收者用户名称
	RecvUserId   string `json:"recv_user_id" varlidate:"required"`   // 红包接收者用户编号
	AccountNo    string `json:"account_no"`
}

type RedEnvelopeActivity struct {
	RedEnvelopeGoodsDTO
	Link string `json:"link"` // 活动链接(发送给收红包的人)
}

func (r RedEnvelopeActivity) CopyTo(target *RedEnvelopeActivity) {
	target.Link = r.Link
	target.EnvelopeNo = r.EnvelopeNo
	target.EnvelopeType = r.EnvelopeType
	target.UserName = r.UserName
	target.UserId = r.UserId
	target.Blessing = r.Blessing
	target.Amount = r.Amount
	target.AmountOne = r.AmountOne
	target.Quantity = r.Quantity
	target.RemainAmount = r.RemainAmount
	target.RemainQuantity = r.RemainQuantity
	target.ExpiredAt = r.ExpiredAt
	target.Status = r.Status
	target.OrderType = r.OrderType
	target.PayStatus = r.PayStatus
	target.CreatedAt = r.CreatedAt
	target.UpdatedAt = r.UpdatedAt
}

type RedEnvelopeGoodsDTO struct {
	EnvelopeNo       string          `json:"envelope_no"`                          // 红包编号,红包唯一标识
	EnvelopeType     int             `json:"envelope_type" validate:"required"`    // 红包类型：普通红包，碰运气红包
	UserName         string          `json:"user_name" validate:"required"`        // 用户名称
	UserId           string          `json:"user_dd" validate:"required"`          // 用户编号, 红包所属用户
	Blessing         string          `json:"blessing"`                             // 祝福语
	Amount           decimal.Decimal `json:"amount" validate:"required,numeric"`   // 红包总金额
	AmountOne        decimal.Decimal `json:"amount_one"`                           // 单个红包金额，碰运气红包无效
	Quantity         int             `json:"quantity" validate:"required,numeric"` // 红包总数量
	RemainAmount     decimal.Decimal `json:"remain_amount"`                        // 红包剩余金额额
	RemainQuantity   int             `json:"remain_quantity"`                      // 红包剩余数量
	ExpiredAt        time.Time       `json:"expired_at" `                          // 过期时间
	Status           OrderStatus     `json:"status"`                               // 红包状态：0红包初始化，1启用，2失效
	OrderType        OrderType       `json:"order_type"`                           // 订单类型：发布单、退款单
	PayStatus        PayStatus       `json:"pay_status"`                           // 支付状态：未支付，支付中，已支付，支付失败
	CreatedAt        time.Time       `json:"created_at"`                           // 创建时间
	UpdatedAt        time.Time       `json:"updated_at"`                           // 更新时间
	AccountNo        string          `json:"account_no"`
	OriginEnvelopeNo string          `json:"origin_envelope_no"`
}

type RedEnvelopeItemDTO struct {
	ItemNo        string          `json:"item_no"`         // 红包订单详情编号
	EnvelopeNo    string          `json:"envelope_no"`     // 订单编号 红包编号,红包唯一标识
	RecvUserName  string          `json:"recv_user_name"`  // 红包接收者用户名称
	RecvUserId    string          `json:"recv_user_id"`    // 红包接收者用户编号
	RecvAccountNo string          `json:"recv_account_no"` // 红包接收者账户编号
	Amount        decimal.Decimal `json:"amount"`          // 收到金额
	Quantity      int             `json:"quantity"`        // 收到数量：对于收红包来说是1
	RemainAmount  decimal.Decimal `json:"remain_amount"`   // 收到后红包剩余金额
	PayStatus     int             `json:"pay_status"`      // 支付状态：未支付，支付中，已支付，支付失败
	CreatedAt     time.Time       `json:"created_at"`      // 创建时间
	UpdatedAt     time.Time       `json:"updated_at"`      // 更新时间
	IsLuckiest    bool            `json:"is_luckiest"`
	Desc          string          `json:"desc"`
}

func (r *RedEnvelopeItemDTO) CopyTo(item *RedEnvelopeItemDTO) {
	item.ItemNo = r.ItemNo
	item.EnvelopeNo = r.EnvelopeNo
	item.RecvUserName = r.RecvUserName
	item.RecvUserId = r.RecvUserId
	item.Amount = r.Amount
	item.Quantity = r.Quantity
	item.RemainAmount = r.RemainAmount
	item.RecvAccountNo = r.RecvAccountNo
	item.PayStatus = r.PayStatus
	item.CreatedAt = r.CreatedAt
	item.UpdatedAt = r.UpdatedAt
	item.Desc = r.Desc
}
