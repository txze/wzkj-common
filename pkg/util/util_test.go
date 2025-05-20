package util_test

import (
	"testing"

	"github.com/txze/wzkj-common/pkg/util"

	"github.com/stretchr/testify/assert"
)

func TestTrimRepeat(t *testing.T) {
	var strs []string

	strs = []string{"a", "a", "b"}
	strs = util.TrimRepeatString(strs)
	assert.Equal(t, 2, len(strs))

	var ints []int
	ints = []int{1, 2, 3, 1}
	ints = util.TrimRepeatInt(ints)
	assert.Equal(t, 3, len(ints))
}
