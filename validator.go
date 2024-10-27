package midtrans

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"reflect"
	"strings"
)

func NewValidator() (*validator.Validate, ut.Translator) {
	v := validator.New()
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		f := field.Tag.Get("json")
		fsplit := strings.Split(f, ",")
		if len(fsplit) == 2 {
			return fsplit[1]
		}
		return f
	})

	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = en_translations.RegisterDefaultTranslations(v, trans)
	return v, trans
}
