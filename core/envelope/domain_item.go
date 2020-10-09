package envelope

import (
	"context"
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
)

type itemDomain struct {
	RedEnvelopeItem
}

// 生成itemNo
func (do *itemDomain) createItemNo() {
	do.ItemNo = ksuid.New().Next().String()
}

// 创建item数据
func (do *itemDomain) Create(item services.RedEnvelopeItemDTO) {
	do.RedEnvelopeItem.FromDTO(&item)
	do.RecvUserName.Valid = true
	do.createItemNo()
}

// 保存item数据
func (do *itemDomain) Save(ctx context.Context) (id int64, err error) {
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		id, err = dao.Insert(&do.RedEnvelopeItem)
		if err != nil {
			return err
		}
		return nil
	})
	return id, err
}

// 通过itemNo查询抢红包明细数据
func (do *itemDomain) GetOne(ctx context.Context, itemNo string) (
	dto *services.RedEnvelopeItemDTO) {
	err := base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		po := dao.GetOne(itemNo)
		if po != nil {
			dto = po.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

func (do *itemDomain) GetByUserIdAndEnvelopeNo(userId, envelopeNo string) (dto *services.RedEnvelopeItemDTO) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		item := dao.GetByUser(userId, envelopeNo)
		if item != nil {
			dto = item.ToDTO()
		}
		return nil
	})
	if err != nil {
		return nil
	}
	return dto
}

// 通过envelopeNo查询已抢到红包列表
func (do *itemDomain) FindItems(envelopeNo string) (itemDtos []*services.RedEnvelopeItemDTO, err error) {
	var items []*RedEnvelopeItem
	err = base.Tx(func(runner *dbx.TxRunner) error {
		dao := RedEnvelopeItemDao{runner: runner}
		items = dao.FindItems(envelopeNo)
		if items == nil {
			return errors.New("该红包编号不存在")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	itemDtos = make([]*services.RedEnvelopeItemDTO, 0)
	var luckyItem *services.RedEnvelopeItemDTO
	for i, po := range items {
		item := po.ToDTO()
		if i == 0 {
			luckyItem = item
		} else {
			if luckyItem.Amount.Cmp(po.Amount) < 0 {
				luckyItem = item
			}
		}
		itemDtos = append(itemDtos, item)
	}
	// luckyItem.IsLuckiest = true
	return itemDtos, nil
}
