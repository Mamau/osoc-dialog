package recoverer

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"osoc-dialog/pkg/log"
)

func TestRecovererOptions(t *testing.T) {
	logger := log.NewDiscardLogger()

	cases := []struct {
		Name     string
		Options  []Option
		Expected map[string]interface{}
	}{
		{
			Name: "check all options",
			Options: []Option{
				Logger(logger),
			},
			Expected: map[string]interface{}{
				"logger": logger,
			},
		},
	}

	for _, v := range cases {
		t.Run(v.Name, func(t *testing.T) {
			o := options{}
			for _, opt := range v.Options {
				opt(&o)
			}
			assert.Equal(t, v.Expected["logger"], o.logger)
		})
	}
}
