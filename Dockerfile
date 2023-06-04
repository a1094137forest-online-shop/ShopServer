# 使用 golang 1.19 的官方映像作为基础
FROM golang:latest

# 设置工作目录
WORKDIR /ShopServer

# 将 ShopServer 项目的代码复制到容器中的工作目录
COPY . .

# 下载并安装依赖包（如果有）
RUN go mod download

# 构建 ShopServer
RUN go build -o main .

# 定义容器运行时需要开放的端口
EXPOSE 1234

# 指定容器启动时运行的命令
CMD ["./main"]