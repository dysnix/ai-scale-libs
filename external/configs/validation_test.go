package configs

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestValidateRequiredIfNotEmpty(t *testing.T) {
	type tmpStruct struct {
		SomeTypeFirst  bool          //`validate:"required"`
		SomeTypeSecond time.Duration `validate:"require_if_not_nil_or_empty=SomeTypeFirst"`
	}

	testValidator := validator.New()

	err := RegisterCustomValidationsTags(context.Background(), testValidator, nil)
	assert.NoError(t, err)

	t.Run("", func(t *testing.T) {
		testData := &tmpStruct{
			//SomeTypeFirst:  false,
			SomeTypeSecond: 10,
		}
		err = testValidator.Struct(testData)
		assert.NoError(t, err)
	})
}
