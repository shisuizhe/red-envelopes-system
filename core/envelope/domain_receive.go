package envelope

import (
	"context"
	"database/sql"
	"errors"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"resk/core/account"
	"resk/infra/algo"
	"resk/infra/base"
	"resk/services"
)

// 收红包业务领域代码
func (g *goodsDomain) Receive(ctx context.Context, dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	// 1.创建收红包的订单明细
	g.preCreateItem(dto)
	// 2.查询出当前红包的剩余数量和剩余金额信息
	goods := g.Get(dto.EnvelopeNo)
	// 3.校验剩余红包和剩余金额
	// - 如果没有剩余，直接返回无可用红包金额
	if goods.RemainQuantity <= 0 || goods.RemainAmount.Cmp(decimal.NewFromFloat(0)) <= 0 {
		return nil, errors.New("没有足够的红包和金额了")
	}
	// 4.使用红包算法计算红包金额
	nextAmount := g.nextAmount(goods)
	err = base.Tx(func(runner *dbx.TxRunner) error {
		// 5.使用乐观锁更新语句，尝试更新剩余数量和剩余金额
		// 如果更新失败，也就是返回0，表示无可用红包数量和金额，抢红包失败
		dao := RedEnvelopeGoodsDao{runner: runner}
		rows, err := dao.UpdateBalance(goods.EnvelopeNo, nextAmount)
		if rows <= 0 || err != nil {
			return errors.New("没有足够的红包和金额了")
		}
		// 如果更新成功，也就是返回1，表示抢到红包
		// 6.保存订单明细数据
		g.item.Quantity = 1
		g.item.PayStatus = int(services.Paid)
		g.item.AccountNo = dto.AccountNo
		g.item.RemainAmount = goods.RemainAmount.Sub(nextAmount)
		g.item.Amount = nextAmount
		g.item.Desc = "抢到" + nextAmount.String() + "元"
		txCtx := base.WithValueContext(ctx, runner)
		_, err = g.item.Save(txCtx)
		if err != nil {
			return err
		}
		// 7.将抢到的红包金额从系统红包中间账户转入当前用户的资金账户
		status, err := g.transfer(txCtx, dto)
		if status == services.TransferStatusSuccess {
			return nil
		}
		return err
	})
	return g.item.ToDTO(), err
}

// 创建收红包的订单明细
func (g *goodsDomain) preCreateItem(dto services.RedEnvelopeReceiveDTO) {
	g.item.AccountNo = dto.AccountNo
	g.item.EnvelopeNo = dto.EnvelopeNo
	g.item.RecvUserId = dto.RecvUserId
	g.item.RecvUserName = sql.NullString{Valid: true, String: dto.RecvUserName}
	g.item.createItemNo()
}

var multiple = decimal.NewFromFloat(100.0)

// 计算红包金额
func (g *goodsDomain) nextAmount(goods *RedEnvelopeGoods) (amount decimal.Decimal) {
	if goods.RemainQuantity == 1 {
		return goods.RemainAmount
	}
	if goods.EnvelopeType == int(services.CommonEnvelope) {
		return goods.AmountOne
	} else if goods.EnvelopeType == int(services.LuckyEnvelope) {
		// 将剩余的金额转换成 分 为单位，IntPart()为取出它的int值
		cent := goods.RemainAmount.Mul(multiple).IntPart()
		next := algo.DoubleAverage(int64(goods.RemainQuantity), cent)
		amount = decimal.NewFromFloat(float64(next)).Div(multiple)
	} else {
		logrus.Error("不支持的红包类型")
	}
	return amount
}

func (g *goodsDomain) transfer(
	ctx context.Context,
	dto services.RedEnvelopeReceiveDTO) (status services.TransferStatus, err error) {
	systemAccount := base.GetSystemAccount()
	body := services.TradeParticipator{
		AccountNo: systemAccount.AccountNo,
		UserId:    systemAccount.UserId,
		UserName:  systemAccount.UserName,
	}
	target := services.TradeParticipator{
		AccountNo: dto.AccountNo,
		UserId:    dto.RecvUserId,
		UserName:  dto.RecvUserName,
	}
	// 从系统红包资金账户扣减
	transfer := services.AccountTransferDTO{
		TradeNo:     dto.EnvelopeNo,
		TradeBody:   body,
		TradeTarget: target,
		Amount:      g.item.Amount,
		ChangeType:  services.EnvelopeOutgoning,
		ChangeFlag:  services.FlagTransferOut,
		Decs:        "红包扣减",
	}
	accountDo := account.NewAccountDomain()
	status, err = accountDo.TransferWithContextTx(ctx, transfer)
	if status != services.TransferStatusSuccess || err != nil {
		return status, err
	}

	// 从系统红包资金账户转入当前用户
	transfer = services.AccountTransferDTO{
		TradeBody:   target,
		TradeTarget: body,
		TradeNo:     dto.EnvelopeNo,
		Amount:      g.item.Amount,
		ChangeType:  services.EnvelopeIncoming,
		ChangeFlag:  services.FlagTransferIn,
		Decs:        "红包收入" + g.item.Amount.String() + "元",
	}
	return accountDo.TransferWithContextTx(ctx, transfer)
}
