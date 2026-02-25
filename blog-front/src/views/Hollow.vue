<script setup lang="ts">
import { getHollow,addHollow } from '../api/hollow';
import AnswerBook from '../components/AnswerBook.vue';
import Mine from '../components/Mine.vue';
import Type from '../components/Type.vue';
import { onMounted, ref } from 'vue';
import { formatTime } from '../utils/format';
import {type ApiResponse } from '../types/articles'
interface hollow{
    date:string,
    word:string
}
const total=ref<number>(0)
const hollows=ref<hollow[]>()
const isEdit = ref<boolean>(false)
const word=ref<string>('')
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
  word.value=''
}
const subbmitHollow=async()=>{
    if (!word.value.trim()) {
        alert('输入内容不能为空')
    }
    const now = new Date()
    const year = now.getFullYear()
    const month = String(now.getMonth() + 1).padStart(2, '0')
    const day = String(now.getDate()).padStart(2, '0')
    const hours = String(now.getHours()).padStart(2, '0')
    const minutes = String(now.getMinutes()).padStart(2, '0')
    const seconds = String(now.getSeconds()).padStart(2, '0')
    const submitDate = `${year}-${month}-${day}T${hours}:${minutes}:${seconds}`
    const res=await addHollow(word.value,submitDate) as unknown as ApiResponse
    word.value=''
    if (res?.code==200) {
        isEdit.value=false
        location.reload()
    }
}
onMounted(async()=>{
    let res=await getHollow()
    hollows.value=res.data
    hollows.value=hollows.value?.reverse()
    total.value=res.data.length
})
</script>

<template>
    <div class="hollow-page">
        <div class="left">
            <div class="title">闲言碎语</div>
            <img 
                src="../assets/edit.png" 
                class="edit-icon" 
                @click="changeEdit"
                alt="编辑"
            >
            <div class="total">共{{ total }}条</div>
            <div class="word-container">
                <div v-for="hollow in hollows" class="word-box">
                    <div>
                        <span>◆</span>
                        <span class="time">{{formatTime(hollow.date)}}</span>
                    </div>
                    <div class="word">{{ hollow.word }}</div>
                </div>
            </div>
        </div>
        <div class="right">
            <Mine></Mine>
            <Type></Type>
            <AnswerBook></AnswerBook>
        </div>
        <!-- 蒙层组件 -->
         <div class="overlay" v-show="isEdit" @click.self="changeEdit">
            <div class="overlay-box">
                <div class="overlay-title">今日感想</div>
                <textarea placeholder="请输入今日感想" class="input-area" v-model="word"></textarea>
                <div class="overlay-footer">
                    <el-button type="primary" @click="subbmitHollow">提交</el-button>
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
    display: flex;
    justify-content: space-around;
    margin: 0 150px;
    padding: 20px;
}
.left {
    width: 75%;
    position: relative;
    background-color: #fff;
    padding: 20px;
    border-radius: 10px;
    box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
}
:root.dark .left{
    background-color: #1A1A1A;
    border:1px solid #fff;
}
.right {
    width: 25%;
    margin-left: 20px;
    display: flex;
    flex-direction: column;
    gap: 20px;
}
.title{
    font-size: 25px;
    margin: 10px 0;
}
.total{
    margin: 10px ;
}
.edit-icon{
    position: absolute;
    right: 20px;
    top: 20px;
}
.word-container{
    display: flex;
    flex-direction: column;
    gap: 5px;
}
.time{
    font-size: 15px;
}
.word{
    border-left: 3px solid #333;
    min-height: 50px;
    display: flex;
    align-items: center;
    padding: 10px;
    margin: 0 6px;
    font-size: 19px;
}
:root.dark .word {
    border-color: #fff;
}
.overlay{
    position:fixed;
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
.overlay-box{
    background-color: #fff;
    width: 400px;
    height: 250px;
    padding: 20px;
    border-radius: 10px;
}
.overlay-title{
    font-size: 22px;
    color: #333;
}
.input-area{
    width: 400px;
    height: 160px;
    resize: none;
    font-size: 18px;
    margin: 10px 0;
}
.overlay-footer{
    display: flex;
    align-items: center;
    justify-content: space-between;
}
</style>
