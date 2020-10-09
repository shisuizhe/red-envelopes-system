package base

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"github.com/tietang/props/kvs"
	"resk/infra"
	log "resk/infra/logrus"
)

var dataBase *dbx.Database

func DataBase() *dbx.Database {
	Check(dataBase)
	return dataBase
}

// dbx数据库starter
type DatabaseStarter struct {
	infra.BaseStarter
}

func (s *DatabaseStarter) Setup(ctx infra.StarterContext) {
	conf := ctx.Props()
	// 数据库配置
	settings := dbx.Settings{}
	err := kvs.Unmarshal(conf, &settings, "mysql")
	if err != nil {
		panic(err)
	}
	logrus.Info("mysql.conn url:", settings.ShortDataSourceName())
	db, err := dbx.Open(settings)
	if err != nil {
		panic(err)
	}
	logrus.Info(db.Ping())
	db.SetLogger(log.NewUpperLogrusLogger())
	dataBase = db
}
