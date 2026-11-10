<script setup lang="ts">

import { ref } from 'vue'
import axios from 'axios'

const answer_zh = ref<string>('')
const answer_en = ref<string>('')
const showAnswer = ref<boolean>(false)
const isLoading = ref(false)

const getAnswer = async () => {
  showAnswer.value = false
  isLoading.value = true
  try {
    const res = await axios.get(`https://api.4qb.cn/api/daan_book`)
    answer_zh.value = res.data.data.zh
    answer_en.value = res.data.data.en
    if (res.data.data) {
      setTimeout(() => {
        showAnswer.value = true
        isLoading.value = false
      }, 800)
    } else {
      isLoading.value = false
    }
  } catch (error) {
    console.error('请求失败', error)
    answer_zh.value = '暂时无法回答，请稍后再试'
    answer_en.value = 'Unable to answer at the moment, please try again later'
    isLoading.value = false
    showAnswer.value = true
  }
}

const handleAgain = () => {
  showAnswer.value = false
}
</script>

<template>
  <div class="answerBook-container">
    <div class="book-header">
      <span class="book-icon" aria-hidden="true">📖</span>
      <div class="book-meta">
        <div class="title">答案之书</div>
        <p class="subtitle">默念问题，翻开属于你的答案</p>
      </div>
    </div>

    <div class="answer-box">
      <div v-if="isLoading" class="loading-state">
        <span class="loading-dot"></span>
        <span class="loading-dot"></span>
        <span class="loading-dot"></span>
        <span class="loading-text">正在翻阅...</span>
      </div>

      <div v-else-if="showAnswer" class="answer-content fade-in">
        <div class="answer-border">
          <p class="answer-label">今日启示</p>
          <p class="answer-text answer-zh">{{ answer_zh }}</p>
          <div class="answer-divider"></div>
          <p class="answer-text answer-en">{{ answer_en }}</p>
        </div>
        <button type="button" class="again-btn" @click="handleAgain">再来一次</button>
      </div>

      <div v-else class="empty-answer">
        <div class="book-cover">
          <div class="cover-overlay"></div>
          <button type="button" class="get-btn" @click="getAnswer">请翻开你的答案之书</button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.answerBook-container {
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.12);
  padding: 16px;
  border-radius: 14px;
  border: 6px solid transparent;
  border-image: url(../assets/border.png) 30 stretch;
  background-color: #fff;
  box-sizing: border-box;
}

:root.dark .answerBook-container {
  background-color: #1a1a1a;
  border: 1px solid #444;
  border-image: none;
}

.answerBook-container:hover {
  scale: 1.01;
}

.book-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-bottom: 12px;
  border-bottom: 2px dotted #e0d6f0;
  margin-bottom: 12px;
}

.book-icon {
  font-size: 26px;
  line-height: 1;
  flex-shrink: 0;
}

.book-meta {
  min-width: 0;
}

.title {
  font-size: 22px;
  font-weight: 700;
  color: #5a3d8c;
  line-height: 1.3;
}

.subtitle {
  margin: 4px 0 0;
  font-size: 13px;
  color: #8b7aa8;
  line-height: 1.4;
}

:root.dark .title {
  color: #c4b5e8;
}

:root.dark .subtitle {
  color: #9a8bb8;
}

:root.dark .book-header {
  border-bottom-color: #444;
}

.answer-box {
  min-height: 220px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.empty-answer {
  width: 100%;
}

.book-cover {
  position: relative;
  width: 100%;
  height: 260px;
  border-radius: 12px;
  overflow: hidden;
  background-image: url(../assets/answerBook.jpg);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  align-items: center;
  justify-content: center;
}

.cover-overlay {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    180deg,
    rgba(58, 38, 89, 0.15) 0%,
    rgba(58, 38, 89, 0.55) 100%
  );
}

.get-btn,
.again-btn {
  position: relative;
  z-index: 1;
  padding: 11px 22px;
  background-color: rgba(90, 61, 140, 0.92);
  color: #fff;
  border: none;
  border-radius: 999px;
  font-size: 14px;
  font-weight: 500;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease, background-color 0.2s ease;
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

.get-btn:hover,
.again-btn:hover {
  background-color: rgba(76, 50, 120, 1);
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(90, 61, 140, 0.35);
}

.loading-state {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 24px 0;
  color: #6b538c;
  font-size: 14px;
}

.loading-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background-color: #8e71c7;
  animation: pulse 1.2s infinite ease-in-out both;
}

.loading-dot:nth-child(2) { animation-delay: 0.15s; }
.loading-dot:nth-child(3) { animation-delay: 0.3s; }

