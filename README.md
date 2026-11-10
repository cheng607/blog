# 下雨天个人博客

前后端分离的个人博客：Vue 3 + Element Plus + Vite，Go (Gin) + MySQL。

## 目录结构

```
Blog/
├── blog-front/     # 前端 SPA
├── blog-back/      # Go API 服务
├── deploy/         # Nginx 等部署示例
├── DEPLOY.md       # 正式部署说明
└── Makefile        # 常用命令
```

## 本地开发

### 环境

- Node.js 18+
- Go 1.21+
- MySQL 5.7 / 8.0（需先导入 `blog-back/scripts/schema.sql`）

### 后端

```bash
cd blog-back
cp .env.example .env
# 编辑 DB_DSN、UPLOAD_PATH
go run ./cmd/server
```

默认 `http://localhost:9090`，健康检查 `GET /health`。

### 前端

```bash
cd blog-front
cp .env.example .env.development
npm install
npm run dev
```

默认 `http://localhost:5173`，通过 Vite 代理访问 `/api`。

### 一键命令（项目根目录）

```bash
make dev-back    # 仅后端
make dev-front   # 仅前端
make build       # 构建前后端
```

## 生产部署

- 首次上线：[BAOTA_DEPLOY.md](./BAOTA_DEPLOY.md)（宝塔逐步操作）
- 通用说明：[DEPLOY.md](./DEPLOY.md)
- **改代码后如何更新**：[BAOTA_DEPLOY.md § 八、代码修改后如何更新部署](./BAOTA_DEPLOY.md#八代码修改后如何更新部署重要)

## 功能模块

| 模块 | 说明 |
|------|------|
| 首页 | 文章列表、分页、封面与摘要 |
| 文章详情 | Markdown 渲染、目录、代码高亮与复制、浏览量 |
| 分类 | 标签统计与按标签筛选 |
| 写文章 | Vditor 编辑器、封面上传、标签 |
| 树洞 | 短句发布与列表 |
| 关于 | 个人介绍与时间轴 |


## 首次发布

1. 配置后端 `ADMIN_TOKEN` 并启动服务，在 `/admin` 登录。
2. 在写文章页面的“博主资料”中填写昵称并保存，签名和头像地址可选。已有资料会自动载入，可在此修改。
3. 输入标签名称，点击“创建并选中标签”；已有标签可直接选择。
4. 填写标题和正文，按需上传封面，然后发布。文章使用保存后的博主 ID，不要求数据库中存在 ID 为 1 的用户。

封面和正文插图统一支持 JPG/JPEG/PNG，单张最多 10MB、1600万像素。后端会解码校验并重新编码图片，拒绝伪装格式或损坏文件。发布失败后保留当前内容，可修正问题后重试。
