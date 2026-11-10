<script setup lang="ts">
import { ElBacktop, ElConfigProvider } from 'element-plus'
import 'element-plus/es/components/backtop/style/css'
import { ref, onMounted, onUnmounted } from 'vue'
import Footer from './components/Footer.vue'
import Header from './components/Header.vue'
import Agent from './components/Agent.vue'
import AdminReauth from './components/AdminReauth.vue'
import { resolveTheme } from './utils/theme'

const isDark = ref(resolveTheme() === 'dark')

let lastTime = 0
const TRAIL_INTERVAL = 20
const trailColors = [
    'rgba(73, 177, 245, 0.8)',
    'rgba(127, 200, 248, 0.6)',
    'rgba(76, 50, 120, 0.7)',
    'rgba(255, 255, 255, 0.5)',
]

const syncDark = () => {
    isDark.value = document.documentElement.classList.contains('dark')
}

const createTrail = (x: number, y: number) => {
    const container = document.getElementById('trail-container')
    if (!container) return

    const trail = document.createElement('div')
    const size = Math.random() * 10 + 5
    const color = trailColors[Math.floor(Math.random() * trailColors.length)]

    trail.style.cssText = `
    position: fixed;
    left: ${x}px;
    top: ${y}px;
    width: ${size}px;
    height: ${size}px;
    background-color: ${color};
    border-radius: 50%;
    pointer-events: none;
    z-index: 999;
    transform: translate(-50%, -50%);
    box-shadow: 0 0 ${size * 2}px ${color};
  `

    container.appendChild(trail)

    const duration = Math.random() * 1000 + 5000
    trail.animate(
        [
            { opacity: 1, transform: `translate(-50%, -50%) scale(1)` },
            {
                opacity: 0,
                transform: `translate(-50%, -50%) scale(0.2) translateY(-${Math.random() * 30 + 20}px)`,
            },
        ],
        { duration, easing: 'cubic-bezier(0.2, 0.8, 0.2, 1)' }
    )

    setTimeout(() => trail.remove(), duration)
}

const handleMouseMove = (e: MouseEvent) => {
    const now = Date.now()
    if (now - lastTime > TRAIL_INTERVAL) {
        createTrail(e.clientX, e.clientY)
        lastTime = now
    }
}

let observer: MutationObserver | undefined

onMounted(() => {
    syncDark()
    observer = new MutationObserver(syncDark)
    observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] })
    window.addEventListener('mousemove', handleMouseMove)
})

onUnmounted(() => {
    observer?.disconnect()
    window.removeEventListener('mousemove', handleMouseMove)
})
</script>

<template>
    <el-config-provider>
        <div class="main" :class="{ 'is-dark': isDark }">
            <Header />
            <router-view />
            <AdminReauth />
            <Footer />
            <div id="trail-container" />
            <Agent class="agent" />
            <el-backtop :right="20" :bottom="24" class="backtop-box">
                <img src="./assets/backTop.png" alt="返回顶部" />
            </el-backtop>
        </div>
    </el-config-provider>
</template>

<style scoped>
.main {
    min-height: 100vh;
    background-color: var(--bg-color);
    color: var(--text-color);
    transition: background-color 0.3s ease, color 0.3s ease;
}

.backtop-box {
    padding: 5px;
    box-shadow: 0 3px 10px rgba(0, 0, 0, 0.2);
}

.agent {
    z-index: 99;
}

#trail-container {
    position: fixed;
    top: 0;
    left: 0;
    pointer-events: none;
    z-index: 999;
    width: 100vw;
    height: 100vh;
}
</style>
