<script setup lang="ts">

import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import lightImg from '../assets/light.png'
import darkImg from '../assets/dark.png'
import {
    applyTheme,
    initTheme,
    watchSystemTheme,
    type ThemeMode,
} from '../utils/theme'

const router = useRouter()
const route = useRoute()

const navItems = [
    { path: '/', name: '首页', index: '0' },
    { path: '/categories', name: '分类', index: '1' },
    { path: '/hollow', name: '树洞', index: '3' },
    { path: '/about', name: '关于', index: '2' },
]

const handleNavClick = (path: string) => {
    router.push(path)
}

const activeIndex = computed(() => {
    const matchedItem = navItems.find((item) => item.path === route.path)
    return matchedItem?.index
})

const currentTheme = ref<ThemeMode>('light')
let unwatchSystem: (() => void) | undefined

onMounted(() => {
    initTheme()
    currentTheme.value = document.documentElement.classList.contains('dark') ? 'dark' : 'light'
    unwatchSystem = watchSystemTheme((t) => {
        currentTheme.value = t
    })
})

onUnmounted(() => {
    unwatchSystem?.()
})

const handleThemeChange = (isDark: boolean) => {
    const theme: ThemeMode = isDark ? 'dark' : 'light'
    currentTheme.value = theme
    applyTheme(theme)
}
</script>

<template>
    <div class="nav-container">
        <div class="nav-box">
            <div
                class="nav-item"
                :class="{ active: activeIndex === item.index }"
                @click="handleNavClick(item.path)"
                v-for="item in navItems"
                :key="item.index"
            >
                {{ item.name }}
            </div>
        </div>
        <div class="change-box">
            <button
                type="button"
                class="theme-btn"
                :title="currentTheme === 'dark' ? '切换浅色' : '切换深色'"
                @click="handleThemeChange(currentTheme !== 'dark')"
            >
                <img :src="lightImg" v-if="currentTheme === 'light'" alt="浅色" />
                <img :src="darkImg" v-else alt="深色" />
            </button>
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
    transition: background-color 0.3s ease, box-shadow 0.3s ease;
}

:root.dark .nav-container {
    background-color: #1a1a1a;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.45);
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

:root.dark .nav-item {
    color: #e0e0e0;
}

.nav-item:hover:not(.active) {
    background-color: #f5f5f5;
    color: #222;
}

:root.dark .nav-item:hover:not(.active) {
    background-color: #2d2d2d;
    color: #fff;
}

.nav-item.active {
    color: #49b1f5;
    font-weight: 600;
}

:root.dark .nav-item.active {
    color: #7fc8f8;
}

.nav-item.active::after {
    content: '';
    position: absolute;
    bottom: 0;
    left: 15px;
    right: 15px;
    height: 2px;
    background-color: #49b1f5;
    border-radius: 1px;
}

:root.dark .nav-item.active::after {
    background-color: #7fc8f8;
}

.change-box {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 5px;
    margin-right: 0;
}

.icon-box img {
    width: 20px;
    height: 20px;
    vertical-align: middle;
}

/* 移动端主题切换纯图标按钮 */
.theme-btn {
    border: none;
    background: transparent;
    padding: 4px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
}

.theme-btn img {
    width: 18px;
    height: 18px;
}

@media (max-width: 768px) {
    .nav-container {
        padding: 4px 8px !important;
        justify-content: center !important;
        align-items: center;
        gap: 2px;
        flex-wrap: wrap;
        min-height: 44px;
        border-bottom: 1px solid #f0f0f0;
    }

    :root.dark .nav-container {
        border-bottom-color: #333;
    }

    .nav-box {
        flex: 1;
        justify-content: center;
        gap: 0;
        display: flex;
        flex-wrap: wrap;
        align-items: center;
    }

    .nav-item {
        font-size: 14px;
        padding: 6px 10px;
        height: 32px;
        line-height: 20px;
        margin: 2px 3px;
        display: flex;
        align-items: center;
    }

    .change-box {
        margin: 0 8px 0 4px; /* 增加右边距，避免紧贴屏幕边缘 */
        padding: 0;
        flex-shrink: 0;
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .theme-btn img {
        width: 17px;
        height: 17px;
    }

    .theme-btn {
        padding: 6px; /* 增大点击区域，提升体验 */
    }
}
</style>
