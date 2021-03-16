/**
 * @Author: Bugzheng
 * @Description:
 * @File:  logger
 * @Version: 1.0.0
 * @Date: 2021/03/15 3:53 下午
 */
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"sync"
	"time"
)

var (
	logger *zap.Logger
	logMap sync.Map
)

type LogConfig struct {
	Path       string `toml:"path"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

// newLogger returns a new logger.
func NewLogger(cfg *LogConfig, debug bool) *zap.Logger {
	//hook := lumberjack.Logger{
	//	Filename:   cfg.Path,
	//	MaxSize:    cfg.MaxSize,
	//	MaxBackups: cfg.MaxBackups,
	//	MaxAge:     cfg.MaxAge,
	//	Compress:   cfg.Compress,
	//}
	//var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	if debug {
		cfg := zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeTime = MyTimeEncoder
		//writes = append(writes, zapcore.AddSync(os.Stdout))
		l, _ := cfg.Build()

		return l
	}
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.Path,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	})
	c := zap.NewProductionEncoderConfig()

	c.TimeKey = "time"
	c.EncodeTime = MyTimeEncoder
	c.EncodeCaller = zapcore.FullCallerEncoder

	core := zapcore.NewCore(zapcore.NewJSONEncoder(c), w, zap.DebugLevel)

	return zap.New(core, zap.AddCaller())
}

//
//func InitLogger2() {
//	// 此处的配置是从我的项目配置文件读取的，读者可以根据自己的情况来设置
//	logPath := config.Cfg.Section("app").Key("logPath").String()
//	name := config.Cfg.Section("app").Key("name").String()
//	debug, err := config.Cfg.Section("app").Key("debug").Bool()
//	if err != nil {
//		debug = false
//	}
//	hook := lumberjack.Logger{
//		Filename:   logPath, // 日志文件路径
//		MaxSize:    128,     // 每个日志文件保存的大小 单位:M
//		MaxAge:     7,       // 文件最多保存多少天
//		MaxBackups: 30,      // 日志文件最多保存多少个备份
//		Compress:   false,   // 是否压缩
//	}
//	encoderConfig := zapcore.EncoderConfig{
//		MessageKey:     "msg",
//		LevelKey:       "level",
//		TimeKey:        "time",
//		NameKey:        "logger",
//		CallerKey:      "file",
//		StacktraceKey:  "stacktrace",
//		LineEnding:     zapcore.DefaultLineEnding,
//		EncodeLevel:    zapcore.LowercaseLevelEncoder,
//		EncodeTime:     zapcore.ISO8601TimeEncoder,
//		EncodeDuration: zapcore.SecondsDurationEncoder,
//		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
//		EncodeName:     zapcore.FullNameEncoder,
//	}
//	// 设置日志级别
//	atomicLevel := zap.NewAtomicLevel()
//	atomicLevel.SetLevel(zap.DebugLevel)
//	var writes = []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
//	// 如果是开发环境，同时在控制台上也输出
//	if debug {
//		writes = append(writes, zapcore.AddSync(os.Stdout))
//	}
//	core := zapcore.NewCore(
//		zapcore.NewJSONEncoder(encoderConfig),
//		zapcore.NewMultiWriteSyncer(writes...),
//		atomicLevel,
//	)
//
//	// 开启开发模式，堆栈跟踪
//	caller := zap.AddCaller()
//	// 开启文件及行号
//	development := zap.Development()
//
//	// 设置初始化字段
//	field := zap.Fields(zap.String("appName", name))
//
//	// 构造日志
//	ZapLogger = zap.New(core, caller, development, field)
//	ZapLogger.Info("log 初始化成功")
//}

// Logger returns a logger
func Logger(name ...string) *zap.Logger {
	if len(name) == 0 {
		return logger
	}

	v, ok := logMap.Load(name[0])

	if !ok {
		return logger
	}
	return v.(*zap.Logger)
}

// MyTimeEncoder zap time encoder.
func MyTimeEncoder(t time.Time, e zapcore.PrimitiveArrayEncoder) {
	e.AppendString(t.Format("2006-01-02 15:04:05"))
}
