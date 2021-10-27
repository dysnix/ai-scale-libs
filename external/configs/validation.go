package configs

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/robfig/cron/v3"
	"github.com/xhit/go-str2duration/v2"
)

const (
	GRPCHostTag             = "grpc_host"
	RequiredIfNotNilOrEmpty = "require_if_not_nil_or_empty"
	CronTag                 = "cron"
	HostIfEnabledTag        = "host_if_enabled"
	PortIfEnabledTag        = "port_if_enabled"
	DurationTag             = "duration"
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

	if err = validator.RegisterValidation(HostIfEnabledTag, ValidateHostIfEnabled(validator)); err != nil {
		return err
	}

	if err = validator.RegisterValidation(PortIfEnabledTag, ValidatePortIfEnabled(validator)); err != nil {
		return err
	}

	if err = validator.RegisterValidation(DurationTag, ValidateDuration); err != nil {
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
	)

	if len(params) >= 1 {
		otherFieldName = params[0]
	} else {
		otherFieldName = fl.Param()
	}

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

	return true
}

// ValidateHostIfEnabled implements validator.Func for validate host string if struct enabled
func ValidateHostIfEnabled(validatorMain *validator.Validate) func(level validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		var otherFieldVal reflect.Value
		if fl.Parent().Kind() == reflect.Ptr {
			otherFieldVal = fl.Parent().Elem().FieldByName("Enabled")
		} else {
			otherFieldVal = fl.Parent().FieldByName("Enabled")
		}

		if !otherFieldVal.IsZero() {
			if err := validatorMain.Var(fl.Field().String(), "hostname"); err == nil {
				return true
			}

			if err := validatorMain.Var(fl.Field().String(), "ip"); err == nil {
				return true
			}

			return false
		}

		return true
	}
}

// ValidatePortIfEnabled implements validator.Func for validate port value if struct enabled
func ValidatePortIfEnabled(validatorMain *validator.Validate) func(level validator.FieldLevel) bool {
	return func(fl validator.FieldLevel) bool {
		var otherFieldVal reflect.Value
		if fl.Parent().Kind() == reflect.Ptr {
			otherFieldVal = fl.Parent().Elem().FieldByName("Enabled")
		} else {
			otherFieldVal = fl.Parent().FieldByName("Enabled")
		}

		if !otherFieldVal.IsZero() {
			result := fl.Field().Uint()
			err := validatorMain.Var(result, "numeric,gte=6060,lte=9999")
			if err == nil {
				return true
			}

			return false
		}

		return true
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

// ValidateDuration implements validator.Func for validate duration values
func ValidateDuration(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if len(field) > 0 {
		if _, err := str2duration.ParseDuration(field); err == nil {
			return true
		}
	}
	return false
}
