package envelope

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"resk/infra/base"
	"resk/services"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.RedEnvelopeServiceInterface = new(redEnvelopeService)
	})
}

type redEnvelopeService struct{}

// 发红包
func (r *redEnvelopeService) SendOut(
	dto services.RedEnvelopeSendOutDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 验证
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	// 获取红包发送人的资金账户信息
	account := services.GetAccountServiceInterface().GetEnvelopeAccountByUserId(dto.UserId)

	if account == nil {
		return nil, errors.New("账户不存在")
	}

	// 如果账户余额小于要发的红包金额，则发送失败
	// envelopeAmount为要发的红包金额
	envelopeAmount := dto.Amount
	if dto.EnvelopeType == int(services.CommonEnvelope) {
		envelopeAmount = dto.Amount.Mul(decimal.NewFromInt(int64(dto.Quantity)))
	}
	if account.Balance.Cmp(envelopeAmount) < 0 {
		return nil, errors.New("账户余额不足，请充值再发")
	}

	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo
	goods.Amount = envelopeAmount

	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == int(services.CommonEnvelope) {
		goods.AmountOne = envelopeAmount.Div(decimal.NewFromInt(int64(dto.Quantity)))
	}
	// 发送红包
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if err != nil {
		logrus.Error(err)
	}
	return activity, err
}

// 收红包
func (r *redEnvelopeService) Receive(dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 参数校验
	if err = base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	// 获取当前收红包用户的账户信息
	a := services.GetAccountServiceInterface().GetEnvelopeAccountByUserId(dto.RecvUserId)
	if a == nil {
		return nil, errors.New("用户不存在")
	}
	dto.AccountNo = a.AccountNo
	// 进行尝试收红包
	gd := goodsDomain{}
	id := itemDomain{}
	item = id.GetByUserIdAndEnvelopeNo(dto.RecvUserId, dto.EnvelopeNo)
	if item != nil {
		return item, err
	}
	item, err = gd.Receive(context.Background(), dto)
	return item, err
}

func (r *redEnvelopeService) Refund(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	return nil
}

func (r *redEnvelopeService) Get(envelopeNo string) (order *services.RedEnvelopeGoodsDTO) {
	gd := goodsDomain{}
	po := gd.GetOne(envelopeNo)
	if po == nil {
		return nil
	}
	return po.ToDTO()
}

func (r *redEnvelopeService) ListSended(userId string, page, size int) (orders []*services.RedEnvelopeGoodsDTO) {
	gd := goodsDomain{}
	pos := gd.FindByUser(userId, page, size)
	orders = make([]*services.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, g := range pos {
		orders = append(orders, g.ToDTO())
	}
	return orders
}

func (r *redEnvelopeService) ListReceived(userId string, page, size int) (items []*services.RedEnvelopeItemDTO) {
	gd := goodsDomain{}
	pos := gd.ListReceived(userId, page, size)
	items = make([]*services.RedEnvelopeItemDTO, 0, len(pos))
	for _, p := range pos {
		items = append(items, p.ToDTO())
	}
	return
}

func (r *redEnvelopeService) ListReceivable(page, size int) (orders []*services.RedEnvelopeGoodsDTO) {
	gd := goodsDomain{}

	pos := gd.ListReceivable(page, size)
	orders = make([]*services.RedEnvelopeGoodsDTO, 0, len(pos))
	for _, p := range pos {
		if p.RemainQuantity > 0 {
			orders = append(orders, p.ToDTO())
		}
	}
	return
}

func (r *redEnvelopeService) ListItems(envelopeNo string) (items []*services.RedEnvelopeItemDTO) {
	id := itemDomain{}
	items, err := id.FindItems(envelopeNo)
	if err != nil {
		return nil
	}
	return items
}
