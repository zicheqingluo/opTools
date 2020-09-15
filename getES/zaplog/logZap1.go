package zaplog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var LogSetting = lumberjack.Logger{
	Filename:   "./logs/getES.log", // 日志文件路径
	MaxSize:    128,                // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups: 30,                 // 日志文件最多保存多少个备份
	MaxAge:     7,                  // 文件最多保存多少天
	Compress:   true,               // 是否压缩
}

var EncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "time",
	LevelKey:       "level",
	NameKey:        "logger",
	CallerKey:      "linenum",
	MessageKey:     "msg",
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
	EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
	EncodeDuration: zapcore.SecondsDurationEncoder, //
	EncodeCaller:   zapcore.FullCallerEncoder,      // 全路径编码器
	EncodeName:     zapcore.FullNameEncoder,
}

// 设置日志级别
var AtomicLevel = zap.NewAtomicLevel()

//AtomicLevel.SetLevel(zap.InfoLevel)

var core = zapcore.NewCore(
	zapcore.NewJSONEncoder(EncoderConfig),                                                 // 编码器配置
	zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&LogSetting)), // 打印到控制台和文件
	AtomicLevel, // 日志级别
)

// 开启开发模式，堆栈跟踪
var caller = zap.AddCaller()

// 开启文件及行号
var development = zap.Development()

// 设置初始化字段
var Filed = zap.Fields(zap.String("serviceName", "getES"))

// 构造日志
var Logger = zap.New(core, caller, development, Filed)

/*
   logger.Info("log 初始化成功")
   logger.Info("无法获取网址",
       zap.String("url", "http://www.baidu.com"),
       zap.Int("attempt", 3),
       zap.Duration("backoff", time.Second))
*/
