package service

import (
	"fmt"
	"testing"

	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/dao"
	. "github.com/ByteDanceCampTeam996/douyin-simple-demo/model"
	. "github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

//基本函数的单元测试
func TestDbHashSalt(t *testing.T) {
	res := DbHashSalt("123456", "test1")
	t.Logf("hash returns %s\n", res)
}

func TestGetRandString(t *testing.T) {
	res := GetRandString()
	t.Logf("get rand string returns %s\n", res)
}

func TestIsLegalUserName(t *testing.T) {
	res := IsLegalUserName("or 1=1") && IsLegalUserName("SELECT *")
	if res == true {
		t.Fatalf("fail judge username")
	} else {
		t.Logf("success judge username")
	}
}

//使用Go monkey对函数打桩进行测试
func TestUserAppend(t *testing.T) {
	//对CreateNewUser函数打桩
	patches := ApplyFunc(CreateNewUser, func(_ DbUserInfo) {
		return
	})
	defer patches.Reset()
	//对CreateNewFollowInfo函数打桩
	patches.ApplyFunc(CreateNewFollowInfo, func(_ UserFollowInfo) {
		return
	})
	//测试打桩后的函数输出
	Convey("TestUserAppend", t, func() {
		userid, token := UserAppend("test1", "123456")
		So(userid, ShouldEqual, 1)
		fmt.Print(token)
		So(len(token), ShouldEqual, 20)
	})

}
