package services

const (
	DefaultBlessing   = "恭喜发财！"
	DefaultTimeFormat = "2006-01-02 15:04:05"
)

// 订单类型：发布单、退款单
type OrderType int

const (
	OrderSendOut OrderType = 1
	OrderRefund  OrderType = 2
)

// 支付状态：未支付，支付中，已支付，支付失败
// 退款状态：未退款，退款中，已退款，退款失败
type PayStatus int

const (
	UnPaid        PayStatus = 1
	Payments      PayStatus = 2
	Paid          PayStatus = 3
	PaymentFailed PayStatus = 4

	UnRefund     PayStatus = 5
	Refunding    PayStatus = 6
	Refunded     PayStatus = 7
	RefundFailed PayStatus = 8
)

// 红包订单状态：创建、发布、过期、失效
type OrderStatus int

const (
	OrderCreated                OrderStatus = 1
	OrderSendOuted              OrderStatus = 2
	OrderExpired                OrderStatus = 3
	OrderInvalid                OrderStatus = 4
	OrderExpiredRefundSucceeded OrderStatus = 5
	OrderExpiredRefundFailed    OrderStatus = 6
)

// 红包类型：普通红包、碰运气红包
type EnvelopeType int

const (
	CommonEnvelope EnvelopeType = 1
	LuckyEnvelope  EnvelopeType = 2
)

var EnvelopeTypes = map[EnvelopeType]string{
	CommonEnvelope: "普通红包",
	LuckyEnvelope:  "碰运气红包",
}
