<script setup lang="ts">
import { ElAlert, ElButton, ElCollapse, ElCollapseItem, ElForm, ElFormItem, ElIcon, ElInput, ElOption, ElSelect, ElUpload } from 'element-plus'
import 'element-plus/es/components/alert/style/css'
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/collapse/style/css'
import 'element-plus/es/components/collapse-item/style/css'
import 'element-plus/es/components/form/style/css'
import 'element-plus/es/components/form-item/style/css'
import 'element-plus/es/components/icon/style/css'
import 'element-plus/es/components/input/style/css'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/option/style/css'
import 'element-plus/es/components/select/style/css'
import 'element-plus/es/components/upload/style/css'
import { ref, onMounted, onUnmounted, watch, computed } from 'vue'
import Vditor from 'vditor'
import 'vditor/dist/index.css'
import { Plus, Close } from '@element-plus/icons-vue'
import { addArticle } from '../api/articles'
import { getTag, createTag } from '../api/tag'
import { uploadImage } from '../api/upload'
import { getUser, saveProfile, type Profile } from '../api/user'
import { IMAGE_ACCEPT, validateImageFile } from '../utils/imageUpload'
import type { UploadFile } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import 'element-plus/es/components/message-box/style/css'
import { type ApiResponse } from '../types/articles'
import { useRouter, onBeforeRouteLeave } from 'vue-router'
import { saveDraft, loadDraft, clearDraft, backupDraft, type ArticleDraft } from '../utils/draft'
import { resolveMediaUrl } from '../utils/media'
import { logoutAdmin, adminHeaders, requireAdminLogin } from '../utils/adminAuth'
import { formatVditorUploadResponse, VDITOR_UPLOAD_URL } from '../utils/vditorUpload'

const router = useRouter()

const handleExitAdmin = async () => {
    try {
        if (!(await confirmLeave())) return
        await logoutAdmin()
        leavingConfirmed = true
        router.push('/')
    } catch {
        ElMessage.error('退出失败，请检查网络后重试')
    }
}
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

const publishing = ref(false)
const savingProfile = ref(false)
const creatingTag = ref(false)
const profileLoaded = ref(false)
const profile = ref<Profile>({ id: 0, username: '', signature: '', avatar: '' })
const newTagName = ref('')
const expandedSettings = ref<string[]>([])

const draftReady = ref(false)
const draftSaving = ref(false)
const draftError = ref('')
const draftNotice = ref('')
const savedAt = ref(0)
const publishError = ref('')
const uploadedCover = ref('')
const hasDraft = computed(() => !!(title.value || content.value || tagIds.value.length || coverFile.value || uploadedCover.value))
let draftTimer: ReturnType<typeof setTimeout> | undefined
let draftVersion = 0
let leavingConfirmed = false
const snapshotDraft = (): ArticleDraft => ({
    title: title.value, content: content.value, tagIds: [...tagIds.value],
    cover: coverFile.value, coverUrl: uploadedCover.value, updatedAt: Date.now(),
})
const persistDraft = async () => {
    if (!draftReady.value) return
    clearTimeout(draftTimer)
    const version = ++draftVersion
    const snapshot = snapshotDraft()
    draftSaving.value = true
    try {
        await saveDraft(snapshot)
        if (version === draftVersion) { draftError.value = ''; savedAt.value = snapshot.updatedAt }
    } catch (error) {
        if (version === draftVersion) draftError.value = error instanceof Error ? error.message : '草稿保存失败，请勿关闭页面'
    } finally { if (version === draftVersion) draftSaving.value = false }
}
watch([title, content, tagIds, coverFile, uploadedCover], () => {
    if (!draftReady.value) return
    // Update text backup immediately; debounce larger IndexedDB writes.
    backupDraft(snapshotDraft())
    clearTimeout(draftTimer)
    draftTimer = setTimeout(persistDraft, 400)
}, { deep: true, flush: 'sync' })
const confirmLeave = async () => {
    if (publishing.value) { ElMessage.warning('正在发布，请等待完成'); return false }
    if (!hasDraft.value) return true
    await persistDraft()
    try {
        await ElMessageBox.confirm(draftError.value || '草稿已保存在此浏览器，下次进入可继续编辑。确定离开吗？', '离开编辑页', {
            confirmButtonText: '离开', cancelButtonText: '继续编辑', type: 'warning',
        })
        return true
    } catch { return false }
}
onBeforeRouteLeave(() => leavingConfirmed || confirmLeave())
const beforeUnload = (event: BeforeUnloadEvent) => {
    if (!hasDraft.value && !publishing.value) return
    backupDraft(snapshotDraft())
    event.preventDefault()
    event.returnValue = ''
}
const discardDraft = async () => {
    if (publishing.value) return
    try { await ElMessageBox.confirm('确定清空当前草稿吗？此操作无法恢复。', '清空草稿', { type: 'warning' }) } catch { return }
    draftReady.value = false
    clearTimeout(draftTimer)
    title.value = ''; content.value = ''; tagIds.value = []
    vditor?.setValue(''); removeCover()
    draftVersion++
    try { await clearDraft(); draftError.value = '' } catch { draftError.value = '草稿存储清理失败，请检查浏览器存储' }
    savedAt.value = 0; draftSaving.value = false; draftNotice.value = ''; publishError.value = ''
    draftReady.value = true
}

