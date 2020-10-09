package envelope

import (
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"resk/services"
	_ "resk/test"
	"strconv"
	"testing"
)

func TestRedEnvelopeService_SendOut(t *testing.T) {
	// 发红包人的红包资金账户
	as := services.GetAccountServiceInterface()
	account := services.AccountCreatedDTO{
		UserId:       ksuid.New().Next().String(),
		UserName:     "发送红包账户",
		AccountName:  "发送红包账户",
		AccountType:  int(services.UserAccount),
		CurrencyCode: "CYN",
		Amount:       "200",
	}
	r := services.GetRedEnvelopeServiceInterface()
	Convey("准备资金账户", t, func() {
		a, err := as.CreateAccount(account)
		So(err, ShouldBeNil)
		So(a, ShouldNotBeNil)
	})
	Convey("发送红包", t, func() {
		Convey("发普通红包", func() {
			goods := services.RedEnvelopeSendOutDTO{
				UserId:       account.UserId,
				UserName:     account.UserName,
				EnvelopeType: services.CommonEnvelope,
				Amount:       decimal.NewFromFloat(10),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			activity, err := r.SendOut(goods)
			So(err, ShouldBeNil)
			So(activity, ShouldNotBeNil)
			So(activity.Link, ShouldNotBeEmpty)
			So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			dto := activity.RedEnvelopeGoodsDTO
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.UserName, ShouldEqual, goods.UserName)
			quantity := decimal.NewFromFloat(float64(dto.Quantity))
			So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(quantity).String())
			So(dto.Quantity, ShouldEqual, goods.Quantity)
		})

		Convey("发碰运气红包", func() {
			goods := services.RedEnvelopeSendOutDTO{
				UserId:       account.UserId,
				UserName:     account.UserName,
				EnvelopeType: services.LuckyEnvelope,
				Amount:       decimal.NewFromFloat(100),
				Quantity:     10,
				Blessing:     services.DefaultBlessing,
			}
			activity, err := r.SendOut(goods)
			So(err, ShouldBeNil)
			So(activity, ShouldNotBeNil)
			So(activity.Link, ShouldNotBeEmpty)
			So(activity.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			// 验证每一个属性
			dto := activity.RedEnvelopeGoodsDTO
			So(dto.UserId, ShouldEqual, goods.UserId)
			So(dto.UserName, ShouldEqual, goods.UserName)
			So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
			So(dto.Quantity, ShouldEqual, goods.Quantity)
		})
	})
}

func TestRedEnvelopeService_Receive(t *testing.T) {
	Convey("准备资金账户", t, func() {
		as := services.GetAccountServiceInterface()

		// 准备发红包资金账户
		account := services.AccountCreatedDTO{
			UserId:       ksuid.New().Next().String(),
			UserName:     "发送红包账户",
			AccountName:  "发送红包账户",
			AccountType:  int(services.UserAccount),
			CurrencyCode: "CYN",
			Amount:       "100",
		}
		fa, err := as.CreateAccount(account)
		So(err, ShouldBeNil)
		So(fa, ShouldNotBeNil)

		// 准备收红包资金账户
		accounts := make([]*services.AccountDTO, 0)
		size := 10
		for i := 0; i < size; i++ {
			account := services.AccountCreatedDTO{
				UserId:       ksuid.New().Next().String(),
				UserName:     "接收账户" + strconv.Itoa(i),
				Amount:       "0",
				AccountName:  "接收账户" + strconv.Itoa(i),
				AccountType:  int(services.UserAccount),
				CurrencyCode: "CNY",
			}
			shou, err := as.CreateAccount(account)
			So(err, ShouldBeNil)
			So(shou, ShouldNotBeNil)
			accounts = append(accounts, shou)
		}
		So(len(accounts), ShouldEqual, size)

		rs := services.GetRedEnvelopeServiceInterface()

		// 收普通红包测试用例
		Convey("收普通红包测试用例", func() {
			goods := services.RedEnvelopeSendOutDTO{
				UserId:       fa.UserId,
				UserName:     fa.UserName,
				EnvelopeType: services.CommonEnvelope,
				Amount:       decimal.NewFromFloat(10),
				Quantity:     10,
				Blessing:     "金额相同",
			}
			// 开始发红包
			at, err := rs.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			g := at.RedEnvelopeGoodsDTO
			So(g.UserName, ShouldEqual, goods.UserName)
			So(g.UserId, ShouldEqual, goods.UserId)
			So(g.Quantity, ShouldEqual, goods.Quantity)
			q := decimal.NewFromFloat(float64(g.Quantity))
			So(g.Amount, ShouldEqual, goods.Amount.Mul(q))
			remainAmount := g.Amount

			Convey("收普通红包", func() {
				for i, account := range accounts {
					recv := services.RedEnvelopeReceiveDTO{
						EnvelopeNo:   at.EnvelopeNo,
						RecvUserId:   account.UserId,
						RecvUserName: account.UserName,
						AccountNo:    account.AccountNo,
					}
					item, err := rs.Receive(recv)
					if item != nil {
						logrus.Infof("index:%d  amount:%+v", i, item.Amount.String())
					}
					So(err, ShouldBeNil)
					So(item, ShouldNotBeNil)
					So(item.Amount, ShouldEqual, at.AmountOne)
					remainAmount = remainAmount.Sub(at.AmountOne)
					So(item.RemainAmount, ShouldEqual, remainAmount)
				}
			})
		})

		// 收幸运红包测试用例
		Convey("收幸运红包测试用例", func() {
			goods := services.RedEnvelopeSendOutDTO{
				UserId:       fa.UserId,
				UserName:     fa.UserName,
				EnvelopeType: services.LuckyEnvelope,
				Amount:       decimal.NewFromFloat(100),
				Quantity:     10,
				Blessing:     "金额随机",
			}
			// 开始发红包
			at, err := rs.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)
			g := at.RedEnvelopeGoodsDTO
			So(g.UserName, ShouldEqual, goods.UserName)
			So(g.UserId, ShouldEqual, goods.UserId)
			So(g.Quantity, ShouldEqual, goods.Quantity)
			So(g.Amount, ShouldEqual, goods.Amount)
			remainAmount := at.Amount

			Convey("收碰运气红包", func() {
				total := decimal.NewFromFloat(0)
				for i, account := range accounts {
					if i > 10 {
						break
					}
					recv := services.RedEnvelopeReceiveDTO{
						EnvelopeNo:   at.EnvelopeNo,
						RecvUserId:   account.UserId,
						RecvUserName: account.UserName,
						AccountNo:    account.AccountNo,
					}
					item, err := rs.Receive(recv)
					if item != nil {
						total = total.Add(item.Amount)
						logrus.Infof("index:%d  amount:%+v", i, item.Amount.String())
					}
					So(err, ShouldBeNil)
					So(item, ShouldNotBeNil)
					remainAmount = remainAmount.Sub(item.Amount)
					So(item.RemainAmount.String(), ShouldEqual, remainAmount.String())
				}
				So(total.String(), ShouldEqual, goods.Amount.String())
			})
		})

		// 余额不足而发红包失败
		Convey("余额不足而发红包失败", func() {
			goods := services.RedEnvelopeSendOutDTO{
				UserId:       fa.UserId,
				UserName:     fa.UserName,
				EnvelopeType: services.LuckyEnvelope,
				Amount:       decimal.NewFromFloat(200),
				Quantity:     10,
				Blessing:     "金额随机",
			}
			// 开始发红包
			at, err := rs.SendOut(goods)
			So(err, ShouldNotBeNil)
			So(at, ShouldBeNil)
		})
	})
}
