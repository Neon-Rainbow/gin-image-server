# ImageServer

一个使用 Gin 框架构建的简单图片上传和获取服务。项目使用 Docker 和 Docker Compose 进行容器化部署。

[English version](README.md)
## 项目结构

```
ImageServer/
├── Dockerfile
├── LICENSE
├── README-zh_cn.md
├── README.md
├── config
│   └── config.go
├── config.yaml
├── controllers
│   └── image_controller.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── models
│   └── image.go
└── uploads


```

config.yaml内容:
```yaml
server:
  host: 127.0.0.1 # 服务器地址
  port: 8080 # 服务器端口
```
返回的url的格式:
```
http://{host}:{post}/image/{image_name}
```
可以直接访问该地址来访问图片

## 功能介绍

- **上传图片**：将图片上传到服务器并返回存储路径。
- **获取图片**：通过给定的图片路径获取图片。

## 准备工作

确保已安装以下工具：

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## 快速开始

### 1. 克隆仓库


### 2. 构建并启动服务

使用 Docker Compose 构建并启动服务：

```sh
docker-compose up --build
```

服务将会在 `http://localhost:8080` 运行。

### 3. 上传图片

使用 `curl` 命令上传图片：

```sh
curl -X POST http://localhost:8080/upload -F "image=@/path/to/your/image.jpg"
```

成功上传后，服务器将返回图片的存储路径。

### 4. 获取图片

使用 `curl` 命令获取图片：

```sh
curl http://localhost:8080/image/your_image_filename.jpg --output downloaded_image.jpg
```

## 代码说明

### `main.go`

设置路由和启动 Gin 服务。

```go
package main

import (
    "github.com/gin-gonic/gin"
    "ImageServer/controllers"
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
```

### `controllers/image_controller.go`

处理图片上传和获取的控制器。

```go
package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "net/http"
    "path/filepath"
    "fmt"
    "os"
)

// UploadImage 处理图片上传
func UploadImage(c *gin.Context) {
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

    // 返回图片存储路径
    c.JSON(http.StatusOK, gin.H{
        "message": "Image uploaded successfully",
        "filepath": filePath,
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
```

### `models/image.go`

图片模型定义及文件检查方法。

```go
package models

import (
    "os"
)

// Image 图片模型
type Image struct {
    Filename string `json:"filename"`
    Filepath string `json:"filepath"`
}

// GetImage 检查图片是否存在
func GetImage(filepath string) (*Image, error) {
    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        return nil, err
    }
    return &Image{
        Filepath: filepath,
    }, nil
}
```

### `Dockerfile`

用于构建 Docker 镜像的文件。

```Dockerfile
# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.20-alpine

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 并下载依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制项目的源代码
COPY . .

# 构建 Go 应用
RUN go build -o main .

# 暴露应用运行的端口
EXPOSE 8080

# 运行 Go 应用
CMD ["./main"]
```

### `docker-compose.yml`

用于定义和运行多容器 Docker 应用的配置文件。

```yaml
version: '3.8'

services:
  imageserver:
    build: .
    ports:
      - "8080:8080"
    volumes:
      - ./uploads:/app/uploads
```

## 持久化存储

使用 Docker Compose 中的 `volumes` 配置，确保上传的文件会保存在主机的 `uploads` 目录中，即使容器被删除或重启，文件数据仍然会保留。

## 贡献

欢迎提交问题和合并请求。如果你有任何建议或发现任何问题，请创建 Issue 或提交 Pull Request。

## 许可证

该项目使用 MIT 许可证。详情请参见 [LICENSE](LICENSE) 文件。
```

这个 `README.md` 文件包含了项目的简介、功能介绍、快速开始指南、代码说明、以及贡献和许可证信息，能够帮助用户快速了解和使用你的项目。