# 宝塔面板部署方案（下雨天个人博客）

适用于已安装 **宝塔 Linux 面板** 的服务器（如阿里云 47.114.122.96），域名示例：`www.raincode.top`。

整体思路：**Nginx 托管前端静态资源**，`/api` 与 `/upload/images` **反代到本机 Go 服务（9090）**，MySQL 由宝塔管理。

---

## 一、部署架构

```text
用户浏览器
    │
    ▼
Nginx (443/80)  ──► /          → blog-front/dist（静态）
              ├──► /api/       → 127.0.0.1:9090（Go）
              └──► /upload/images/ → 127.0.0.1:9090（图片）
    │
    ▼
Go blog-back (:9090) ──► MySQL (blog 库)
                      └──► /www/wwwroot/rain/upload/（图片文件）
```

---

## 二、宝塔环境准备

### 1. 安装软件（软件商店）

| 软件 | 用途 |
|------|------|
| Nginx | 网站与反向代理 |
| MySQL 5.7/8.0 | 数据库 |
| PM2 管理器（可选） | 也可用「Supervisor 管理器」守护 Go 进程 |

Go **不必**在软件商店安装：在本地或服务器编译好 Linux 二进制后上传即可。

### 2. 创建目录（文件管理或 SSH）

```bash
mkdir -p /www/wwwroot/rain/blog-front
mkdir -p /www/wwwroot/rain/blog-back
mkdir -p /www/wwwroot/rain/upload/images
chown -R www:www /www/wwwroot/rain
chmod -R 755 /www/wwwroot/rain
```

---

## 三、数据库（宝塔 → 数据库）

1. 打开 **数据库** → **添加数据库**
   - 数据库名：`blog`
   - 用户名：`blog`（或自定义）
   - 密码：自行设置（记下，写入后端 `.env`）
   - 访问权限：本地服务器

2. 点击数据库右侧 **管理**（phpMyAdmin）→ **导入**  
   上传项目中的：`blog-back/scripts/schema.sql`

3. 若已有历史数据，可跳过建表，仅确认表名与字段一致。

---

## 四、后端 Go 服务

### 方式 A：在 Windows 本机交叉编译（推荐）

在项目根目录 PowerShell：

```powershell
cd D:\project\Blog\blog-back
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o blog-server ./cmd/server
```

得到 `blog-server`（无扩展名），用宝塔 **文件** 上传到：

`/www/wwwroot/rain/blog-back/blog-server`

### 方式 B：在服务器上编译

SSH 安装 Go 1.21+ 后：

```bash
cd /www/wwwroot/rain/blog-back
# 上传整个 blog-back 源码目录后
go build -o blog-server ./cmd/server
```

### 配置 `.env`

在 `/www/wwwroot/rain/blog-back/.env` 新建：

```env
SERVER_PORT=9090
GIN_MODE=release

DB_DSN=blog:你的数据库密码@tcp(127.0.0.1:3306)/blog?parseTime=true&loc=Asia%2FShanghai&charset=utf8mb4

UPLOAD_PATH=/www/wwwroot/rain/upload/
UPLOAD_URL_PREFIX=/upload/images/

CORS_ORIGINS=https://www.raincode.top,https://raincode.top,http://www.raincode.top,http://raincode.top,https://raincode.top

# 管理员 Token（保护发布文章等接口，必须配置）
# 管理密钥仅配置在后端 ADMIN_TOKEN，前端不配置密钥。
ADMIN_TOKEN=请设置一串复杂随机字符串

# AI 助手（阿里百炼 DashScope，密钥仅在服务端配置，勿写入前端）
DASHSCOPE_API_KEY=你的百炼API密钥
DASHSCOPE_MODEL=qwen-plus,qwen-turbo
DASHSCOPE_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
```

### 只上传 `blog-server` 为什么还会报错？

上传二进制只是其中一步，缺一不可：

| 步骤 | 说明 |
|------|------|
| 1. **Linux 二进制** | 在 Windows 上必须用 `$env:GOOS="linux"` 交叉编译，得到**无 `.exe` 后缀**的 `blog-server`。若上传 `blog-server.exe`，Linux **无法执行**，Supervisor 往往仍在跑**旧进程**。 |
| 2. **覆盖路径一致** | 文件必须覆盖 Supervisor 里配置的同一路径，例如 `/www/wwwroot/rain/blog-back/blog-server`。 |
| 3. **重启进程** | 替换文件后必须在 Supervisor 点 **重启**（或 `supervisorctl restart blog-back`），否则内存里仍是旧程序。 |
| 4. **配置 `.env`** | `DASHSCOPE_API_KEY` 写在 `.env` 里，**不会**打进二进制；没配置会提示「AI 未配置」，而不是 404。 |
| 5. **前端 dist** | AI 请求走 `/api/chat`，需重新 `npm run build` 并上传 `dist`（若仍是旧前端，也可能请求错路径）。 |

**先确认版本是否已更新**（新程序 `health` 会带 `version` 字段）：