@keyframes pulse {
  0%, 80%, 100% { opacity: 0.3; transform: scale(0.85); }
  40% { opacity: 1; transform: scale(1); }
}

.answer-content {
  width: 100%;
  text-align: center;
  opacity: 0;
  transform: translateY(16px);
  transition: opacity 0.5s ease, transform 0.5s ease;
}

.fade-in {
  opacity: 1;
  transform: translateY(0);
}

.answer-border {
  padding: 18px 16px;
  background: linear-gradient(145deg, #faf7ff 0%, #f3ecff 100%);
  border-radius: 12px;
  border: 1px solid #e8dff5;
  box-shadow:
    inset 0 0 20px rgba(142, 113, 199, 0.08),
    0 2px 8px rgba(0, 0, 0, 0.04);
}

:root.dark .answer-border {
  background: linear-gradient(145deg, #2a2438 0%, #1f1a2e 100%);
  border-color: #443a5c;
}

.answer-label {
  margin: 0 0 12px;
  font-size: 12px;
  letter-spacing: 0.12em;
  color: #8e71c7;
  text-transform: uppercase;
}

.answer-text {
  margin: 0;
  line-height: 1.7;
  font-style: italic;
}

.answer-zh {
  font-size: 19px;
  color: #3a2659;
  font-weight: 500;
}

.answer-en {
  font-size: 15px;
  color: #6b538c;
  opacity: 0.92;
}

:root.dark .answer-zh {
  color: #e8dff5;
}

:root.dark .answer-en {
  color: #b8a8d4;
}

.answer-divider {
  width: 40px;
  height: 2px;
  margin: 14px auto;
  background: linear-gradient(90deg, transparent, #c4b5e8, transparent);
  border-radius: 1px;
}

.again-btn {
  margin-top: 16px;
  background-color: #fff;
  color: #5a3d8c;
  border: 1px solid #d8ccef;
}

.again-btn:hover {
  background-color: #f9f5ff;
  color: #4a3278;
  box-shadow: 0 4px 12px rgba(90, 61, 140, 0.15);
}

:root.dark .again-btn {
  background-color: #2a2438;
  color: #c4b5e8;
  border-color: #443a5c;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .answerBook-container {
    max-width: 100%;
    padding: 14px;
    border: none;
    border-image: none;
    border-radius: 14px;
    box-shadow: 0 2px 10px rgba(90, 61, 140, 0.08);
    background: linear-gradient(180deg, #fff 0%, #faf8ff 100%);
  }

  :root.dark .answerBook-container {
    background: linear-gradient(180deg, #1e1e1e 0%, #1a1724 100%);
    box-shadow: 0 2px 10px rgba(0, 0, 0, 0.25);
  }

  .answerBook-container:hover {
    scale: 1;
  }

  .book-header {
    gap: 12px;
    padding-bottom: 10px;
    margin-bottom: 10px;
    border-bottom: 1px solid #ede8f5;
  }

  :root.dark .book-header {
    border-bottom-color: #333;
  }

  .book-icon {
    font-size: 22px;
  }

  .title {
    font-size: 17px;
    font-weight: 600;
  }

  .subtitle {
    font-size: 12px;
    margin-top: 2px;
  }

  .answer-box {
    min-height: auto;
  }

  .book-cover {
    height: 180px;
    border-radius: 10px;
  }

  .get-btn {
    padding: 10px 18px;
    font-size: 13px;
    background-color: rgba(255, 255, 255, 0.95);
    color: #5a3d8c;
    box-shadow: 0 4px 14px rgba(0, 0, 0, 0.12);
  }

  .get-btn:hover {
    background-color: #fff;
    color: #4a3278;
  }

  .answer-border {
    padding: 16px 14px;
    border-radius: 10px;
  }

  .answer-label {
    font-size: 11px;
    margin-bottom: 10px;
  }

  .answer-zh {
    font-size: 17px;
    line-height: 1.65;
  }

  .answer-en {
    font-size: 14px;
    line-height: 1.55;
  }

  .answer-divider {
    margin: 12px auto;
  }

  .again-btn {
    width: 100%;
    max-width: 200px;
    margin-top: 14px;
    padding: 10px 16px;
    font-size: 13px;
    background-color: #f3ecff;
    color: #5a3d8c;
    border: 1px solid #e0d4f5;
  }

  .again-btn:hover {
    background-color: #ebe0ff;
  }

  :root.dark .again-btn {
    background-color: #2a2438;
    border-color: #443a5c;
  }
}
</style>
