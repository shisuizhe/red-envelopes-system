package account

import (
	"context"
	"errors"
	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
)

type accountDomain struct {
	account    Account
	accountLog AccountLog
}

func NewAccountDomain() *accountDomain {
	return new(accountDomain)
}

// 创建账户编号
func (ad *accountDomain) createAccountNo() {
	ad.account.AccountNo = ksuid.New().Next().String()
}

// 创建账户流水编号
func (ad *accountDomain) createAccountLogNo() {
	ad.accountLog.LogNo = ksuid.New().Next().String()
}

// 创建流水的记录
func (ad *accountDomain) CreateAccountLog() {
	// 通过account来创建流水
	ad.accountLog = AccountLog{}
	ad.createAccountLogNo()
	ad.accountLog.TradeNo = ad.accountLog.LogNo
	// 流水中的交易主体信息
	ad.accountLog.AccountNo = ad.account.AccountNo
	ad.accountLog.UserId = ad.account.UserId
	ad.accountLog.UserName = ad.account.UserName.String
	// 交易对象信息
	ad.accountLog.TargetAccountNo = ad.account.AccountNo
	ad.accountLog.TargetUserId = ad.account.UserId
	ad.accountLog.TargetUserName = ad.account.UserName.String
	// 交易金额
	ad.accountLog.Amount = ad.account.Balance
	ad.accountLog.Balance = ad.account.Balance
	// 交易变化信息
	ad.accountLog.Decs = "账户创建"
	ad.accountLog.ChangeType = services.AccountCreated
	ad.accountLog.ChangeFlag = services.FlagAccountCreated
}

// 创建账户
func (ad *accountDomain) Create(adto services.AccountDTO) (*services.AccountDTO, error) {
	// 创建账户持久化对象
	ad.account = Account{}
	ad.account.FromDTO(&adto)
	ad.createAccountNo()
	ad.account.UserName.Valid = true
	// 创建账户流水持久化对象
	ad.CreateAccountLog()
	adao := AccountDao{}
	ldao := AccountLogDao{}
	var ndto *services.AccountDTO
	err := base.Tx(func(runner *dbx.TxRunner) error {
		adao.runner = runner
		ldao.runner = runner
		// 插入账户数据
		id, err := adao.Insert(&ad.account)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户失败")
		}
		// 如果插入成功就插入流水数据
		id, err = ldao.Insert(&ad.accountLog)
		if err != nil {
			return err
		}
		if id <= 0 {
			return errors.New("创建账户流水失败")
		}
		ad.account = *adao.GetByAccountNo(ad.account.AccountNo)
		return nil
	})
	ndto = ad.account.ToDTO()
	return ndto, err
}

// 转账
func (ad *accountDomain) Transfer(dto services.AccountTransferDTO) (status services.TransferStatus, err error) {
	// 验证输入参数
	if err := base.ValidateStruct(&dto); err != nil {
		return services.TransferStatusFailed, err
	}

	// 类型转换
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferStatusFailed, err
	}
	dto.Amount = amount

	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeFlag > 0 {
			return services.TransferStatusFailed, errors.New("如果changeFlag为支出，那么changeType必须小于0")
		}
		// 如果余额小于要转账的金额，则无法进行转账操作
		bodyUserId := dto.TradeBody.UserId
		a := ad.GetAccountByUserId(bodyUserId)
		if a.Balance.LessThan(dto.Amount) {
			return services.TransferStatusSufficientFunds, errors.New("余额不足")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferStatusFailed, errors.New("如果changeFlag为收入，那么changeType必须大于0")
		}
	}

	// 进行转账操作
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		status, err = ad.TransferWithContextTx(ctx, dto)

		if err == nil {
			ldao := AccountLogDao{runner: runner}
			err = ldao.UpdateDecs()
			if err != nil {
				return err
			}
		}
		return err
	})
	return status, err
}

// 必须在base.Tx事务块中运行，不能单独运行
func (ad *accountDomain) TransferWithContextTx(
	ctx context.Context,
	dto services.AccountTransferDTO) (status services.TransferStatus, err error) {
	// 如果交易变化是支出，修正amount
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = amount.Mul(decimal.NewFromFloat(-1))
	}
	// 创建账户流水记录
	ad.accountLog = AccountLog{}
	ad.accountLog.FromTransferDTO(&dto)
	ad.createAccountLogNo()
	// 检查余额是否足够和更新余额：通过乐观锁来验证，更新余额的同时来验证余额是否足够
	// 更新成功后，写入流水记录
	err = base.ExecuteContext(ctx, func(runner *dbx.TxRunner) error {
		adao := AccountDao{runner: runner}
		ldao := AccountLogDao{runner: runner}
		rows, err := adao.UpdateBalance(dto.TradeBody.AccountNo, amount)
		if err != nil || rows <= 0 {
			status = services.TransferStatusFailed
			return err
		}
		// if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
		// 	status = services.TransferStatusSufficientFunds
		// 	return errors.New("余额不足")
		// }
		account := adao.GetByAccountNo(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("账户不存在")
		}
		ad.account = *account
		ad.accountLog.Balance = ad.account.Balance
		id, err := ldao.Insert(&ad.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferStatusFailed
			return errors.New("账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
	} else {
		status = services.TransferStatusSuccess
	}
	return status, err
}

// 根据账户编号来查询账户信息
func (ad *accountDomain) GetAccountByAccountNo(accountNo string) *services.AccountDTO {
	adao := AccountDao{}
	var account *Account
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		adao.runner = runner
		account = adao.GetByAccountNo(accountNo)
		return nil
	})
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据用户ID来查询账户信息
func (ad *accountDomain) GetAccountByUserId(userId string) *services.AccountDTO {
	adao := AccountDao{}
	var account *Account
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		adao.runner = runner
		account = adao.GetByUserId(userId)
		return nil
	})
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

// 根据流水编号查询账户流水
func (ad *accountDomain) GetAccountLogByLogNo(logNo string) *services.AccountLogDTO {
	ldao := AccountLogDao{}
	var accountLog *AccountLog
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		ldao.runner = runner
		accountLog = ldao.GetByLogNo(logNo)
		return nil
	})
	if accountLog == nil {
		return nil
	}
	return accountLog.ToDTO()
}

// 根据交易编号来查询账户流水
func (ad *accountDomain) GetAccountLogByTradeNo(tradeNo string) *services.AccountLogDTO {
	ldao := AccountLogDao{}
	var accountLog *AccountLog
	_ = base.Tx(func(runner *dbx.TxRunner) error {
		ldao.runner = runner
		accountLog = ldao.GetByTradeNo(tradeNo)
		return nil
	})
	if accountLog == nil {
		return nil
	}
	return accountLog.ToDTO()
}