```bash
curl -s http://127.0.0.1:9090/health
# 期望含： "version":"2026.06.04-chat1" 且 features 含 chat
curl -s http://127.0.0.1:9090/chat/status
# 期望： {"code":200,"data":{"configured":true,...}}
```

若 `health` 没有 `version` 或 `chat/status` 返回 404，说明**当前运行的仍是旧 blog-server**，请按上文重新编译、上传、chmod +x、Supervisor 重启。

更新代码后若 AI 助手提示「请求地址不存在」，说明线上仍是旧版 `blog-server`（无 `/chat` 接口）。需重新 `go build` 上传并 Supervisor 重启，再用 SSH 验证：

```bash
curl -s -X POST http://127.0.0.1:9090/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"你好"}'
```

应返回 JSON 且 `code` 为 200；若 `AI 服务未配置` 则检查 `.env` 中 `DASHSCOPE_API_KEY`。

> `DB_DSN` 中密码若含特殊字符需 URL 编码。`127.0.0.1` 表示 MySQL 与网站同机。

### 赋予执行权限

```bash
chmod +x /www/wwwroot/rain/blog-back/blog-server
chown -R www:www /www/wwwroot/rain/upload
```

### 用 Supervisor 守护进程（宝塔 → 软件商店 → Supervisor）

1. 添加守护进程：
   - 名称：`blog-back`
   - 运行目录：`/www/wwwroot/rain/blog-back`
   - 启动命令：`/www/wwwroot/rain/blog-back/blog-server`
   - 进程数量：1

2. 保存并启动，日志中应出现：`blog-back listening on :9090`

3. SSH 自测：

```bash
curl http://127.0.0.1:9090/health
# 期望：{"status":"ok","version":"2026.06.04-chat1","features":["chat",...]}
```

---

## 五、前端构建与上传

### 1. 本地构建（Windows）

```powershell
cd D:\project\Blog\blog-front
# 确认 .env.production 为：
# VITE_API_BASE_URL=/api/
npm install
npm run build
```

### 2. 上传 dist

将 `blog-front/dist` 内**全部文件**上传到：

`/www/wwwroot/rain/blog-front/`

（即该目录下直接有 `index.html` 和 `assets/` 文件夹）

---

## 六、宝塔添加网站与 Nginx

### 1. 添加站点

**网站** → **添加站点**

- 域名：`www.raincode.top`（可同时填 `raincode.top`）
- 根目录：`/www/wwwroot/rain/blog-front`
- PHP：纯静态，**不需要** PHP

### 2. 配置 SSL

站点设置 → **SSL** → Let's Encrypt 申请证书 → 开启 **强制 HTTPS**

### 3. 修改 Nginx 配置

站点 → **设置** → **配置文件**，在 `server { ... }` 内找到 `location /` 块，**在其前面或替换为**以下完整逻辑（保留你原有的 SSL、listen 等行）：

```nginx
    # API 反代到 Go
    location /api/ {
        proxy_pass http://127.0.0.1:9090/;
        proxy_connect_timeout 5s;
        proxy_send_timeout 20s;
        proxy_read_timeout 85s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        client_max_body_size 12m;
    }

    # 上传图片（封面、编辑器）
    location /upload/images/ {
        proxy_pass http://127.0.0.1:9090/upload/images/;
        proxy_set_header Host $host;
        expires 30d;
    }

    location / {
        try_files $uri $uri/ /index.html;
    }
```

保存后 **重载 Nginx**。

> 本项目使用 **History 路由**（如 `/article?id=1`），必须配置 `try_files` 回退到 `index.html`，否则刷新子页面会 404。

### 4. 放行端口（安全）

- 对外只需 **80 / 443**
- **9090 不要**在防火墙对公网开放，仅本机 Nginx 访问

---

## 七、部署后验证清单

在浏览器依次检查：

| 步骤 | 地址 | 预期 |
|------|------|------|
| 1 | `https://www.raincode.top/` | 首页文章列表 |
| 2 | `https://www.raincode.top/api/health` 经反代 | 若未配置 /api/health 可 SSH curl 本机 |
| 3 | 点开一篇文章 | 正文、目录、浏览量正常 |
| 4 | 分类页 → 某分类 | 筛选文章正确 |
| 5 | 管理后台 `/admin` | 密钥登录后可发布文章 |
| 6 | 树洞 | 发布与列表刷新 |
| 7 | F12 网络 | 接口为 `https://域名/api/...`，无跨域报错 |

SSH 快速测 API（本机）：

```bash
curl http://127.0.0.1:9090/article/list
curl http://127.0.0.1:9090/tag
```

---

## 八、代码修改后如何更新部署（重要）

上线后改代码，**不必重装环境**，按改动范围操作即可。

### 判断改动了什么

| 改动位置 | 需要更新的部分 |
|----------|----------------|
| 只改了 `blog-front/`（页面、样式、组件） | **仅前端** |
| 只改了 `blog-back/`（Go 接口、逻辑） | **仅后端** |
| 前后端都改了 | **先后端再前端**（或只更新你改的那一侧） |
| 改了 `scripts/schema.sql`（表结构） | **数据库** + 可能后端 |

