<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { Plus, Close } from '@element-plus/icons-vue'
import { addArticle } from '../api/articles'
import { getTag } from '../api/tag'
import type { UploadFile } from 'element-plus'
import { ElMessage } from 'element-plus'
import { type ApiResponse } from '../types/articles'
// 定义标签项接口
interface TagItem {
    name: string
    id: number
}

// 核心变量
const content = ref('')
const vditorRef = ref<HTMLDivElement>()
let vditor: Vditor | null = null
const tags = ref<TagItem[]>([])
const tagIds = ref<number[]>([])
const title = ref<string>('')
const imageUrl = ref('')
const coverFile = ref<File | null>(null)

// 1. 提前定义粘贴处理函数（方便销毁）
let handlePaste: (event: ClipboardEvent) => void

// 2. Markdown 预处理（最小化）：仅对 HTML 实体与行尾做规范，不修改代码围栏或语言标识
const preprocessMarkdown = (mdContent: string): string => {
    if (!mdContent) return ''
    // 仅解码常见实体，规范化 CRLF 为 LF，保留所有原始标记（避免修改代码围栏或语言）
    return mdContent
        .replace(/\r\n/g, '\n')
        .replace(/&grave;|&#96;/g, '`')
        .replace(/&lt;/g, '<')
        .replace(/&gt;/g, '>')
}

// 初始化编辑器和获取标签列表
onMounted(async () => {
    // 获取标签列表
    try {
        const res = await getTag()
        tags.value = res.data || []
    } catch (error) {
        console.error('获取标签列表失败:', error)
        ElMessage.error('加载标签失败')
    }

    // 3. 自定义粘贴处理：强制保留代码块语言标识
    handlePaste = (event: ClipboardEvent) => {
        const clipboardData = event.clipboardData
        if (!clipboardData || !vditor) return

        const pasteText = clipboardData.getData('text/plain')
        if (!pasteText) return

        event.preventDefault()
        event.stopPropagation()

        const processedText = pasteText
            .replace(/\r\n/g, '\n')
            .replace(/&grave;|&#96;/g, '`')
            .replace(/&lt;/g, '<')
            .replace(/&gt;/g, '>')

        const currentValue = (vditor.getValue && vditor.getValue()) || ''
        // 保证与现有内容之间至少有一个换行
        const sep = currentValue.endsWith('\n') || currentValue === '' ? '' : '\n'
        vditor.setValue(currentValue + sep + processedText)

        // 手动同步到content变量（兜底）
        content.value = vditor.getValue() || ''
    }

    // 初始化Vditor（核心：确保input事件正确同步）
    if (vditorRef.value) {
        vditor = new Vditor(vditorRef.value, {
            height: 600,
            width: '100%',
            mode: 'ir',
            upload: {
                url: '/api/upload/image',
                filename(name) {
                    return `${Date.now()}-${name}`
                }
            },
            preview: {
                markdown: {
                    toc: true
                }
            },
            value: content.value,
            cache: {
                id: 'article-editor',
                enable: true
            },
            input: (value) => {
                content.value = value
            },
        })

        // 延迟绑定粘贴事件
        setTimeout(() => {
            if (vditorRef.value) {
                vditorRef.value.addEventListener('paste', handlePaste, true)
            }
        }, 0)
    }
})

onUnmounted(() => {
    vditor?.destroy()
    if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value)
    }
    if (vditorRef.value && handlePaste) {
        vditorRef.value.removeEventListener('paste', handlePaste, true)
    }
})

// 选择图片逻辑
const handleCoverChange = (uploadFile: UploadFile) => {
    const file = uploadFile.raw
    if (!file) return
    if (!file.type.startsWith('image/')) {
        ElMessage.error('请上传图片类型文件')
        imageUrl.value = ''
        coverFile.value = null
        return
    }
    if (file.size > 10 * 1024 * 1024) {
        ElMessage.error('图片大小不能超过10MB')
        imageUrl.value = ''
        coverFile.value = null
        return
    }
    imageUrl.value = URL.createObjectURL(file)
    coverFile.value = file
}

// 移除封面逻辑
const removeCover = () => {
    if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value)
    }
    imageUrl.value = ''
    coverFile.value = null
}

