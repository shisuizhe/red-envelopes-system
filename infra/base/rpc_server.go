package base

import (
	"github.com/sirupsen/logrus"
	"net"
	"net/rpc"
	"reflect"
	"resk/infra"
)

var rpcServer *rpc.Server

func RpcServer() *rpc.Server {
	Check(rpcServer)
	return rpcServer
}

func RpcServerRegister(i interface{}) {
	tf := reflect.TypeOf(i)
	logrus.Infof("RpcServer Register: %s", tf.String())
	_ = RpcServer().Register(i)
}

type RpcServerStarter struct {
	infra.BaseStarter
	server *rpc.Server
}

func (s *RpcServerStarter) Init(infra.StarterContext) {
	s.server = rpc.NewServer()
	rpcServer = s.server
}

func (s *RpcServerStarter) Start(ctx infra.StarterContext) {
	port, _ := ctx.Props().Get("app.rpc.port")
	// 监听网络端口
	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		logrus.Panic(err)
	}
	logrus.Infof("rpc is listening on port: %s", port)
	// 处理网络连接和请求
	go s.server.Accept(listener)
}
