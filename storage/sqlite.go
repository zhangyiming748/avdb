package storage

import (
	"log"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var gormDB *gorm.DB

// https://www.sqlite.org/download.html
func SetSqlite() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("创建本地sqlite数据库目录失败:%s", err.Error())
	}
	location := filepath.Join(home, "AV.db")
	//location := "duplicate.db"
	db, err := gorm.Open(sqlite.Open(location), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
			NameReplacer:  nil,  // 不使用名称替换器，保持驼峰命名
		},
	})
	if err != nil {
		log.Fatalf("打开本地sqlite数据库失败:%s", err.Error())
	}

	gormDB = db
	log.Println("本地sqlite数据库初始化完成")

	// 自动同步表结构，避免后期插入数据库失败
	if err := gormDB.AutoMigrate(&AVDB{}); err != nil {
		log.Fatalf("同步数据库表结构失败: %s", err.Error())
	}
	log.Println("数据库表结构同步完成")
}

func GetSqlite() *gorm.DB {
	return gormDB
}
