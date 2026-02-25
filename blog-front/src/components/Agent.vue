<script setup lang="ts">
import {ref,onMounted,onUnmounted} from 'vue'
import axios from 'axios'
const question=ref<string>('')
const answer=ref<string>('')
const isLoading=ref<boolean>(false)
const isAgent=ref<boolean>(false)
const defQuestion =ref<string[]>([])

// 拖拽相关变量
const isDragging = ref(false)
const startX = ref(0)
const startY = ref(0)
const initialLeft = ref(30)
const initialTop = ref(window.innerHeight - 200)

const isDragged=ref<boolean>(false)

const handleQA = async () => {
    if (!question.value.trim()) return
    isLoading.value = true
    answer.value = ''

    try {
        const res = await axios.get(`https://api.fohok.xin/API/GLM_ai.php?msg=${question.value}`,)
        answer.value = res.data?.data?.choices?.[0]?.message?.content || '未获取到答案'
        question.value=''
    } catch (err) {
        console.error(err)
        answer.value = '请求失败，请稍后再试'
    } finally {
        isLoading.value = false
    }
}

const handleAgent = () => {
    if (isDragged.value) return
    isAgent.value = true
}
const cancleAgent=()=>{
    isAgent.value=false
    answer.value=''
}
const handleQuestion=(q:string)=>{
    question.value=q
    handleQA()
}

// 鼠标按下：开始拖拽
const handleMouseDown = (e: MouseEvent) => {
    isDragging.value=true
    startX.value = e.clientX
    startY.value = e.clientY
    e.preventDefault()
}

// 鼠标移动：更新位置
const handleMouseMove = (e: MouseEvent) => {
    if (!isDragging.value) return
    const dx = e.clientX - startX.value
    const dy = e.clientY - startY.value
    isDragged.value = true

    initialLeft.value += dx
    initialTop.value += dy

    startX.value = e.clientX
    startY.value = e.clientY
}

// 鼠标松开：结束拖拽
const handleMouseUp = () => {
    isDragging.value = false
    setTimeout(() => { isDragged.value = false }, 0)
}

// 监听全局鼠标事件
onMounted(() => {
    defQuestion.value = ['你是谁？', '介绍下中南民族大学？', '能给我个简单的python程序吗？']
    window.addEventListener('mousemove', handleMouseMove)
    window.addEventListener('mouseup', handleMouseUp)
})

onUnmounted(() => {
    window.removeEventListener('mousemove', handleMouseMove)
    window.removeEventListener('mouseup', handleMouseUp)
})
</script>

<template>
    <div 
        class="agent-container" 
        @mousedown.stop="handleMouseDown"
        :style="{
            left: `${initialLeft}px`,
            top: `${initialTop}px`,
            cursor: isDragging ? 'grabbing' : 'grab'
        }"
      >
        <div class="agent-trigger" @click="handleAgent" v-show="!isAgent">
          <div class="agent-icon">
            <img src="../assets/agent.svg" alt="AI助手">
          </div>
          <div class="agent-bubble">
            你好，我是你的贴心小助手
          </div>
        </div>
        <div class="agent-show" v-show="isAgent">
            <div class="agent-title">
                <span>AI 助手</span>
                <span class="cancleIcon" @click="cancleAgent">⨉</span>
            </div>
            <!-- 对话区域 -->
            <div class="chat-area">
                <!-- 答案显示区域 -->
                <div class="answer-box" v-if="answer">
                    <img src="../assets/agent.svg" class="ai-avatar">
                    <div class="answer-content">{{ answer }}</div>
                </div>
                <!-- 默认问题区域 -->
                <div class="default-box" v-else>
                    <div class="default-item"
                     @click="handleQuestion(item)"
                     v-for="(item,index) in defQuestion" :key="index">
                        {{ item }}
                    </div>
                </div>
                <!-- 加载状态 -->
                <div class="loading" v-if="isLoading">
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span class="dot"></span>
                    <span>正在思考...</span>
                </div>
            </div>

            <!-- 输入区域 -->
            <div class="input-area">
                <input type="text" v-model="question" placeholder="请输入你的问题..." class="question-input" @keyup.enter="handleQA"
                >
                <button @click="handleQA" class="submit-btn" :disabled="!question.trim()">
                    发送
                </button>
            </div>
        </div>
    </div>
</template>

<style scoped>
/* 容器基础样式 */
.agent-container {
    position: fixed;
    z-index: 99;
}

