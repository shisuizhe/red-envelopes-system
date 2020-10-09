package account

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

// 根据用户编号查询用户
func (adao *AccountDao) GetByAccountNo(accoutNo string) *Account {
	a := &Account{AccountNo: accoutNo}
	ok, err := adao.runner.GetOne(a)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

// 通过用户ID和账户类型来查询账户信息
func (adao *AccountDao) GetByUserId(userId string) *Account {
	a := &Account{}
	sql := "select * from account where user_id=?"
	ok, err := adao.runner.Get(a, sql, userId)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

// 账户数据的插入
func (adao *AccountDao) Insert(a *Account) (int64, error) {
	res, err := adao.runner.Insert(a)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// 账户余额的更新(amount 如果是负数，就是扣减；如果是正数，就是增加)
func (adao *AccountDao) UpdateBalance(accountNo string, amount decimal.Decimal) (int64, error) {
	sql := "update account" +
		" set balance = balance + CAST(? AS DECIMAL(30, 6))" +
		" where account_no = ?" +
		" and balance >= -1*CAST(? AS DECIMAL(30, 6))"
	res, err := adao.runner.Exec(sql, amount.String(), accountNo, amount.String())
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

// 账户状态更新
func (adao *AccountDao) UpdateStatus(accountNo string, status int) (int64, error) {
	sql := "update account set status = ? where account_no = ?"
	res, err := adao.runner.Exec(sql, status, accountNo)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
