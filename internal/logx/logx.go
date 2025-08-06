package logx

import (
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	logger       = log.New(os.Stdout, "[influxdb-sync] ", log.LstdFlags|log.Lshortfile)
	currentLevel = "info" // 默认级别
)

// 设置日志级别
func SetLevel(level string) {
	currentLevel = strings.ToLower(level)
}

// 检查是否应该输出指定级别的日志
func shouldLog(level string) bool {
	levels := map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}

	currentLevelInt, ok1 := levels[currentLevel]
	targetLevelInt, ok2 := levels[level]

	if !ok1 || !ok2 {
		return true // 如果级别不存在，默认输出
	}

	return targetLevelInt >= currentLevelInt
}

func Debug(v ...interface{}) {
	if shouldLog("debug") {
		logger.SetPrefix("[DEBUG] ")
		logger.Output(2, sprint(v...))
	}
}

func Info(v ...interface{}) {
	logger.SetPrefix("[INFO] ")
	logger.Output(2, sprint(v...))
}

func Warn(v ...interface{}) {
	logger.SetPrefix("[WARN] ")
	logger.Output(2, sprint(v...))
}

func Error(v ...interface{}) {
	logger.SetPrefix("[ERROR] ")
	logger.Output(2, sprint(v...))
}

func Fatal(v ...interface{}) {
	logger.SetPrefix("[FATAL] ")
	logger.Output(2, sprint(v...))
	os.Exit(1)
}

func sprint(v ...interface{}) string {
	return fmt.Sprint(v...)
}
