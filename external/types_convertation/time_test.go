package types_convertation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseMillisecondUnixTimestamp(t *testing.T) {
	testData := 1639473822

	ts, err := ParseMillisecondUnixTimestamp(testData)
	assert.NoError(t, err)
	t.Log(ts)
}
