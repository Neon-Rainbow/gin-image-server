package main

import (
	"ImageServer/controllers"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// 创建用于存储上传图片的目录
	router.Static("/uploads", "./uploads")

	// 上传图片接口
	router.POST("/upload", controllers.UploadImage)

	// 获取图片接口
	router.GET("/image/:filename", controllers.GetImage)

	router.Run(":8080")
}
