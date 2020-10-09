// 初始化测试环境需要的代码
package test

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"resk/infra"
	"resk/infra/base"
)

func init() {
	// 获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("../brun/config.ini", 1)
	// 加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)

	infra.Register(&base.PropsStarter{})
	infra.Register(&base.DatabaseStarter{})
	infra.Register(&base.ValidatorStarter{})
	// infra.Register(&base.IrisServerStarter{})

	app := infra.NewBootApplication(conf)
	app.Start()
}
