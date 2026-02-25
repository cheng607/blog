<script setup lang="ts">
import { ref, onMounted, nextTick, watch, watchEffect } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Marked } from 'marked'
import { markedHighlight } from 'marked-highlight'
import { getArticle, addViewsCount } from '../api/articles'
import hljs from 'highlight.js'
import 'highlight.js/styles/github.css'
import Sidebar from '../components/Sidebar.vue'
import { ElMessage } from 'element-plus'
const route = useRoute()
const article = ref<any>(null)
const loading = ref(true)
const error = ref('')
const markdownHtml = ref('')
const toc = ref<Array<{ id: string; text: string; level: number }>>([])
// 检测当前是否为暗色模式
const isDarkMode = ref(document.documentElement.classList.contains('dark'))

// 初始化marked
const marked = new Marked()
marked.use(
    markedHighlight({
        emptyLangClass: 'hljs',
        langPrefix: 'hljs language-',
            highlight(code, lang) {
                const requested = (lang || '').toLowerCase();
                // 优先使用明确的语言（且 highlight.js 支持）
                if (requested && hljs.getLanguage(requested)) {
                    return hljs.highlight(code, { language: requested }).value;
                }
                // 否则尝试自动检测语言，尽量避免把 `plaintext` 字样显示在页面中
                const detected = hljs.highlightAuto(code);
                if (detected && detected.value) return detected.value;
                // 最后兜底为 plaintext 的高亮渲染（不会把 `plaintext` 文本插入页面）
                return hljs.highlight(code, { language: 'plaintext' }).value;
            }
    }),
    { breaks: true }
);

// 刷新代码高亮
const refreshCodeHighlight = () => {
    nextTick(() => {
        document.querySelectorAll('.article-body pre code.hljs').forEach((el) => {
            // 先移除hljs类，再重新高亮，避免样式残留
            el.classList.remove('hljs');
            hljs.highlightElement(el as HTMLElement);
        });
    });
}

// 新增：监听暗色模式切换
watchEffect(() => {
    isDarkMode.value = document.documentElement.classList.contains('dark')
    refreshCodeHighlight() // 模式变化时刷新代码高亮
})

// 核心：为代码块添加复制功能
const addCopyToCodeBlocks = () => {
    nextTick(() => {
        const codeBlocks = document.querySelectorAll('.article-body pre');

        codeBlocks.forEach((pre) => {
            // 避免重复添加复制按钮
            if (pre.querySelector('.code-copy-btn')) return;
            const copyBtn = document.createElement('button');
            copyBtn.className = 'code-copy-btn';
            copyBtn.innerHTML = `
                <!-- 引用SVG精灵图中的图标 -->
                <svg class="copy-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <use xlink:href="#icon-copy"></use>
                </svg>
                <svg class="check-icon" width="16" height="16" viewBox="0 0 24 24" fill="currentColor">
                    <use xlink:href="#icon-check"></use>
                </svg>
            `;
            copyBtn.title = '点击复制代码';

            // 3. 绑定复制事件
            copyBtn.addEventListener('click', () => {
                const codeElement = pre.querySelector('code');
                if (!codeElement) return;

                const codeContent = codeElement.textContent || '';
                // 复制到剪贴板
                navigator.clipboard.writeText(codeContent)
                    .then(() => {
                        copyBtn.classList.add('copied');
                        ElMessage.success('代码复制成功！');

                        // 2秒后恢复初始状态
                        setTimeout(() => {
                            copyBtn.classList.remove('copied');
                        }, 2000);
                    })
                    .catch((err) => {
                        console.error('复制失败:', err);
                        ElMessage.error('复制失败，请手动复制');
                    });
            });

            // 直接添加到pre标签
            pre.appendChild(copyBtn);
        });
    });
}

onMounted(async () => {
    window.scrollTo({
        top: 500,
        behavior: "smooth"
    })
    const id = route.query.id
    if (!id) {
        error.value = '文章ID不存在'
        loading.value = false
        return
    }
    await watchViewsCount(Number(id))
    try {
        const res = await getArticle(Number(id))
        article.value = res.data
        markdownHtml.value = await marked.parse(article.value.content)
        // 关键：解析完成后，等待DOM渲染再添加复制按钮
        await nextTick()
        addCopyToCodeBlocks()
        // 新增：初始化时刷新一次高亮
        refreshCodeHighlight()
            // 生成目录（TOC）并注入到页面标题
            // 选择渲染后的文章内的所有标题，生成唯一 id 并收集
            toc.value = []
            const headers = document.querySelectorAll('.article-body h1, .article-body h2, .article-body h3, .article-body h4, .article-body h5, .article-body h6')
            const usedIds = new Set<string>()
            headers.forEach((h) => {
                const el = h as HTMLElement
                let text = (el.textContent || '').trim()
                if (!text) return
                // 生成 slug
                let slug = text.toLowerCase().replace(/\s+/g, '-').replace(/[^a-z0-9\-]/g, '')
                if (!slug) slug = 'section'
                let uniq = slug
                let idx = 1
                while (usedIds.has(uniq)) {
                    uniq = `${slug}-${idx++}`
                }
                usedIds.add(uniq)
                el.id = uniq
                const level = Number(el.tagName.replace('H', '')) || 1
                toc.value.push({ id: uniq, text, level })
            })
    } catch (err) {
        error.value = '获取文章失败'
        console.error(err)
    } finally {
        loading.value = false
    }
    // 清理过期记录
    Object.keys(localStorage).forEach(key => {
        if (key.startsWith('article_viewed_')) {
            const item = JSON.parse(localStorage.getItem(key) || '{}');
            if (item.expireTime && Date.now() > item.expireTime) {
                localStorage.removeItem(key);
            }
        }
    });
})

