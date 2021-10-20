package configs

import (
	"context"
	//"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
)

const (
	GRPCHostTag             = "grpc_host"
	RequiredIfNotNilOrEmpty = "require_if_not_nil_or_empty"
	CronTag                 = "cron"
)

// RegisterCustomValidationsTags registers all custom validation tags
func RegisterCustomValidationsTags(ctx context.Context, validator *validator.Validate, in map[string]func(fl validator.FieldLevel) bool) (err error) {
	for tag, _ := range in {
		newTag := tag
		newCallback := in[tag]
		if err = validator.RegisterValidation(newTag, newCallback); err != nil {
			return err
		}
	}

	if err = validator.RegisterValidation(GRPCHostTag, ValidateGRPCHost(validator)); err != nil {
		return err
	}

	if err = validator.RegisterValidation(RequiredIfNotNilOrEmpty, ValidateRequiredIfNotEmpty); err != nil {
		return err
	}

	if err = validator.RegisterValidation(CronTag, ValidateCronString); err != nil {
		return err
	}

	return err
}

// ValidateGRPCHost implements validator.Func
func ValidateGRPCHost(validatorMain *validator.Validate) func(level validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		field := fl.Field().String()
		if len(field) > 0 {
			if err := validatorMain.Var(field, "hostname"); err == nil {
				return true
			}

			return false
		}
		return true
	}
}

// ValidateRequiredIfNotEmpty implements validator.Func for check some field not empty or nil
// Example for usage:
//	Type        string	`json:"type" validate:"required,oneof=flat_off percent_off"`
//	MaxValue	uint	`json:"max_value" validate:"required_if=Type percent_off"`
func ValidateRequiredIfNotEmpty(fl validator.FieldLevel) bool {
	params := strings.Split(fl.Param(), " ")
	if len(params) == 0 {
		params = strings.Split(fl.Param(), ",")
	}

	var (
		otherFieldName string
		//otherFieldValCheck string
	)

	if len(params) >= 1 {
		otherFieldName = params[0]
	} else {
		otherFieldName = fl.Param()
	}

	//if len(params) >= 2 {
	//	otherFieldValCheck = params[1]
	//} else {
	//	otherFieldValCheck = fl.Param()
	//}

	var otherFieldVal reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		otherFieldVal = fl.Parent().Elem().FieldByName(otherFieldName)
	} else {
		otherFieldVal = fl.Parent().FieldByName(otherFieldName)
	}

	if !otherFieldVal.IsZero() {
		result := !fl.Field().IsZero()
		return result
	}
	//if otherFieldValCheck != otherFieldVal.String() {
	//	fmt.Println(otherFieldVal.IsZero().String(), fl.Field().String(), fl.Field().IsZero())
	//	if isNilOrZeroValue(fl.Field()) {
	//		return true
	//	}
	//	return false
	//}

	return true
}

func isNilOrZeroValue(field reflect.Value) bool {
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsZero()
	}
}

// ValidateCronString implements validator.Func for validate cron string format
func ValidateCronString(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if len(field) > 0 {
		if _, err := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor).
			Parse(field); err == nil {
			return true
		}

		return false
	}
	return true
}