const loadProfile = async () => {
    profileLoaded.value = false
    try {
        const res = await getUser()
        profile.value = (res.data ?? [])[0] ?? { id: 0, username: '', signature: '', avatar: '' }
        if (!profile.value.id) expandedSettings.value = ['profile']
        profileLoaded.value = true
    } catch {
        ElMessage.error('加载博主资料失败，请点击重试')
    }
}
const handleSaveProfile = async () => {
    if (savingProfile.value || !profileLoaded.value || publishing.value) return
    if (!profile.value.username.trim()) { ElMessage.warning('请输入博主昵称'); return }
    savingProfile.value = true
    try {
        const res = await saveProfile({ ...profile.value, username: profile.value.username.trim() })
        profile.value = res.data
        ElMessage.success('博主资料已保存')
    } catch (e) {
        ElMessage.error(e instanceof Error ? e.message : '保存失败')
    } finally { savingProfile.value = false }
}
const handleCreateTag = async () => {
    if (creatingTag.value || publishing.value) return
    const name = newTagName.value.trim()
    if (!name) { ElMessage.warning('请输入标签名称'); return }
    const existing = tags.value.find(tag => tag.name === name)
    if (existing) {
        if (!tagIds.value.includes(existing.id)) tagIds.value.push(existing.id)
        newTagName.value = ''
        return
    }
    creatingTag.value = true
    try {
        const res = await createTag(name)
        tags.value.push(res.data)
        tagIds.value.push(res.data.id)
        newTagName.value = ''
        ElMessage.success('标签已创建并选中')
    } catch (e) {
        ElMessage.error(e instanceof Error ? e.message : '创建标签失败')
    } finally { creatingTag.value = false }
}
let disposed = false

// 初始化编辑器和获取标签列表
onMounted(async () => {
    window.addEventListener('beforeunload', beforeUnload)
    const restored = await loadDraft()
    if (disposed) return
    if (restored.draft) {
        title.value = restored.draft.title
        content.value = restored.draft.content
        tagIds.value = restored.draft.tagIds
        coverFile.value = restored.draft.cover
        uploadedCover.value = restored.draft.coverUrl
        imageUrl.value = coverFile.value ? URL.createObjectURL(coverFile.value) : resolveMediaUrl(uploadedCover.value)
        savedAt.value = restored.draft.updatedAt
        if (hasDraft.value) draftNotice.value = restored.missingCover ? '草稿已恢复，封面文件未保存成功，请重新选择封面。' : '已恢复上次未发布的草稿。'
    }
    draftReady.value = true
    const profileTask = loadProfile()
    // 获取标签列表
    try {
        const res = await getTag()
        tags.value = res.data || []
        const available = new Set(tags.value.map(tag => tag.id))
        const filtered = tagIds.value.filter(id => available.has(id))
        if (filtered.length !== tagIds.value.length) draftNotice.value = '草稿中的部分标签已删除，请重新选择标签。'
        tagIds.value = filtered
    } catch (error) {
        console.error('获取标签列表失败:', error)
        ElMessage.error('加载标签失败')
    }

    await profileTask
    if (disposed) return

    // 初始化Vditor（核心：确保input事件正确同步）
    if (vditorRef.value) {
        vditor = new Vditor(vditorRef.value, {
            height: 600,
            width: '100%',
            mode: 'ir',
            upload: {
                url: VDITOR_UPLOAD_URL,
                setHeaders: adminHeaders,
                fieldName: 'file',
                multiple: false,
                max: 10 * 1024 * 1024,
                accept: IMAGE_ACCEPT,
                format: formatVditorUploadResponse,
                error(message: string) {
                    try {
                        const result = JSON.parse(message)
                        if (result.code === 401) {
                            requireAdminLogin()
                            return
                        }
                        if (typeof result.msg === 'string') { ElMessage.error(result.msg); return }
                    } catch {}
                    ElMessage.error('正文图片上传失败')
                },
                filename(name: string) {
                    return `${Date.now()}-${name}`
                },
            },
            preview: {
                markdown: {
                    toc: true
                }
            },
            value: content.value,
            cache: {
                id: 'article-editor',
                enable: false
            },
            input: (value) => {
                content.value = value
            },
        })


    }
})

