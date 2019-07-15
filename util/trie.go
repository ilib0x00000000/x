// Copyright 2019 ilib0x00000000.  All rights reserved.
// 由于 github 上开源的字典树使用 rune 最为节点的值
// 源码参见: https://github.com/derekparker/trie
// 而本项目只需要记录 URL，使用 rune 浪费内存且效率相对低下
// 故模仿源码实现了以 byte 为节点值的字典树
//
// 但是在开发过程中，发现使用字典树存储 GFW 黑名单存在 bug
// 望后来人可以修复此问题
package util

import "sync"

type node struct {
	val      byte
	term     bool
	depth    int
	meta     interface{}
	mask     uint64
	parent   *node
	children map[byte]*node
}

const null = 0x0

// NewChild 创建子节点
func (n *node) NewChild(val byte, bitmask uint64, meta interface{}, term bool) *node {
	node := &node{
		val:      val,
		mask:     bitmask,
		term:     term,
		meta:     meta,
		parent:   n,
		children: make(map[byte]*node),
		depth:    n.depth + 1,
	}
	n.children[val] = node
	n.mask |= bitmask

	return node
}

func (n *node) RemoveChild(r byte) {
	delete(n.children, r)

	for nd := n.parent; nd != nil; nd = nd.parent {
		nd.mask ^= nd.mask
		nd.mask |= uint64(1) << uint64(nd.val-'a')

		for _, c := range nd.children {
			nd.mask |= c.mask
		}
	}
}

func (n *node) Parent() *node {
	return n.parent
}

func (n *node) Meta() interface{} {
	return n.meta
}

func (n *node) Children() map[byte]*node {
	return n.children
}

func (n *node) Terminating() bool {
	return n.term
}

func (n *node) Value() byte {
	return n.val
}

func (n *node) Depth() int {
	return n.depth
}

func (n *node) Mask() uint64 {
	return n.mask
}

// findNode 字典树查找
func findNode(node *node, domain []byte) *node {
	if node == nil {
		return nil
	}

	if len(domain) == 0 {
		return node
	}

	n, ok := node.Children()[domain[0]]
	if !ok {
		return nil
	}

	var offset []byte
	if len(domain) > 1 {
		offset = domain[1:]
	} else {
		offset = domain[0:0]
	}

	return findNode(n, offset)
}

// findTopNode 最多匹配到二级域名
func findTopNode(root *node, domain []byte, flag bool) *node {
	if root == nil {
		return nil
	}

	if len(domain) == 0 {
		return root
	}

	key := domain[0]
	if flag && key == '.' {
		return root
	}

	if !flag && key == '.' {
		flag = true
	}

	n, ok := root.Children()[domain[0]]
	if !ok {
		return nil
	}

	var offset []byte
	if len(domain) > 1 {
		offset = domain[1:]
	} else {
		offset = domain[0:0]
	}

	return findTopNode(n, offset, flag)
}

// maskByteSlice 计算字符串的掩码
func maskByteSlice(bs []byte) uint64 {
	var m uint64
	for _, r := range bs {
		m |= uint64(1) << uint64(r-'a')
	}
	return m
}

type Trie struct {
	mu   sync.Mutex
	size int
	root *node
}

// Root 返回跟节点信息
func (t *Trie) Root() *node {
	return t.root
}

// Add 添加一个字符串
func (t *Trie) Add(key string, meta interface{}) *node {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.size++
	domain := []byte(key)
	bitmask := maskByteSlice(domain)

	node := t.root
	node.mask |= bitmask

	for i := range domain {
		b := domain[i]
		bitmask = maskByteSlice(domain[i:])
		if n, ok := node.children[b]; ok {
			node = n
			node.mask |= bitmask
		} else {
			node = node.NewChild(b, bitmask, nil, false)
		}
	}

	node = node.NewChild(null, 0, meta, true)
	return node
}

func (t *Trie) Find(key string) (*node, bool) {
	node := findNode(t.Root(), []byte(key))
	if node == nil {
		return nil, false
	}

	node, ok := node.Children()[null]
	if !ok || !node.term {
		return nil, false
	}

	return node, true
}

// HasKeysWithPrefix 前缀匹配
func (t *Trie) HasKeysWithPrefix(key string) bool {
	node := findTopNode(t.Root(), []byte(key), false)
	return node != nil
}

// Remove 删除
func (t *Trie) Remove(key string) {
	var (
		i    int
		bs   = []byte(key)
		node = findNode(t.Root(), []byte(key))
	)

	t.mu.Lock()
	defer t.mu.Unlock()

	t.size--
	for n := node.Parent(); n != nil; n = n.Parent() {
		i++
		if len(n.Children()) > 1 {
			b := bs[len(bs)-i]
			n.RemoveChild(b)
			break
		}
	}
}

func NewTrie() *Trie {
	return &Trie{
		size: 0,
		root: &node{
			children: make(map[byte]*node),
			depth:    0,
		},
	}
}
