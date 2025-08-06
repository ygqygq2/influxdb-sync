package utils

import "fmt"

// PrintProgress 进度条式输出
func PrintProgress(db string, done, total int, m string) {
	fmt.Printf("\r[%-20s] %d/%d 正在同步: %s", db, done, total, m)
}

// PrintProgressDone 输出换行结束进度条显示
func PrintProgressDone() {
	fmt.Println()
}
