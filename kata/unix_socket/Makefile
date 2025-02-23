# Makefile

# 定义变量
MODULE_NAME := github.com/walterfan/kata-go/unix_socket
SERVER_BINARY := server
CLIENT_BINARY := client

# 默认目标
all: server client

# 构建 server
server:
	go build -o $(SERVER_BINARY) server.go common.go

# 构建 client
client:
	go build -o $(CLIENT_BINARY) client.go common.go

# 清理构建文件
clean:
	rm -f $(SERVER_BINARY) $(CLIENT_BINARY)

# 运行 server
run-server:
	./$(SERVER_BINARY)

# 运行 client
run-client:
	./$(CLIENT_BINARY)