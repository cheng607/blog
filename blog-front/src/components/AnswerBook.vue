<script setup lang="ts">
import { ref } from 'vue'
import axios from 'axios'

const answer_zh = ref<string>('')
const answer_en = ref<string>('')
const showAnswer = ref<boolean>(false)

// 获取答案的方法
const getAnswer = async () => {
  showAnswer.value = false
  try {
    const res = await axios.get(`https://api.4qb.cn/api/daan_book`)
    answer_zh.value = res.data.data.zh
    answer_en.value = res.data.data.en
    console.log(res.data)
    if (res.data.data) {
      // showAnswer.value = true
      setTimeout(() => {
        showAnswer.value = true
      }, 1000)
    }
  } catch (error) {
    console.error('请求失败', error)
    answer_zh.value = '暂时无法回答，请稍后再试'
    answer_en.value = 'Unable to answer at the moment, please try again later'
    showAnswer.value = true
  }
}
const handleAgain = () => {
  showAnswer.value = false
}
</script>

<template>
  <div class="answerBook-container">
    <div class="title">答案之书</div>
    <div class="answer-box">
      <div v-if="showAnswer" class="answer-content" :class="{ 'fade-in': showAnswer }">
        <div class="answer-border">
          <p class="answer-text">{{ answer_zh }}</p>
          <p class="answer-text">{{ answer_en }}</p>
        </div>
        <el-button class="again-btn" @click="handleAgain">再来一次</el-button>
      </div>
      <div v-else class="empty-answer">
        <el-button class="get-btn" @click="getAnswer">请翻开你的答案之书</el-button>
      </div>
    </div>
  </div>
</template>

<style scoped>
.answerBook-container {
  width: 100%;
  max-width: 500px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
  padding: 10px;
  border-radius: 10px;
  border: 8px solid transparent;
  border-image: url(../assets/border.png) 30 stretch;
  background-color: #fff;
}

:root.dark .answerBook-container {
  background-color: #1A1A1A;
  border: 1px solid #fff;
}

.answerBook-container:hover {
  scale: 1.01;
}

.title {
  text-align: center;
  font-size: 28px;
  font-weight: bold;
  color: #5a3d8c;
  padding-bottom: 10px;
  border-bottom: 2px dotted #e0d6f0;
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}

.answer-box {
  min-height: 250px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 10px;
  position: relative;
}

.empty-answer {
  width: 100%;
  height: 300px;
  background-image: url(../assets/answerBook.jpg);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 自定义按钮样式 */
.get-btn,
.again-btn {
  padding: 10px;
  background-color: rgba(90, 61, 140, 0.9);
  color: white;
  border: none;
  border-radius: 30px;
  font-size: 14px;
  transition: all 0.3s ease;
  margin-top: 20px;
}

.get-btn {
  margin-top: 100px;
}

.get-btn:hover,
.again-btn:hover {
  background-color: rgba(76, 50, 120, 1);
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.25);
}

.answer-content {
  width: 100%;
  text-align: center;
  opacity: 0;
  transform: translateY(20px);
  transition: opacity 0.6s ease, transform 0.6s ease;
}

.fade-in {
  opacity: 1;
  transform: translateY(0);
}

.answer-border {
  padding: 10px;
  background-color: #f9f5ff;
  border-radius: 10px;
  box-shadow:
    inset 0 0 15px rgba(142, 113, 199, 0.15),
    0 2px 5px rgba(0, 0, 0, 0.05);
}

.answer-text:first-child {
  font-size: 20px;
  line-height: 1.8;
  color: #3a2659;
  font-style: italic;
  margin-bottom: 15px;
}

/* 英文答案样式 */
.answer-text:last-child {
  font-size: 16px;
  line-height: 1.6;
  color: #6b538c;
  font-style: italic;
  opacity: 0.9;
}
</style>