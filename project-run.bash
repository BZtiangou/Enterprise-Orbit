# 方法1: 使用 Makefile（推荐）
make deps    # 下载依赖
make build   # 编译项目
make run     # 运行项目

# # 方法2: 直接运行
# go run cmd/server/main.go

# # 方法3: 开发模式（需要先安装 air）
# go install github.com/cosmtrek/air@latest
# make dev