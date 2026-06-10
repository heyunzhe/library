# ========== 构建阶段 ==========
FROM golang:1.23-alpine AS builder

WORKDIR /build

# 先复制依赖文件，利用 Docker 缓存
COPY go.mod go.sum ./
RUN go mod download

# 复制源码并编译
COPY . .
RUN CGO_ENABLED=0 go build -o library-server .

# ========== 运行阶段 ==========
FROM alpine:latest

WORKDIR /app

# 复制编译好的二进制
COPY --from=builder /build/library-server .

# 复制静态资源
COPY --from=builder /build/html ./html
COPY --from=builder /build/css ./css
COPY --from=builder /build/js ./js
COPY --from=builder /build/images ./images
COPY --from=builder /build/font ./font
COPY --from=builder /build/userphoto ./userphoto

EXPOSE 8080

CMD ["./library-server"]
