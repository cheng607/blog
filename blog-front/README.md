# blog-front

Vue 3 + TypeScript + Vite + Element Plus 个人博客前端。

## 开发

```bash
cp .env.example .env.development
npm install
npm run dev
```

访问 http://localhost:5173 。接口通过 `/api` 代理到后端（见 `vite.config.ts`）。

## 构建

```bash
cp .env.example .env.production
npm run build
```

产物在 `dist/`，由 Nginx 托管。生产环境需配置 `VITE_API_BASE_URL`（默认同域 `/api/`）。

第四阶段的分页接口变更、草稿行为与验证说明见 [STAGE4.md](../STAGE4.md)。本版本需与对应后端同步部署。

## 目录说明

| 路径 | 说明 |
|------|------|
| `src/views/` | 页面：首页、文章、分类、树洞、编辑、关于 |
| `src/components/` | 公共组件：导航、列表、侧栏等 |
| `src/api/` | 后端接口封装 |
| `src/utils/` | 请求、格式化、Markdown 摘要、媒体 URL |
| `src/config.ts` | API 基址配置 |
