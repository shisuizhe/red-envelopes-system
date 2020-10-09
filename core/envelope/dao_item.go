package envelope

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type RedEnvelopeItemDao struct {
	runner *dbx.TxRunner
}

// 查询
func (dao *RedEnvelopeItemDao) GetOne(itemNo string) *RedEnvelopeItem {
	po := &RedEnvelopeItem{ItemNo: itemNo}
	ok, err := dao.runner.GetOne(po)
	if err != nil || !ok {
		return nil
	}
	return po
}

// 红包订单详情数据的插入
func (dao *RedEnvelopeItemDao) Insert(po *RedEnvelopeItem) (int64, error) {
	res, err := dao.runner.Insert(po)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (dao *RedEnvelopeItemDao) FindItems(envelopeNo string) []*RedEnvelopeItem {
	var items = make([]*RedEnvelopeItem, 0)
	sql := "select * from red_envelope_item where envelope_no = ?"
	err := dao.runner.Find(&items, sql, envelopeNo)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return items
}

func (dao *RedEnvelopeItemDao) GetByUser(userId, envelopeNo string) *RedEnvelopeItem {
	item := &RedEnvelopeItem{}
	sql := "select * from red_envelope_item where envelope_no=? and recv_user_id=?"
	ok, err := dao.runner.Get(item, sql, envelopeNo, userId)
	if !ok {
		return nil
	}
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return item
}

func (dao *RedEnvelopeItemDao) ListReceivedItems(userId string, page, size int) []*RedEnvelopeItem {
	items := make([]*RedEnvelopeItem, 0)
	sql := "select * from red_envelope_item where recv_user_id=? order by created_at desc limit ?,?"
	err := dao.runner.Find(&items, sql, userId, page, size)
	if err != nil {
		logrus.Error(err)
		return items
	}
	return items
}
