package envelope

import (
	"context"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
	"time"
)

type goodsDomain struct {
	RedEnvelopeGoods
	item itemDomain
}

// 生成一个红包编号
func (g *goodsDomain) createEnvelopeNo() {
	g.EnvelopeNo = ksuid.New().Next().String()
}

// 创建一个红包商品对象
func (g *goodsDomain) Create(goods services.RedEnvelopeGoodsDTO) {
	g.FromDTO(&goods)
	g.RemainQuantity = goods.Quantity
	g.UserName.Valid = true
	g.Blessing.Valid = true
	if g.EnvelopeType == int(services.CommonEnvelope) {
		g.Amount = goods.AmountOne.Mul(
			decimal.NewFromFloat(float64(goods.Quantity)))
	}
	if g.EnvelopeType == int(services.LuckyEnvelope) {
		g.AmountOne = decimal.NewFromFloat(0)
	}
	g.RemainAmount = g.Amount
	// 过期时间
	g.ExpiredAt = time.Now().Add(24 * time.Hour)

	g.Status = services.OrderCreated
	g.createEnvelopeNo()
}

// 保存到红包商品表
func (g *goodsDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		id, err = dao.Insert(&g.RedEnvelopeGoods)
		return err
	})
	return id, err
}

// 创建并保存红包商品
func (g *goodsDomain) CreateAndSave(ctx context.Context, goods services.RedEnvelopeGoodsDTO) (id int64, err error) {
	g.Create(goods)
	return g.Save(ctx)
}

// 查询红包商品信息
func (g *goodsDomain) Get(envelopeNo string) (goods *RedEnvelopeGoods) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
