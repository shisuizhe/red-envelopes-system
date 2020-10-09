package services

const (
	// 默认常熟
	DefaultCurrencyCode = "CNY"
)

// 转账状态
type TransferStatus int8

const (
	// 转账失败
	TransferStatusFailed TransferStatus = -1
	// 余额不足
	TransferStatusSufficientFunds TransferStatus = 0
	// 转账成功
	TransferStatusSuccess TransferStatus = 1
)

// 转账类型：0 创建账户，>0 收入类型，<0 支出类型
type ChangeType int8

const (
	// 创建账户
	AccountCreated ChangeType = 0
	// 进账
	AccountStoreValue ChangeType = 1
	// 红包资金的支出
	EnvelopeOutgoning ChangeType = -2
	// 红包资金的收入
	EnvelopeIncoming ChangeType = 2
	// 红包过期退还
	EnvelopeExpiredRefund ChangeType = 3
)

// 资金交易变化标识
type ChangeFlag int8

const (
	// 创建账户
	FlagAccountCreated ChangeFlag = 0
	// 支出
	FlagTransferOut ChangeFlag = -1
	// 收入
	FlagTransferIn ChangeFlag = 1
)

type AccountType int8

const (
	UserAccount   AccountType = 1
	SystemAccount AccountType = 2
)
