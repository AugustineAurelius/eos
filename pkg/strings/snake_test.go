package strings_test

import (
	"testing"

	"github.com/AugustineAurelius/eos/pkg/strings"
	"github.com/stretchr/testify/assert"
)

func Test_Snake(t *testing.T) {
	assert.Equal(t, "user_time", strings.ToSnakeCase("UserTime"))
}
