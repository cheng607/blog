<script setup lang="ts">
import { ref, onMounted } from 'vue'
interface activity {
    content: string,
    timestamp: string
}
const activities = ref<activity[]>()
onMounted(() => {
    activities.value = [
        {
            content: '出生',
            timestamp: '2005-08-01'
        },
        {
            content: '进入大学生活',
            timestamp: '2023-09-01'
        },
        {
            content: '实现一个自己的技术博客',
            timestamp: '2025-11'
        },
        {
            content: '网站正式通过备案',
            timestamp: '2025-12-25（圣诞节）'
        },
    ]
})
</script>

<template>
    <div class="about-page">
        <article class="word">
            在我看来，作为一名程序员，必须要有一个属于自己的技术博客。没别的原因，单纯觉得挺酷的😎。<br><br>
            于是乎在2025年11月初便开始了个人技术博客的搭建，从最开始的前端搭建，到后来的后端基础学习，再到后面的服务器，域名的使用，每一步都出乎我自己的意料，只有当自己真正的开始，才会发现很多东西并没有自己想象的那么困难，只有迈出那一步才能自己行不行。<br><br>
            从最初自己考虑是从0开始自己手搓一个个人博客还是用Hexo这些框架，到自己独自进行ICP备案和公安备案，最难的从来不是技术，而是个人的懒惰和爱玩，总是不能静下心来认真弄个东西，不是以最近课多没时间为借口，就是认为技术太垃圾了不愿意上线。<br><br>
            这个博客用于记录我的学习历程、技术心得和项目经验，希望能为同样在技术道路上探索的同学提供一些参考和帮助。<br><br>
            博客内容主要围绕前端开发技术和个人所学的一些有趣的小知识，包括但不限于Vue、React、HTML等技术的实践与思考，同时也会有个人的树洞用于记录自己的生活和偶尔的生活感想的感想。<br><br>
            前端技术栈：Vue3+Element Plus+Vditor<br><br>
            后端技术栈：SpringBoot
        </article>
        <div class="timeline">
            <div class="title">个人时间轴</div>
            <el-timeline>
                <el-timeline-item v-for="(activity, index) in activities" type="primary" size="large" hollow=true
                    :key="index" placement="top" :timestamp="activity.timestamp">
                    <span>{{ activity.content }}</span>
                </el-timeline-item>
            </el-timeline>
        </div>
    </div>
</template>

<style scoped>
.about-page {
    width: 100%;
    min-height: 400px;
    display: flex;
    gap: 50px;
    justify-content: center;
    padding: 30px 0;
}

.word {
    width: 700px;
    background-color: #fff;
    padding: 30px 50px;
    border-radius: 30px;
    font-size: 18px;
    color: #333;
}

.timeline {
    width: 300px;
    background-color: #fff;
    padding: 10px 20px;
    border-radius: 20px;
}

.timeline .title {
    margin: 20px;
    font-size: 20px;
    text-align: center;
    color: #333;
}

.timeline span {
    font-size: 18px;
}
:root.dark .word{
    background-color: #353739;
    color: #fff;
}
:root.dark .timeline{
    background-color: #353739;
}
:root.dark .timeline span{
    color: #fff;
}
</style>
