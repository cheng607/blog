<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import lightImg from '../assets/light.png'
import darkImg from '../assets/dark.png'

const router = useRouter()
const route = useRoute()

// 导航配置数据
const navItems = [
    { path: '/', name: '首页', index: '0' },
    { path: '/categories', name: '分类', index: '1' },
    { path: '/hollow', name: '树洞', index: '3' },
    { path: '/about', name: '关于', index: '2' }
]

// 处理导航点击
const handleNavClick = (path: string) => {
    router.push(path)
}

const activeIndex = computed(() => {
    const matchedItem = navItems.find(item => item.path === route.path);
    if (matchedItem) {
        return matchedItem.index;
    }
})

const currentTheme = ref<'light' | 'dark'>('light')
onMounted(() => {
    const savedTheme = localStorage.getItem('app-theme') as 'light' | 'dark' | null
    if (savedTheme) {
        currentTheme.value = savedTheme
    }
    document.documentElement.className = currentTheme.value
})

const handleThemeChange = (isDark: boolean) => {
    currentTheme.value = isDark ? 'dark' : 'light'
    document.documentElement.className = currentTheme.value
    localStorage.setItem('app-theme', currentTheme.value)
}
</script>

<template>
    <div class="nav-container">
        <div class="nav-box">
            <div class="nav-item" :class="{ active: activeIndex === item.index }" @click="handleNavClick(item.path)"
                v-for="item in navItems" :key="item.index">
                {{ item.name }}
            </div>
        </div>
        <div class="change-box">
            <div class="icon-box">
                <img :src="lightImg" v-if="currentTheme === 'light'">
                <img :src="darkImg" v-else>
            </div>
            <el-switch :model-value="currentTheme === 'dark'" @change="handleThemeChange" />
        </div>
    </div>
</template>

<style scoped>
.nav-container {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 80px;
    background-color: #fff;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
}

.nav-box {
    display: flex;
    align-items: center;
}

.nav-item {
    font-size: 15px;
    cursor: pointer;
    color: #333;
    padding: 10px 15px;
    border-radius: 6px;
    transition: all 0.3s ease;
    position: relative;
}

.nav-item:hover:not(.active) {
    background-color: #f5f5f5;
    color: #222;
}

.nav-item.active {
    color: #49B1F5;
    font-weight: 600;
}

.nav-item.active:hover {
    background-color: transparent;
}

.nav-item.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 15px;
    right: 15px;
    height: 2px;
    background-color: #49B1F5;
    border-radius: 1px;
}

.change-box {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 5px;
    margin-right: 150px;
}

.icon-box img {
    width: 20px;
    height: 20px;
    z-index: 1;
    vertical-align: middle;
}
</style>