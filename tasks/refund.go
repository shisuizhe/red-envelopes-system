package tasks

import (
	"fmt"
	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	"resk/core/envelope"
	"resk/infra"
	"time"
)

type RefundExpiredTaskStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (r *RefundExpiredTaskStarter) Init(ctx infra.StarterContext) {
	d := ctx.Props().GetDurationDefault("tasks.refund.interval", 60*time.Second)
	r.ticker = time.NewTicker(d)

	ip, err := utils.GetExternalIP()
	if err != nil {
		ip = "127.0.0.1"
	}

	maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	idleTimeout := ctx.Props().GetDurationDefault("redis.idleTimeout", 20*time.Second)
	host := ctx.Props().GetDefault("redis.host", "127.0.0.1")
	port := ctx.Props().GetDefault("redis.port", "6379")
	pwd := ctx.Props().GetDefault("redis.pwd", "Pd940810")

	pools := make([]redsync.Pool, 0)
	pool := &redis.Pool{
		MaxIdle:     maxIdle,     // 最大空闲连接数
		MaxActive:   maxActive,   // 最大连接数
		IdleTimeout: idleTimeout, // 超时时间
		Dial: func() (conn redis.Conn, err error) {
			// 1.打开连接
			conn, err = redis.Dial("tcp", host + ":" + port)
			if err != nil {
				return nil, err
			}
			// 2.访问认证
			if _, err = conn.Do("auth", pwd); err != nil {
				conn.Close()
				return nil, err
			}
			return conn, nil
		},
	}
	pools = append(pools, pool)
	rsync := redsync.New(pools)
	r.mutex = rsync.NewMutex("lock:RefundExpired",
		redsync.SetExpiry(60*time.Second),
		// 重试次数
		redsync.SetRetryDelay(3),
		// 设置此key的value
		redsync.SetGenValueFunc(func() (s string, err error) {
			now := time.Now()
			logrus.Infof("节点%s正在执行过期红包的退款任务", ip)
			return fmt.Sprintf("%d:%s", now.Unix(), ip), nil
		}),
	)
}

func (r *RefundExpiredTaskStarter) Start(ctx infra.StarterContext) {
	go func() {
		for {
			c := <-r.ticker.C
			err := r.mutex.Lock()
			if err == nil {
				logrus.Info("过期红包退款开始...", c)
				// 红包过期退款业务逻辑代码
				domain := envelope.ExpiredEnvelopeDomain{}
				_ = domain.Expired()
			}else {
				logrus.Info("已经有节点在运行该任务了")
			}
			r.mutex.Unlock()
		}
	}()
}

func (r *RefundExpiredTaskStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}
