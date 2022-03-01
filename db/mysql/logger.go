package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/chuan-fu/Common/log"
	"github.com/pkg/errors"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type Writer interface {
	Printf(logger.LogLevel, ...interface{})
}

type logWriter struct{}

func NewLogWriter() Writer {
	return &logWriter{}
}

func (w logWriter) Printf(level logger.LogLevel, data ...interface{}) {
	switch level {
	case logger.Info:
		log.Info(data)
	case logger.Warn:
		log.Warn(data)
	case logger.Error:
		log.Error(data)
	default:
		log.Info(data)
	}
}

type gormLogger struct {
	Writer
	Config
}

type Config struct {
	SlowThreshold             time.Duration // 慢sql
	IgnoreRecordNotFoundError bool          // 是否忽略记录未找到bug
	LogLevel                  logger.LogLevel
}

const (
	DefaultSlowThreshold = 200 * time.Millisecond // 默认慢sql
)

// New initialize logger
func NewGormLogger(writer Writer, config Config) logger.Interface {
	if config.SlowThreshold == 0 {
		config.SlowThreshold = DefaultSlowThreshold
	}
	return &gormLogger{
		Writer: writer,
		Config: config,
	}
}

// LogMode log mode
func (l *gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l gormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.Printf(logger.Info, fmt.Sprintf(msg, data...))
	}
}

// Warn print warn messages
func (l gormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.Printf(logger.Warn, fmt.Sprintf(msg, data...))
	}
}

// Error print error messages
func (l gormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.Printf(logger.Error, fmt.Sprintf(msg, data...))
	}
}

// Trace print sql message
func (l gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, logger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(logger.Error, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(logger.Error, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(logger.Warn, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(logger.Warn, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(logger.Info, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(logger.Info, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}
