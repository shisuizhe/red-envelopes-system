package rpc

import (
	"resk/infra"
	"resk/infra/base"
)

type RpcApiStarter struct {
	infra.BaseStarter
}

func (r *RpcApiStarter) Init(ctx infra.StarterContext) {
	base.RpcServerRegister(new(EnvelopeRpc))
}
