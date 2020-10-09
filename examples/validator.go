package main

import (
	"fmt"
	"github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
)

type User struct {
	FirstName string `validate:"required"`
	Age       uint8  `validate:"gte=0,lte=130"`
	Email     string `validate:"required,email"`
}

func main() {
	user := &User{
		FirstName: "",
		Age:       200,
		Email:     "1234@.com",
	}
	validate := validator.New()
	// 创建消息国际化通用翻译器
	cn := zh_Hans.New()
	uni := ut.New(cn, cn)
	translator, found := uni.GetTranslator("zh_Hans")
	if found {
		err := zh.RegisterDefaultTranslations(validate, translator)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("not found")
	}
	err := validate.Struct(user)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			fmt.Println(err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, err := range errs {
				fmt.Println(err.Translate(translator))
			}
		}
	}
}