// 监控用户停留时间并修改浏览量
async function watchViewsCount(articleId: number) {
    const storageKey = `article_viewed_${articleId}`;
    // 检查本地存储：如果已记录，则不重复调用接口
    if (localStorage.getItem(storageKey)) {
        return;
    }
    try {
        await addViewsCount(articleId)
        const expireTime = Date.now() + 24 * 60 * 60 * 1000;
        localStorage.setItem(storageKey, JSON.stringify({ viewed: true, expireTime }));
    } catch (error) {
        console.error(error)
    }
}

// 监听markdownHtml变化
watch(markdownHtml, () => {
    addCopyToCodeBlocks()
    refreshCodeHighlight()
})

const router = useRouter()
const handleTag = (tagId: number) => {
    router.push(`/categoryDetail?typeId=${tagId}`)
}
</script>

<template>
    <div class="article-page">
        <svg xmlns="http://www.w3.org/2000/svg" style="display: none;">
            <!-- 复制图标 -->
            <symbol id="icon-copy" viewBox="0 0 24 24">
                <path
                    d="M16 1H4c-1.1 0-2 .9-2 2v14h2V3h12V1zm3 4H8c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h11c1.1 0 2-.9 2-2V7c0-1.1-.9-2-2-2zm0 16H8V7h11v14z" />
            </symbol>
            <!-- 复制成功图标 -->
            <symbol id="icon-check" viewBox="0 0 24 24">
                <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
            </symbol>
        </svg>

        <div class="article-container" v-loading="loading" element-loading-text="加载中..."
            element-loading-background="rgba(255, 255, 255, 0.8)">
            <div v-if="error" class="error-message">{{ error }}</div>

            <!-- 文章内容 -->
            <div v-else-if="article" class="article-content">
                <h1 class="article-title">{{ article.title }}</h1>
                <div class="article-meta">
                    <div class="article-box">
                        <span class="publish-time">发布时间：{{ $formatTime(article.publishTime) }}</span>
                        <span class="views-count">浏览量：{{ article.viewsCount }}</span>
                    </div>
                    <div class="tags">
                        <span>标签：</span>
                        <el-tag v-for="tag in article.tags" :key="tag.id" type="primary" class="tag"
                            @click="handleTag(tag.id)">
                            {{ tag.name }}
                        </el-tag>
                    </div>
                </div>
                <el-divider />
                <!-- 文章封面 -->
                <div v-if="article.coverImage" class="article-cover">
                    <img :src="article.coverImage" alt="文章封面">
                </div>
                <div class="article-body-container">
                    <div class="article-body" v-html="markdownHtml"></div>
                </div>
            </div>
        </div>
        <div class="right">
            <Sidebar :toc="toc" />
        </div>
    </div>
</template>

<style scoped>
.article-page {
    width: 100%;
    padding: 20px 0;
    display: flex;
    justify-content: center;
}

.article-container {
    width: 1000px;
    padding: 0 20px;
}

.right {
    width: 260px;
    margin-left: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}

.article-content {
    background-color: #fff;
    padding: 20px;
}

:root.dark .article-content {
    background-color: #1A1A1A;
    border-radius: 20px;
    border: 1px solid #333;
}

.article-title {
    font-size: 28px;
    text-align: center;
    margin: 20px 0;
}

.article-meta {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin: 20px 0;
    padding-bottom: 10px;
    border-bottom: 1px solid #eee;
}

.article-box {
    display: flex;
    align-items: center;
    gap: 20px;
}

.publish-time,
.views-count,
.tags span {
    color: #666;
    font-size: 14px;
}

.tags {
    display: flex;
    gap: 10px;
}

.tag {
    cursor: pointer;
}

.article-cover {
    text-align: center;
    margin: 30px 0;
}

.article-cover img {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
}

.article-body-container {
    padding: 0 10px;
}

.article-body {
    min-height: 300px;
    line-height: 1.8;
}

