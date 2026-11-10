# 下雨天个人博客

一个前后端分离的个人博客项目，用于记录技术文章与生活随笔。前端基于 Vue 3 和 TypeScript，后端使用 Go、Gin 和 MySQL，支持 Markdown 写作、分类浏览、树洞和可选的 AI 助手。

[在线访问](https://www.raincode.top/) · [快速开始](#快速开始) · [部署指南](#部署) · [接口说明](#接口概览)

## 功能特性

- **文章阅读**：服务端分页、文章摘要、分类筛选、Markdown 渲染、代码高亮、目录与浏览量。
- **文章发布**：管理员登录后使用 Vditor 编辑文章，支持标签、封面与正文图片上传。
- **草稿保护**：自动保存到当前浏览器，支持手动保存、恢复、离开提醒及发布失败后保留内容。
- **个人空间**：博主资料、分类统计、树洞、关于页面和深色模式。
- **AI 助手**：可选接入阿里云百炼，提供请求频率、并发与每日调用量限制。
- **加载优化**：页面按需加载，列表使用分页和摘要，首页统计通过独立接口获取。

> 草稿保存在当前浏览器中，不会自动同步到其他设备。当前管理页面支持新文章发布；已有文章编辑入口尚未完善。

## 技术栈

| 层次 | 技术 |
| --- | --- |
| 前端 | Vue 3、TypeScript、Vite、Vue Router、Element Plus |
| 内容编辑与展示 | Vditor、Marked、highlight.js、DOMPurify |
| 后端 | Go、Gin、MySQL |
| 部署 | Nginx 静态托管与反向代理，可使用 Supervisor 管理后端进程 |

## 项目结构

```text
.
├── blog-front/              # Vue 前端
│   ├── src/api/             # API 请求
│   ├── src/components/      # 公共组件
│   ├── src/views/           # 页面
│   ├── src/utils/           # 草稿、认证与媒体处理等工具
│   └── tests/               # 前端测试
├── blog-back/               # Go 后端
│   ├── cmd/server/          # 服务入口
│   ├── internal/            # 路由、处理器、数据访问和中间件
│   └── scripts/schema.sql   # 数据库结构
├── deploy/                  # Nginx 配置示例
├── DEPLOY.md                # 通用部署说明
├── BAOTA_DEPLOY.md          # 宝塔部署说明
└── STAGE4.md                # 分页接口及验证说明
```

## 快速开始

### 环境要求

- Node.js **20.19.x 或 22.12 及以上版本**，以及 npm（Vite 7 的运行要求）。
- Go 1.21 或更高版本。
- MySQL 5.7 / 8.0。

### 1. 获取代码与初始化数据库

```bash
git clone https://github.com/cheng607/blog.git
cd blog
```

创建数据库并导入表结构：

```bash
mysql -u root -p -e "CREATE DATABASE IF NOT EXISTS blog DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
mysql -u root -p blog < blog-back/scripts/schema.sql
```

建表脚本中的数据库名为 `blog`。已有数据的环境请先备份并检查表结构，不要通过重建数据库的方式更新项目。

### 2. 配置并启动后端

```bash
cd blog-back
cp .env.example .env
# 编辑 .env：至少配置 DB_DSN、ADMIN_TOKEN、UPLOAD_PATH 和 CORS_ORIGINS
go run ./cmd/server
```

`DB_DSN` 使用 Go MySQL DSN 格式，需包含 `parseTime=true`；`UPLOAD_PATH` 应指向有写入权限的目录。开发环境的 `CORS_ORIGINS` 可包含 `http://localhost:5173`。后端默认监听 `9090`，可访问 `http://localhost:9090/health` 检查服务。

### 3. 启动前端

新开一个终端：

```bash
cd blog-front
cp .env.example .env.development
npm ci
npm run dev
```

默认访问 `http://localhost:5173`。开发服务器将 `/api` 请求代理到本地后端；可通过 `VITE_DEV_API_TARGET` 调整代理目标。

## 配置说明

后端配置示例见 [`blog-back/.env.example`](./blog-back/.env.example)，前端配置示例见 [`blog-front/.env.example`](./blog-front/.env.example)。

| 配置 | 用途 |
| --- | --- |
| `DB_DSN` | MySQL 连接地址及数据库配置 |
| `ADMIN_TOKEN` | 管理后台登录密钥，仅配置在后端 |
| `SERVER_PORT` | 后端端口，默认 `9090` |
| `UPLOAD_PATH` / `UPLOAD_URL_PREFIX` | 图片存储路径及访问前缀 |
| `CORS_ORIGINS` | 允许访问 API 的前端来源 |
| `DASHSCOPE_API_KEY` / `DASHSCOPE_MODEL` | 可选的 AI 服务密钥与模型配置 |
| `AI_RATE_PER_MINUTE` 等 | AI 请求频率、并发与每日额度限制 |
| `TRUSTED_PROXIES` | 可信反向代理的 IP 或网段 |
| `VITE_API_BASE_URL` | 前端 API 基址，同域反代通常使用 `/api/` |

请勿将真实密码、管理员密钥或 API Key 提交到仓库；前端 `VITE_` 变量会进入构建产物，不应存放密钥。

## 写作与发布

访问 `/admin`，使用后端配置的管理员密钥登录。首次使用时，在写文章页面展开「博主资料」并保存昵称；随后填写标题、正文，选择或创建分类，按需上传封面并发布。

封面及正文图片支持 JPG、JPEG、PNG，单张最大 10 MB、1600 万像素。草稿和封面文件依赖浏览器本地存储；请留意页面中的保存失败提示。

## 构建与测试

```bash
# 前端
cd blog-front
npm test
npm run build

# 后端（从 blog-front 目录切换）
cd ../blog-back
go test ./...
go build -o bin/blog-server ./cmd/server
```

前端构建产物位于 `blog-front/dist/`。后端编译命令中的输出路径相对于 `blog-back`；Windows 环境可按需调整可执行文件名。

## 部署

推荐使用 Nginx 托管前端，并将 `/api/` 反向代理到后端。站点使用 History 路由，Nginx 需要将非文件路径回退到 `index.html`；上传图片的访问路径也需按实际配置代理。

- [通用部署指南](./DEPLOY.md)
- [宝塔部署指南](./BAOTA_DEPLOY.md) · [宝塔快速开始](./BAOTA_QUICKSTART.md)
- [Nginx 配置示例](./deploy/nginx.conf.example)
- [分页接口变更与验证说明](./STAGE4.md)

**前后端需要配套更新**：新版文章列表接口返回分页对象，旧版后端返回数组，两者混用可能导致首页显示为空。更新前请备份数据库、上传文件和现有配置。

## 接口概览

下列为后端路径；使用默认 Nginx 代理时，浏览器请求路径通常带 `/api` 前缀。

| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | `/health` | 服务健康检查 |
| GET | `/site/summary` | 博主资料及站点统计 |
| GET | `/article/list` | 文章摘要分页，支持 `page`、`pageSize` |
| GET | `/article/tag/:tagId` | 按分类分页 |
| GET | `/article/:id` | 文章详情 |
| GET | `/tag` | 分类列表 |
| GET | `/hollow` | 树洞列表 |
| POST | `/auth/login` | 管理员登录 |
| POST | `/article` | 发布文章，需管理员会话 |
| POST | `/upload/image` | 上传图片，需管理员会话 |

分页响应中的 `data` 包含 `items`、`total`、`page` 和 `pageSize`。更多接口及具体路由可参考 [`router.go`](./blog-back/internal/router/router.go)。

## 参与与反馈

欢迎通过 [Issues](https://github.com/cheng607/blog/issues) 反馈问题或提出建议，也欢迎提交 Pull Request。提交问题时，请尽量附上复现步骤、运行环境和经过脱敏的日志。

## 开源许可

当前仓库尚未提供许可证文件；在作者明确选择并添加许可证前，请先与作者确认代码使用和分发权限。
