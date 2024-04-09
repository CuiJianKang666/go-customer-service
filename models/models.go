package models

import (
	"fmt"
	gorm2 "github.com/jinzhu/gorm"
	"github.com/taoshihan1991/imaptool/common"
	"gorm.io/gorm"
	"log"
	"time"
)

var OldDB *gorm2.DB
var DB *gorm.DB

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func init() {
	Connect()
}
func Connect() error {
	mysql := common.GetMysqlConf()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", mysql.Username, mysql.Password, mysql.Server, mysql.Port, mysql.Database)
	var err error
	OldDB, err = gorm2.Open("mysql", dsn)
	if err != nil {
		log.Println(err)
		panic("数据库连接失败!")
		return err
	}
	OldDB.SingularTable(true)
	OldDB.LogMode(true)
	OldDB.DB().SetMaxIdleConns(10)
	OldDB.DB().SetMaxOpenConns(100)
	OldDB.DB().SetConnMaxLifetime(59 * time.Second)
	InitConfig()
	return nil
}
func Execute(sql string) error {
	return OldDB.Exec(sql).Error
}
func CloseDB() {
	defer OldDB.Close()
}
