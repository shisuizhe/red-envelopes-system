package resk

import (
	// "resk/apis/rpc"
	// _ "resk/apis/rpc"
	_ "resk/apis/web"
	_ "resk/core/account"
	_ "resk/core/envelope"
	"resk/infra"
	"resk/infra/base"
	_ "resk/public/ui"
	"resk/tasks"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	// infra.Register(&base.RpcServerStarter{})
	// infra.Register(&rpc.RpcApiStarter{})
	infra.Register(&tasks.RefundExpiredTaskStarter{})
	infra.Register(&base.IrisServerStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.HookStarter{})
}
