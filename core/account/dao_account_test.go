// 测试通过

package account

import (
	"database/sql"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
	"resk/infra/base"
	_ "resk/test"
	"testing"
)

func TestAccountDao_GetByAccountNo(t *testing.T) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}

		Convey("通过账户编号编号查询账户数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "通过账户编号编号查询账户数据",
				UserId:      ksuid.New().Next().String(),
				UserName:    sql.NullString{String: "通过账户编号编号查询账户数据", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			a1 := dao.GetByAccountNo(a.AccountNo)
			So(a1, ShouldNotBeNil)
			So(a1.Balance.String(), ShouldEqual, a.Balance.String())
			So(a1.UserId, ShouldEqual, a.UserId)
			So(a1.UserName.String, ShouldEqual, a.UserName.String)
		})
		return nil
	})
}

func TestAccountDao_GetByUserId(t *testing.T) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}

		Convey("通过用户ID和账户类型查询账户数据", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "通过用户ID和账户类型查询账户数据",
				UserId:      ksuid.New().Next().String(),
				UserName:    sql.NullString{String: "通过用户ID和账户类型查询账户数据", Valid: true},
				AccountType: 1,
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			a1 := dao.GetByUserId(a.UserId)
			So(a1, ShouldNotBeNil)
			So(a1.Balance.String(), ShouldEqual, a.Balance.String())
			So(a1.UserId, ShouldEqual, a.UserId)
			So(a1.UserName.String, ShouldEqual, a.UserName.String)
		})
		return nil
	})
}

func TestAccountDao_UpdateBalance(t *testing.T) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := &AccountDao{runner: runner}

		Convey("测试增加余额", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试增加余额",
				UserId:      ksuid.New().Next().String(),
				UserName:    sql.NullString{String: "测试增加余额", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			// 增加余额
			Convey("增加余额", func() {
				amount := decimal.NewFromFloat(100)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)

				a1 := dao.GetByAccountNo(a.AccountNo)
				So(a1, ShouldNotBeNil)
				So(a1.Balance.Sub(decimal.NewFromFloat(200)), ShouldEqual, decimal.NewFromFloat(0))
			})
		})

		Convey("测试扣减余额", t, func() {
			a := &Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试扣减余额",
				UserId:      ksuid.New().Next().String(),
				UserName:    sql.NullString{String: "测试扣减余额", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			// 扣减余额
			Convey("扣减余额", func() {
				amount := decimal.NewFromFloat(-100)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldBeGreaterThan, 0)

				a1 := dao.GetByAccountNo(a.AccountNo)
				So(a1, ShouldNotBeNil)
				So(a1.Balance.Add(decimal.NewFromFloat(200)), ShouldEqual, decimal.NewFromFloat(200))
			})
		})

		return nil
	})
}
