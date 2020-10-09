package rpc

import (
	"resk/services"
)

type EnvelopeRpc struct{}

// Go内置的RPC接口有一些规范：
// 1.入参和出参都要作为方法参数
// 2.方法必须有2个参数，并且是可导出类型
// 3.第二个参数（返回值）必须是指针类型
// 4.方法返回值要返回error类型
// 5.方法必须是可导出的
func (e *EnvelopeRpc) SendOut(
	in services.RedEnvelopeSendOutDTO,
	out *services.RedEnvelopeActivity) error {
	s := services.GetRedEnvelopeServiceInterface()
	a, err := s.SendOut(in)
	a.CopyTo(out)
	return err
}

func (e *EnvelopeRpc) Receive(
	in services.RedEnvelopeReceiveDTO,
	out *services.RedEnvelopeItemDTO) error {
	s := services.GetRedEnvelopeServiceInterface()
	a, err := s.Receive(in)
	a.CopyTo(out)
	return err
}
