package controllers

import (
	"ImageServer/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
)

// UploadImage 处理图片上传
func UploadImage(c *gin.Context) {
	var serveBaseURL = fmt.Sprintf(
		"http://%v:%v/image",
		config.GetConfig().Server.Host,
		config.GetConfig().Server.Port,
	)

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to upload image",
		})
		return
	}

	// 获取文件扩展名
	ext := filepath.Ext(file.Filename)
	// 生成新的文件名
	newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	// 创建文件存储路径
	filePath := filepath.Join("uploads", newFileName)

	// 创建uploads目录
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create uploads directory",
		})
		return
	}

	// 保存图片
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save image",
		})
		return
	}

	fileURL := fmt.Sprintf("%s/%s", serveBaseURL, newFileName)

	// 返回图片存储路径
	c.JSON(http.StatusOK, gin.H{
		"message":  "Image uploaded successfully",
		"filepath": fileURL,
	})
}

// GetImage 处理图片获取
func GetImage(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("uploads", filename)

	// 检查图片是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Image not found",
		})
		return
	}

	// 返回图片
	c.File(filePath)
}
