<script setup lang="ts">

import { ref, onMounted, onUnmounted, computed } from 'vue'
import { chatWithAI, getChatStatus } from '../api/chat'

/** 站点 favicon / 网页图标（public/logo.png） */
const siteLogo = '/logo.png'

const question = ref('')
const answer = ref('')
const isLoading = ref(false)
const isAgent = ref(false)
const defQuestion = ref<string[]>([])

const isDragging = ref(false)
const hasMoved = ref(false)
const dragStartX = ref(0)
const dragStartY = ref(0)
const posLeft = ref(0)
const posTop = ref(0)

const DRAG_THRESHOLD = 8

const isMobile = ref(false)
const serviceHint = ref('')

const PANEL_WIDTH = 400
const PANEL_HEIGHT = 520

const clampPosition = () => {
    const margin = 8
    const w = isAgent.value ? Math.min(PANEL_WIDTH, window.innerWidth - margin * 2) : 56
    const maxL = window.innerWidth - w - margin
    // 不限制底部，允许拖到底部
    posLeft.value = Math.min(Math.max(margin, posLeft.value), Math.max(margin, maxL))
    // 顶部至少留 margin
    posTop.value = Math.max(margin, posTop.value)
}

const initPosition = () => {
    isMobile.value = window.innerWidth <= 768
    // 悬浮按钮与聊天窗默认靠左
    posLeft.value = 16
    const h = isAgent.value ? Math.min(PANEL_HEIGHT, window.innerHeight - 100) : 56
    posTop.value = Math.max(70, window.innerHeight - h - 20)
    clampPosition()
}

const panelStyle = computed(() => ({
    left: `${posLeft.value}px`,
    top: `${posTop.value}px`,
}))

const handleQA = async () => {
    const q = question.value.trim()
    if (!q || isLoading.value) return
    isLoading.value = true
    answer.value = ''

    try {
        const res = await chatWithAI(q)
        answer.value = res.data?.content || '未获取到答案'
        question.value = ''
    } catch (err) {
        console.error(err)
        const msg = err instanceof Error ? err.message : '请求失败，请稍后再试'
        answer.value = msg
    } finally {
        isLoading.value = false
    }
}

const openAgent = () => {
    isAgent.value = true

    // 移动端强制居中显示，不使用拖拽位置
    if (isMobile.value) {
        posLeft.value = Math.max(12, (window.innerWidth - 340) / 2)
        posTop.value = 80
    } else {
        clampPosition()
    }
}

const closeAgent = () => {
    isAgent.value = false
    answer.value = ''
    clampPosition()
}

const handleQuestion = (q: string) => {
    question.value = q
    handleQA()
}

const startDrag = (e: PointerEvent) => {
    isDragging.value = true
    hasMoved.value = false
    dragStartX.value = e.clientX
    dragStartY.value = e.clientY

    const target = e.currentTarget as HTMLElement
    target?.setPointerCapture?.(e.pointerId)
    e.preventDefault()
}

/** 悬浮按钮：支持拖拽，轻点打开 */
const onTriggerPointerDown = (e: PointerEvent) => {
    if (isAgent.value) return
    startDrag(e)
}

/** 桌面端：标题栏拖拽移动面板 */
const onPanelDragStart = (e: PointerEvent) => {
    if (isMobile.value) return
    startDrag(e)
}

const onDragMove = (e: PointerEvent) => {
    if (!isDragging.value) return

    const dx = e.clientX - dragStartX.value
    const dy = e.clientY - dragStartY.value

    if (!hasMoved.value && Math.hypot(dx, dy) < DRAG_THRESHOLD) return

    hasMoved.value = true
    posLeft.value += dx
    posTop.value += dy
    dragStartX.value = e.clientX
    dragStartY.value = e.clientY
    clampPosition()
}

const onDragEnd = (e: PointerEvent) => {
    if (!isDragging.value) return

    const moved = hasMoved.value
    const shouldOpen = !isAgent.value && !moved
    isDragging.value = false
    hasMoved.value = false

    const target = e.currentTarget as HTMLElement
    try {
        target?.releasePointerCapture?.(e.pointerId)
    } catch {}

    if (shouldOpen) {
        openAgent()
    }
}

const onResize = () => initPosition()

