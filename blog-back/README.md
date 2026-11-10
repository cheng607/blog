# blog-back

Go (Gin) + MySQL 个人博客 API，与 `blog-front` 配套。

## 快速开始

```bash
# 1. 初始化数据库
mysql -u root -p < scripts/schema.sql

# 2. 配置环境
cp .env.example .env

# 3. 启动
make run
# 或 go run ./cmd/server
```

- 服务地址：`http://localhost:9090`
- 健康检查：`GET /health`

## 环境变量

| 变量 | 说明 |
|------|------|
| `SERVER_PORT` | 监听端口，默认 9090 |
| `DB_DSN` | MySQL DSN（需 `parseTime=true`） |
| `UPLOAD_PATH` | 上传文件根目录 |
| `UPLOAD_URL_PREFIX` | 图片 URL 前缀，默认 `/upload/images/` |
| `CORS_ORIGINS` | 允许的前端来源，逗号分隔 |
| `GIN_MODE` | 设为 `debug` 可开启 Gin 调试 |

## 项目结构

```
blog-back/
├── cmd/server/       # 入口
├── internal/
│   ├── config/       # 配置
│   ├── router/       # 路由注册
│   ├── handler/      # HTTP 处理
│   ├── store/        # 数据访问
│   ├── model/        # 模型
│   ├── middleware/   # CORS 等
│   ├── upload/       # 图片上传
│   └── response/     # 统一响应 {code,msg,data}
└── scripts/schema.sql
```

## API 约定

与历史 Java 版保持一致，统一 JSON：

```json
{ "code": 200, "msg": "成功", "data": {} }
```

主要路径：`/article`、`/tag`、`/user`、`/hollow`、`/upload/image`。