/* 悬浮触发按钮 */
.agent-trigger {
    z-index: 99;
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    transition: all 0.2s;
}

/* AI图标容器 */
.agent-icon {
    width: 56px;
    height: 56px;
    border-radius: 50%;
    background-color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 16px rgba(64, 150, 255, 0.3);
    transition: transform 0.2s, box-shadow 0.2s;
}

.agent-icon img {
    width: 32px;
    height: 32px;
    object-fit: contain;
}

/* 气泡提示框 */
.agent-bubble {
    opacity: 0;
    padding: 8px 16px;
    background-color: #fff;
    border-radius: 20px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
    color: #333;
    font-size: 14px;
    min-width: 170px;
    left: 65px;
    position: absolute;
    transition: box-shadow 0.2s;
}


/* 触发按钮hover效果 */
.agent-trigger:hover .agent-icon {
    transform: scale(1.05);
    box-shadow: 0 6px 20px rgba(64, 150, 255, 0.4);
}

.agent-trigger:hover .agent-bubble {
    opacity: 1;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

/* AI助手主体样式 */
.agent-show {
    z-index: 98;
    width: 350px;
    max-width: 90vw;
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
    padding: 20px;
    border-radius: 12px;
    background-color: #fff;
}

/* 标题栏样式 */
.agent-title {
    font-size: 18px;
    font-weight: 600;
    color: #333;
    margin-bottom: 16px;
    padding-bottom: 8px;
    border-bottom: 1px solid #f0f0f0;
    display: flex;
    align-items: center;
    justify-content: space-between;
    position: relative;
}

/* 关闭按钮样式 */
.cancleIcon {
    font-size: 18px;
    cursor: pointer;
    color: #999;
    transition: color 0.2s;
}

.cancleIcon:hover {
    color: #333;
}

/* 聊天区域样式 */
.chat-area {
    min-height: 120px;
    max-height: 250px;
    margin-bottom: 20px;
    padding: 10px;
    border-radius: 8px;
    background-color: #f9f9f9;
    overflow-y: auto;
}
.default-box{
    display: flex;
    flex-direction: column;
    gap: 12px;
    padding: 15px 10px;
}
.default-item {
  background-color: #fff;
  padding: 10px 14px;
  border-radius: 8px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  font-size: 14px;
  color: #555;
  cursor: pointer;
  transition: all 0.2s ease; 
}
.default-item:hover {
  background-color: #f0f7ff;
  color: #4096ff;
  box-shadow: 0 2px 6px rgba(64, 150, 255, 0.15);
}

/* 答案显示样式 */
.answer-box {
    display: flex;
    align-items: flex-start;
    max-height: 180px;
    overflow-y: auto;
    gap: 10px;
    padding: 10px 0;
}

.ai-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background-color: #e6f7ff;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    object-fit: contain;
    padding: 4px;
}

.answer-content {
    flex: 1;
    color: #333;
    line-height: 1.6;
    padding: 8px 12px;
    border-radius: 8px;
    background-color: #fff;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
    white-space: pre-wrap;
}

/* 加载动画样式 */
.loading {
    display: flex;
    align-items: center;
    gap: 8px;
    color: #666;
    padding: 10px;
}

.dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background-color: #4096ff;
    animation: bounce 1.4s infinite ease-in-out both;
}

.dot:nth-child(2) {
    animation-delay: 0.2s;
}

.dot:nth-child(3) {
    animation-delay: 0.4s;
}

@keyframes bounce {

    0%,
    80%,
    100% {
        transform: scale(0);
    }

    40% {
        transform: scale(1);
    }
}

/* 输入区域样式 */
.input-area {
    display: flex;
    gap: 10px;
}

.question-input {
    flex: 1;
    padding: 12px 16px;
    border: 1px solid #ddd;
    border-radius: 8px;
    font-size: 14px;
    transition: border-color 0.2s, box-shadow 0.2s;
}

.question-input:focus {
    outline: none;
    border-color: #4096ff;
    box-shadow: 0 0 0 2px rgba(64, 150, 255, 0.1);
}

.question-input::placeholder {
    color: #999;
}

.submit-btn {
    padding: 0 20px;
    background-color: #4096ff;
    color: white;
    border: none;
    border-radius: 8px;
    font-size: 14px;
    font-weight: 500;
    cursor: pointer;
    transition: background-color 0.2s;
}

.submit-btn:hover {
    background-color: #1890ff;
}

.submit-btn:disabled {
    background-color: #b3d8ff;
    cursor: not-allowed;
}

</style>