onMounted(async () => {
    defQuestion.value = ['你是谁？', '介绍下中南民族大学？', '能给我个简单的 Python 程序吗？']
    initPosition()
    try {
        const st = await getChatStatus()
        if (!st.data?.configured) {
            serviceHint.value = 'AI 暂时不可用'
        }
    } catch {
        serviceHint.value = 'AI 暂时不可用，请稍后重试'
    }
    window.addEventListener('resize', onResize)
    window.addEventListener('pointermove', onDragMove as any)
    window.addEventListener('pointerup', onDragEnd as any)
    window.addEventListener('pointercancel', onDragEnd as any)
})

onUnmounted(() => {
    window.removeEventListener('resize', onResize)
    window.removeEventListener('pointermove', onDragMove as any)
    window.removeEventListener('pointerup', onDragEnd as any)
    window.removeEventListener('pointercancel', onDragEnd as any)
})
</script>

<template>
    <div
        class="agent-container"
        :class="{ 'is-open': isAgent, 'is-mobile': isMobile }"
        :style="panelStyle"
    >
        <button
            v-show="!isAgent"
            type="button"
            class="agent-trigger"
            aria-label="打开 AI 助手"
            @pointerdown="onTriggerPointerDown"
            @pointerup="onDragEnd"
            @pointercancel="onDragEnd"
        >
            <span class="agent-icon">
                <img :src="siteLogo" alt="AI助手" class="site-logo">
            </span>
            <span class="agent-bubble" v-if="!isMobile">你好，我是你的贴心小助手</span>
        </button>

        <div class="agent-show" v-show="isAgent">
            <div
                class="agent-title"
                @pointerdown="onPanelDragStart"
                @pointerup="onDragEnd"
                @pointercancel="onDragEnd"
            >
                <span>AI 助手</span>
                <button
                    type="button"
                    class="cancleIcon"
                    aria-label="关闭"
                    @pointerdown.stop
                    @click.stop="closeAgent"
                >⨉</button>
            </div>
            <p v-if="serviceHint" class="service-hint">{{ serviceHint }}</p>
            <div class="chat-area">
                <div class="answer-box" v-if="answer">
                    <img :src="siteLogo" class="ai-avatar" alt="">
                    <div class="answer-content">{{ answer }}</div>
                </div>
                <div class="default-box" v-else-if="!isLoading">
                    <button
                        type="button"
                        class="default-item"
                        v-for="(item, index) in defQuestion"
                        :key="index"
                        @click="handleQuestion(item)"
                    >
                        {{ item }}
                    </button>
                </div>
                <div class="loading" v-if="isLoading">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span>正在思考...</span>
                </div>
            </div>
            <div class="input-area">
                <input
                    type="text"
                    v-model="question"
                    placeholder="请输入你的问题..."
                    class="question-input"
                    @keyup.enter="handleQA"
                >
                <button type="button" @click="handleQA" class="submit-btn" :disabled="!question.trim() || isLoading">
                    发送
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
.agent-container {
    position: fixed;
    z-index: 10001;
    pointer-events: none;
}

.agent-container > * {
    pointer-events: auto;
}

.agent-trigger {
    display: flex;
    align-items: center;
    gap: 8px;
    border: none;
    padding: 0;
    background: transparent;
    cursor: grab;
    touch-action: none;
    -webkit-tap-highlight-color: transparent;
    position: relative;
}

.agent-trigger:active {
    cursor: grabbing;
}

.agent-icon {
    width: 56px;
    height: 56px;
    border-radius: 50%;
    background-color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 16px rgba(73, 177, 245, 0.35);
    overflow: hidden;
}

.site-logo {
    width: 40px;
    height: 40px;
    object-fit: cover;
    border-radius: 50%;
}

.agent-bubble {
    opacity: 0;
    padding: 8px 14px;
    background-color: #fff;
    border-radius: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    color: #333;
    font-size: 14px;
    position: absolute;
    left: 64px;
    white-space: nowrap;
    pointer-events: none;
    transition: opacity 0.2s;
}

.agent-trigger:hover .agent-bubble,
.agent-trigger:focus-visible .agent-bubble {
    opacity: 1;
}

