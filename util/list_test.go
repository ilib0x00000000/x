package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrrev(t *testing.T) {
	s := "abc"
	r := strrev(s)
	assert.Equal(t, r, "cba")

	s = "12345"
	r = strrev(s)
	assert.Equal(t, r, "54321")
}

func BenchmarkStrrev(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = strrev("abcd")
	}
}

// TestBlackList
func TestBlackList(t *testing.T) {
	black := NewBlackList()
	assert.NotNil(t, black)

	host := "google.com"
	if !black.IsBlacked(host) {
		black.Add(host)
	}

	assert.Equal(t, black.IsBlacked(host), true)

	host = "baidu.com"
	assert.Equal(t, black.IsBlacked(host), false)
	black.Add(host)
	assert.Equal(t, black.IsBlacked(host), true)
}

// BenchmarkBlackList 压力测试
func BenchmarkBlackList(b *testing.B) {
	host := "github.com"
	black := NewBlackList()

	for i := 0; i < b.N; i++ {
		black.Add(host)
		black.IsBlacked(host)
	}
}