onUnmounted(() => {
    disposed = true
    window.removeEventListener('beforeunload', beforeUnload)
    clearTimeout(draftTimer)
    if (draftReady.value && hasDraft.value) void persistDraft()
    vditor?.destroy()
    if (imageUrl.value) {
        URL.revokeObjectURL(imageUrl.value)
    }

})

// 选择图片逻辑
const handleCoverChange = (uploadFile: UploadFile) => {
    const file = uploadFile.raw
    if (!file) return
    const error = validateImageFile(file)
    if (error) { ElMessage.error(error); return }
    if (imageUrl.value) URL.revokeObjectURL(imageUrl.value)
    uploadedCover.value = ''
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
    uploadedCover.value = ''
}

// 发布逻辑
const handlePublish = async () => {
    if (publishing.value || savingProfile.value || creatingTag.value) return
    if (!profileLoaded.value || !profile.value.id) {
        ElMessage.warning('请先保存博主资料')
        return
    }
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

    publishError.value = ''
    publishing.value = true
    vditor?.disabled()
    const selectedTags = [...tagIds.value]
    const authorId = profile.value.id
    try {
        let coverImage: string | undefined = uploadedCover.value || undefined
        if (coverFile.value && !coverImage) {
            const res = await uploadImage(coverFile.value)
            coverImage = res.data?.imageUrl
            if (!coverImage) throw new Error('封面上传失败：未返回图片地址')
            uploadedCover.value = coverImage
        }
        const res = await addArticle({
            title: latestTitle,
            content: latestContent.replace(/\r\n/g, '\n'),
            tagIds: selectedTags,
            userId: authorId,
            ...(coverImage ? { coverImage } : {}),
        }) as unknown as ApiResponse
        if (res.code !== 200) throw new Error(res.msg || res.message || '发布失败')
        ElMessage.success('文章发布成功')
        draftReady.value = false
        clearTimeout(draftTimer)
        draftVersion++
        title.value = ''
        content.value = ''
        vditor?.setValue('')
        tagIds.value = []
        removeCover()
        vditor?.clearCache()
        try { await clearDraft(); draftError.value = '' } catch { draftError.value = '文章已发布，但本地草稿清理失败' }
        savedAt.value = 0; draftNotice.value = ''; draftSaving.value = false
        draftReady.value = true
    } catch (error) {
        publishError.value = error instanceof Error ? error.message : '发布失败，请重试'
        if (/超时|网络/.test(publishError.value)) publishError.value += '；发布结果可能尚未确认，请先查看首页，避免重复发布。'
        await persistDraft()
        ElMessage.error(publishError.value)
    } finally {
        publishing.value = false
        if (!disposed) vditor?.enable()
    }

}
</script>

