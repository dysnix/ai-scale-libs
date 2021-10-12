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
