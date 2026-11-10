<script setup lang="ts">
import { ElButton, ElLoading, ElPagination } from 'element-plus'
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/loading/style/css'
import 'element-plus/es/components/message/style/css'
import 'element-plus/es/components/pagination/style/css'
const vLoading = ElLoading.directive
import { ref, onUnmounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import loadingImg from '../assets/loading.gif'
import { getList, getArticlesByTag } from '../api/articles'
import { resolveMediaUrl } from '../utils/media'
import { ElMessage } from 'element-plus'

const currentPage = ref<number>(1)
const pageSize = ref<number>(5)
const isLoading = ref<boolean>(true)
interface tagItem {
    id: number
    name: string
}
interface listItem {
    id: number
    coverImage: string
    title: string
    excerpt: string
    tags: tagItem[]
    publishTime: string
}
const total = ref(0)
const error = ref('')
const listData = ref<listItem[]>([])
const router = useRouter()
const route = useRoute()

const getCoverSrc = (item: listItem) => {
    const url = resolveMediaUrl(item.coverImage)
    return url || loadingImg
}


let requestId = 0
onUnmounted(() => { requestId++ })
const fetchList = async () => {
    const current = ++requestId
    listData.value = []
    error.value = ''
    isLoading.value = true
    const typeId = route.query.typeId
    try {
        const res = typeId
            ? await getArticlesByTag(Number(typeId), currentPage.value, pageSize.value)
            : await getList(currentPage.value, pageSize.value)
        if (current !== requestId) return
        listData.value = res.data?.items ?? []
        total.value = res.data?.total ?? 0
        currentPage.value = res.data?.page ?? currentPage.value
    } catch (err) {
        if (current !== requestId) return
        console.error('获取文章列表失败', err)
        error.value = err instanceof Error ? err.message : '加载文章列表失败'
        ElMessage.error(error.value)
        listData.value = []
    } finally {
        if (current === requestId) isLoading.value = false
    }
}

watch(
    () => route.query.typeId,
    () => {
        currentPage.value = 1
        fetchList()
    },
    { immediate: true }
)

const handleChange = (val: number) => {
    if (val === currentPage.value) return
    currentPage.value = val
    fetchList()
}

const handleDetail = (id: number) => {
    router.push(`/article?id=${id}`)
}
</script>

<template>
    <div class="list-container" v-loading="isLoading">
        <div v-if="error" class="list-empty" role="alert">{{ error }} <el-button @click="fetchList">重试</el-button></div>
        <div v-else-if="!isLoading && listData.length === 0" class="list-empty">暂无文章</div>
        <div class="item" v-for="item in listData" :key="item.id" @click="handleDetail(item.id)">
            <img :src="getCoverSrc(item)" class="item-img" alt="封面" @error="($event.target as HTMLImageElement).src = loadingImg">
            <div class="item-box">
                <div class="item-title">{{ item.title }}</div>
                <div class="item-tag">
                    <div v-for="tag in item.tags" :key="tag.id">
                        <span>{{ tag.name }}</span>
                    </div>
                </div>
                <div class="item-content">{{ item.excerpt }}</div>
                <div class="item-time">发布时间：{{ $formatTime(item.publishTime) }}</div>
            </div>
        </div>
        <el-pagination
            v-if="!error && total > pageSize"
            class="controller"
            background
            layout="prev, pager, next"
            :total="total"
            :disabled="isLoading"
            :current-page="currentPage"
            :page-size="pageSize"
            @current-change="handleChange"
        />
    </div>
</template>

<style scoped>
.list-container {
    width: 100%;
    height: auto;
    min-height: 120px;
}

.list-empty {
    text-align: center;
    padding: 60px 20px;
    color: #999;
    font-size: 16px;
}

.item {
    max-width: 100%;
    width: 100%;
    height: 250px;
    display: flex;
    align-items: center;
    margin: 10px 0;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
    border-radius: 10px;
    background-color: #fff;
    box-sizing: border-box;
    overflow: hidden;
    cursor: pointer;
}

.item:hover {
    scale: 1.01;
}

.item-img {
    width: 220px;
    height: 220px;
    padding: 15px;
    flex-shrink: 0;
    object-fit: cover;
}

/* 移动端取消封面图，让列表更紧凑 */
@media (max-width: 768px) {
  .item-img {
    display: none;
  }

  .item {
    flex-direction: column !important;
    height: auto !important;
    min-height: auto;
    padding: 16px 18px !important;
    border-bottom: 1px solid #f0f0f0;
  }

  .item:last-child {
    border-bottom: none;
  }

  .item-box {
    padding: 0 !important;
  }

  .item-title {
    font-size: 17px !important;
    line-height: 1.4 !important;
    font-weight: 600;
  }

  .item-content {
    font-size: 14px !important;
    color: #666 !important;
    max-height: none !important;
    -webkit-line-clamp: 4 !important;
    line-height: 1.6 !important;
    margin: 8px 0;
  }

  .item-tag {
    margin: 8px 0 !important;
  }

  .item-tag span {
    font-size: 12px !important;
    padding: 2px 8px !important;
  }

  .item-time {
    font-size: 12px !important;
    color: #999 !important;
    padding: 4px 0 0 !important;
  }
}

.item-box {
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    flex-direction: column;
    box-sizing: border-box;
    padding: 0 10px;
}

.item-title {
    font-size: 24px;
}

.item-tag {
    display: flex;
    gap: 10px;
    margin: 15px 0;
}

.item-tag span {
    font-size: 14px;
    background-color: #49B1F5;
    color: #fff;
    padding: 5px 10px;
    border-radius: 20px;
}

.item-content {
    font-size: 20px;
    height: auto;
    max-height: 90px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    line-height: 1.5;
    width: 100%;
    word-break: break-all;
    box-sizing: border-box;
    color: #666;
}

.controller {
    justify-content: center;
    margin: 30px;
    margin-top: 50px;
}

.item-time {
    margin-left: auto;
    font-size: 14px;
    padding: 20px;
}

:root.dark .item {
    border: 1px solid #f5f5f5;
    background-color: #1A1A1A;
}

:root.dark .item-content {
    color: #aaa;
}

@media (max-width: 768px) {
    .item {
        flex-direction: column;
        height: auto;
        align-items: stretch;
    }

    .item-img {
        width: 100%;
        height: 180px;
        padding: 12px 12px 0;
        box-sizing: border-box;
    }

    .item-box {
        padding: 12px;
    }

    .item-title {
        font-size: 18px;
        line-height: 1.4;
    }

    .item-content {
        font-size: 15px;
        max-height: none;
        -webkit-line-clamp: 4;
    }

    .item-time {
        margin-left: 0;
        padding: 8px 0 0;
    }

    .controller {
        margin: 20px 0;
    }
}
</style>
