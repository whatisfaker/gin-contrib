package validatoroverriding

import (
	"reflect"
	"sync"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"gopkg.in/go-playground/validator.v9"

	//"gopkg.in/go-playground/validator.v9"
	zhtranslations "gopkg.in/go-playground/validator.v9/translations/zh"
)

var ins *TransBinding
var once sync.Once

type TransBinding struct {
	uni    *ut.UniversalTranslator
	trans  ut.Translator
	locale string
}

func (c *TransBinding) SetLocale(locale string) {
	if c.locale != locale {
		c.locale = locale
		c.trans, _ = c.uni.GetTranslator(c.locale)
	}
}

func (c *TransBinding) translateFirst(err error) string {
	if es, ok := err.(validator.ValidationErrors); !ok {
		return err.Error()
	} else {
		return es[0].Translate(c.trans)
	}
}

//TranslateFirst 把第一条错误翻译输出
func TranslateFirst(err error) string {
	return getTransBinding().translateFirst(err)
}

func (c *TransBinding) bindValidator(validate *validator.Validate) error {
	err := zhtranslations.RegisterDefaultTranslations(validate, c.trans)
	if err != nil {
		return err
	}
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		tag := fld.Name
		name := fld.Tag.Get("trans")
		if name == "" {
			return tag
		}
		return name
	})
	return nil
}

//BindValidator 使验证器支持中文
func BindValidator(validate *validator.Validate) error {
	return getTransBinding().bindValidator(validate)
}

func getTransBinding() *TransBinding {
	once.Do(func() {
		ins = new(TransBinding)
		zhlocale := zh.New()
		enlocale := en.New()
		ins.uni = ut.New(enlocale, enlocale, zhlocale)
		ins.SetLocale("zh")
	})
	return ins
}

func RegisterCustomTranslate(tag string, text string, validate *validator.Validate) error {
	return getTransBinding().registerCustomTranslate(tag, text, validate)
}

func (c *TransBinding) registerCustomTranslate(tag string, text string, validate *validator.Validate) error {
	return validate.RegisterTranslation(tag, c.trans, func(ut ut.Translator) error {
		return ut.Add(tag, text, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(tag, fe.Field())
		return t
	})
}
