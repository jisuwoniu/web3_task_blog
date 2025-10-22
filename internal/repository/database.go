package repository

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
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

func GetDB() (*gorm.DB, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
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
		dateBaseConfig.User, dateBaseConfig.Password, dateBaseConfig.Host, dateBaseConfig.Port, dateBaseConfig.DBName, dateBaseConfig.Charset, dateBaseConfig.ParseTime, dateBaseConfig.Loc)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func GetTestDB() (*gorm.DB, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath("../../configs")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
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
		dateBaseConfig.User, dateBaseConfig.Password, dateBaseConfig.Host, dateBaseConfig.Port, dateBaseConfig.DBName, dateBaseConfig.Charset, dateBaseConfig.ParseTime, dateBaseConfig.Loc)
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