.article-body h1,
.article-body h2,
.article-body h3 {
    margin: 1.5em 0 1em;
    font-weight: 600;
}

.article-body p {
    margin: 1em 0;
}

.article-body strong {
    font-weight: 600;
}

.article-body ul {
    margin: 1em 0;
    padding-left: 2em;
}

.article-body li {
    margin: 0.5em 0;
}

::v-deep .article-body pre {
    border-radius: 6px;
    margin: 1.5em 0;
    position: relative;
    overflow-x: auto;
    padding: 1em;
    border: 1px solid #f5f5f5;
    font-size: large;
}
::v-deep .article-body code{
    font-size: large;
}
::v-deep .article-body pre code.hljs {
    font-size: 14px !important;
}

/* 暗模式代码高亮背景*/
:root.dark ::v-deep .article-body pre code.hljs {
    background: #333 !important;
    color: #e0e0e0 !important;
}

::v-deep .article-body img {
    width: 100%;
}

/* 代码块复制按钮样式 */
::v-deep .code-copy-btn {
    position: absolute;
    top: 20px;
    right: 20px;
    z-index: 10;
    background-color: #fff;
    border: 1px solid #eee;
    border-radius: 4px;
    padding: 6px;
    cursor: pointer;
    transition: all 0.2s ease;
    display: inline-flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    width: 32px;
    height: 32px;
    display: none;
}

/* 鼠标悬浮代码块时显示按钮（原有逻辑不变） */
::v-deep .article-body pre:hover .code-copy-btn {
    display: flex;
}

::v-deep .code-copy-btn svg {
    transition: all 0.2s ease;
}

::v-deep .code-copy-btn:hover {
    background-color: #f8f9fa;
    border-color: #ddd;
}

::v-deep .code-copy-btn .copy-icon {
    display: block;
    color: #666;
}

::v-deep .code-copy-btn .check-icon {
    display: none;
    color: #fff;
}

::v-deep .code-copy-btn.copied {
    background-color: #67c23a;
    border-color: #67c23a;
}

::v-deep .code-copy-btn.copied .copy-icon {
    display: none;
}

::v-deep .code-copy-btn.copied .check-icon {
    display: block;
}

::v-deep .article-body table {
    width: 100%;
    border-collapse: collapse;
    margin: 1.5em 0;
    font-size: 14px;
}

::v-deep .article-body table-container {
    overflow-x: auto;
}

/* 表头样式（原有逻辑不变） */
::v-deep .article-body th {
    font-weight: 600;
    text-align: left;
    padding: 10px 12px;
    border: 2px solid #ddd;
}

/* 单元格样式（原有逻辑不变） */
::v-deep .article-body td {
    padding: 10px 12px;
    border: 2px solid #ddd;
}

::v-deep:root.dark .article-body td {
    color: #fff;
}

/* 暗色模式适配优化（原有逻辑不变） */
:root.dark ::v-deep .article-content {
    background-color: #1A1A1A;
    border-color: #333;
}

:root.dark ::v-deep .article-body pre {
    background-color: #2d2d2d !important;
    border-color: #444 !important;
}

:root.dark ::v-deep .article-body code {
    color: #e0e0e0 !important;
}

:root.dark ::v-deep .code-copy-btn {
    background-color: #3a3a3a;
    border-color: #555;
}

:root.dark ::v-deep .code-copy-btn:hover {
    background-color: #444;
    border-color: #666;
}

:root.dark ::v-deep .code-copy-btn .copy-icon {
    color: #ccc;
}

:root.dark ::v-deep .code-copy-btn.copied {
    background-color: #67c23a;
    border-color: #67c23a;
}

.error-message {
    text-align: center;
    padding: 50px;
    color: #f56c6c;
}

:root.dark .article-meta {
    border-bottom-color: #333;
}

:root.dark .publish-time,
:root.dark .tags span {
    color: #aaa;
}

::v-deep .article-body .hljs-keyword {
    color: #007acc;
}

::v-deep .article-body .hljs-string {
    color: #ce9178;
}

::v-deep .article-body .hljs-comment {
    color: #6a9955;
}

::v-deep .article-body .hljs-function {
    color: #dcdcaa;
}

::v-deep .article-body .hljs-number {
    color: #b5cea8;
}

/* 暗色模式下代码高亮颜色适配（原有逻辑不变） */
:root.dark ::v-deep .article-body .hljs-keyword {
    color: #569cd6;
}

:root.dark ::v-deep .article-body .hljs-string {
    color: #ce9178;
}

:root.dark ::v-deep .article-body .hljs-comment {
    color: #6a9955;
}

:root.dark ::v-deep .article-body .hljs-function {
    color: #dcdcaa;
}

:root.dark ::v-deep .article-body .hljs-number {
    color: #b5cea8;
}
</style>