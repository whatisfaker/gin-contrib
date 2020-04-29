package validatoroverriding

import (
	"regexp"

	"gopkg.in/go-playground/validator.v9"
)

var (
	alphanumericUnderscore = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	chinaMobile            = regexp.MustCompile(`^1\d{10}$`)
)

//NewValidateForAlphanumericUnderscore
func NewValidateForAlphanumericUnderscore() (func(validator.FieldLevel) bool, string) {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			return alphanumericUnderscore.MatchString(value)
		}
		return false
	}, "{0}必须由英数字和下划线组成"
}

//NewValidateChinaMobile
func NewValidateChinaMobile() (func(validator.FieldLevel) bool, string) {
	return func(fl validator.FieldLevel) bool {
		if value, ok := fl.Field().Interface().(string); ok {
			return chinaMobile.MatchString(value)
		}
		return false
	}, "{0}手机号码不正确"
}

func RegisterCustomValidateAndTranslate(tag string, text string, fnc func(validator.FieldLevel) bool, validate *validator.Validate) error {
	err := validate.RegisterValidation(tag, fnc)
	if err != nil {
		return err
	}
	return getTransBinding().registerCustomTranslate(tag, text, validate)
}
