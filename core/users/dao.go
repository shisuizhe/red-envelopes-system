package users

// import "github.com/tietang/dbx"
//
// type UserDao struct {
// 	runner *dbx.TxRunner
// }
//
// func (dao *UserDao) GetOne(password string) *User {
// 	form := &User{Password: password}
// 	ok, err := dao.runner.GetOne(form)
// 	if err != nil {
// 		return nil
// 	}
// 	if !ok {
// 		return nil
// 	}
// 	return form
// }
//
// func (dao *UserDao) Insert(form *User) (int64, error) {
// 	rs, err := dao.runner.Insert(form)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return rs.LastInsertId()
// }
//
// func (dao *UserDao) Update(po *User) (int64, error) {
// 	rs, err := dao.runner.Update(po)
// 	if err != nil {
// 		return 0, err
// 	}
// 	return rs.RowsAffected()
// }

import "github.com/tietang/dbx"

type UserDao struct {
	runner *dbx.TxRunner
}

func (dao *UserDao) GetOne(mobile string) *User {
	form := &User{Mobile: mobile}
	ok, err := dao.runner.GetOne(form)
	if err != nil {
		return nil
	}
	if !ok {
		return nil
	}
	return form
}

func (dao *UserDao) Insert(form *User) (int64, error) {
	rs, err := dao.runner.Insert(form)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

func (dao *UserDao) Update(po *User) (int64, error) {
	rs, err := dao.runner.Update(po)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}