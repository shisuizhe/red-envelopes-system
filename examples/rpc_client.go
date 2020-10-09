package main

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"net/rpc"
	"resk/services"
)

func main() {
	cli, err := rpc.Dial("tcp", ":8082")
	if err != nil {
		logrus.Panic(err)
	}
	// sendOut(cli)
	receive(cli)
}

func sendOut(cli *rpc.Client) {
	in := services.RedEnvelopeSendOutDTO{
		Amount: decimal.NewFromFloat(90),
		UserId: "1fSJJYtPJBQFItcEV0NmvWgOjYi",
		UserName: "测试用户",
		EnvelopeType: int(services.LuckyEnvelope),
		Quantity: 2,
		Blessing: "快抢呀",
	}
	out := &services.RedEnvelopeActivity{}
	err := cli.Call("EnvelopeRpc.SendOut", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(out)
}

func receive(cli *rpc.Client) {
	in := services.RedEnvelopeReceiveDTO{
		EnvelopeNo:   "1fSRcy9zqRZ2Ix0arlfzp2GUPVu",
		RecvUserId:   "1fSJJYtPJBQFItcEV0NmvWgOjYi",
		RecvUserName: "测试账户1",
		// AccountNo:    "",
	}
	out := &services.RedEnvelopeItemDTO{}
	err := cli.Call("EnvelopeRpc.Receive", in, out)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Info(out)
}
