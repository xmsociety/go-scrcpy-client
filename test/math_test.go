package test

import (
	"fmt"
	"sort"
	"testing"
)

// minus 获取差集
func minus(a []string, b []string) []string {
	var inter []string
	mp := make(map[string]bool)
	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			delete(mp, s)
		}
	}
	for key := range mp {
		inter = append(inter, key)
	}
	return inter
}

func TestMinux(t *testing.T) {
	newSn := []string{"5555", "5557", "5559", "542"}
	oldSn := []string{"5555", "5558"}
	fmt.Println(minus(newSn, oldSn))
	keep := newSn
	sort.Slice(keep, func(i, j int) bool {
		return i < j
	})
	fmt.Println(keep)
}
