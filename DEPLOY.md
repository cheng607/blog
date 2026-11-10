# 个人博客正式部署指南

项目为 **Vue3 前端**（`blog-front`）+ **Go 后端**（`blog-back`，Gin）。推荐在同一域名下用 Nginx 反代 API，避免跨域问题。

## 一、环境要求

| 组件 | 版本建议 |
|------|----------|
| Node.js | 18+ |
| Java | 8+ |
| MySQL | 5.7 / 8.0 |
| Nginx | 任意稳定版 |

## 二、数据库

1. 创建数据库 `blog`，导入表结构（article、tag、article_tag、hollow、user 等）。
2. 生产环境通过环境变量配置，勿在仓库中提交真实密码：

```bash
export DB_URL="jdbc:mysql://127.0.0.1:3306/blog?useSSL=false&serverTimezone=Asia/Shanghai&characterEncoding=utf8"
export DB_USER=blog
export DB_PASSWORD=你的密码
export UPLOAD_PATH=/www/wwwroot/rain/upload/
```

## 三、后端部署（Go）

```bash
cd blog-back
cp .env.example .env
# 编辑 .env：DB_DSN、UPLOAD_PATH、CORS_ORIGINS

make build
./bin/blog-server
```

首次部署请执行 `scripts/schema.sql` 初始化数据库。

或使用宝塔 / systemd 守护进程，确保：

- 端口 `9090` 可访问（或仅本机监听由 Nginx 反代）
- 上传目录 `UPLOAD_PATH` 存在且可写，例如 `/www/wwwroot/rain/upload/images/`
- `CORS_ORIGINS` 包含你的站点域名（如 `https://www.raincode.top`）

## 四、前端构建

```bash
cd blog-front
npm install
npm run build
```

产物在 `blog-front/dist/`。

- 默认 `.env.production` 使用 `VITE_API_BASE_URL=/api/`，需配合 Nginx 反代。
- 若前后端不同域且无反代，在 `.env.production` 中设置 `VITE_API_BASE_URL=http://你的服务器:9090/` 后重新构建。

## 五、Nginx 示例（推荐）

见 `deploy/nginx.conf.example`。核心配置：

- `/` → 静态站点 `dist`
- `/api/` → 反代到 `http://127.0.0.1:9090/`
- `/upload/images/` → 反代到后端（封面与正文图片）

站点使用 **History 路由**，Nginx 需配置 `try_files $uri $uri/ /index.html;`（见 `deploy/nginx.conf.example`）。

## 六、本地开发

```bash
# 终端 1：Go 后端（需 MySQL 可连）
cd blog-back && go run ./cmd/server

# 终端 2：前端
cd blog-front && npm run dev
```

访问 `http://localhost:5173`，API 经 `/api` 代理到后端。

## 七、功能自检清单

- [ ] 首页文章列表、封面、摘要正常
- [ ] 文章详情 Markdown、代码高亮、目录、浏览量 +1
- [ ] 分类页 / 分类详情筛选
- [ ] 写文章：标题、标签、正文、封面上传、发布
- [ ] 树洞列表与发布
- [ ] 关于页、页脚备案链接
- [ ] 暗色模式切换持久化

## 八、代码修改后如何更新

