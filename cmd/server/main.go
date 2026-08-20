package main

import (
	"blog/internal/app"
)

func main() {
	// 创建应用实例
	application := app.NewApp()

	// 启动应用
	application.Run()
}