// 发布逻辑
const handlePublish = async () => {
    // 1. 优先从Vditor实例获取内容
    const latestContent = vditor?.getValue() || ''
    const latestTitle = title.value.trim()

    // 2. 表单验证
    if (!latestTitle) {
        ElMessage.warning('请输入文章标题')
        return
    }
    if (!latestContent.trim()) {
        ElMessage.warning('请输入文章内容')
        return
    }
    if (tagIds.value.length === 0) {
        ElMessage.warning('请至少选择一个标签')
        return
    }

    // 3. 预处理内容（仅最小化处理），优先保留编辑器原始 Markdown
    const processedContent = preprocessMarkdown(latestContent)

    // 4. 构建提交数据
    const articleData = {
        title: latestTitle,
        content: processedContent,
        tagIds: tagIds.value,
        userId: 1
    }

    try {
        const res = await addArticle(articleData) as unknown as ApiResponse
        console.log('接口返回:', res)
        if (res.code === 200) {
            ElMessage.success('文章发布成功')
            // 重置所有状态
            title.value = ''
            content.value = ''
            vditor?.setValue('') // 清空编辑器
            tagIds.value = []
            imageUrl.value = ''
            coverFile.value = null
            vditor?.clearCache()
        } else {
            ElMessage.error('发布失败: ' + (res.message || res.data.message || '未知错误'))
        }
    } catch (error) {
        console.error('发布文章错误:', error)
        ElMessage.error('发布失败，请重试')
    }
}
</script>

<template>
    <div class="edit-page">
        <!-- 左侧编辑器 -->
        <div class="editor-container">
            <div ref="vditorRef"></div>
        </div>
        <!-- 右侧标签+封面 -->
        <div class="tag-container">
            <div class="tag-item">
                <div class="title">标题</div>
                <el-input v-model="title" size="large" placeholder="请输入文章标题" clearable />
            </div>
            <div class="tag-item">
                <div class="title">分类标签</div>
                <el-select v-model="tagIds" placeholder="请选择文章标签" multiple style="width: 100%" clearable>
                    <el-option v-for="item in tags" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
            </div>
            <div class="tag-item">
                <div class="title">封面</div>
                <el-upload class="avatar-uploader" :action="''" :show-file-list="false" :on-change="handleCoverChange"
                    :auto-upload="false">
                    <div class="cover-preview-container" v-if="imageUrl">
                        <img :src="imageUrl" class="avatar" alt="封面预览" />
                        <el-icon class="cover-remove-btn" @click.stop="removeCover">
                            <Close />
                        </el-icon>
                    </div>
                    <el-icon v-else class="avatar-uploader-icon">
                        <Plus />
                    </el-icon>
                </el-upload>
                <div class="upload-hint">支持jpg、png、jpeg格式，最大10MB</div>
            </div>
            <el-button class="publish-btn" size="large" type="primary" @click="handlePublish">发布</el-button>
        </div>
    </div>
</template>

<style scoped>
.edit-page {
    width: 90%;
    min-height: 100vh;
    background-color: #f9f9f9;
    display: flex;
    gap: 20px;
    padding: 20px;
    margin: 0 auto;
}

:root.dark .edit-page {
    background-color: #1A1A1A;
}

.editor-container {
    width: 1000px;
}

.tag-container {
    flex: 1;
    height: 500px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 10px;
    background-color: #fff;
}

:root.dark .tag-container {
    background-color: #333;
}

.tag-item {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.tag-item .title {
    font-weight: 500;
}

.tag-item .title::before {
    content: '*';
    color: #C26968;
}

.tag-hint,
.upload-hint {
    font-size: 12px;
    color: #8c939d;
    margin-top: -5px;
}

.cover-preview-container {
    position: relative;
}

.avatar-uploader .avatar {
    width: 178px;
    height: 178px;
    display: block;
    object-fit: cover;
    border-radius: 6px;
}

.cover-remove-btn {
    position: absolute;
    top: 8px;
    right: 8px;
    width: 24px;
    height: 24px;
    line-height: 24px;
    text-align: center;
    background-color: rgba(0, 0, 0, 0.5);
    color: white;
    border-radius: 50%;
    cursor: pointer;
    transition: background-color 0.2s;
}

.cover-remove-btn:hover {
    background-color: rgba(0, 0, 0, 0.7);
}

.avatar-uploader .el-upload {
    border: 1px dashed var(--el-border-color);
    border-radius: 6px;
    cursor: pointer;
    position: relative;
    overflow: hidden;
    transition: var(--el-transition-duration-fast);
}

.avatar-uploader .el-upload:hover {
    border-color: var(--el-color-primary);
}

.el-icon.avatar-uploader-icon {
    font-size: 28px;
    color: #8c939d;
    background-color: #f5f5f5;
    width: 178px;
    height: 178px;
    text-align: center;
    line-height: 178px;
}

.publish-btn {
    width: 80%;
    margin: 10px auto;
}
</style>