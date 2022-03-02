package logger

import (
	"encoding/json"
	"fmt"
	"gohub/pkg/app"
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 全局logger 对象
var Logger *zap.Logger

// 日志初始化
func InitLogger(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string, level string) {
	// 获取日志写入介质
	writeSyncer := getLogWriter(filename, maxSize, maxBackup, maxAge, compress, logType)

	// 设置日志等级, 具体见 config/log.go 文件
	logLevel := new(zapcore.Level)
	if err := logLevel.UnmarshalText([]byte(level)); err != nil {
		fmt.Println("日志初始化错误, 日志级别设置有误")
	}

	// 初始化 core
	core := zapcore.NewCore(getEncoder(), writeSyncer, logLevel)

	// 初始化 logger
	Logger = zap.New(core,
		zap.AddCaller(),                   // 调用文件和行号 内部使用 runtime.Caller
		zap.AddCallerSkip(1),              // 封装了一层, 调用文件去除一层(runtime.Caller(1))
		zap.AddStacktrace(zap.ErrorLevel), // Error时才会显示 stacktrace
	)

	// 将自定义的 logger 替换为全局的 logger
	// zap.L().Fatal() 调用时, 就会使用我们自定义的 logger
	zap.ReplaceGlobals(Logger)
}

// 设置日志存储格式
func getEncoder() zapcore.Encoder {
	// 日志格式规则
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", // 代码调用 如 paginator/paginator.go: 148
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,      // 每行日志的结尾添加 "\n"
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 日志级别名称大写 如 ERROR
		EncodeTime:     customTimeEncoder,              // 时间格式, 2006-01-02 15:04:05
		EncodeDuration: zapcore.SecondsDurationEncoder, // 执行时间 以秒为单位
		EncodeCaller:   zapcore.ShortCallerEncoder,     // Caller 短格式
	}

	// 本地环境
	if app.IsLocal() {
		// 终端输出关键词高亮
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		// 本地设置内置的 Console 解码器 (支持 stacktrace 换行)
		return zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 线上环境使用 JSON 编码器
	return zapcore.NewJSONEncoder(encoderConfig)
}

// 自定义友好时间格式
func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool, logType string) zapcore.WriteSyncer {
	// 如果配置了按照日期记录日志文件
	if logType == "daily" {
		logname := time.Now().Format("2006-01-02.log")
		filename = strings.ReplaceAll(filename, "logs.log", logname)
	}

	// 滚动日志
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
		Compress:   compress,
	}

	// 配置输出介质
	if app.IsLocal() {
		// 本地开发终端打印和记录文件
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(lumberJackLogger))
	} else {
		// 生产环境只记录文件
		return zapcore.AddSync(lumberJackLogger)
	}

}

// 调试专用, 不会中断程序, 会在中断打印出 warning 消息
func Dump(value interface{}, msg ...string) {
	valueString := jsonString(value)
	// 判断第二个参数是否传参 msg
	if len(msg) > 0 {
		Logger.Warn("Dump", zap.String(msg[0], valueString))
	} else {
		Logger.Warn("Dump", zap.String("data", valueString))
	}
}

func LogIf(err error) {
	if err != nil {
		Logger.Error("Error Occurred:", zap.Error(err))
	}
}

func LogWarnIf(err error) {
	if err != nil {
		Logger.Warn("Error Occurred:", zap.Error(err))
	}
}

func LogInfoIf(err error) {
	if err != nil {
		Logger.Info("Error occurred:", zap.Error(err))
	}
}

// 调试日志
func Debug(mouduleName string, fields ...zap.Field) {
	Logger.Debug(mouduleName, fields...)
}

func Info(moduleName string, fields ...zap.Field) {
	Logger.Info(moduleName, fields...)
}

func Warn(moduleName string, fields ...zap.Field) {
	Logger.Warn(moduleName, fields...)
}

// Error 错误时记录，不应该中断程序，查看日志时重点关注
func Error(moduleName string, fields ...zap.Field) {
	Logger.Error(moduleName, fields...)
}

// Fatal 级别同 Error(), 写完 log 后调用 os.Exit(1) 退出程序
func Fatal(moduleName string, fields ...zap.Field) {
	Logger.Fatal(moduleName, fields...)
}

// 记录一条字符串类型的 debug 日志
func DebugString(moduleName, name, msg string) {
	Logger.Debug(moduleName, zap.String(name, msg))
}

func InfoString(moduleName, name, msg string) {
	Logger.Info(moduleName, zap.String(name, msg))
}

func WarnString(moduleName, name, msg string) {
	Logger.Warn(moduleName, zap.String(name, msg))
}

func ErrorString(moduleName, name, msg string) {
	Logger.Error(moduleName, zap.String(name, msg))
}

func FatalString(moduleName, name, msg string) {
	Logger.Fatal(moduleName, zap.String(name, msg))
}

func DebugJSON(moduleName, name string, value interface{}) {
	Logger.Debug(moduleName, zap.String(name, jsonString(value)))
}

func InfoJSON(moduleName, name string, value interface{}) {
	Logger.Info(moduleName, zap.String(name, jsonString(value)))
}

func WarnJSON(moduleName, name string, value interface{}) {
	Logger.Warn(moduleName, zap.String(name, jsonString(value)))
}

func ErrorJSON(moduleName, name string, value interface{}) {
	Logger.Error(moduleName, zap.String(name, jsonString(value)))
}

func FatalJSON(moduleName, name string, value interface{}) {
	Logger.Fatal(moduleName, zap.String(name, jsonString(value)))
}

func jsonString(value interface{}) string {
	b, err := json.Marshal(value)
	if err != nil {
		Logger.Error("logger", zap.String("JSON marshal error", err.Error()))
	}

	return string(b)
}
