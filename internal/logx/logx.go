package logx

import (
	"fmt"
	"log"
	"os"
)

var (
	logger = log.New(os.Stdout, "[influxdb-sync] ", log.LstdFlags|log.Lshortfile)
)

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

