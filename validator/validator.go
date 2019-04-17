package validator

import (
	v "gopkg.in/go-playground/validator.v9"
	"reflect"
	"strings"
)

// NewValidator create a validator with custom validations
func NewValidator() *v.Validate {
	validate := v.New()
	validate.RegisterTagNameFunc(JsonFieldNames)

	return validate
}

// JsonFieldNames ignore hidden json fields
func JsonFieldNames(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

	if name == "-" {
		return ""
	}

	return name
}
