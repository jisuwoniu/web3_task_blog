package utils

import (
	"github.com/spf13/viper"
)

// InitConfig 初始化配置
func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs") // 为测试添加额外的路径
	viper.AddConfigPath("../../configs") // 为测试添加额外的路径
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic("Failed to read config file: " + err.Error())
	}
	JwtKey = []byte(viper.GetString("jwt.secret"))
}