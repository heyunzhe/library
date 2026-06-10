FROM nginx:latest

WORKDIR /app

# 复制 Go 后端
COPY library-server .

# 复制静态文件
COPY html ./html
COPY css ./css
COPY js ./js
COPY images ./images
COPY font ./font
COPY userphoto ./userphoto

# 暴露端口
EXPOSE 8080

# 运行 Go 后端
CMD ["./library-server"]
