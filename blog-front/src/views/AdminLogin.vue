<script setup lang="ts">
import { ElButton, ElCard, ElInput } from 'element-plus'
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/card/style/css'
import 'element-plus/es/components/input/style/css'
import 'element-plus/es/components/message/style/css'
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { checkAdminSession, loginAdmin } from '../utils/adminAuth'

const route = useRoute()
const router = useRouter()
const password = ref('')
const loading = ref(false)
const redirect = () => {
    const to = route.query.redirect
    router.replace(typeof to === 'string' && ['/admin/write', '/hollow'].includes(to) ? to : '/admin/write')
}
onMounted(async () => {
    // Remove legacy secret-bearing URLs without using their contents.
    if (route.query.key || route.query.token || route.query.secret) {
        await router.replace({ path: '/admin', query: { redirect: route.query.redirect } })
    }
    if (await checkAdminSession()) redirect()
})
const handleLogin = async () => {
    if (loading.value) return
    if (!password.value) { ElMessage.warning('请输入管理密钥'); return }
    loading.value = true
    try {
        await loginAdmin(password.value)
        ElMessage.success('验证成功')
        redirect()
    } catch (error) {
        ElMessage.error(error instanceof Error ? error.message : '登录失败')
    } finally {
        password.value = ''
        loading.value = false
    }
}
</script>

<template>
    <div class="admin-login-page">
        <el-card class="login-card" shadow="hover">
            <h2>管理入口</h2>
            <p class="hint">请输入管理密钥，由服务器验证后进入后台。</p>
            <el-input
                v-model="password"
                type="password"
                show-password
                placeholder="编辑密钥"
                size="large"
                @keyup.enter="handleLogin"
            />
            <el-button type="primary" size="large" class="submit-btn" :loading="loading" @click="handleLogin">
                进入编辑器
            </el-button>
            <el-button text type="info" @click="router.push('/')">返回首页</el-button>
        </el-card>
    </div>
</template>

<style scoped>
.admin-login-page {
    min-height: 60vh;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 40px 20px;
}

.login-card {
    width: 100%;
    max-width: 400px;
    padding: 10px;
}

.login-card h2 {
    margin: 0 0 12px;
    text-align: center;
    color: var(--text-color);
}

.hint {
    font-size: 14px;
    color: #888;
    margin-bottom: 20px;
    line-height: 1.6;
}

:root.dark .hint {
    color: #aaa;
}

.submit-btn {
    width: 100%;
    margin: 16px 0 8px;
}
</style>
