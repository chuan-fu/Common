package log

const (
	defaultConfFilePath = "./conf/log.yaml"
	defaultConfFile     = `
lumberjack:
  filename: "./log/sys.log"
  maxsize: 1
  maxage: 0
  maxbackups: 0
  localtime: true
  compress: false
# zap config
zapConfig:
  level: debug
  development: true
  disableCaller: false
  disableStacktrace: false
  sampling:
  encoding: consul
  outputPaths:
    - stdout
    - "./log/sys.log"
  errorOutputPaths:
    - stderr
  initialFields:
    sysname: Common
  encoderConfig:
    messageKey: msg
    levelKey: level
    timeKey: logtime
    nameKey: logger
    callerKey: gofile
    stacktraceKey: stack
    lineEnding: "\n"
    levelEncoder: capital
    timeEncoder: ISO8601
    durationEncoder:
    callerEncoder: full
    nameEncoder: full
`
)
