package configs

import (
	"context"

	"github.com/go-playground/validator/v10"
	"golang.org/x/sync/errgroup"
)

const (
	GRPCHostTag = "grpc_host"
)

// RegisterCustomValidationsTags registers all custom validation tags
func RegisterCustomValidationsTags(ctx context.Context, validator *validator.Validate, in map[string]func(fl validator.FieldLevel) bool) error {
	eg, _ := errgroup.WithContext(ctx)
	for tag, callback := range in {
		eg.Go(func() error {
			return validator.RegisterValidation(tag, callback)
		})
	}

	return eg.Wait()
}

// ValidateGRPCHost implements validator.Func
func ValidateGRPCHost(validator *validator.Validate, fl validator.FieldLevel) bool {
	field := fl.Field().String()
	if len(field) > 0 {
		if err := validator.Var(field, "hostname"); err == nil {
			return true
		}
	}
	return true
}
