package lib

import (
	"gopkg.in/go-playground/validator.v9"
	"fmt"
)

func ValidateRequest(requestInfo interface{}) bool {
	//初步验证，使用validate
	validate := validator.New()
	err := validate.Struct(requestInfo)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			//错误的结构体的字段，未通过的标签，参数的值
			fmt.Println(err.Namespace())
			fmt.Println(err.Tag())
			fmt.Println(err.Value())
			return false
		}
	}
	return true
}
