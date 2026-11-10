<script setup lang="ts">
import { ElButton, ElEmpty, ElLoading } from 'element-plus'
import 'element-plus/es/components/button/style/css'
import 'element-plus/es/components/empty/style/css'
import 'element-plus/es/components/loading/style/css'
const vLoading = ElLoading.directive
import 'element-plus/es/components/message/style/css'
import { getHollow, addHollow } from '../api/hollow'
import AnswerBook from '../components/AnswerBook.vue'
import Mine from '../components/Mine.vue'
import Type from '../components/Type.vue'
import { onMounted, onUnmounted, ref } from 'vue'
import { formatTime } from '../utils/format'
import { type ApiResponse } from '../types/articles'
import { ElMessage } from 'element-plus'

interface hollow {
    date: string
    word: string
}
const total = ref<number>(0)
const hollows = ref<hollow[]>([])
const isEdit = ref<boolean>(false)
const word = ref<string>('')
const submitting = ref(false)

const preventScroll = (e: Event) => {
    e.preventDefault()
}

const changeEdit = () => {
    isEdit.value = !isEdit.value

    if (isEdit.value) {
        document.addEventListener('wheel', preventScroll, { passive: false })
        document.addEventListener('touchmove', preventScroll, { passive: false })
    } else {
        document.removeEventListener('wheel', preventScroll)
        document.removeEventListener('touchmove', preventScroll)
    }
    word.value = ''
}

const loading = ref(true)
const loadError = ref('')
const loadHollows = async () => {
    loading.value = true
    loadError.value = ''
    try {
        const res = await getHollow()
        const list = res.data ?? []
        hollows.value = [...list].reverse()
        total.value = list.length
    } catch (e) {
        console.error(e)
        loadError.value = e instanceof Error ? e.message : '加载树洞失败'
    } finally { loading.value = false }
}

const subbmitHollow = async () => {
    if (!word.value.trim()) {
        ElMessage.warning('输入内容不能为空')
        return
    }
    if (submitting.value) return
    submitting.value = true
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const day = String(now.getDate()).padStart(2, '0')
    const hours = String(now.getHours()).padStart(2, '0')
    const minutes = String(now.getMinutes()).padStart(2, '0')
    const seconds = String(now.getSeconds()).padStart(2, '0')
    const submitDate = `${year}-${month}-${day}T${hours}:${minutes}:${seconds}`
    try {
        const res = await addHollow(word.value.trim(), submitDate) as unknown as ApiResponse
        if (res?.code === 200) {
            word.value = ''
            isEdit.value = false
            document.removeEventListener('wheel', preventScroll)
            document.removeEventListener('touchmove', preventScroll)
            await loadHollows()
            ElMessage.success('发布成功')
        } else {
            ElMessage.error(res?.msg || res?.message || '发布失败')
        }
    } catch (e) {
        console.error(e)
        ElMessage.error(e instanceof Error ? e.message : '发布失败，请重试')
    } finally {
        submitting.value = false
    }
}

onMounted(() => {
    loadHollows()
})

onUnmounted(() => {
    document.removeEventListener('wheel', preventScroll)
    document.removeEventListener('touchmove', preventScroll)
})
</script>

