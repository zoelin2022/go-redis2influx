package utils

import (
	"go-redis2influx/global"

	"os"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func InitLogger() {
	var logger *zap.Logger
	cfg := zap.NewDevelopmentConfig()

	c := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大寫帶顏色
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	cfg.EncoderConfig = c

	level := global.EnvConfig.Log.Level

	core := zapcore.NewTee(
		// 1. console & db
		zapcore.NewCore(CustomLogConsole(), zapcore.AddSync(os.Stdout), DefaultLogLevel(level)),

		// 2. file
		zapcore.NewCore(
			CustomLogFile(),
			DefaultRotateWriteSyncer(),
			DefaultLogLevel(level)),
	)

	// caller 顯示文件名、行號和zap調用者的函數名
	if global.EnvConfig.Log.Level == "debug" {
		logger = zap.New(core, zap.AddCaller())
	} else {
		logger = zap.New(core)
	}

	global.Logger = logger

}

// *** 螢幕輸出 ***//
func CustomLogConsole() zapcore.Encoder {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大寫帶顏色
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	})

	return encoder
}

// *** 寫到檔案 ***//
func CustomLogFile() zapcore.Encoder {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// 在日誌文件中使用大寫字母記錄日誌級別
		EncodeLevel: zapcore.CapitalLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	})
	return encoder
}

// error > warn > info > debug
func DefaultLogLevel(level string) zap.LevelEnablerFunc {
	var priority zap.LevelEnablerFunc

	switch level {
	case "error": // error
		highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= zap.ErrorLevel
		})
		priority = highPriority

	case "info": // warn + info
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev <= zap.ErrorLevel && lev > zap.DebugLevel
		})
		priority = lowPriority
	case "debug": // warn + info + debug
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev <= zap.ErrorLevel && lev >= zap.DebugLevel
		})
		priority = lowPriority
	}

	return priority
}

// *** 日誌切割 ***//
// lumberjack：如果 MaxBackups 和 MaxAge均為 0，則不會刪除任何舊的日誌檔。
func DefaultRotateWriteSyncer() zapcore.WriteSyncer {
	cfg := global.EnvConfig.Log
	var fileWriteSyncer zapcore.WriteSyncer
	logPath := global.EnvConfig.Log.Path
	fileWriteSyncer = zapcore.AddSync(&lumberjack.Logger{
		Filename: logPath,     // 日誌文件存放目錄，如果文件夾不存在會自動創建
		MaxSize:  cfg.MaxSize, // 文件大小限制,單位MB
		MaxAge:   cfg.MaxAge,  // 日誌文件保留天數
		//Compress:   false,   // 是否壓縮處理
	})

	return fileWriteSyncer
}

// *** Custom ***//
func CustomLogLevel(level string) zap.LevelEnablerFunc {
	var priority zap.LevelEnablerFunc

	switch level {
	case "error": // error
		highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev >= zap.ErrorLevel
		})
		priority = highPriority

	case "info": // warn + info
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev < zap.ErrorLevel && lev > zap.DebugLevel
		})
		priority = lowPriority
	case "debug": // warn + info + debug
		lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
			return lev <= zap.ErrorLevel && lev >= zap.DebugLevel
		})
		priority = lowPriority
	}

	return priority
}

// *** Custom ***//
// lumberjack：如果 MaxBackups 和 MaxAge均為 0，則不會刪除任何舊的日誌檔。
func CustomRotateWriteSyncer(level string) zapcore.WriteSyncer {
	cfg := global.EnvConfig.Log
	var fileWriteSyncer zapcore.WriteSyncer
	switch level {
	case "error":
		filePath := filepath.Join(cfg.Path, "error.json")
		fileWriteSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename: filePath,    // 日誌文件存放目錄，如果文件夾不存在會自動創建
			MaxSize:  cfg.MaxSize, // 文件大小限制,單位MB
			MaxAge:   cfg.MaxAge,  // 日誌文件保留天數
			//Compress:   false,   // 是否壓縮處理
		})
	case "info", "debug":
		filePath := filepath.Join(cfg.Path, "info.json")
		fileWriteSyncer = zapcore.AddSync(&lumberjack.Logger{
			Filename: filePath,    // 日誌文件存放目錄，如果文件夾不存在會自動創建
			MaxSize:  cfg.MaxSize, // 文件大小限制,單位MB
			MaxAge:   cfg.MaxAge,  // 日誌文件保留天數
			//Compress:   false,   // 是否壓縮處理
		})
	}
	return fileWriteSyncer
}

// *** Custom ***//
func CustomInitLogger() {
	var logger *zap.Logger
	cfg := zap.NewDevelopmentConfig()

	c := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "message",
		StacktraceKey:  "stack",
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 大寫帶顏色
		EncodeCaller:   zapcore.ShortCallerEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
	}

	cfg.EncoderConfig = c

	level := global.EnvConfig.Log.Level

	switch level {
	case "error": // error
		core := zapcore.NewTee(
			//************ error ************//
			// 1. console
			zapcore.NewCore(
				CustomLogConsole(),
				zapcore.Lock(zapcore.AddSync(os.Stdout)),
				CustomLogLevel("error")),
			// 2. file
			zapcore.NewCore(
				CustomLogFile(),
				CustomRotateWriteSyncer("error"),
				CustomLogLevel("error")),
		)
		logger = zap.New(core)
	case "info": // error + warn + info
		core := zapcore.NewTee(
			//************ error ************//
			// 1. console
			// zapcore.NewCore(
			// 	CustomLogConsole(),
			// 	zapcore.Lock(zapcore.AddSync(os.Stdout)),
			// 	CustomLogLevel("error")),
			// 2. file
			zapcore.NewCore(
				CustomLogFile(),
				CustomRotateWriteSyncer("error"),
				CustomLogLevel("error")),

			//************ info ************//
			// 1. console
			// zapcore.NewCore(
			// 	CustomLogConsole(),
			// 	zapcore.Lock(zapcore.AddSync(os.Stdout)),
			// 	CustomLogLevel(level)),
			// 2, file
			zapcore.NewCore(
				CustomLogFile(),
				CustomRotateWriteSyncer(level),
				CustomLogLevel(level)),
		)

		// caller 顯示文件名、行號和zap調用者的函數名
		logger = zap.New(core, zap.AddCaller())

	case "debug": // error + warn + info + debug
		core := zapcore.NewTee(
			//************ error ************//
			// 1. console
			// zapcore.NewCore(
			// 	CustomLogConsole(),
			// 	zapcore.Lock(zapcore.AddSync(os.Stdout)),
			// 	CustomLogLevel("error")),
			// 2. file
			zapcore.NewCore(
				CustomLogFile(),
				CustomRotateWriteSyncer("error"),
				CustomLogLevel("error")),
			// 3. db
			//&HookCore{Core: hook.Core()},

			//************ info ************//
			// 1. console
			// zapcore.NewCore(
			// 	CustomLogConsole(),
			// 	zapcore.Lock(zapcore.AddSync(os.Stdout)),
			// 	CustomLogLevel(level)),
			// 2, file
			zapcore.NewCore(
				CustomLogFile(),
				CustomRotateWriteSyncer(level),
				CustomLogLevel(level)),
		)

		// caller 顯示文件名、行號和zap調用者的函數名
		logger = zap.New(core, zap.AddCaller())

	}
	global.Logger = logger
}
