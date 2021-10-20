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
		Profiling      Profiling
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
		{
			name: "4. invalid validation of 'host_if_enabled' tag",
			opts: []*tmpStruct{
				{
					SomeTypeFirst:  false,
					SomeTypeSecond: 0,
					SomeCron:       "* * * * *",
					Profiling: Profiling{
						Enabled: true,
						Host:    "",
						Port:    8080,
					},
				},
			},
			want: `Key: 'tmpStruct.Profiling.Host' Error:Field validation for 'Host' failed on the 'host_if_enabled' tag`,
		},
		{
			name: "5. invalid validation of 'port_if_enabled' tag",
			opts: []*tmpStruct{
				{
					SomeTypeFirst:  false,
					SomeTypeSecond: 0,
					SomeCron:       "* * * * *",
					Profiling: Profiling{
						Enabled: true,
						Host:    "localhost",
						Port:    0,
					},
				},
			},
			want: `Key: 'tmpStruct.Profiling.Port' Error:Field validation for 'Port' failed on the 'port_if_enabled' tag`,
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