<template>
    <div class="hollow-page">
        <Type class="page-type" />
        <div class="left">
            <div class="title">闲言碎语</div>
            <img src="../assets/edit.png" class="edit-icon" @click="changeEdit" alt="编辑">
            <div class="total">共{{ total }}条</div>
            <div class="word-container" v-loading="loading">
                <el-empty v-if="loadError" :description="loadError"><el-button @click="loadHollows">重试</el-button></el-empty>
                <el-empty v-else-if="!loading && hollows.length === 0" description="暂无树洞内容" />
                <div v-for="(hollow, index) in hollows" :key="`${hollow.date}-${index}`" class="word-box">
                    <div>
                        <span>◆</span>
                        <span class="time">{{ formatTime(hollow.date) }}</span>
                    </div>
                    <div class="word">{{ hollow.word }}</div>
                </div>
            </div>
        </div>
        <div class="right">
            <Mine class="page-mine" />
            <AnswerBook class="page-answer" />
        </div>
        <div class="overlay" v-show="isEdit" @click.self="changeEdit">
            <div class="overlay-box">
                <div class="overlay-title">今日感想</div>
                <textarea placeholder="请输入今日感想" class="input-area" v-model="word" maxlength="1000"></textarea>
                <div class="overlay-footer">
                    <el-button type="primary" :loading="submitting" @click="subbmitHollow">提交</el-button>
                    <span class="wordLength">共{{ word.length }}字</span>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.hollow-page {
    height: auto;
    position: relative;
    display: grid;
    grid-template-columns: 1fr 280px;
    grid-template-areas:
        "left mine"
        "left type"
        "left answer";
    justify-content: center;
    max-width: 1400px;
    margin: 0 auto;
    padding: 20px 24px;
    gap: 20px;
    box-sizing: border-box;
}

.page-type {
    grid-area: type;
    align-self: start;
}

.left {
    grid-area: left;
    min-width: 0;
    position: relative;
    background-color: #fff;
    padding: 24px;
    border-radius: 12px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
}

:root.dark .left {
    background-color: #1f1f1f;
    border: 1px solid #333;
}

.right {
    display: contents;
}

.page-mine {
    grid-area: mine;
    align-self: start;
}

.page-answer {
    grid-area: answer;
    align-self: start;
}

.title {
    font-size: 22px;
    font-weight: 600;
    margin: 0 0 8px;
    color: #333;
}

.total {
    font-size: 14px;
    color: #888;
    margin-bottom: 16px;
}

.edit-icon {
    position: absolute;
    right: 24px;
    top: 24px;
    cursor: pointer;
    width: 22px;
    height: 22px;
}

.word-container {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.word-box {
    background-color: #fafafa;
    border-radius: 10px;
    padding: 14px 16px;
    border: 1px solid #f0f0f0;
}

:root.dark .word-box {
    background-color: #2a2a2a;
    border-color: #333;
}

.time {
    font-size: 13px;
    color: #888;
    margin-left: 6px;
}

.word {
    border-left: 3px solid #49b1f5;
    min-height: 40px;
    display: flex;
    align-items: center;
    padding: 8px 0 8px 12px;
    margin: 6px 0 0;
    font-size: 15px;
    line-height: 1.6;
    color: #333;
}

:root.dark .word {
    border-color: #7fc8f8;
    color: #ddd;
}
.overlay {
    position: fixed;
    width: 100%;
    height: 100%;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.7);
    z-index: 9999;
    display: flex;
    align-items: center;
    justify-content: center;
    overflow: auto;
}
.overlay-box {
    background-color: #fff;
    width: 400px;
    height: 250px;
    padding: 20px;
    border-radius: 10px;
}
.overlay-title {
    font-size: 22px;
    color: #333;
}
.input-area {
    width: 100%;
    max-width: 400px;
    height: 160px;
    resize: none;
    font-size: 18px;
    margin: 10px 0;
    box-sizing: border-box;
}
.overlay-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
}

/* 移动端优化 - 和首页一致的现代风格 */
@media (max-width: 992px) {
  .hollow-page {
    grid-template-columns: 1fr;
    grid-template-areas:
      "type"
      "left"
      "mine"
      "answer";
    padding: 16px;
    margin: 0;
    gap: 12px;
  }

  .page-type {
    margin: 0 -4px;
  }

  .left {
    width: 100% !important;
    padding: 20px 16px !important;
  }

  .title {
    font-size: 19px !important;
  }

  .edit-icon {
    right: 16px !important;
    top: 20px !important;
  }

  .word-box {
    padding: 12px 14px !important;
  }

  .word {
    font-size: 14.5px !important;
    padding-left: 10px !important;
  }
}
</style>
