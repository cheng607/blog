<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router';
import { getTag } from '../api/tag';

interface TypeItem {
    name: string,
    count: number,
    id: number
}

const typeData = ref<TypeItem[]>([])
const router = useRouter()

// 计算属性：最多显示前4条数据
const displayTypeData = computed(() => {
    const top5Data = [...typeData.value]
    top5Data.sort(function (a, b) {
        return b.count - a.count
    })
    return top5Data.slice(0, 5)
})

onMounted(async () => {
    try {
        const res = await getTag()
        typeData.value = res.data
    } catch (error) {
        console.error('获取标签', error)
    }
})

const handleType = (typeId: number) => {
    router.push(`/categoryDetail?typeId=${typeId}`)
}

const handelMore = () => {
    router.push('/categories')
}
</script>

<template>
    <div class="type-container">
        <div class="type-title">
            <span>分类</span>
        </div>

        <div v-if="typeData.length > 0">
            <div class="type-list">
                <div class="type-item" v-for="item in displayTypeData" :key="item.id" @click="handleType(item.id)">
                    <div class="type-name">{{ item.name }}</div>
                    <div class="type-count">{{ item.count }}</div>
                </div>
            </div>
            <div v-if="typeData.length > 4" class="more-bottom" @click="handelMore">
                查看更多 >
            </div>
        </div>

        <div class="type-empty" v-else>
            暂无内容
        </div>
    </div>
</template>

<style scoped>
.type-container {
    width: 100%;
    height: auto;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
    padding: 15px;
    border-radius: 10px;
    background-color: #fff;
    transition: all 0.3s ease;
}

:root.dark .type-container {
    background-color: #1A1A1A;
    border: 1px solid #333;
}

.type-container:hover {
    scale: 1.01;
}

.type-title {
    padding: 10px;
    font-size: 18px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-weight: 500;
}

.more {
    width: 15px;
    height: 15px;
    cursor: pointer;
    transition: all 0.2s ease;
}

.more:hover {
    width: 18px;
    height: 18px;
    transform: translateX(2px);
}

.type-list {
    display: flex;
    flex-direction: column;
    gap: 8px;
    /* 增加项目之间的间距 */
    margin-top: 10px;
}

.type-item {
    height: 40px;
    /* 增加高度，提升点击区域 */
    padding: 0 15px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: #666261;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s ease;
    background-color: #f8f9fa;
}

:root.dark .type-item {
    color: #f5f5f5;
    background-color: #2d2d2d;
}

.type-item:hover {
    background-color: #7FC8F8;
    color: #fff;
    transform: translateX(3px);
}

:root.dark .type-item:hover {
    background-color: #7FC8F8;
    color: #fff;
    transform: translateX(3px);
}

.type-name {
    font-size: 15px;
}

.type-count {
    font-size: 14px;
    padding: 2px 8px;
    border-radius: 12px;
}


.type-empty {
    text-align: center;
    padding: 50px 10px;
    color: #999;
    font-size: 16px;
}

/* 底部更多按钮样式 */
.more-bottom {
    margin-top: 15px;
    text-align: center;
    color: #7FC8F8;
    font-size: 14px;
    cursor: pointer;
    padding: 8px 0;
    border-radius: 8px;
    transition: all 0.2s ease;
}

.more-bottom:hover {
    color: #4aa8e0;
}

:root.dark .more-bottom:hover {
    background-color: rgba(127, 200, 248, 0.05);
}

/* 移除重复的dark模式样式 */
:root.dark .type-container {
    border: 1px solid #333;
}
</style>