---

### 仅更新前端（最常见）

在**你自己电脑**上：

```powershell
cd D:\project\Blog\blog-front
# 生产构建前在 .env.production 中设置（勿提交到 Git）：
# 管理密钥仅配置在后端 ADMIN_TOKEN，前端不配置密钥。
npm run build
```

> 写文章入口：`/admin` 或 `/edit`。不支持 URL 携带密钥登录。

在**宝塔 → 文件**：

1. 进入 `/www/wwwroot/rain/blog-front/`
2. **删除** 该目录下旧的 `assets` 文件夹和 `index.html`（或整目录清空后上传，避免残留旧 JS）
3. 将本地 `blog-front\dist\` 里的**全部内容**上传进去（含新的 `index.html` 和 `assets/`）

无需重启 Nginx、无需动 Go 服务。浏览器 **Ctrl + F5** 强制刷新。

> 若只改了一两个文件仍建议全量覆盖 `dist`，避免 `index.html` 引用的 JS 文件名（带 hash）不一致导致白屏。

---

### 仅更新后端

在**你自己电脑**上交叉编译 Linux 可执行文件：

```powershell
cd D:\project\Blog\blog-back
$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o blog-server ./cmd/server
```

在**宝塔**：

1. **文件** → `/www/wwwroot/rain/blog-back/`
2. 上传新的 `blog-server`，**覆盖**原文件
3. 若改过环境变量，编辑服务器上的 `.env`（**不要**覆盖成开发机密码，除非你有意修改）
4. **Supervisor** → 找到 `blog-back` → **重启**

验证：

```bash
curl http://127.0.0.1:9090/health
```

返回 `{"status":"ok"}` 即正常。

> **不要**上传 `blog-server.exe`（Windows 版），服务器只能运行 Linux 版 `blog-server`。

---

### 前后端都更新

1. 按上一节 **先更新后端** 并 Supervisor 重启
2. 再按 **仅更新前端** 覆盖 `dist`
3. 浏览器强刷后全站走一遍（首页、文章、发布、树洞）

---

### 改了数据库结构

1. 宝塔 → **数据库** → 对应库 → 管理 → 执行增量 SQL（备份后再改）
2. 若接口逻辑也变了，再 **更新后端**
3. 前端一般不用动

---

### 更新检查清单（每次发版勾一下）

- [ ] Supervisor 里 `blog-back` 状态为运行中
- [ ] `curl 127.0.0.1:9090/health` 正常
- [ ] 网站首页能打开，F12 无大量 404
- [ ] 接口请求地址为 `https://你的域名/api/...`
- [ ] 新功能路径手动点一遍

---

## 九、常见问题

### 1. 接口 502 Bad Gateway

- Supervisor 中 `blog-back` 是否运行
- `curl 127.0.0.1:9090/health` 是否通
- Nginx `proxy_pass` 末尾是否有 `/`（与本文一致）

### 2. 跨域错误

- 检查 `.env` 中 `CORS_ORIGINS` 是否包含完整前端域名（含 `https://`）
- 修改后重启 Go 进程

### 3. 图片/封面 404

- `UPLOAD_PATH` 目录是否存在且 `www` 可写
- Nginx 是否配置了 `/upload/images/` 反代
- 数据库里 `cover_image` 是否为 `/upload/images/...` 格式

### 4. 数据库连接失败

- 宝塔数据库账号密码与 `DB_DSN` 一致
- MySQL 是否仅允许 `localhost`（同机部署用 `127.0.0.1`）

### 5. 更新后页面没变化

- 是否覆盖上传了完整 `dist`（含 `index.html`）
- 浏览器是否强刷（Ctrl+F5）或清缓存
- 是否误传到了别的网站根目录

### 6. 更新后端后 502

- 是否上传的是 **Linux** 版 `blog-server` 且已 `chmod +x`
- Supervisor 是否已重启，日志有无报错

---

## 十、备案与安全建议

管理员在 `/admin` 输入后端 `ADMIN_TOKEN` 登录；前端不保存或构建管理密钥。
2. 数据库密码、`.env` 不要提交到 Git。
3. 定期在宝塔备份 **数据库** + **`/www/wwwroot/rain/upload`** 目录。

---

## 十一、与项目文件对应关系

| 宝塔路径 | 项目来源 |
|----------|----------|
| `/www/wwwroot/rain/blog-front/` | `blog-front/dist` |
| `/www/wwwroot/rain/blog-back/blog-server` | `go build` 产物 |
| `/www/wwwroot/rain/blog-back/.env` | 参考 `blog-back/.env.example` |
| `/www/wwwroot/rain/upload/` | 运行时上传目录 |
| MySQL `blog` 库 | `blog-back/scripts/schema.sql` |

更通用的说明见 [DEPLOY.md](./DEPLOY.md)。
