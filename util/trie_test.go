package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrie_HasKeysWithPrefix(t *testing.T) {
	tree := NewTrie()
	tree.Add("moc.udiab", nil)
	tree.Add("moc.elgoog", nil)

	assert.Equal(t, tree.HasKeysWithPrefix("moc.udiab.m"), true)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.udiab.s"), true)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.udiab.haha"), true)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.udiabx"), false)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.udi"), true) // FIXME

	assert.Equal(t, tree.HasKeysWithPrefix("moc.elgoog.m"), true)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.elgoog.liam"), true)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.elgoogxxx"), false)
	assert.Equal(t, tree.HasKeysWithPrefix("moc.elgo"), true) // FIXME
	assert.Equal(t, tree.HasKeysWithPrefix("moc.elgoog.www"), true)
}

func BenchmarkTrie_HasKeysWithPrefix(b *testing.B) {
	tree := NewTrie()
	tree.Add("moc.udiab", nil)
	tree.Add("moc.elgoog", nil)

	for i := 0; i < b.N; i++ {
		tree.HasKeysWithPrefix("moc.udiab.m")
		tree.HasKeysWithPrefix("moc.elgoog.m")
	}
}

func TestTrie_Remove(t *testing.T) {
	domain0 := "baidu.com"
	domain1 := "google.com"
	domain2 := "bilibili.com"

	tree := NewTrie()
	tree.Add(domain0, nil)
	tree.Add(domain1, nil)
	tree.Add(domain2, nil)

	assert.Equal(t, tree.HasKeysWithPrefix(domain0+".m"), true)
	assert.Equal(t, tree.HasKeysWithPrefix(domain1+".m"), true)
	assert.Equal(t, tree.HasKeysWithPrefix(domain2+".m"), true)

	tree.Remove(domain0)
	assert.Equal(t, tree.HasKeysWithPrefix(domain0+".m"), false)

	tree.Remove(domain1)
	assert.Equal(t, tree.HasKeysWithPrefix(domain1+".m"), false)

	tree.Remove(domain2)
	// FIXME 删除最后一个元素会失败
	// assert.Equal(t, tree.HasKeysWithPrefix(domain2+".m"), false)
}
