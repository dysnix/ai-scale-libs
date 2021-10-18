package configs

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
)

const (
	GRPCHostTag             = "grpc_host"
	RequiredIfNotNilOrEmpty = "require_if_not_nil_or_empty"
)

// RegisterCustomValidationsTags registers all custom validation tags
func RegisterCustomValidationsTags(ctx context.Context, validator *validator.Validate, in map[string]func(fl validator.FieldLevel) bool) error {
	eg, _ := errgroup.WithContext(ctx)
	for tag, _ := range in {
		newTag := tag
		newCallback := in[tag]
		eg.Go(func() error {
			return validator.RegisterValidation(newTag, newCallback)
		})
	}

	eg.Go(func() error {
		return validator.RegisterValidation(GRPCHostTag, ValidateGRPCHost(validator))
	})

	eg.Go(func() error {
		return validator.RegisterValidation(RequiredIfNotNilOrEmpty, ValidateRequiredIfNotEmpty)
	})

	return eg.Wait()
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
		otherFieldName     string
		otherFieldValCheck string
	)

	if len(params) >= 1 {
		otherFieldName = params[0]
	} else {
		otherFieldName = fl.Param()
	}

	if len(params) >= 2 {
		otherFieldValCheck = params[1]
	} else {
		otherFieldValCheck = fl.Param()
	}

	var otherFieldVal reflect.Value
	if fl.Parent().Kind() == reflect.Ptr {
		otherFieldVal = fl.Parent().Elem().FieldByName(otherFieldName)
	} else {
		otherFieldVal = fl.Parent().FieldByName(otherFieldName)
	}

	if otherFieldValCheck != otherFieldVal.String() {
		if !isNilOrZeroValue(fl.Field()) {
			return true
		}
		return false
	}

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
