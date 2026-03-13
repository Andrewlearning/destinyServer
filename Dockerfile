# 基础镜像：带 Go 编译器的 Alpine Linux
FROM golang:1-alpine

# go-sqlite3 需要 CGO，安装 C 编译工具链
RUN apk add --no-cache gcc musl-dev

# 设置容器内工作目录为 /app
WORKDIR /app

# 把你本地项目所有文件拷进容器的 /app
COPY . .

# 腾讯 Go 模块代理（国内加速）
ENV GOPROXY=https://mirrors.tencent.com/go

# SQLite 需要 CGO 支持
ENV CGO_ENABLED=1

RUN go mod download

# 编译 cmd/ 下的 main 包，输出到 bin/
RUN go build -o bin/destinyServer ./cmd

# 容器启动时执行这个二进制
CMD ["bin/destinyServer"]

# 声明容器对外暴露 8080 端口
EXPOSE 8080
