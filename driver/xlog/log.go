package xlog

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"context"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"math/rand"
)

// 常量定义
const (
	LogCtxID    = "logger_log_id" // context中log_id的字段名
	RandomRange = 10000           // 随机数范围
)

// 全局日志实例
var log = logrus.New()

// 初始化随机数生成器
func init() {
	rand.Seed(time.Now().UnixNano())
}

// Init 初始化日志配置
func Init(c *Config) {
	// 设置日志格式为JSON，包含log_id字段
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05.000",
	})

	// 设置日志级别
	log.SetLevel(logrus.DebugLevel)

	// 配置日志分割和级别分离
	writerMap := setupLogRotation(c)
	if writerMap != nil {
		log.AddHook(lfshook.NewHook(writerMap, log.Formatter))
	}

	// 配置标准输出
	if c.Stdout == STDOUT_YES {
		log.AddHook(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.DebugLevel: os.Stdout,
				logrus.InfoLevel:  os.Stdout,
				logrus.WarnLevel:  os.Stdout,
				logrus.ErrorLevel: os.Stdout,
				logrus.FatalLevel: os.Stdout,
				logrus.PanicLevel: os.Stdout,
			},
			log.Formatter,
		))
	}
}

// setupLogRotation 配置日志轮转
func setupLogRotation(c *Config) lfshook.WriterMap {
	if c.Dir == "" {
		return nil
	}

	// 创建日志目录
	if err := os.MkdirAll(c.Dir, 0755); err != nil {
		log.Errorf("创建日志目录失败: %v", err)
		return nil
	}

	// 配置保留时间
	maxAge := 7 * 24 * time.Hour
	if c.MaxAge > 0 {
		maxAge = time.Duration(c.MaxAge) * 24 * time.Hour
	}

	// 配置轮转时间
	rotationTime := 1 * time.Hour
	if c.RotationHour > 0 {
		rotationTime = time.Duration(c.RotationHour) * time.Hour
	}

	// 定义各级别日志路径
	logLevels := map[logrus.Level]string{
		logrus.DebugLevel: filepath.Join(c.Dir, c.Name, "debug", "debug.log"),
		logrus.InfoLevel:  filepath.Join(c.Dir, c.Name, "info", "info.log"),
		logrus.WarnLevel:  filepath.Join(c.Dir, c.Name, "warn", "warn.log"),
		logrus.ErrorLevel: filepath.Join(c.Dir, c.Name, "error", "error.log"),
		logrus.FatalLevel: filepath.Join(c.Dir, c.Name, "fatal", "fatal.log"),
		logrus.PanicLevel: filepath.Join(c.Dir, c.Name, "panic", "panic.log"),
	}

	writerMap := lfshook.WriterMap{}

	// 为每个级别配置轮转
	for level, path := range logLevels {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			log.Errorf("创建级别日志目录失败: %v", err)
			continue
		}

		rotator, err := rotatelogs.New(
			path+".%Y%m%d%H",
			rotatelogs.WithRotationTime(rotationTime),
			rotatelogs.WithMaxAge(maxAge),
			rotatelogs.WithRotationSize(0),
		)

		if err != nil {
			log.Errorf("创建日志轮转器失败 %s: %v", level, err)
			continue
		}

		writerMap[level] = rotator
	}

	return writerMap
}

// 生成自动log_id: 毫秒时间戳+0~10000随机数
func generateLogID() string {
	millisecond := time.Now().UnixNano() / 1e6
	random := rand.Intn(RandomRange)
	return fmt.Sprintf("%d%d", millisecond, random)
}

// WithLogID 为context设置log_id
func WithLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, LogCtxID, logID)
}

// GetLogID 从context获取log_id，不存在则生成新的
func GetLogID(ctx context.Context) string {
	if ctx == nil {
		return generateLogID()
	}

	logID, ok := ctx.Value(LogCtxID).(string)
	if !ok || logID == "" {
		return generateLogID()
	}
	return logID
}

// 带context的日志方法，会自动处理log_id
func Infof(ctx context.Context, format string, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Logf(logrus.InfoLevel, format, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Logf(logrus.WarnLevel, format, args...)
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Logf(logrus.ErrorLevel, format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.InfoLevel, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.WarnLevel, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.ErrorLevel, args...)
}

// 新增其他级别日志方法
func Debugf(ctx context.Context, format string, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Logf(logrus.DebugLevel, format, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.DebugLevel, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Logf(logrus.FatalLevel, format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.FatalLevel, args...)
}
