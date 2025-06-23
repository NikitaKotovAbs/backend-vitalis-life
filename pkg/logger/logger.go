package logger

import (
	"os"
	"path/filepath"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var log *zap.Logger

func Init(level string) {
	// Создаем директорию для логов
	logDir := filepath.Join("..", "pkg", "logger", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		panic(err)
	}

	// Базовая конфигурация энкодера
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller", 
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, 
	}

	// Настраиваем цветной вывод для debug
	if level == "debug" {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	// Устанавливаем уровень логирования
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Создаем список ядер
	var cores []zapcore.Core

	// Файловый вывод только для info
	if level == "info" {
		infoCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   filepath.Join(logDir, "info.log"),
				MaxSize:    30,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   false,
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.InfoLevel && lvl < zapcore.ErrorLevel
			}),
		)

		errorCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(&lumberjack.Logger{
				Filename:   filepath.Join(logDir, "error.log"),
				MaxSize:    20,
				MaxBackups: 5,
				MaxAge:     30,
				Compress:   false,
			}),
			zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				return lvl >= zapcore.ErrorLevel
			}),
		)

		cores = append(cores, infoCore, errorCore)
	}

	// Консольный вывод (всегда включен)
	consoleCore := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapLevel
		}),
	)
	cores = append(cores, consoleCore)

	// Создаем логгер с дополнительными опциями
	log = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),      // Добавляем информацию о вызове для всех уровней
		zap.AddCallerSkip(1), // Пропускаем 1 уровень стека (сами функции логгера)
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
}

// Оберточные функции с добавлением caller информации
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}