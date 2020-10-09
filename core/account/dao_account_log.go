package account

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountLogDao struct {
	runner *dbx.TxRunner
}

// 通过流水编号查询流水记录
func (dao *AccountLogDao) GetByLogNo(logNo string) *AccountLog {
	al := &AccountLog{LogNo: logNo}
	ok, err := dao.runner.GetOne(al)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return al
}

// 通过交易编号来查询流水记录
func (dao *AccountLogDao) GetByTradeNo(tradeNo string) *AccountLog {
	sql := "select * from account_log where trade_no=?"
	al := &AccountLog{}
	ok, err := dao.runner.Get(al, sql, tradeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return al
}

// 流水记录的插入
func (dao *AccountLogDao) Insert(al *AccountLog) (int64, error) {
	res, err := dao.runner.Insert(al)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (dao *AccountLogDao) UpdateDecs() error {
	sql := "update account_log set `decs`=concat('入账', `amount`, '元') where change_flag in (1,2)"
	_, err := dao.runner.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}
