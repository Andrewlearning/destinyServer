# 使用基于 alpine 的 golang 官方镜像
FROM golang:1-alpine

# go-sqlite3 需要 CGO，安装 C 编译工具链
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY . .

ENV GOPROXY=https://mirrors.tencent.com/go
# SQLite 需要 CGO 支持
ENV CGO_ENABLED=1

RUN go mod download

RUN go build -o wxcloudrun-files/main

CMD ["wxcloudrun-files/main"]

EXPOSE 8080