package mysql

import (
	"fmt"
	"time"

	"github.com/chuan-fu/Common/log"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type MysqlConf struct {
	DataSourceName            string `required:"true"` // 主库
	MaxOpenConns              int    `default:"200"`   // 最大连接数
	MaxIdleConns              int    `default:"50"`    // 最大空闲连接数
	ConnMaxLifetime           int64  `default:"21600"` // 连接可以重用的最长时间
	Debug                     bool   `default:"false"` // 是否开启debug
	LogLevel                  string `default:"error"` // 日志等级
	IgnoreRecordNotFoundError bool   `default:"false"` // 是否忽略记录未找到bug
}

var gormDb *gorm.DB

func GetGormDB() *gorm.DB {
	return gormDb
}

func ConnectGORM(conf MysqlConf) error {
	if err := connectGORM(conf); err != nil {
		err = errors.Wrap(err, "GORM连接错误")
		log.Error(err)
		return err
	}
	return nil
}

func ReloadGORM(conf MysqlConf) error {
	oldDb := gormDb
	if err := connectGORM(conf); err != nil {
		err = errors.Wrap(err, "GORM重连错误")
		log.Error(err)
		return err
	}
	CloseGORM(oldDb)
	return nil
}

func connectGORM(conf MysqlConf) error {
	baseDB, err := gorm.Open(mysql.Open(conf.DataSourceName), &gorm.Config{
		Logger: NewGormLogger(NewLogWriter(), Config{
			IgnoreRecordNotFoundError: conf.IgnoreRecordNotFoundError,
			LogLevel:                  GetLoggerLevel(conf.LogLevel),
		}),
		AllowGlobalUpdate: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	})
	if err != nil {
		log.Error(err)
		return err
	}

	sqlDb, err := baseDB.DB()
	if err != nil {
		log.Error(err)
		return err
	}
	sqlDb.SetMaxOpenConns(conf.MaxOpenConns)
	sqlDb.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDb.SetConnMaxLifetime(time.Duration(conf.ConnMaxLifetime) * time.Second)

	err = sqlDb.Ping()
	if err != nil {
		log.Error(err)
		return err
	}

	gormDb = baseDB
	if conf.Debug {
		gormDb = gormDb.Debug()
	}
	return nil
}

func CloseGORM(db *gorm.DB) {
	if db != nil {
		sqlDb, err := db.DB()
		if err != nil {
			err = fmt.Errorf("数据库链接错误 %w ", err)
			log.Error(err)
			return
		}
		err = sqlDb.Close()
		if err != nil {
			log.Error(err)
		} else {
			log.Info(" gorm closed")
		}
	}
}

func GetLoggerLevel(level string) logger.LogLevel {
	var loggerLevel logger.LogLevel
	switch level {
	case "error":
		loggerLevel = logger.Error
	case "warn":
		loggerLevel = logger.Warn
	case "info":
		loggerLevel = logger.Info
	default:
		loggerLevel = logger.Info
	}

	return loggerLevel
}
