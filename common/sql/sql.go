package sql

import (
	"fmt"
	gmysql "github.com/go-sql-driver/mysql"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"

	"gorm.io/gorm"
	gormloger "gorm.io/gorm/logger"
)

var useTls = false

func Dail(dsn string) *gorm.DB {
	if useTls {
		myconfig, err := gmysql.ParseDSN(dsn)
		if err != nil {
			panic(fmt.Sprintf("ParseDSN(%s) failed : %s", dsn, err.Error()))
		}

		myconfig.TLSConfig = "custom"
		dsn = myconfig.FormatDSN()
		logx.Infof("%s", dsn)
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: NewLogger(gormloger.Config{
			LogLevel: gormloger.Info,
		}),
	})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	return db
}
