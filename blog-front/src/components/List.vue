<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router';
import loadingImg from '../assets/loading.gif'
import { getList } from '../api/articles'
const currentPage = ref<number>(1)
const pageSize = ref<number>(5)
const isLoading = ref<boolean>(false)
interface tagItem {
    id: number,
    name: string
}
interface listItem {
    id: number,
    coverImage: string,
    title: string,
    content: string,
    tags: tagItem[],
    publishTime: string
}
const listData = ref<listItem[]>([])
const router = useRouter()
const route = useRoute()
const fetchList = async () => {
    const typeId = route.query.typeId
    try {
        const res = await getList()
        if (typeId) {
            listData.value = res.data
            listData.value = listData.value.filter((item) => {
                return item.tags.some(tag => tag.id === Number(typeId));
            })
        } else {
            listData.value = res.data
        }
    } catch (err) {
        console.error('获取文章列表失败', err);
        listData.value = [];
    } finally {
        isLoading.value = false;
    }
}

onMounted(() => {

    fetchList()
})

const displayData = computed(() => {
    const startIndex = (currentPage.value - 1) * pageSize.value
    const endIndex = startIndex + pageSize.value
    return listData.value.slice(startIndex, endIndex)
})
// 处理页码变化
const handleChange = (val: number) => {
    currentPage.value = val
}

// 跳转文章详情页
const handleDetail = (id: number) => {
    router.push(`/article?id=${id}`)
}
</script>

<template>
    <div class="list-container">
        <div class="item" v-for="item in displayData" :key="item.id" @click="handleDetail(item.id)">
            <!-- <img :src="'http://47.114.122.96:9090' + item.coverImage || loadingImg" class="item-img"> -->
            <img :src="loadingImg" class="item-img">
            <div class="item-box">
                <div class="item-title">{{ item.title }}</div>
                <div class="item-tag">
                    <div v-for="tag in item.tags" :key="tag.id">
                        <span>{{ tag.name }}</span>
                    </div>
                </div>
                <div class="item-content">{{ item.content }}</div>
                <div class="item-time">发布时间：{{ $formatTime(item.publishTime) }}</div>
            </div>
        </div>
        <el-pagination class="controller" background layout="prev, pager, next" :total="listData.length"
            :current-page="currentPage" :page-size="pageSize" @current-change="handleChange" />
    </div>
</template>

<style scoped>
.list-container {
    width: 100%;
    height: auto;
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
}

.item:hover {
    scale: 1.01;
}

.item-img {
    width: 220px;
    height: 220px;
    padding: 15px;
    flex-shrink: 0;
}

.item-box {
    width: 100%;
    /* width: 620px; */
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
    padding-right: 20px;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    line-height: 1.5;
    width: 100%;
    word-break: break-all;
    padding-right: 0;
    box-sizing: border-box;
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
</style>
