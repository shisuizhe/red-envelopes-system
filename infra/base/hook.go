package base

import (
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"reflect"
	"resk/infra"
	"syscall"
)

var callbacks []func()

func Register(f func()) {
	callbacks = append(callbacks, f)
}

type HookStarter struct {
	infra.BaseStarter
}

func (s *HookStarter) Init(ctx infra.StarterContext) {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGQUIT, syscall.SIGTERM)
	go func() {
		for {
			s := <-ch
			logrus.Error("notify:", s)
			for _, f := range callbacks {
				f()
			}
			break
			os.Exit(0)
		}
	}()
}

func (s *HookStarter) Start(ctx infra.StarterContext) {
	starters := infra.GetStarters()
	for _, s := range starters {
		tf := reflect.TypeOf(s)
		logrus.Info("[Register Notify Stop]:%s.Stop()", tf.String())
		Register(func() {
			s.Stop(ctx)
		})
	}
}
