<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElCard, ElDivider } from 'element-plus'
import { getTag } from '../api/tag';
// 定义分类数据接口
interface Category {
    id: number
    name: string
    description: string
    count: number
}

// 模拟分类数据
const categories = ref<Category[]>([])
const activeCategory = ref<number | null>(null)

const route = useRoute()
const router = useRouter()

onMounted(async () => {
    const res = await getTag()
    categories.value = res.data

    const type = route.query.type
    if (type && typeof type === 'string') {
        const matched = categories.value.find(cat => cat.name === type)
        if (matched) {
            activeCategory.value = matched.id
        }
    }
})

// 处理分类点击
const handleCategoryClick = (id: number) => {
    activeCategory.value = id
    router.push(`/categoryDetail?typeId=${id}`)
}

</script>

<template>
    <div class="category-page">
        <!-- 页面标题 -->
        <div class="cat-header">
            <h1>博客分类</h1>
            <p>探索不同类别的技术文章和教程</p>
        </div>
        <ElDivider />
        <div class="category-content">
            <div class="all-categories">
                <ElCard v-for="category in categories" :key="category.id" class="category-card"
                    @click="handleCategoryClick(category.id)">
                    <div class="category-card-header">
                        <h3 class="category-card-title">{{ category.name }}</h3>
                        <span class="category-count">{{ category.count }} 篇文章</span>
                    </div>
                    <div class="category-card-desc">{{ category.description }}</div>
                    <ElDivider />
                    <div class="category-card-preview">
                        <div class="more-articles">
                            查看全部 {{ category.count }} 篇文章 →
                        </div>
                    </div>
                </ElCard>
            </div>
        </div>
    </div>
</template>

<style scoped>
.category-page {
    width: 100%;
    padding: 40px 20px;
    box-sizing: border-box;
    max-width: 1300px;
    margin: 0 auto;
}

.cat-header {
    text-align: center;
    margin-bottom: 40px;
}

.cat-header h1 {
    font-size: 2.5rem;
    margin-bottom: 15px;
    color: var(--text-color);
}

.cat-header p {
    font-size: 1.1rem;
    color: #666;
}

.category-content {
    margin-top: 30px;
}

.all-categories {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 25px;
}

.category-card {
    cursor: pointer;
    transition: all 0.3s ease;
    height: 100%;
}

.category-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.category-card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 15px;
}

.category-card-title {
    margin: 0;
    font-size: 1.3rem;
    flex-grow: 1;
}

.category-count {
    background-color: #f0f7ff;
    color: #4096ff;
    padding: 3px 10px;
    border-radius: 12px;
    font-size: 0.9rem;
}

.category-card-desc {
    height: 70px;
    color: #666;
    line-height: 1.5;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
}

.more-articles {
    color: #4096ff;
    margin-top: 10px;
    font-size: 0.9rem;
}

:root.dark .category-card-desc,
:root.dark .cat-header p {
    color: #aaa;
}
</style>