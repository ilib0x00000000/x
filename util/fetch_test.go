package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFetchBlackList 测试黑名单拉取不出错
// go test
func TestFetchBlackList(t *testing.T) {
	black, err := fetchBlockedDomains()
	assert.Nil(t, err)

	assert.Equal(t, len(black) > 0, true)

	// for _, url := range black {
	// 	fmt.Println(url)
	// }
}
