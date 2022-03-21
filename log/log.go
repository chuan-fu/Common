package log

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/chuan-fu/Common/util"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

type LogConfig struct {
	Lumberjack *lumberjack.Logger `yaml:"lumberjack"`
	ZapConfig  zap.Config         `yaml:"zapConfig"`
}

func (lc *LogConfig) Build() (logger *zap.Logger) {
	return lc.buildRollLog()
}

func (lc *LogConfig) buildRollLog() (logger *zap.Logger) {
	hook := lc.Lumberjack
	if hook == nil {
		hook = &lumberjack.Logger{}
		hook.Filename = "./log/sys.log"
		hook.MaxSize = 128
		hook.MaxBackups = 0
		hook.MaxAge = 0
		hook.LocalTime = true
		hook.Compress = false
	}

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(lc.ZapConfig.Level.Level())

	encoderConfig := lc.ZapConfig.EncoderConfig
	var encoder zapcore.Encoder
	switch lc.ZapConfig.Encoding {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default: // case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	core := zapcore.NewCore(
		encoder, // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)

	opts := []zap.Option{zap.AddStacktrace(zap.DPanicLevel), zap.AddCallerSkip(1)}
	if lc.ZapConfig.Development {
		opts = append(opts, zap.Development())
	}
	if !lc.ZapConfig.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}
	stackLevel := zap.ErrorLevel
	if lc.ZapConfig.Development {
		stackLevel = zap.WarnLevel
	}
	if !lc.ZapConfig.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(stackLevel))
	}
	if len(lc.ZapConfig.InitialFields) != 0 {
		var fields []zap.Field
		for k, v := range lc.ZapConfig.InitialFields {
			fields = append(fields, zap.Any(k, v))
		}
		opts = append(opts, zap.Fields(fields...))
	}
	if len(lc.ZapConfig.ErrorOutputPaths) != 0 {
		errSink, _, err := zap.Open(lc.ZapConfig.ErrorOutputPaths...)
		if err == nil {
			opts = append(opts, zap.ErrorOutput(errSink))
		}
	}

	return zap.New(core, opts...)
}

// 命令行应用log
func (lc *LogConfig) buildLog4Cmd() (logger *zap.Logger) {
	opts := []zap.Option{zap.AddStacktrace(zap.DPanicLevel), zap.AddCallerSkip(1)}
	var err error
	logger, err = lc.ZapConfig.Build(opts...)
	if err != nil {
		panic(err)
	}
	return logger
}

var logger *zap.Logger

func init() {
	if CheckExist(defaultConfFilePath) {
		if bs, err := ioutil.ReadFile(defaultConfFilePath); err == nil {
			logger = initZapLog(bs)
			return
		}
	}
	logger = initZapLog(util.StringToBytes(defaultConfFile))
}

func initZapLog(bs []byte) (logger *zap.Logger) {
	lc := LogConfig{}
	err := yaml.Unmarshal(bs, &lc)
	if err != nil {
		panic(err)
	}
	return lc.Build()
}

func ReplaceLoggerFromString(cfg string) {
	logger = initZapLog(util.StringToBytes(cfg))
}

func GetLogger() *zap.Logger {
	return logger
}

func Debug(v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Debug(fmt.Sprint(ps...), fields...)
}

func Debugf(format string, v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Debug(fmt.Sprintf(format, ps...), fields...)
}

func Info(v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Info(fmt.Sprint(ps...), fields...)
}

func Infof(format string, v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Info(fmt.Sprintf(format, ps...), fields...)
}

func Warn(v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Warn(fmt.Sprint(ps...), fields...)
}

func Warnf(format string, v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Warn(fmt.Sprintf(format, ps...), fields...)
}

func Error(v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Error(fmt.Sprint(ps...), fields...)
}

func Errorf(format string, v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Error(fmt.Sprintf(format, ps...), fields...)
}

func Fatal(v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Fatal(fmt.Sprint(ps...), fields...)
}

func Fatalf(format string, v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Fatal(fmt.Sprintf(format, ps...), fields...)
}

func Panic(v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Panic(fmt.Sprint(ps...), fields...)
}

func Panicf(format string, v ...interface{}) {
	fields, ps := splitLogValue(v)
	logger.Panic(fmt.Sprintf(format, ps...), fields...)
}

func splitLogValue(v []interface{}) ([]zap.Field, []interface{}) {
	var fields []zap.Field
	var ps []interface{}
	for _, p := range v {
		field, ok := p.(zap.Field)
		if ok {
			fields = append(fields, field)
		} else {
			ps = append(ps, p)
		}
	}
	return fields, ps
}

// 使用前必须 defer log.Flush()
func Flush() {
	_ = logger.Sync()
}

func CheckExist(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
