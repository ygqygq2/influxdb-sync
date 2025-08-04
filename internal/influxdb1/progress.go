package influxdb1

import "fmt"

// 进度条式输出
func PrintProgress(db string, done, total int, m string) {
	   fmt.Printf("\r[%-20s] %d/%d 正在同步: %s", db, done, total, m)
}

func PrintProgressDone() {
	fmt.Println()
}
