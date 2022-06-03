package service

import (
	"testing"

	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

//基本函数的单元测试
func TestGetFavoriteList(t *testing.T) {
	res := GetFavoriteList(1)
	t.Logf("hash returns %s\n", res)
}
