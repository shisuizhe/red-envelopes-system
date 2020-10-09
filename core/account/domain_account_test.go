package account

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"resk/services"
	"testing"
)

func TestAccountDomain_Create(t *testing.T) {
	a := services.AccountDTO{
		UserName: "测试创建账户",
		UserId:   ksuid.New().Next().String(),
		Balance:  decimal.NewFromFloat(999999999999),
		Status:   1,
	}

	Convey("账户创建", t, func() {
		ad := accountDomain{}
		ndto, err := ad.Create(a)
		So(err, ShouldBeNil)
		So(ndto, ShouldNotBeNil)
		So(ndto.Balance.String(), ShouldEqual, a.Balance.String())
		So(ndto.UserId, ShouldEqual, a.UserId)
		So(ndto.UserName, ShouldEqual, a.UserName)
		So(ndto.Status, ShouldEqual, a.Status)
	})
}

func TestAccountDomain_Transfer(t *testing.T) {
	// 准备2个账户，交易主体账户要有余额
	a1 := &services.AccountDTO{
		UserName: "测试转账1",
		UserId:   ksuid.New().Next().String(),
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}
	a2 := &services.AccountDTO{
		UserName: "测试转账2",
		UserId:   ksuid.New().Next().String(),
		Balance:  decimal.NewFromFloat(100),
		Status:   1,
	}

	ad := accountDomain{}

	Convey("转账测试", t, func() {
		// 创建账户a1
		a11, err := ad.Create(*a1)
		So(err, ShouldBeNil)
		So(a11, ShouldNotBeNil)
		So(a11.Balance.String(), ShouldEqual, a1.Balance.String())
		So(a11.UserId, ShouldEqual, a1.UserId)
		So(a11.UserName, ShouldEqual, a1.UserName)
		So(a11.Status, ShouldEqual, a1.Status)
		a1 = a11
		// 创建账户a2
		a22, err := ad.Create(*a2)
		So(err, ShouldBeNil)
		So(a22, ShouldNotBeNil)
		So(a22.Balance.String(), ShouldEqual, a2.Balance.String())
		So(a22.UserId, ShouldEqual, a2.UserId)
		So(a22.UserName, ShouldEqual, a2.UserName)
		So(a22.Status, ShouldEqual, a2.Status)
		a2 = a22

		// 余额充足，金额转入其他账户
		Convey("余额充足，金额转入其他账户", func() {
			amount := decimal.NewFromFloat(100)
			body := services.TradeParticipator{
				AccountNo: a1.AccountNo,
				UserId:    a1.UserId,
				UserName:  a1.UserName,
			}
			target := services.TradeParticipator{
				AccountNo: a2.AccountNo,
				UserId:    a2.UserId,
				UserName:  a2.UserName,
			}
			transfer := services.AccountTransferDTO{
				TradeNo:     ksuid.New().Next().String(),
				TradeBody:   body,
				TradeTarget: target,
				AmountStr:   amount.String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := ad.Transfer(transfer)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferStatusSuccess)

			na := ad.GetAccountByAccountNo(a1.AccountNo)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a1.Balance.Sub(amount).String())
		})

		// 余额不足，金额转出
		Convey("余额不足，金额转出", func() {
			amount := decimal.NewFromFloat(101)
			body := services.TradeParticipator{
				AccountNo: a1.AccountNo,
				UserId:    a1.UserId,
				UserName:  a1.UserName,
			}
			target := services.TradeParticipator{
				AccountNo: a2.AccountNo,
				UserId:    a2.UserId,
				UserName:  a2.UserName,
			}
			transfer := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				Amount:      amount,
				ChangeType:  services.ChangeType(-1),
				ChangeFlag:  services.FlagTransferOut,
				Decs:        "转账",
			}
			status, err := ad.Transfer(transfer)
			So(err, ShouldNotBeNil)
			So(status, ShouldEqual, services.TransferStatusSufficientFunds)

			na := ad.GetAccountByAccountNo(a1.AccountNo)
			So(na, ShouldNotBeNil)
			So(transfer.Amount.Sub(na.Balance), ShouldEqual, decimal.NewFromFloat(1))
		})

		// 充值
		Convey("充值", func() {
			amount := decimal.NewFromFloat(50)
			body := services.TradeParticipator{
				AccountNo: a1.AccountNo,
				UserId:    a1.UserId,
				UserName:  a1.UserName,
			}
			target := services.TradeParticipator{
				AccountNo: a2.AccountNo,
				UserId:    a2.UserId,
				UserName:  a2.UserName,
			}
			transfer := services.AccountTransferDTO{
				TradeBody:   body,
				TradeTarget: target,
				TradeNo:     ksuid.New().Next().String(),
				AmountStr:   amount.String(),
				Amount:      amount,
				ChangeType:  services.AccountStoreValue,
				ChangeFlag:  services.FlagTransferIn,
				Decs:        "充值",
			}
			status, err := ad.Transfer(transfer)
			So(err, ShouldBeNil)
			So(status, ShouldEqual, services.TransferStatusSuccess)

			na := ad.GetAccountByAccountNo(a1.AccountNo)
			So(na, ShouldNotBeNil)
			So(na.Balance.Sub(amount), ShouldEqual, a1.Balance)
		})
	})
}
