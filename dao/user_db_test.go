package dao

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestUserExistByToken(t *testing.T) {
	//新建一个mock数据库
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	//将mock数据库连接到原有数据库
	Db, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Init Db failed")
	}
	//在mock数据库中新建一行内容
	rows := sqlmock.NewRows([]string{"id", "user_name", "password_hash", "token"}).AddRow(1, "test1", "123456", "1234567")

	//声明预期查询应该返回的结果
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	// now we execute our method
	if _, exist := UserExistByToken("1234567"); exist != true {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	// 确保所有查询都满足了
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
func TestFindUserByToken(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	Db, err = gorm.Open(mysql.New(mysql.Config{
		SkipInitializeWithVersion: true,
		Conn:                      db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("Init Db failed")
	}
	rows := sqlmock.NewRows([]string{"id", "user_name", "password_hash", "token"}).AddRow(1, "test1", "123456", "1234567")
	mock.ExpectQuery("SELECT").WillReturnRows(rows)
	if _, findrow := FindUserByToken("1234567"); findrow.UserName == "123456" {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
