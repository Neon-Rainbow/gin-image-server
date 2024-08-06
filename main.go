package main

import (
	"ImageServer/config"
	"ImageServer/controllers"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	err := config.InitConfig()
	if err != nil {
		fmt.Println("Init config failed, err:", err)
	}

	router := gin.Default()

	// 创建用于存储上传图片的目录
	router.Static("/uploads", "./uploads")

	// 上传图片接口
	router.POST("/upload", controllers.UploadImage)

	// 获取图片接口
	router.GET("/image/:filename", controllers.GetImage)

	router.Run(":8080")
}
