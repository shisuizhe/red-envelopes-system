package envelope

import (
	"github.com/tietang/dbx"
	"resk/infra/base"
)

func (*goodsDomain) Find(po *RedEnvelopeGoods) (regs []RedEnvelopeGoods) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		regs = dao.Find(po)
		return nil
	})
	return regs
}

func (*goodsDomain) GetOne(envelopeNo string) (goods *RedEnvelopeGoods) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.GetOne(envelopeNo)
		return nil
	})
	return goods
}

func (*goodsDomain) FindByUser(userId string, offset, limit int) (goods []RedEnvelopeGoods) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.FindByUser(userId, offset, limit)
		return nil
	})
	return goods
}

func (*goodsDomain) ListReceived(userId string, offset, limit int) (items []*RedEnvelopeItem) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		items = dao.ListReceivedItems(userId, offset, limit)
		return nil
	})
	return items
}

func (*goodsDomain) ListReceivable(offset, limit int) (goods []RedEnvelopeGoods) {
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeGoodsDao{runner: runner}
		goods = dao.ListReceivable(offset, limit)
		return nil
	})
	return goods
}