.agent-show {
    width: min(400px, calc(100vw - 24px));
    min-height: 420px;
    max-height: min(520px, calc(100vh - 48px));
    display: flex;
    flex-direction: column;
    box-shadow: 0 8px 28px rgba(0, 0, 0, 0.16);
    padding: 18px 18px 16px;
    border-radius: 14px;
    background-color: #fff;
    box-sizing: border-box;
}

:root.dark .agent-show {
    background-color: #1e1e1e;
    border: 1px solid #444;
}

.agent-title {
    font-size: 17px;
    font-weight: 600;
    color: #333;
    margin-bottom: 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: grab;
    user-select: none;
}

/* 移动端顶部拖拽条 */
@media (max-width: 768px) {
  .agent-title::before {
    content: '';
    position: absolute;
    top: 6px;
    left: 50%;
    transform: translateX(-50%);
    width: 36px;
    height: 4px;
    background: #ddd;
    border-radius: 2px;
  }
}

:root.dark .agent-title {
    color: #eee;
    border-bottom-color: #444;
}

.service-hint {
    margin: -6px 0 10px;
    font-size: 12px;
    color: #e6a23c;
    line-height: 1.4;
}

.cancleIcon {
    font-size: 20px;
    line-height: 1;
    border: none;
    background: none;
    cursor: pointer;
    color: #999;
    padding: 4px 8px;
    flex-shrink: 0;
    touch-action: manipulation;
    -webkit-tap-highlight-color: transparent;
}

.cancleIcon:hover {
    color: #666;
}

.chat-area {
    flex: 1;
    min-height: 260px;
    max-height: min(380px, calc(100vh - 220px));
    margin-bottom: 14px;
    padding: 12px;
    border-radius: 10px;
    background-color: #f5f8fc;
    overflow-y: auto;
    -webkit-overflow-scrolling: touch;
}

:root.dark .chat-area {
    background-color: #2a2a2a;
}

.default-box {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.default-item {
    text-align: left;
    background-color: #fff;
    padding: 10px 12px;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    color: #555;
    cursor: pointer;
    line-height: 1.5;
    touch-action: manipulation;
}

:root.dark .default-item {
    background-color: #333;
    color: #ddd;
}

.answer-box {
    display: flex;
    gap: 10px;
    align-items: flex-start;
}

.ai-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    object-fit: cover;
    flex-shrink: 0;
}

.answer-content {
    flex: 1;
    color: #333;
    line-height: 1.65;
    font-size: 14px;
    padding: 8px 10px;
    border-radius: 8px;
    background-color: #fff;
    white-space: pre-wrap;
    word-break: break-word;
}

:root.dark .answer-content {
    background-color: #333;
    color: #eee;
}

.loading {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #666;
    font-size: 14px;
}

.dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: #49b1f5;
    animation: bounce 1.4s infinite ease-in-out both;
}

.dot:nth-child(2) { animation-delay: 0.2s; }
.dot:nth-child(3) { animation-delay: 0.4s; }

@keyframes bounce {
    0%, 80%, 100% { transform: scale(0); }
    40% { transform: scale(1); }
}

.input-area {
    display: flex;
    gap: 10px;
    flex-shrink: 0;
}

.question-input {
    flex: 1;
    min-width: 0;
    padding: 12px 14px;
    border: 1px solid #ddd;
    border-radius: 10px;
    font-size: 15px;
    line-height: 1.4;
}

.submit-btn {
    flex-shrink: 0;
    min-height: 44px;
    padding: 0 20px;
    background-color: #49b1f5;
    color: #fff;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    cursor: pointer;
    touch-action: manipulation;
}

.submit-btn:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.is-mobile .agent-icon {
    width: 52px;
    height: 52px;
}

.is-mobile .site-logo {
    width: 36px;
    height: 36px;
}

.is-mobile.is-open .agent-show {
    width: 92vw;
    max-width: 340px;
    min-height: 360px;
    max-height: min(520px, calc(100vh - 100px));
    padding: 16px 14px;
    box-shadow: 0 8px 30px rgba(0, 0, 0, 0.2);
    border-radius: 16px;
}

.is-mobile.is-open {
    left: 50% !important;
    top: 80px !important;
    transform: translateX(-50%);
}

.is-mobile .chat-area {
    min-height: 200px;
    max-height: 280px;
    flex: 1;
}

.is-mobile .agent-icon {
    width: 42px;
    height: 42px;
}
</style>
