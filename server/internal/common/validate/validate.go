package validate

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// v is a single shared validate - the library caches struct info, so we reuse
// one instance
var v = validator.New()

func init() {
	// Report errors using the json field name (e.g. "new_password") not Go name ("NewPassword")
	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		// take the part before the comma in `json:"new_password,omitempty"`
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Struct checks a struct's `validate` tags - returns "" if valid, or one
// readable message
func Struct(s any) string {
	err := v.Struct(s)
	if err == nil {
		return  ""
	}
	// error is a list of field faiures - we surface the first one as a friendly message
	for _, fe := range err.(validator.ValidationErrors) {
		return message(fe)
	}

	return  "invalid request"
}

// message turns one field error into human text like "email must be a valid email"
func message(fe validator.FieldError) string {
	field := fe.Field()
	switch fe.Tag() {
		
	case "required":
		return  field + " is required"

	case "email":
		return field + " must be a valid email"
	
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
	
	default:
		return field + " is invalid"
	}
}
