# ImageServer

A simple image upload and retrieval service built using the Gin framework. The project is containerized using Docker and Docker Compose.

[简体中文版](./README-zh_cn.md)

## Project Structure

```
ImageServer/
├── Dockerfile
├── README.md
├── controllers
│   └── image_controller.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
└── models
    └── image.go
```

## Features

- **Image Upload**: Upload images to the server and return the storage path.
- **Image Retrieval**: Retrieve images by the given image path.

## Prerequisites

Ensure the following tools are installed:

- [Docker](https://www.docker.com/get-started)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Quick Start

### 1. Clone the Repository

```sh
git clone <repository-url>
cd ImageServer
```

### 2. Build and Start the Service

Use Docker Compose to build and start the service:

```sh
docker-compose up --build
```

The service will run on `http://localhost:8080`.

### 3. Upload an Image

Use the `curl` command to upload an image:

```sh
curl -X POST http://localhost:8080/upload -F "image=@/path/to/your/image.jpg"
```

Upon successful upload, the server will return the storage path of the image.

### 4. Retrieve an Image

Use the `curl` command to retrieve an image:

```sh
curl http://localhost:8080/image/your_image_filename.jpg --output downloaded_image.jpg
```

## Code Explanation

### `main.go`

Sets up routes and starts the Gin server.

```go
package main

import (
    "github.com/gin-gonic/gin"
    "ImageServer/controllers"
)

func main() {
    router := gin.Default()

    // Serve static files from the uploads directory
    router.Static("/uploads", "./uploads")

    // Image upload endpoint
    router.POST("/upload", controllers.UploadImage)

    // Image retrieval endpoint
    router.GET("/image/:filename", controllers.GetImage)

    router.Run(":8080")
}
```

### `controllers/image_controller.go`

Controller handling image upload and retrieval.

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

// UploadImage handles image uploads
func UploadImage(c *gin.Context) {
    file, err := c.FormFile("image")
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Failed to upload image",
        })
        return
    }

    // Get file extension
    ext := filepath.Ext(file.Filename)
    // Generate new file name
    newFileName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
    // Create file storage path
    filePath := filepath.Join("uploads", newFileName)

    // Create uploads directory if not exists
    if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create uploads directory",
        })
        return
    }

    // Save the image
    if err := c.SaveUploadedFile(file, filePath); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save image",
        })
        return
    }

    // Return the image storage path
    c.JSON(http.StatusOK, gin.H{
        "message": "Image uploaded successfully",
        "filepath": filePath,
    })
}

// GetImage handles image retrieval
func GetImage(c *gin.Context) {
    filename := c.Param("filename")
    filePath := filepath.Join("uploads", filename)

    // Check if the image exists
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Image not found",
        })
        return
    }

    // Return the image
    c.File(filePath)
}
```

### `models/image.go`

Defines the image model and file check method.

```go
package models

import (
    "os"
)

// Image model
type Image struct {
    Filename string `json:"filename"`
    Filepath string `json:"filepath"`
}

// GetImage checks if the image exists
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

File used to build the Docker image.

```Dockerfile
# Use the official Golang image as the base image
FROM golang:1.20-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the project source code
COPY . .

# Build the Go application
RUN go build -o main .

# Expose the port the app runs on
EXPOSE 8080

# Run the Go application
CMD ["./main"]
```

### `docker-compose.yml`

Configuration file to define and run the multi-container Docker application.

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

## Persistent Storage

By using the `volumes` configuration in Docker Compose, uploaded files will be saved in the host's `uploads` directory. This ensures that the file data remains even if the container is deleted or restarted.

## Contribution

Feel free to submit issues and merge requests. If you have any suggestions or find any problems, please create an issue or submit a pull request.

## License

This project is licensed under the MIT License. For details, please see the [LICENSE](LICENSE) file.

---

This `README.md` provides an overview, feature description, quick start guide, code explanation, and contribution and license information to help users quickly understand and use your project.