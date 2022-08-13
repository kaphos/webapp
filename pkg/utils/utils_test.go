package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandStr(t *testing.T) {
	str1 := RandStr(32)
	str2 := RandStr(32)

	assert.NotEmpty(t, str1)
	assert.NotEmpty(t, str2)
	assert.NotEqual(t, str1, str2)
}

func TestRandInt(t *testing.T) {
	int1 := RandInt()
	int2 := RandInt()

	assert.NotEmpty(t, int1)
	assert.NotEmpty(t, int2)
	assert.NotEqual(t, int1, int2)
}
