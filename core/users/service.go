package users

import (
	log "github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
	"resk/infra/base"
	"resk/services"
	"strconv"
)

// type UserService struct{}
//
// func (u *UserService) Login(Password, UserName string) (user *User) {
// 	as := services.GetAccountServiceInterface()
// 	err := base.Tx(func(runner *dbx.TxRunner) error {
// 		dao := UserDao{runner: runner}
// 		user = dao.GetOne(Password)
// 		// 创建用户
// 		if user == nil {
// 			user = &User{
// 				Password: Password,
// 				UserName: UserName,
// 			}
// 			id, err := dao.Insert(user)
// 			if err != nil {
// 				return err
// 			}
// 			user.Id = id
// 			user.UserId = strconv.Itoa(int(id))
// 		}
// 		if user.UserName != UserName {
// 			user.UserName = UserName
// 			rows, err := dao.Update(user)
// 			if err != nil {
// 				logrus.Error(err)
// 				return err
// 			}
// 			if rows <= 0 {
// 				logrus.Warn("non updated:", user)
// 			}
// 		}
// 		return nil
//
// 	})
// 	if err != nil {
// 		logrus.Error(err)
// 		return nil
// 	}
// 	// 创建资金账户
// 	a := as.GetEnvelopeAccountByUserId(user.UserId)
// 	if a == nil {
// 		dto := services.AccountCreatedDTO{
// 			UserId:       user.UserId,
// 			UserName:     user.UserName,
// 			AccountName:  user.UserName,
// 			AccountType:  int(services.UserAccount),
// 			CurrencyCode: services.DefaultCurrencyCode,
// 			Amount:       "1000000000000",
// 		}
// 		_, err := as.CreateAccount(dto)
// 		if err != nil {
// 			logrus.Error(err)
// 			return
// 		}
// 	} else {
// 		logrus.Info("已经存在红包资金账户：", user)
// 	}
// 	return
// }

type UserService struct {}

func (u *UserService) Login(mobile, username string) (user *User) {
	as := services.GetAccountServiceInterface()
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := UserDao{runner: runner}
		user = dao.GetOne(mobile)
		//创建用户
		if user == nil {
			user = &User{
				Mobile:   mobile,
				Username: username,
			}
			id, err := dao.Insert(user)
			if err != nil {
				log.Error(err)
				return err
			}
			user.Id = id
			user.UserId = strconv.Itoa(int(id))
		}
		if user.Username != username {
			user.Username = username
			rows, err := dao.Update(user)
			if err != nil {
				log.Error(err)
				return err
			}
			if rows <= 0 {
				log.Warn("non updated:", user)
			}

		}
		return nil

	})
	if err != nil {
		log.Error(err)
		return nil
	}
	user.UserId = strconv.Itoa(int(user.Id))
	//创建资金账户
	a := as.GetEnvelopeAccountByUserId(user.UserId)
	if a == nil {
		dto := services.AccountCreatedDTO{
			UserId:       user.UserId,
			UserName:     user.Username,
			AccountName:  user.Username,
			AccountType:  int(services.UserAccount),
			CurrencyCode: services.DefaultCurrencyCode,
			Amount:       "10000000000000",
		}
		_, err := as.CreateAccount(dto)
		if err != nil {
			log.Error(err)
			return
		}
	} else {
		log.Info("已经存在红包资金账户：", user)
	}
	return
}