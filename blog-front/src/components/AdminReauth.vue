<script setup lang="ts">
import { ElButton, ElDialog, ElInput } from 'element-plus'
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/dialog/style/css'
import 'element-plus/es/components/input/style/css'
import 'element-plus/es/components/message/style/css'
import { ref } from 'vue'
import { ElMessage } from 'element-plus'
import { loginAdmin, reauthRequired } from '../utils/adminAuth'

const password = ref('')
const loading = ref(false)
const login = async () => {
    if (loading.value || !password.value) return
    loading.value = true
    try {
        await loginAdmin(password.value)
        ElMessage.success('登录成功，请重新执行刚才的发布或上传操作')
    } catch (error) {
        ElMessage.error(error instanceof Error ? error.message : '登录失败')
    } finally {
        password.value = ''
        loading.value = false
    }
}
</script>

<template>
    <el-dialog v-model="reauthRequired" title="请重新登录" width="min(420px, 92vw)"
        :close-on-click-modal="false" @closed="password = ''">
        <p>当前编辑内容已保留。登录后可继续发布或上传。</p>
        <el-input v-model="password" type="password" show-password autocomplete="current-password"
            placeholder="管理密钥" @keyup.enter="login" />
        <template #footer>
            <el-button @click="reauthRequired = false">继续编辑</el-button>
            <el-button type="primary" :loading="loading" @click="login">登录</el-button>
        </template>
    </el-dialog>
</template>
