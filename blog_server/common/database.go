package common

import (
	"blog_server/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"net/url"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	user := "你的mysql用户名"
	password := "你的mysql用户名"
	host := "localhost"
	port := "3306"
	database := "blog"
	charset := "utf8"
	loc := "Asia/Shanghai"
	args := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=%s&parseTime=true&loc=%s",
		user,
		password,
		host,
		port,
		database,
		charset,
		url.QueryEscape(loc))
	// 连接数据库
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to open database: " + err.Error())
	}
	// 迁移数据表
	db.AutoMigrate(&model.User{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}