已上线后的增量发布（只改前端 / 只改后端 / 全量更新），见 **[BAOTA_DEPLOY.md 第八节](./BAOTA_DEPLOY.md#八代码修改后如何更新部署重要)**。

简要流程：

- **前端**：本地 `npm run build` → 覆盖服务器 `blog-front` 目录（dist 内容）→ 浏览器强刷
- **后端**：`GOOS=linux go build` → 覆盖 `blog-server` → Supervisor 重启

## 九、安全建议

1. 不要将 `.env` 中的数据库密码提交到公开仓库。
2. `/edit` 写文章入口目前无登录校验，生产环境建议在 Nginx 层限制 `/edit` 或增加鉴权。
3. 定期备份 MySQL 与 `upload` 目录。


## 管理登录升级与密钥轮换

- 管理密钥只配置在后端 `.env` 的 `ADMIN_TOKEN`。删除所有前端环境文件中的旧编辑密钥变量，重新构建前端。
- 将后端 `ADMIN_TOKEN` 更换为全新的高强度随机值（例如使用 `openssl rand -hex 32` 生成），不要复用曾进入前端构建的旧值。
- 同时部署新的前后端，重启后端；清理服务器/CDN 上旧版前端资源缓存。在 HTTPS 站点的 `/admin` 输入新密钥登录。
- 后端未配置密钥时，公开阅读仍可使用，但登录和管理写操作关闭。
- 会话有效期 8 小时，保存在当前标签页的 sessionStorage；退出会在服务器撤销会话。后端重启会使所有会话失效，目前适用于单进程部署。
- 会话过期时页面会弹窗重新登录，当前文章及所选封面保留在页面中；登录后手动重试操作，不要刷新页面。正文上传、封面上传、文章和树洞发布都使用会话认证。
- 密钥不再通过 URL、X-Admin-Token 或直接 Bearer 密钥传递。跨域部署必须将前端来源加入 CORS_ORIGINS。
- 验收：旧密钥登录失败；无会话写入返回 401；新登录可发布/上传；退出后原会话访问 `/auth/session` 返回 401。


## AI 限流、额度和超时

后端 `.env` 可设置 `AI_RATE_PER_MINUTE=6`（每 IP 滚动一分钟最多6次）、`AI_MAX_CONCURRENT=3`（全站3个进行中的请求）、`AI_DAILY_LIMIT=500`（全站每日上游调用次数）、`AI_MAX_TOKENS=1024`（单次输出上限）。未设置、非数字或非正值使用安全默认值；允许上限依次为60、20、100000、4096。额度是调用次数与单次输出限制，不是人民币预算；具体费用取决于模型和输入输出。

- 每个问题最多尝试两个不同模型；网络错误、超时、HTTP 408/429/5xx、无效或空响应才尝试备用模型。认证/参数等确定性错误不重复尝试。
- 每一次上游尝试（包括失败和备用模型）都先扣除一次额度。无效输入、限流及并发拒绝不调用模型。两条聊天路由共用限制，返回429时携带 `Retry-After`。
- 每日额度文件为 `UPLOAD_PATH/.ai-quota.json`，不在公开图片目录中。必须保留该文件及上传目录的持久化存储，部署不要删除；按北京时间00:00重置，重启不会重置已使用次数。文件损坏或无法写入时AI接口停止调用模型并返回503。
- 当前额度、并发与会话机制面向单进程部署。不要让多个后端实例共用此额度文件；多实例部署需换为数据库/Redis原子计数。
- `TRUSTED_PROXIES=127.0.0.1,::1` 只信任本机 Nginx；有其他代理时填写准确IP/CIDR。不要信任所有公网地址，否则来源IP限制可被伪造。
- 单模型30秒、后端总预算65秒、Go写响应80秒、Nginx读取85秒、前端90秒。更新现有 Nginx `/api/` 配置为 `proxy_connect_timeout 5s; proxy_send_timeout 20s; proxy_read_timeout 85s;`，检查配置后重载。客户端断开时取消上游请求。

参数错误/不存在/冲突/内容过大/限流分别返回HTTP 400/404/409/413/429，内部错误500、上游失败502、不可用503、超时504。响应保留 `{code,msg,data}` 结构，前端展示安全的 `msg`；详细数据库和上游错误仅用于服务端诊断。

正文最多1MB（UTF-8字节），文章JSON请求最多2MB，标题255字、标签1至20个且ID不重复；标签名称64字、描述255字；昵称64字、签名255字；树洞1000字；头像/封面只接受站内绝对路径或HTTP(S)地址（512字）。旧客户端须兼容真实HTTP错误状态码。
