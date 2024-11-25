package database

import (
	"campus_forum_cloud/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	logging "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	//"gorm.io/plugin/dbresolver"
	"gorm.io/plugin/opentelemetry/tracing"
	//"strings"
	"time"
)

var Client *gorm.DB

func init() {
	var err error
	var cfg gorm.Config

	cfg = gorm.Config{
		PrepareStmt: true,
		Logger:      logging.Default.LogMode(logging.Info),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: config.Configs.MySQLPrefix,
		},
	}

	if Client, err = gorm.Open(
		mysql.Open(config.Configs.MysqlDSN),
		&cfg,
	); err != nil {
		panic(err)
	}
	// TODO: mysql 主从同步
	//if config.EnvCfg.MySQLReplicaState == "enable" {
	//	var replicas []gorm.Dialector
	//	for _, addr := range strings.Split(config.EnvCfg.MySQLReplicaAddress, ",") {
	//		pair := strings.Split(addr, ":")
	//		if len(pair) != 2 {
	//			continue
	//		}
	//
	//		replicas = append(replicas, mysql.Open(
	//			fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s parseTime=true",
	//				config.EnvCfg.MySQLReplicaUser,
	//				config.EnvCfg.MySQLReplicaPassword,
	//				config.EnvCfg.MySQLDatabase,
	//				pair[0],
	//				pair[1])))
	//	}
	//
	//	err := Client.Use(dbresolver.Register(dbresolver.Config{
	//		Replicas: replicas,
	//		Policy:   dbresolver.RandomPolicy{},
	//	}))
	//	if err != nil {
	//		panic(err)
	//	}
	//}

	sqlDB, err := Client.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(200)
	sqlDB.SetConnMaxLifetime(24 * time.Hour)
	sqlDB.SetConnMaxIdleTime(time.Hour)
	if err := Client.Use(tracing.NewPlugin()); err != nil {
		panic(err)
	}
}
