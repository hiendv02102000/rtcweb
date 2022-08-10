package validators

import (
	"regexp"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func CustomValidate() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("password", passwordCheck)
	}
}
func passwordCheck(fl validator.FieldLevel) bool {
	password, ok := fl.Field().Interface().(string)
	if ok {
		if len(password) < 8 {
			return false
		}
		if m, _ := regexp.MatchString("[0-9]+", password); !m {
			return false
		}
		if m, _ := regexp.MatchString("[a-z]+", password); !m {
			return false
		}
		if m, _ := regexp.MatchString("[A-Z]+", password); !m {
			return false
		}
	}
	return true
}
