package zlog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

const (
	defaultMaxSize   = 800
	defaulMaxBackups = 5
	defaultMaxAge    = 28
)

var (
	defaultZapLogger *ZapLogger
	logOnce          sync.Once
)

type ZConfig struct {
	FileName    string
	Level       string
	MaxSize     int
	MaxAge      int
	BackUpCount int
	Console     bool
	Compress    bool
}

type ZapLogger struct {
	logger   *zap.Logger
	outer    zapcore.WriteSyncer
	inlogger *lumberjack.Logger
}

func init() {
	resetLogger(ZConfig{Console: true})
}

func ResetOnce(c ZConfig) {
	logOnce.Do(func() {
		_ = defaultZapLogger.inlogger.Close()
		resetLogger(c)
	})
}

func setDefaultConf(c *ZConfig) {
	if len(c.Level) < 1 {
		c.Level = "info"
	}
	if c.MaxSize < 1 {
		c.MaxSize = defaultMaxSize
	}
	if c.BackUpCount < 1 {
		c.BackUpCount = defaulMaxBackups
	}
	if c.MaxAge < 1 {
		c.MaxAge = defaultMaxAge
	}
	now := time.Now()
	if len(c.FileName) < 1 {
		c.FileName = fmt.Sprintf("log/server-%04d%02d%02d-%02d-%02d-%02d.log",
			now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	}
}

func resetLogger(c ZConfig) {
	setDefaultConf(&c)
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeTime = zapcore.TimeEncoderOfLayout(time.DateTime + ".000")
	encoder := zapcore.NewJSONEncoder(encodeConf)

	fileW := &lumberjack.Logger{
		Filename:   c.FileName,
		MaxSize:    c.MaxSize,
		MaxBackups: c.BackUpCount,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}

	outer := zapcore.AddSync(fileW)
	ws := []zapcore.WriteSyncer{outer}
	if c.Console {
		ws = append(ws, zapcore.AddSync(os.Stdout))
	}
	ll := getLevel(c.Level)
	core := zapcore.NewCore(
		encoder, zapcore.NewMultiWriteSyncer(ws...), ll)

	opts := []zap.Option{
		zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(zap.NewAtomicLevelAt(zap.ErrorLevel)),
	}

	zapLogger := zap.New(core, opts...)
	if defaultZapLogger == nil {
		defaultZapLogger = &ZapLogger{}
	}
	defaultZapLogger.logger = zapLogger
	defaultZapLogger.outer = outer
	defaultZapLogger.inlogger = fileW
}

func SetLogger(logger *zap.Logger) {
	defaultZapLogger.logger = logger
}
func GetLogger() *ZapLogger {
	return defaultZapLogger
}
func GetOuter() io.Writer {
	return defaultZapLogger.outer
}

func getLevel(level string) zap.AtomicLevel {
	ll := strings.ToLower(level)
	zl := zap.NewAtomicLevel()
	switch ll {
	case "debug":
		zl = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		zl = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		zl = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		zl = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		zl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	return zl
}

func Sync() {
	defaultZapLogger.logger.Sync()
}

func Info(msg string, fields ...zap.Field) {
	defaultZapLogger.logger.Info(msg, fields...)
}
func Warn(msg string, fields ...zap.Field) {
	defaultZapLogger.logger.Warn(msg, fields...)
}
func Error(msg string, fields ...zap.Field) {
	defaultZapLogger.logger.Error(msg, fields...)
}
func Debug(msg string, fields ...zap.Field) {
	defaultZapLogger.logger.Debug(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	defaultZapLogger.logger.Fatal(msg, fields...)
}
