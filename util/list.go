package util

import (
	"log"
	"time"
)

// BlackList 黑名单[字典树]
type BlackList struct {
	list *Trie // GFW 名单上的列表
	temp *Trie // 不在 GFW 名单上，但是直接访问不了的会加入
}

// LoadGFWList 加载 GFW 名单上数据
func (s *BlackList) LoadGFWList() {
	domains, err := fetchBlockedDomains()
	if err != nil {
		log.Println("load gfw list faild: ", err)
	}

	for _, d := range domains {
		s.list.Add(strrev(d), nil)
	}
}

// Add 动态发现被墙的域名时，加入黑名单
func (s *BlackList) Add(host string) {
	// TODO 定时清理
	// 如果已经不能直接访问，加入 list
	// 如果可以访问 解除黑名单
	s.temp.Add(strrev(host), time.Now())
}

// IsBlacked 判断域名是否被墙
func (s *BlackList) IsBlacked(host string) bool {
	host = strrev(host)

	if s.temp.HasKeysWithPrefix(host) {
		return true
	}

	if s.list.HasKeysWithPrefix(host) {
		return true
	}

	return false
}

// strrev 字符串倒序 abc  -->  cba
func strrev(s string) string {
	b := []byte(s)
	length := len(s)

	for start, end := 0, length-1; start < end; start++ {
		b[start], b[end] = s[end], s[start]
		end--
	}

	return string(b)
}

// NewBlackList 返回黑名单句柄
func NewBlackList() *BlackList {
	return &BlackList{
		list: NewTrie(),
		temp: NewTrie(),
	}
}
