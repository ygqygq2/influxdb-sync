package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("用法: influxdb-sync <subcmd> [config.yaml]")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "influxdb1":
		os.Args = append([]string{os.Args[0]}, os.Args[2:]...)
		influxdb1main()
	default:
		fmt.Println("未知子命令: ", os.Args[1])
		os.Exit(1)
	}
}

func influxdb1main() {
	// 直接调用 cmd/influxdb1.go 的 main
	// go run . influxdb1 config.yaml
	// 这里采用 import _ "github.com/ygqygq2/influxdb-sync/cmd/influxdb1" 方式可进一步优化
}
