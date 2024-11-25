package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Host        string `json:"host" yaml:"host"`
	Port        int    `json:"port" yaml:"port"`
	MysqlDSN    string `json:"MysqlDSN" yaml:"MysqlDSN"`
	MySQLPrefix string `json:"MySQLPrefix" yaml:"MySQLPrefix"`
}

var (
	configFile = flag.String("f", "config/config.yaml", "the config file")
	Configs    *Config
)

func init() {
	// 解析命令行参数
	flag.Parse()
	// 初始化 Viper 配置
	viper.SetConfigFile(*configFile)
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("无法读取配置文件: %w", err))
	}
	// 解析配置到 Configs 变量
	if err := viper.Unmarshal(&Configs); err != nil {
		panic(fmt.Errorf("无法解析配置: %w", err))
	}
}
