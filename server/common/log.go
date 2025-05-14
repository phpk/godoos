package common

import (
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// InitLogger 初始化日志系统
func InitLogger() (slog.Handler, error) {
	// 根据 Debug 字段判断是否开启文件日志记录
	handler, err := NewLogHandler(&Config.Log)
	if err != nil {
		return nil, err
	}

	return handler, nil
}

// NewLogHandler 创建一个新的日志处理器
func NewLogHandler(log *Log) (slog.Handler, error) {
	if Config.Log.WriteFile {
		logFilePath := filepath.Join("data", log.Path, log.Filename)
		lumberjackLogger := &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    log.MaxSize,
			MaxBackups: log.MaxBackups,
			MaxAge:     log.MaxAge,
		}

		// 创建一个 slog.TextHandler，它将使用 lumberjackLogger 作为输出
		textHandler := slog.NewTextHandler(lumberjackLogger, nil)

		return textHandler, nil
	} else {
		// 如果不是 debug 模式，使用控制台输出
		return slog.NewJSONHandler(os.Stdout, nil), nil
	}
}