<template>
    <div class="edit-page">
        <div class="admin-bar">
            <span class="admin-label">后台 · 写文章</span>
            <div class="admin-actions">
                <el-button text type="primary" @click="router.push('/')">查看站点</el-button>
                <el-button text type="danger" @click="handleExitAdmin">退出后台</el-button>
            </div>
        </div>
        <div class="draft-status" aria-live="polite">
            <span>{{ draftSaving ? '正在保存草稿…' : savedAt ? '草稿已自动保存在此浏览器' : '草稿将在输入后自动保存' }}</span>
            <el-button text :disabled="publishing || !hasDraft" @click="persistDraft">保存草稿</el-button>
            <el-button text type="danger" :disabled="publishing || !hasDraft" @click="discardDraft">清空草稿</el-button>
            <el-alert v-if="draftNotice" :title="draftNotice" type="info" @close="draftNotice = ''" />
            <el-alert v-if="draftError" :title="draftError" type="warning" :closable="false" />
            <el-alert v-if="publishError" :title="publishError" type="error" :closable="false" />
        </div>
        <!-- 左侧编辑器 -->
        <div class="editor-container">
            <div ref="vditorRef"></div>
        </div>
        <!-- 右侧标签+封面 -->
        <div class="tag-container">
            <el-collapse v-model="expandedSettings">
                <el-collapse-item title="博主资料（首次发布前填写）" name="profile">
                    <el-form label-position="top" :disabled="publishing || savingProfile || !profileLoaded">
                        <el-form-item label="昵称"><el-input v-model="profile.username" maxlength="64" /></el-form-item>
                        <el-form-item label="个性签名"><el-input v-model="profile.signature" maxlength="255" /></el-form-item>
                        <el-form-item label="头像地址（可选）"><el-input v-model="profile.avatar" maxlength="512" placeholder="https://…" /></el-form-item>
                        <el-button :loading="savingProfile" @click="handleSaveProfile">保存博主资料</el-button>
                    </el-form>
                    <el-button v-if="!profileLoaded" text @click="loadProfile">重新加载资料</el-button>
                </el-collapse-item>
            </el-collapse>
            <el-alert v-if="profileLoaded && !profile.id" title="请展开博主资料，填写昵称并保存后发布" type="info" :closable="false" />
            <div class="tag-item">
                <div class="title">标题</div>
                <el-input v-model="title" :disabled="publishing" maxlength="255" size="large" placeholder="请输入文章标题" clearable />
            </div>
            <div class="tag-item">
                <div class="title">分类标签</div>
                <el-select v-model="tagIds" :disabled="publishing" placeholder="请选择文章标签" multiple style="width: 100%" clearable>
                    <el-option v-for="item in tags" :key="item.id" :label="item.name" :value="item.id" />
                </el-select>
                <el-input v-model="newTagName" :disabled="publishing || creatingTag" maxlength="64" placeholder="新标签名称" @keyup.enter="handleCreateTag" />
                <el-button :loading="creatingTag" :disabled="publishing" @click="handleCreateTag">创建并选中标签</el-button>
            </div>
            <div class="tag-item">
                <div class="title">封面</div>
                <el-upload class="avatar-uploader" :action="''" :show-file-list="false" :on-change="handleCoverChange"
                    :auto-upload="false" :accept="IMAGE_ACCEPT" :disabled="publishing">
                    <div class="cover-preview-container" v-if="imageUrl">
                        <img :src="imageUrl" class="avatar" alt="封面预览" />
                        <el-icon class="cover-remove-btn" @click.stop="!publishing && removeCover()">
                            <Close />
                        </el-icon>
                    </div>
                    <el-icon v-else class="avatar-uploader-icon">
                        <Plus />
                    </el-icon>
                </el-upload>
                <div class="upload-hint">支持jpg、png、jpeg格式，最大10MB、1600万像素</div>
            </div>
            <el-button class="publish-btn" size="large" type="primary" :loading="publishing" :disabled="savingProfile || creatingTag || !profileLoaded || !profile.id" @click="handlePublish">发布</el-button>
        </div>
    </div>
</template>

<style scoped>
.draft-status { width: 100%; color: var(--text-color); }

.edit-page {
    width: 90%;
    min-height: 100vh;
    background-color: var(--bg-color);
    display: flex;
    flex-wrap: wrap;
    gap: 20px;
    padding: 20px;
    margin: 0 auto;
    box-sizing: border-box;
}

.admin-bar {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 8px 12px;
    background: var(--card-bg);
    border: 1px solid var(--border-color);
    border-radius: 8px;
}

.admin-label {
    font-weight: 600;
    color: var(--text-color);
}

.admin-actions {
    display: flex;
    gap: 8px;
}

.editor-container {
    width: 1000px;
}

.tag-container {
    flex: 1;
    min-height: 500px;
    display: flex;
    flex-direction: column;
    gap: 20px;
    padding: 10px;
    background-color: var(--card-bg);
    border: 1px solid var(--border-color);
    border-radius: 8px;
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