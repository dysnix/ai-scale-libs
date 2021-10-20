package configs

import (
	"context"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidateRequiredIfNotEmpty(t *testing.T) {
	type tmpStruct struct {
		SomeTypeFirst  bool          //`validate:"required"`
		SomeTypeSecond time.Duration `validate:"require_if_not_nil_or_empty=SomeTypeFirst"`
		SomeCron       CronStr       `validate:"cron"`
	}

	var cases = []struct {
		name string
		opts []*tmpStruct
		want string
	}{
		{
			name: "1. valid validation",
			opts: []*tmpStruct{
				{
					SomeTypeFirst:  false,
					SomeTypeSecond: 0,
					SomeCron:       "* * * * *",
				},
				{
					SomeTypeFirst:  true,
					SomeTypeSecond: 15,
					SomeCron:       "@every 1m",
				},
			},
			want: "",
		},
		{
			name: "2. invalid validation of 'require_if_not_nil_or_empty' tag",
			opts: []*tmpStruct{
				{
					SomeTypeFirst:  true,
					SomeTypeSecond: 0,
					SomeCron:       "* * * * *",
				},
			},
			want: `Key: 'tmpStruct.SomeTypeSecond' Error:Field validation for 'SomeTypeSecond' failed on the 'require_if_not_nil_or_empty' tag`,
		},
		{
			name: "3. invalid validation of 'cron' tag",
			opts: []*tmpStruct{
				{
					SomeTypeFirst:  false,
					SomeTypeSecond: 0,
					SomeCron:       "* & * * *",
				},
			},
			want: `Key: 'tmpStruct.SomeCron' Error:Field validation for 'SomeCron' failed on the 'cron' tag`,
		},
	}

	testValidator := validator.New()

	err := RegisterCustomValidationsTags(context.Background(), testValidator, nil)
	assert.NoError(t, err)

	for i := range cases {
		tc := cases[i]
		t.Run(tc.name, func(t *testing.T) {
			for j := range tc.opts {
				err = testValidator.Struct(tc.opts[j])
				if valErr, ok := err.(validator.ValidationErrors); ok {
					for _, err := range valErr {
						assert.Equal(t, tc.want, err.Error())
					}
				}
			}
		})
	}
}
