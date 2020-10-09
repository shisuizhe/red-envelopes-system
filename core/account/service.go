package account

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"resk/infra/base"
	"resk/services"
	"sync"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.AccountServiceInterface = new(accountService)
	})
}

type accountService struct{}

func (as *accountService) CreateAccount(dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	// 验证输入参数
	if err := base.ValidateStruct(&dto); err != nil {
		return nil, err
	}
	// 验证账户 是否存在 幂等性
	ad := accountDomain{}

	a := ad.GetAccountByUserId(dto.UserId)
	if a != nil {
		return nil, errors.New(fmt.Sprintf("该类型账户已存在"))
	}
	// 账户创建
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	accountDto := services.AccountDTO{
		AccountType:  dto.AccountType,
		AccountName:  dto.AccountName,
		CurrencyCode: dto.CurrencyCode,
		UserId:       dto.UserId,
		UserName:     dto.UserName,
		Balance:      amount,
		Status:       1,
	}
	res, err := ad.Create(accountDto)
	return res, err
}

func (as *accountService) Transfer(dto services.AccountTransferDTO) (services.TransferStatus, error) {
	ad := accountDomain{}
	status, ok := ad.Transfer(dto)

	// 转账成功，并且交易主体和交易目标不是同一个人，而且交易类型不是储值，则进行反向操作
	if status == services.TransferStatusSuccess && dto.TradeBody.AccountNo != dto.TradeTarget.AccountNo && dto.ChangeType != services.AccountStoreValue {
		backwardDto := dto
		backwardDto.TradeBody = dto.TradeTarget
		backwardDto.TradeTarget = dto.TradeBody
		backwardDto.ChangeType = -dto.ChangeType
		backwardDto.ChangeFlag = -dto.ChangeFlag
		status, err := ad.Transfer(backwardDto)
		return status, err
	}
	return status, ok
}

func (as *accountService) StoreValue(dto services.AccountTransferDTO) (services.TransferStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return as.Transfer(dto)
}

func (as *accountService) GetEnvelopeAccountByUserId(userId string) *services.AccountDTO {
	ad := accountDomain{}
	return ad.GetAccountByUserId(userId)
}

func (as *accountService) GetAccount(accountNo string) *services.AccountDTO {
	ad := accountDomain{}
	return ad.GetAccountByAccountNo(accountNo)
}
