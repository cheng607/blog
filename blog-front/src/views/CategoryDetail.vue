<script setup lang="ts">
import { ElButton, ElDivider, ElEmpty, ElLoading } from 'element-plus'
import 'element-plus/es/components/empty/style/css'
import 'element-plus/es/components/divider/style/css'
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/loading/style/css'
import 'element-plus/es/components/message/style/css'
const vLoading = ElLoading.directive
import { ref, watch, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import List from '../components/List.vue'
import { getTagByid } from '../api/tag'
import { ElMessage } from 'element-plus'

interface typeObj {
    id: number
    name: string
    description: string
}
const route = useRoute()
const typeData = ref<typeObj>()
const loading = ref(true)

const error = ref('')
let requestId = 0
onUnmounted(() => { requestId++ })
const loadCategory = async (typeId: unknown = route.query.typeId) => {
    const current = ++requestId
    typeData.value = undefined
    loading.value = true
    error.value = ''
    if (typeof typeId !== 'string' || !/^[1-9]\d*$/.test(typeId) || !Number.isSafeInteger(Number(typeId))) {
        error.value = '分类不存在'
        loading.value = false
        return
    }
    try {
        const res = await getTagByid(Number(typeId))
        if (current !== requestId) return
        typeData.value = res.data ?? undefined
        if (!typeData.value) error.value = '分类不存在'
    } catch (e) {
        if (current !== requestId) return
        error.value = e instanceof Error ? e.message : '加载分类信息失败，请稍后重试'
        ElMessage.error(error.value)
    } finally {
        if (current === requestId) loading.value = false
    }
}
watch(() => route.query.typeId, id => loadCategory(id), { immediate: true })

</script>

<template>
    <div class="detail-page">
        <div class="cat-header">
            <h1>{{ typeData?.name }}</h1>
            <p>{{ typeData?.description }}</p>
        </div>
        <ElDivider />
        <div class="cat-content" v-loading="loading">
            <el-empty v-if="error" :description="error"><el-button @click="loadCategory()">重试</el-button><el-button @click="$router.push('/categories')">返回分类</el-button></el-empty>
            <List v-else-if="typeData" :key="typeData.id" />
        </div>
    </div>
</template>

<style scoped>
.detail-page {
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

.cat-content {
    margin-top: 30px;
    padding: 0 20px;
}
</style>
