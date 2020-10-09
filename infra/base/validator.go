package base

import (
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
	"github.com/sirupsen/logrus"
	"resk/infra"
)

var validate *validator.Validate
var translator ut.Translator

func Validate() *validator.Validate {
	return validate
}

func Translator() ut.Translator {
	Check(translator)
	return translator
}

type ValidatorStarter struct {
	infra.BaseStarter
}

func (s *ValidatorStarter) Init(infra.StarterContext) {
	validate = validator.New()
	// 创建消息国际化通用翻译器
	cn := zh_Hans.New()
	uni := ut.New(cn, cn)
	var found bool
	translator, found = uni.GetTranslator("zh_Hans")
	if found {
		err := zh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			logrus.Error(err)
		}
	}else {
		logrus.Error("Not found translator: zh_Hans")
	}
}

func ValidateStruct(i interface{}) error {
	err := Validate().Struct(i)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(Translator()))
			}
		}
		return err
	}
	return nil
}
