package repository

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"sync"
)

type DatabaseConfig struct {
	Host      string
	Port      int
	User      string
	Password  string
	DBName    string
	Charset   string
	ParseTime bool
	Loc       string
}

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// 初始化数据库连接（单例模式）
func InitDB() {
	var err error
	once.Do(func() {
		dbInstance, err = createDBConnection()
	})
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

// GetDB 获取数据库实例
func GetDB() (*gorm.DB, error) {
	if dbInstance == nil {
		return nil, fmt.Errorf("database not initialized. Call InitDB() first")
	}
	return dbInstance, nil
}

// createDBConnection 创建数据库连接
func createDBConnection() (*gorm.DB, error) {
	dateBaseConfig := DatabaseConfig{
		Host:      viper.GetString("database.host"),
		Port:      viper.GetInt("database.port"),
		User:      viper.GetString("database.user"),
		Password:  viper.GetString("database.password"),
		DBName:    viper.GetString("database.dbname"),
		Charset:   viper.GetString("database.charset"),
		ParseTime: viper.GetBool("database.parseTime"),
		Loc:       viper.GetString("database.loc"),
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		dateBaseConfig.User, dateBaseConfig.Password, dateBaseConfig.Host,
		dateBaseConfig.Port, dateBaseConfig.DBName, dateBaseConfig.Charset,
		dateBaseConfig.ParseTime, dateBaseConfig.Loc)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
