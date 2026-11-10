.PHONY: dev dev-back dev-front build build-front build-back

# 本地同时开发：先启动后端，再启动前端（两个终端分别执行 dev-back / dev-front 亦可）
dev-back:
	cd blog-back && go run ./cmd/server

dev-front:
	cd blog-front && npm run dev

build-back:
	cd blog-back && go build -o bin/blog-server ./cmd/server

# 宝塔 Linux 服务器必须用 Linux 二进制，不能把 Windows 的 .exe 直接上传
build-back-linux:
	cd blog-back && set GOOS=linux&& set GOARCH=amd64&& go build -ldflags="-s -w" -o bin/blog-server-linux ./cmd/server

build-front:
	cd blog-front && npm run build

build: build-back build-front
