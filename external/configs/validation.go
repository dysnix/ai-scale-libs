package configs

import (
	"context"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
)

var Validate = validator.New()

// RegisterCustomValidationsTags registers all custom validation tags
func RegisterCustomValidationsTags(ctx context.Context, in map[string]func(fl validator.FieldLevel) bool) error {
	eg, _ := errgroup.WithContext(ctx)
	for tag, callback := range in {
		eg.Go(func() error {
			return Validate.RegisterValidation(tag, callback)
		})
	}

	return eg.Wait()
}

// ValidateGRPCHost implements validator.Func
func ValidateGRPCHost(fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if len(field) > 0 {
		if err := Validate.Var(field, "hostname"); err == nil {
			return true
		}
	}
	return true
}
