# 第四阶段：性能与编辑体验

## 已完成

- 文章列表、按分类/作者筛选统一为服务端分页，列表只返回摘要；先分页再关联标签，避免多标签文章重复或被拆到不同页。
- `/site/summary` 一次返回公开博主资料与文章、分类、树洞计数，侧栏不再下载完整列表计算数量。
- 页面路由懒加载，编辑器和 Markdown 高亮随对应页面加载；Element Plus 组件及样式按需引入。
- 草稿自动保存/手动保存、重新进入恢复、离开提醒、清空确认、发布失败提示；文本使用 localStorage 备份，封面文件使用 IndexedDB。浏览器存储不可用时显示明确提示，草稿不跨设备同步。
- 分类、树洞、文章列表和详情补充错误/空状态及重试，未知路由显示 404；清理旧 `::v-deep` 写法。

## 接口兼容性

以下接口的 `data` 从数组变为 `{ items, total, page, pageSize }`，必须同时更新前后端：

- `GET /article/list`
- `GET /article/tag/:tagId`
- `GET /article/user/:userId`

查询参数：`page` 默认 1，范围 1–1000000；`pageSize` 默认 5，范围 1–50。超过末页会回到末页，空列表返回 `items: []`、`total: 0`、`page: 1`。摘要项使用 `excerpt`，不返回 `content`；正文仍从 `GET /article/:id` 获取。

`GET /site/summary` 的 `data` 为 `{ profile, articleCount, tagCount, hollowCount }`，未设置博主时 `profile` 为 `null`。上述路径均为后端路径，前端默认经 `/api` 代理。

## 验证与体积

在 `blog-back` 执行 `go test ./... -timeout 30s`；在 `blog-front` 执行 `npm test`、`npm run build`。

构建后运行 `node tests/bundle-report.mjs`，按 manifest 递归统计首页入口和全部静态 JS 依赖、去除重复模块。当前合计 378214 字节，逐文件 gzip 合计 144894 字节；不包含 CSS、图片、外部媒体，不能等同于首屏总流量或真实加载时间。编辑器和文章高亮不在此依赖集合内。文章详情的全语言高亮包仍触发大 chunk 提示，可后续根据实际语言范围进一步缩减。

本地浏览器验证使用 `node tests/preview-server.mjs`（仅监听 `127.0.0.1:4178`，使用模拟数据和模拟登录，不连接真实后端，不用于部署）：

- 下一页显示对应的 5 篇文章。
- 离开编辑器弹出确认，重新进入恢复标题。
- 发布返回模拟 503 时显示错误，保留正文、标题、分类，并可继续保存草稿。

后端测试使用 SQL 驱动桩验证分页参数、摘要、空数据和计数；尚未连接真实 MySQL 做集成验证，也未部署线上。
