package account

import (
	"database/sql"
	"github.com/shopspring/decimal"
	"resk/services"
	"time"
)

// 持久化对象是ORM映射的基础
// 1. dbx支持自动映射名称，默认是把驼峰命名转换成下划线命名
// 2. 表名称默认是结构体名称转换成下划线命名来映射
// 3. 字段名称默认是field name转换成下划线命名来映射，字段映射时使用db作为tag
// 4. 使用uni|unique的tag值来标识字段为唯一索引字段
// 5. 使用id|pk的tag值来标识字段为主键
// 6. 使用omitempty的tag值来标识字段更新和写入时会被忽略
// 7. 使用-中划线的tag值来标识字段在更新，写入、查询时会被忽略

// 账户持久化对象
type Account struct {
	Id           int64           `db:"id,omitempty"`         // 账户ID
	AccountNo    string          `db:"account_no,uni"`       // 账户编号,账户唯一标识
	AccountName  string          `db:"account_name"`         // 账户名称,用来说明账户的简短描述,账户对应的名称或者命名，比如xx积分、xx零钱
	AccountType  int             `db:"account_type"`         // 账户类型，用来区分不同类型的账户：积分账户、会员卡账户、钱包账户、红包账户
	CurrencyCode string          `db:"currency_code"`        // 货币类型编码：CNY人民币，EUR欧元，USD美元 。。。
	UserId       string          `db:"user_id"`              // 用户编号, 账户所属用户
	UserName     sql.NullString  `db:"user_name"`            // 用户名称
	Balance      decimal.Decimal `db:"balance"`              // 账户可用余额
	Status       int             `db:"status"`               // 账户状态，账户状态：0账户初始化，1启用，2停用
	CreatedAt    time.Time       `db:"created_at,omitempty"` // 创建时间
	UpdatedAt    time.Time       `db:"updated_at,omitempty"` // 更新时间
}

func (po *Account) ToDTO() *services.AccountDTO {
	dto := &services.AccountDTO{}
	dto.AccountNo = po.AccountNo
	dto.AccountName = po.AccountName
	dto.AccountType = po.AccountType
	dto.CurrencyCode = po.CurrencyCode
	dto.UserId = po.UserId
	dto.UserName = po.UserName.String
	dto.Balance = po.Balance
	dto.Status = po.Status
	dto.CreatedAt = po.CreatedAt
	dto.UpdatedAt = po.UpdatedAt
	return dto
}

func (po *Account) FromDTO(dto *services.AccountDTO) {
	po.AccountNo = dto.AccountNo
	po.AccountName = dto.AccountName
	po.AccountType = dto.AccountType
	po.CurrencyCode = dto.CurrencyCode
	po.UserId = dto.UserId
	po.UserName = sql.NullString{Valid: true, String: dto.UserName}
	po.Balance = dto.Balance
	po.Status = dto.Status
	po.CreatedAt = dto.CreatedAt
	po.UpdatedAt = dto.UpdatedAt
}
