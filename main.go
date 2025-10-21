package main

import (
	"fmt"
	"os"

	"github.com/ygqygq2/influxdb-sync/cmd"
)

func main() {
	if len(os.Args) < 2 {
		cmd.ShowUsage()
		os.Exit(1)
	}

	cfgPath := os.Args[1]

	// 使用 cmd.Run 执行同步，自动识别版本
	if err := cmd.Run(cfgPath); err != nil {
		fmt.Println("同步失败:", err)
		os.Exit(1)
	}

	fmt.Println("同步完成")
}
