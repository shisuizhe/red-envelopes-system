package base

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"resk/infra"
	"sync"
)

var props kvs.ConfigSource

func Props() kvs.ConfigSource {
	return props
}

// 系统启动后，最先执行
type PropsStarter struct {
	infra.BaseStarter
}

func (s *PropsStarter) Init(ctx infra.StarterContext) {
	props = ctx.Props()
	log.Info("初始化配置中...")
	GetSystemAccount()
}

type SystemAccount struct {
	AccountNo string
	AccountName string
	UserId string
	UserName string
}

var systemAccount *SystemAccount
var once sync.Once

func GetSystemAccount() *SystemAccount {
	once.Do(func() {
		systemAccount = new(SystemAccount)
		err := kvs.Unmarshal(Props(), systemAccount, "system.account")
		if err != nil {
			panic(err)
		}
	})
	return systemAccount
}

func GetEnvelopeActivityLink() string {
	link, _ := Props().Get("envelope.link")
	return link
}

func GetEnvelopeDomain() string {
	domain, _ := Props().Get("envelope.domain")
	return domain
}
