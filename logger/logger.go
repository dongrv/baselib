package logger

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

type FlowLevel uint8 // 数据流向等级类型

const (
	Console          FlowLevel = 1 // 控制台输出
	ConsoleFile      FlowLevel = 2 // 控制台输出 & 磁盘文件
	ConsoleFileKafka FlowLevel = 4 // 控制台输出 & 磁盘文件 & Kafka

	// 输出格式
	PrintConsoleMode uint = 1
	PrintJsonMode    uint = 2
)

// logger 配置结构
type Setting struct {
	Debug      uint8     `json:"debug"`
	StoreLevel FlowLevel `json:"store_level"`
	Filename   string    `json:"filename,omitempty"`
	PrintMode  uint      `json:"print_mode"`
}

var (
	logger *zap.Logger
	Sugar  *zap.SugaredLogger
	writer []zapcore.WriteSyncer

	Config *Setting
)

// 初始化logger
func InitLogger(setting *Setting) {
	var (
		level   zapcore.Level   // 日志级别
		core    zapcore.Core    // 更小更快的核心接口
		encoder zapcore.Encoder // 编码器，用于设置编码格式
	)

	Config = setting // 设置为全局变量

	// debug 模式
	if setting.Debug == 1 {
		level = zapcore.DebugLevel // 调试
	} else {
		level = zapcore.ErrorLevel //| zapcore.DPanicLevel | zapcore.PanicLevel | zapcore.FatalLevel // 错误
	}
	encoder = getEncoder(setting.PrintMode)
	switch setting.StoreLevel {
	case Console:
		core = getConsoleCore(&encoder, level)
		break
	case ConsoleFile, ConsoleFileKafka:
		core = getMixedCore(&encoder, level, setting)
		break
	default:
		fmt.Println("[Init logger] Invalid setting parameter `StoreLevel`.")
		return
	}
	logger = zap.New(core, zap.AddCaller())
	Sugar = logger.Sugar()
	return
}

// 编码器
func getEncoder(printMode uint) (zapEncoder zapcore.Encoder) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	switch printMode {
	case PrintJsonMode:
		zapEncoder = zapcore.NewJSONEncoder(config)
		break
	case PrintConsoleMode:
		zapEncoder = zapcore.NewConsoleEncoder(config)
	}
	return
}

// 控制台输出
func getConsoleCore(encoder *zapcore.Encoder, level zapcore.Level) (core zapcore.Core) {
	core = zapcore.NewCore(*encoder, zapcore.AddSync(os.Stdout), level)
	return
}

// 控制台输出 & 存文件
func getMixedCore(encoder *zapcore.Encoder, level zapcore.Level, setting *Setting) (core zapcore.Core) {
	lumberjack := zapcore.AddSync(&lumberjack.Logger{
		Filename:   setting.Filename, // 文件路径和文件名
		MaxSize:    10,               // 切割文件的标准最大大小：xx M
		MaxBackups: 10,               // 保留历史文件的最大个数
		MaxAge:     7,                // 保留历史文件的最大天数
		Compress:   false,            // 是否压缩/归档历史文件
	})
	writer = append(writer, lumberjack, zapcore.AddSync(os.Stdout))
	core = zapcore.NewCore(*encoder, zapcore.NewMultiWriteSyncer(writer...), level)
	return
}
