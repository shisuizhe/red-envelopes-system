package envelope

import (
	"context"
	"github.com/tietang/dbx"
	"path"
	"resk/core/account"
	"resk/infra/base"
	"resk/services"
)

// 发红包业务领域代码
func (g *goodsDomain) SendOut(
	goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	// 创建红包商品
	g.Create(goods)
	// 创建活动
	activity = new(services.RedEnvelopeActivity)
	// 红包链接 http://域名/v1/envelope/link/{id}
	link := base.GetEnvelopeActivityLink()
	dm := base.GetEnvelopeDomain()
	activity.Link = path.Join(dm, link, g.EnvelopeNo)

	ad := account.NewAccountDomain()

	err = base.Tx(func(runner *dbx.TxRunner) (err error) {
		// 事务：保存红包商品和红包金额的支付必须要保证全部成功或者全部失败
		ctx := base.WithValueContext(context.Background(), runner)
		// 保存红包商品
		id, err := g.Save(ctx)
		if id <= 0 || err != nil {
			return err
		}
		// 红包金额支付
		// 1.需要红包中间商的红包资金账户，定义在配置文件中，事先初始化到资金账户表中
		// 2.从红包发送人的资金账户中扣减红包金额

		// 发红包账户
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserId:    goods.UserId,
			UserName:  goods.UserName,
		}
		// target为系统账户
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserId:    systemAccount.UserId,
			UserName:  systemAccount.UserName,
		}
		// 转账transfer
		transfer := services.AccountTransferDTO{
			TradeNo:     g.EnvelopeNo,
			TradeBody:   body,
			TradeTarget: target,
			Amount:      g.Amount,
			ChangeType:  services.EnvelopeOutgoning,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		status, err := ad.TransferWithContextTx(ctx, transfer)
		if status == services.TransferStatusSuccess {
			return nil
		}
		// 3.将扣减的红包总金额转入红包中间商的红包资金账户
		transfer = services.AccountTransferDTO{
			TradeNo:     g.EnvelopeNo,
			TradeBody:   target,
			TradeTarget: body,
			Amount:      g.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		status, err = ad.TransferWithContextTx(ctx, transfer)
		if status == services.TransferStatusSuccess {
			return nil
		} else {
			return err
		}
	})
	if err != nil {
		return nil, err
	}
	// 扣减金额没有问题，返回活动
	activity.RedEnvelopeGoodsDTO = *g.RedEnvelopeGoods.ToDTO()
	return activity, err
}
