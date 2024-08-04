# 使用官方的 Golang 镜像作为基础镜像
FROM golang:1.22-alpine

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