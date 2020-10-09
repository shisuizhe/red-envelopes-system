package users

import (
	"time"
)

// type User struct {
// 	Id        int64     `db:"id,omitempty"`
// 	UserId    string    `db:"-"`
// 	UserName  string    `db:"user_name"`
// 	Password  string    `db:"password,uni"`
// 	CreatedAt time.Time `db:"created_at,omitempty"` // 创建时间
// 	UpdatedAt time.Time `db:"updated_at,omitempty"` // 更新时间
// }

type User struct {
	Id        int64     `db:"id,omitempty"`
	UserId    string    `db:"-"`
	Mobile    string    `db:"mobile,uni"`
	Username  string    `db:"username"`
	CreatedAt time.Time `db:"created_at,omitempty"` // 创建时间
	UpdatedAt time.Time `db:"updated_at,omitempty"` // 更新时间
}
