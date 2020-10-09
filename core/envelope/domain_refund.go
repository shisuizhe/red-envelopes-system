package envelope

import (
	"context"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
)

const (
	pageSize = 100
)

type ExpiredEnvelopeDomain struct {
	expriedGoods []RedEnvelopeGoods
	offest       int
}

// 查询出过期红包
func (d *ExpiredEnvelopeDomain) Next() (ok bool) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		d.expriedGoods = dao.FindExpired(d.offest, pageSize)
		if len(d.expriedGoods) > 0 {
			d.offest += len(d.expriedGoods)
			ok = true
		}
		return nil
	})
	return ok
}

func (d *ExpiredEnvelopeDomain) Expired() (err error) {
	for d.Next() {
		for _, goods := range d.expriedGoods {
			logrus.Debugf("过期红包退款开始：%+v", goods)
			err := d.ExpiredOne(goods)
			if err != nil {
				logrus.Error(err)
			}
			logrus.Debugf("过期红包退款结束：%+v", goods)
		}
	}
	return err
}

// 发起退款流程
func (d *ExpiredEnvelopeDomain) ExpiredOne(goods RedEnvelopeGoods) (err error) {
	// 创建一个退款订单
	refund := goods
	refund.OrderType = services.OrderRefund
	refund.RemainAmount = goods.RemainAmount.Mul(decimal.NewFromFloat(-1))
	refund.RemainQuantity = -goods.RemainQuantity
	refund.Status = services.OrderExpired
	refund.PayStatus = services.Refunding
	refund.OriginEnvelopeNo = goods.EnvelopeNo
	refund.EnvelopeNo = ""

	domain := goodsDomain{RedEnvelopeGoods: refund}
	domain.createEnvelopeNo()

	err = base.Tx(func(runner *dbx.TxRunner) error {
		txCtx := base.WithValueContext(context.Background(), runner)
		id, err := domain.Save(txCtx)
		if err != nil || id == 0 {
			return errors.New("创建退款订单失败")
		}
		// 修改原订单状态
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpired)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	// 调用资金账户接口进行转账
	systemAccount := base.GetSystemAccount()
	account := services.GetAccountServiceInterface().GetEnvelopeAccountByUserId(goods.UserId)
	if account == nil {
		return errors.New("没有找到该用户的资金红包账户：" + goods.UserId)
	}
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		UserName:  systemAccount.UserName,
	}
	target := services.TradeParticipator{
		AccountNo: account.AccountNo,
		UserId:    account.UserId,
		UserName:  account.UserName,
	}
	transfer := services.AccountTransferDTO{
		TradeNo:     domain.RedEnvelopeGoods.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      goods.RemainAmount,
		ChangeType:  services.EnvelopeExpiredRefund,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包过期退款：" + goods.EnvelopeNo,
	}
	status, err := services.GetAccountServiceInterface().Transfer(transfer)
	if status != services.TransferStatusSuccess {
		return err
	}

	// 修改订单状态
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		// 修改原订单状态
		rows, err := dao.UpdateOrderStatus(goods.EnvelopeNo, services.OrderExpiredRefundSucceeded)
		if err != nil || rows == 0 {
			return errors.New("更新原订单状态失败")
		}
		// 修改退款订单状态
		rows, err = dao.UpdateOrderStatus(refund.EnvelopeNo, services.OrderExpiredRefundSucceeded)
		if err != nil || rows == 0 {
			return errors.New("更新退款订单状态失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return err
	}
	return nil
}
