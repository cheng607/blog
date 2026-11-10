# 宝塔面板快速部署指南（最简化版）

> 目标：**本地构建 → 上传文件 → 宝塔配置**，5 分钟完成部署。

---

## 一、准备工作（宝塔面板）

### 1. 安装软件
- **Nginx**（网站）
- **MySQL 5.7/8.0**（数据库）
- **Supervisor 管理器**（守护 Go 进程）

### 2. 创建目录（文件管理 → 新建目录）
```
/www/wwwroot/rain/blog-front/
/www/wwwroot/rain/blog-back/
/www/wwwroot/rain/upload/images/
```
设置权限：`www:www`，755

### 3. 创建数据库
宝塔 → 数据库 → 添加数据库
- 库名：`blog`
- 用户：`blog`
- 密码：**记下来**（后面写 `.env` 用）
- 权限：本地服务器

导入表结构：上传 `blog-back/scripts/schema.sql`

---

## 二、本地构建（Windows 电脑）

### 1. 前端构建
```powershell
cd D:\project\Blog\blog-front

# 确认 .env.production 内容为：
# VITE_API_BASE_URL=/api/
# 管理密钥仅配置在后端 ADMIN_TOKEN，前端不配置密钥。

npm install
npm run build
```
产物：`blog-front/dist/`（里面有 `index.html` + `assets/`）

### 2. 后端交叉编译（Linux 二进制）
```powershell
cd D:\project\Blog\blog-back

$env:GOOS="linux"
$env:GOARCH="amd64"
go build -o blog-server ./cmd/server
```
得到 `blog-server`（**无扩展名**）。

---

## 三、上传文件到宝塔

### 必须上传的文件清单

| 本地文件 | 宝塔目标路径 | 说明 |
|---------|-------------|------|
| `blog-front/dist/*`（全部） | `/www/wwwroot/rain/blog-front/` | 前端静态资源 |
| `blog-back/blog-server` | `/www/wwwroot/rain/blog-back/blog-server` | Go 可执行文件 |
| `blog-back/.env`（新建） | `/www/wwwroot/rain/blog-back/.env` | 环境变量（见下方模板） |
| `blog-back/scripts/schema.sql` | 导入到 MySQL | 数据库表结构（首次） |

### 后端 `.env` 模板（宝塔上新建文件）
路径：`/www/wwwroot/rain/blog-back/.env`

```env
SERVER_PORT=9090
GIN_MODE=release

DB_DSN=blog:你的数据库密码@tcp(127.0.0.1:3306)/blog?parseTime=true&loc=Asia%2FShanghai&charset=utf8mb4

UPLOAD_PATH=/www/wwwroot/rain/upload/
UPLOAD_URL_PREFIX=/upload/images/

CORS_ORIGINS=https://你的域名,https://你的域名,http://你的域名

# 管理员 Token（保护发布文章等接口，必须配置）
# 管理密钥仅配置在后端 ADMIN_TOKEN，前端不配置密钥。
ADMIN_TOKEN=请设置一串复杂随机字符串

# AI 助手（可选）
DASHSCOPE_API_KEY=你的百炼密钥
DASHSCOPE_MODEL=qwen-plus,qwen-turbo
DASHSCOPE_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1
```

管理员在 `/admin` 输入后端 `ADMIN_TOKEN` 登录；前端不保存或构建管理密钥。

### 赋予执行权限（SSH 或宝塔终端）
```bash
chmod +x /www/wwwroot/rain/blog-back/blog-server
chown -R www:www /www/wwwroot/rain/upload
```

---

## 四、宝塔配置

### 1. Supervisor 守护 Go 进程
宝塔 → Supervisor 管理器 → 添加进程
- 名称：`blog-back`
- 运行目录：`/www/wwwroot/rain/blog-back`
- 启动命令：`/www/wwwroot/rain/blog-back/blog-server`
- 进程数量：1

启动后，日志应出现：`blog-back listening on :9090`

### 2. 添加网站（Nginx）
宝塔 → 网站 → 添加站点
- 域名：`你的域名`
- 根目录：`/www/wwwroot/rain/blog-front`
- PHP：纯静态

### 3. 配置 SSL（可选但推荐）
站点设置 → SSL → Let's Encrypt 申请 → 开启**强制 HTTPS**

### 4. 修改 Nginx 配置
站点 → 设置 → 配置文件，**替换** `location /` 相关内容为：

```nginx
# API 反代
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

# 上传图片
location /upload/images/ {
    proxy_pass http://127.0.0.1:9090/upload/images/;
    proxy_set_header Host $host;
    expires 30d;
}

# History 路由回退
location / {
    try_files $uri $uri/ /index.html;
}
```

保存 → **重载 Nginx**

### 5. 安全设置
- 防火墙仅放行 **80/443**
- **9090 端口不要对公网开放**

---

## 五、验证部署

### 浏览器访问
- `https://你的域名/` → 首页文章列表
- `https://你的域名/api/health` → `{"status":"ok"}`
- 点开文章 → 正文、目录、浏览量正常
管理员在 `/admin` 输入后端 `ADMIN_TOKEN` 登录；前端不保存或构建管理密钥。

### SSH 自检
```bash
curl http://127.0.0.1:9090/health
curl http://127.0.0.1:9090/article/list
```

---

## 六、更新部署（改代码后）

### 只改前端
1. 本地 `cd blog-front && npm run build`
2. 宝塔文件 → 删除旧 `blog-front` 目录内容 → 上传新 `dist/*`
3. 浏览器 **Ctrl+F5** 强制刷新

### 只改后端
1. 本地交叉编译 `blog-server`
2. 宝塔上传覆盖 `/www/wwwroot/rain/blog-back/blog-server`
3. Supervisor → 重启 `blog-back`

### 前后端都改
先更新后端 → 再更新前端

---

## 七、常见问题

| 问题 | 解决 |
|------|------|
| 502 Bad Gateway | Supervisor 检查 `blog-back` 是否运行，`curl 127.0.0.1:9090/health` 是否通 |
| 接口 404 | Nginx `proxy_pass` 末尾是否有 `/`（与本文一致） |
| 图片 404 | `UPLOAD_PATH` 目录存在且可写，Nginx 配置了 `/upload/images/` 反代 |
| 无法发布文章 | 检查后端 ADMIN_TOKEN 已配置，并在 /admin 重新登录 |
| 页面刷新 404 | Nginx 配置了 `try_files $uri $uri/ /index.html;` |

---

## 八、文件对应关系

| 宝塔路径 | 来源 |
|---------|------|
| `/www/wwwroot/rain/blog-front/` | `blog-front/dist` |
| `/www/wwwroot/rain/blog-back/blog-server` | `go build` 产物 |
| `/www/wwwroot/rain/blog-back/.env` | 参考 `blog-back/.env.example` |
| `/www/wwwroot/rain/upload/` | 运行时上传目录 |
| MySQL `blog` 库 | `blog-back/scripts/schema.sql` |

---

**完成！** 部署后如有问题，优先检查：
1. Supervisor 日志
2. Nginx 配置是否重载
3. `.env` 环境变量是否正确
4. 后端 `ADMIN_TOKEN` 是否已配置，并已重新登录

更多细节见 [BAOTA_DEPLOY.md](./BAOTA_DEPLOY.md) 和 [DEPLOY.md](./DEPLOY.md)。
