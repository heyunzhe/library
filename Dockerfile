FROM golang:1.23-alpine

# 安装 gcc、musl-dev 等依赖支持 cgo
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /app

# 拷贝依赖
COPY go.mod go.sum ./
RUN go mod download

# 拷贝源码
COPY . .

# 编译时启用 CGO
ENV CGO_ENABLED=1
RUN go build -o server main.go



EXPOSE 8080
CMD ["./server"]

