// DTO是数据传输对象(Data Transfer Object)的缩写
package services

import (
	"github.com/shopspring/decimal"
	"resk/infra/base"
	"time"
)

var AccountServiceInterface AccountService

func GetAccountServiceInterface() AccountService {
	base.Check(AccountServiceInterface)
	return AccountServiceInterface
}

type AccountService interface {
	CreateAccount(dto AccountCreatedDTO) (*AccountDTO, error)
	Transfer(dto AccountTransferDTO) (TransferStatus, error)
	StoreValue(dto AccountTransferDTO) (TransferStatus, error)
	GetEnvelopeAccountByUserId(userId string) *AccountDTO
	GetAccount(accountNo string) *AccountDTO
}

// 转账对象
type AccountTransferDTO struct {
	TradeNo     string            // `validate:"required"`         // 交易单号 全局不重复字符或数字，唯一性标识
	TradeBody   TradeParticipator `validate:"required"` // 交易主体
	TradeTarget TradeParticipator `validate:"required"` // 交易对象
	AmountStr   string            // `validate:"required,numeric"` // 交易金额,该交易涉及的金额
	Amount      decimal.Decimal   ``                            // 交易金额,该交易涉及的金额
	ChangeType  ChangeType        `validate:"required,numeric"` // 流水交易类型，0 创建账户，>0 为收入类型，<0 为支出类型
	ChangeFlag  ChangeFlag        `validate:"required,numeric"` // 交易变化标识：-1 出账 1为进账，枚举
	Decs        string            ``                            // 交易描述
}

// 交易参与者
type TradeParticipator struct {
	AccountNo string `validate:"required"` // 账户编号 账户ID
	UserId    string `validate:"required"` // 用户编号
	UserName  string `validate:"required"` // 用户姓名
}

// 账户创建对象
type AccountCreatedDTO struct {
	UserId       string `validate:"required"`
	UserName     string `validate:"required"`
	AccountName  string `validate:"required"`
	AccountType  int    ``
	CurrencyCode string ``
	Amount       string `validate:"required"` // 金额用string类型，因为在go中float32、float64在计算时会丢失精度
}

// 账户对象
type AccountDTO struct {
	AccountNo    string          // 账户编号,账户唯一标识
	AccountName  string          // 账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xxx积分、xxx零钱
	AccountType  int             // 账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户
	CurrencyCode string          // 货币类型编码：CNY人民币，EUR欧元，USD美元 。。。
	UserId       string          // 用户编号, 账户所属用户
	UserName     string          // 用户名称
	Balance      decimal.Decimal // 账户可用余额
	Status       int             // 账户状态，账户状态：0账户初始化，1启用，2停用
	CreatedAt    time.Time       // 创建时间
	UpdatedAt    time.Time       // 更新时间
}

func (a *AccountDTO) CopyTo(trgt *AccountDTO) {
	trgt.AccountNo = a.AccountNo
	trgt.AccountName = a.AccountName
	trgt.AccountType = a.AccountType
	trgt.CurrencyCode = a.CurrencyCode
	trgt.UserId = a.UserId
	trgt.UserName = a.UserName
	trgt.Balance = a.Balance
	trgt.Status = a.Status
	trgt.CreatedAt = a.CreatedAt
	trgt.UpdatedAt = a.UpdatedAt
}

// 账户流水对象
type AccountLogDTO struct {
	LogNo           string          // 流水编号 全局不重复字符或数字，唯一性标识
	TradeNo         string          // 交易单号 全局不重复字符或数字，唯一性标识
	AccountNo       string          // 账户编号 账户ID
	UserId          string          // 用户编号
	UserName        string          // 用户名称
	TargetAccountNo string          // 账户编号 账户ID
	TargetUserId    string          // 目标用户编号
	TargetUserName  string          // 目标用户名称
	Amount          decimal.Decimal // 交易金额,该交易涉及的金额
	Balance         decimal.Decimal // 交易后余额,该交易后的余额
	ChangeType      ChangeType      // 流水交易类型：0 创建账户，>0 为收入类型，<0 为支出类型
	ChangeFlag      ChangeFlag      // 交易变化标识：-1 出账 1为进账
	Status          int             // 交易状态：
	Decs            string          // 交易描述
	CreatedAt       time.Time       // 创建时间
}
