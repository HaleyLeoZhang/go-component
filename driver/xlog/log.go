package xlog

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
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
	// 设置日志格式为JSON，指定中国时区时间格式
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05", // 目标时间格式
	})

	// 设置日志级别
	log.SetLevel(logrus.DebugLevel)

	// 禁用logrus默认的stdout输出
	log.SetOutput(io.Discard)  // 添加这一行，禁用默认输出

	// 配置日志分割和级别分离
	writerMap := setupLogRotation(c)
	if writerMap != nil {
		log.AddHook(lfshook.NewHook(writerMap, log.Formatter))
	}

	// 配置标准输出
	if c.Stdout == StdoutYes {
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
		logrus.DebugLevel: filepath.Join(c.Dir, "debug.log"),
		logrus.InfoLevel:  filepath.Join(c.Dir, "info.log"),
		logrus.WarnLevel:  filepath.Join(c.Dir, "warn.log"),
		logrus.ErrorLevel: filepath.Join(c.Dir, "error.log"),
		logrus.FatalLevel: filepath.Join(c.Dir, "fatal.log"),
		logrus.PanicLevel: filepath.Join(c.Dir, "panic.log"),
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

// GenerateLogID 生成自动log_id: 毫秒时间戳+0~10000随机数
func GenerateLogID() string {
	now := time.Now()

	// 提取小时、分钟、秒和毫秒
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()
	millisecond := now.Nanosecond() / 1e6 // 转换纳秒到毫秒

	// 组合时间部分：小时(2位)+分钟(2位)+秒(2位)+毫秒(3位)
	timePart := hour*10000000 + minute*100000 + second*1000 + millisecond

	// 生成随机数
	random := rand.Intn(RandomRange)

	// 组合并返回结果
	return fmt.Sprintf("%d%d", timePart, random)
}

// WithLogID 为context设置log_id
func WithLogID(ctx context.Context, logID string) context.Context {
	return context.WithValue(ctx, LogCtxID, logID)
}

// GetLogID 从context获取log_id，不存在则生成新的
func GetLogID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	logID, ok := ctx.Value(LogCtxID).(string)
	if !ok || logID == "" {
		return ""
	}
	return logID
}

// 带context的日志方法，自动处理log_id，兼容纯字符串format
func Infof(ctx context.Context, format string, args ...interface{}) {
	entry := log.WithField("log_id", GetLogID(ctx))
	if len(args) == 0 {
		// 没有格式化参数时，直接打印字符串
		entry.Log(logrus.InfoLevel, format)
	} else {
		// 有格式化参数时，使用格式化打印
		entry.Logf(logrus.InfoLevel, format, args...)
	}
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	entry := log.WithField("log_id", GetLogID(ctx))
	if len(args) == 0 {
		entry.Log(logrus.WarnLevel, format)
	} else {
		entry.Logf(logrus.WarnLevel, format, args...)
	}
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	entry := log.WithField("log_id", GetLogID(ctx))
	if len(args) == 0 {
		entry.Log(logrus.ErrorLevel, format)
	} else {
		entry.Logf(logrus.ErrorLevel, format, args...)
	}
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	entry := log.WithField("log_id", GetLogID(ctx))
	if len(args) == 0 {
		entry.Log(logrus.DebugLevel, format)
	} else {
		entry.Logf(logrus.DebugLevel, format, args...)
	}
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	entry := log.WithField("log_id", GetLogID(ctx))
	if len(args) == 0 {
		entry.Log(logrus.FatalLevel, format)
	} else {
		entry.Logf(logrus.FatalLevel, format, args...)
	}
}

// 非格式化日志方法保持不变
func Info(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.InfoLevel, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.WarnLevel, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.ErrorLevel, args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.DebugLevel, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	log.WithField("log_id", GetLogID(ctx)).Log(logrus.FatalLevel, args...)
}
