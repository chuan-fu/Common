package mysql

import (
	"fmt"
	"time"

	"github.com/chuan-fu/Common/zlog"
	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

/*
type Config struct {
	Mysql mysql.MysqlConf `json:"mysql" yaml:"mysql"`
}

mysql:
	dataSourceName: root:123456@tcp(0.0.0.0:3306)/user?charset-utf8mb4
	maxOpenConns: 2000
	maxIdleConns: 1000
	connMaxLifetime: -1
	debug: true
	logLevel: info
	ignoreRecordNotFoundError: true

*/
type MysqlConf struct {
	DataSourceName            string `required:"true" json:"dataSourceName" yaml:"dataSourceName"`                       // 主库
	MaxOpenConns              int    `default:"200" json:"maxOpenConns" yaml:"maxOpenConns"`                             // 最大连接数
	MaxIdleConns              int    `default:"50" json:"maxIdleConns" yaml:"maxIdleConns"`                              // 最大空闲连接数
	ConnMaxLifetime           int64  `default:"21600" json:"connMaxLifetime" yaml:"connMaxLifetime"`                     // 连接可以重用的最长时间
	Debug                     bool   `default:"false" json:"debug" yaml:"debug"`                                         // 是否开启debug
	LogLevel                  string `default:"error" json:"logLevel" yaml:"logLevel"`                                   // 日志等级
	IgnoreRecordNotFoundError bool   `default:"false" json:"ignoreRecordNotFoundError" yaml:"ignoreRecordNotFoundError"` // 是否忽略记录未找到bug
}

var gormDb *gorm.DB

func GetGorm() *gorm.DB {
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
