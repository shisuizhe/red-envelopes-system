// 测试通过

package account

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
	"testing"
)

func TestAccountLogDao(t *testing.T) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountLogDao{runner: runner}

		Convey("test", t, func() {
			a := &AccountLog{
				LogNo:      ksuid.New().Next().String(),
				TradeNo:    ksuid.New().Next().String(),
				Status:     1,
				AccountNo:  ksuid.New().Next().String(),
				UserId:     ksuid.New().Next().String(),
				UserName:   "通过logNo来查询",
				Amount:     decimal.NewFromFloat(100),
				Balance:    decimal.NewFromFloat(100),
				ChangeFlag: services.FlagAccountCreated,
				ChangeType: services.AccountCreated,
			}
			// 通过logNo来查询
			Convey("通过logNo来查询", func() {
				id, err := dao.Insert(a)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)

				a1 := dao.GetByLogNo(a.LogNo)
				So(a1, ShouldNotBeNil)
				So(a1.Balance.String(), ShouldEqual, a.Balance.String())
				So(a1.Amount.String(), ShouldEqual, a.Amount.String())
				So(a1.CreatedAt, ShouldNotBeNil)
			})
		})

		Convey("test", t, func() {
			a := &AccountLog{
				LogNo:      ksuid.New().Next().String(),
				TradeNo:    ksuid.New().Next().String(),
				Status:     1,
				AccountNo:  ksuid.New().Next().String(),
				UserId:     ksuid.New().Next().String(),
				UserName:   "通过tradeNo来查询",
				Amount:     decimal.NewFromFloat(200),
				Balance:    decimal.NewFromFloat(200),
				ChangeFlag: services.FlagAccountCreated,
				ChangeType: services.AccountCreated,
			}
			// 通过tradeNo来查询
			Convey("通过tradeNo来查询", func() {
				id, err := dao.Insert(a)
				So(err, ShouldBeNil)
				So(id, ShouldBeGreaterThan, 0)

				a1 := dao.GetByTradeNo(a.TradeNo)
				So(a1, ShouldNotBeNil)
				So(a1.Balance.String(), ShouldEqual, a.Balance.String())
				So(a1.Amount.String(), ShouldEqual, a.Amount.String())
				So(a1.CreatedAt, ShouldNotBeNil)
			})
		})

		return nil
	})
}
