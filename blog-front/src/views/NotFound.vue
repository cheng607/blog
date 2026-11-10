<script setup lang="ts">
import { ElButton } from 'element-plus'
import 'element-plus/es/components/button/style/css'
import { useRouter } from 'vue-router'
const router = useRouter()
</script>

<template>
    <main class="not-found">
        <h1>404</h1>
        <p>这个页面不存在，可能已被移动或删除。</p>
        <el-button type="primary" @click="router.push('/')">返回首页</el-button>
    </main>
</template>

<style scoped>
.not-found { min-height: 50vh; display: grid; place-content: center; text-align: center; gap: 16px; padding: 32px; }
h1 { font-size: 64px; margin: 0; }
</style>
