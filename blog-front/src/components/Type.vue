<script setup lang="ts">
import { ElButton } from 'element-plus'
import 'element-plus/es/components/button/style/css'

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

const error = ref('')
const loading = ref(true)
const loadTags = async () => {
    error.value = ''
    loading.value = true
    try {
        const res = await getTag()
        typeData.value = res.data ?? []
    } catch (e) {
        error.value = e instanceof Error ? e.message : '加载分类失败'
    } finally { loading.value = false }
}
onMounted(loadTags)

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

        <div v-if="error" role="alert">{{ error }} <el-button text @click="loadTags">重试</el-button></div>
        <div v-else-if="loading">加载中…</div>
        <div v-else-if="typeData.length > 0">
            <div class="type-list">
                <!-- 全部标签 -->
                <div class="type-item all-item" @click="handelMore">
                    <div class="type-name">全部</div>
                </div>
                <div class="type-item" v-for="item in displayTypeData" :key="item.id" @click="handleType(item.id)">
                    <div class="type-name">{{ item.name }}</div>
                </div>
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
    margin-top: 10px;
}

/* 移动端分类 - 参考掘金风格：小巧横向标签 */
@media (max-width: 768px) {
  .type-container {
    padding: 8px 10px !important;
    box-shadow: none !important;
    background-color: transparent !important;
    border: none !important;
  }

  .type-title {
    display: none;
  }

  .type-list {
    flex-direction: row !important;
    flex-wrap: nowrap !important;
    overflow-x: auto !important;
    gap: 6px !important;
    padding: 2px 0 8px;
    -webkit-overflow-scrolling: touch;
    scrollbar-width: none;
  }

  .type-list::-webkit-scrollbar {
    display: none;
  }

  .type-item {
    flex-shrink: 0 !important;
    height: 30px !important;
    padding: 0 11px !important;
    font-size: 12.5px !important;
    border-radius: 15px !important;
    background-color: #f0f4f8 !important;
    color: #555 !important;
    transition: all 0.2s ease;
  }

  .type-item:hover,
  .type-item:active {
    background-color: #e6f0ff !important;
    color: #49b1f5 !important;
  }

  .all-item {
    background-color: #fff7e6 !important;
    color: #d46b08 !important;
    font-weight: 500;
  }

  :root.dark .type-item {
    background-color: #2d2d2d !important;
    color: #ddd !important;
  }

  :root.dark .all-item {
    background-color: #3a2a1a !important;
    color: #f0a040 !important;
  }

  .type-count {
    display: none;
  }
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