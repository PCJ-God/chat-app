# ---- build ----
FROM golang:1.27-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
ENV GOOS=linux
# 改用国内 Go 模块代理
ENV GOPROXY=https://goproxy.cn,direct
# 清空代理环境变量，避免容器内走宿主机 127.0.0.1 代理（不可达）导致下载失败
ENV HTTP_PROXY=
ENV HTTPS_PROXY=
ENV ALL_PROXY=
ENV http_proxy=
ENV https_proxy=
ENV all_proxy=
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/chat-app ./cmd/server

# ---- runtime ----
FROM alpine:3.20
ENV http_proxy=
ENV https_proxy=
ENV all_proxy=
ENV HTTP_PROXY=
ENV HTTPS_PROXY=
ENV ALL_PROXY=
# apk 走阿里云镜像，安装 CA 证书和时区数据
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/chat-app ./chat-app
# 打包一份默认配置；部署时由 compose 用 config/config.docker.yaml 覆盖
COPY --from=builder /app/config ./config
EXPOSE 8080
ENTRYPOINT ["./chat-app"]
