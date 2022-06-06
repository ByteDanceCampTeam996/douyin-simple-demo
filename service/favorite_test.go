package service

import (
	"testing"
)

//基本函数的单元测试
func TestGetFavoriteList(t *testing.T) {
	_, res := GetFavoriteList(1)
	t.Logf("hash returns %s\n", res)
}
