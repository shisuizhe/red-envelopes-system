package infra

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"reflect"
)

// 应用程序启动管理器
type BootApplication struct {
	IsTest         bool
	conf           kvs.ConfigSource
	starterContext StarterContext
}
// 构造系统
func NewBootApplication(conf kvs.ConfigSource) *BootApplication {
	b := &BootApplication{
		conf:           conf,
		starterContext: StarterContext{},
	}
	b.starterContext[KeyProps] = conf
	return b
}

func (b *BootApplication) Start() {
	// 1. 初始化 starter
	b.init()
	// 2. 安装 starter
	b.setup()
	// 3. 启动 starter
	b.start()
}

// 程序初始化
func (b *BootApplication) init() {
	logrus.Info("Initializing starters...")
	for _, v := range GetStarters() {
		tf := reflect.TypeOf(v)
		logrus.Debugf("Initializing: PriorityGroup=%d,Priority=%d,type=%s\n", v.PriorityGroup(), v.Priority(), tf.String())
		v.Init(b.starterContext)
	}
}

// 程序安装
func (b *BootApplication) setup() {
	logrus.Info("Setup starters...")
	for _, v := range GetStarters() {
		tf := reflect.TypeOf(v)
		logrus.Debug("Setup: ", tf.String())
		v.Setup(b.starterContext)
	}
}

// 程序开始运行，开始接受调用
func (b *BootApplication) start() {
	logrus.Info("Starting starters...")
	for i, v := range GetStarters() {
		tf := reflect.TypeOf(v)
		logrus.Debug("Starting: ", tf.String())
		// 启动器是否可阻塞
		if v.StartBlocking() {
			// 如果是最后一个可阻塞starter，直接启动并阻塞
			if i+1 == len(GetStarters()) {
				v.Start(b.starterContext)
				// 使用携程异步启动，防止阻塞后面的starter
			} else {
				go v.Start(b.starterContext)
			}
		} else {
			v.Start(b.starterContext)
		}
	}
}

func (b *BootApplication) Stop() {
	logrus.Info("Stoping starters...")
	for _, v := range GetStarters() {
		tf := reflect.TypeOf(v)
		logrus.Debug("Stoping: ", tf.String())
		v.Stop(b.starterContext)
	}
}
