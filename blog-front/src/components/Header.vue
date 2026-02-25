<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { getMusic, getWord } from '../api/other'
import Nav from './Nav.vue'
const isScroll = ref(false)
const lastScrollTop = ref(0)
interface musicItem {
    name: string,
    image: string,
    url: string,
    singer: string
}
const word = ref<string>('')
const music = ref<musicItem>()
const handleScroll = () => {
    const currentScrollTop = window.scrollY
    if (currentScrollTop > 100) {
        if (currentScrollTop > lastScrollTop.value) {
            isScroll.value = true
        } else {
            isScroll.value = false
        }
    } else {
        isScroll.value = false
    }
    lastScrollTop.value = currentScrollTop
}

onMounted(async () => {

    window.addEventListener('scroll', handleScroll)
    try {
        const wordRes = await getWord()
        word.value = wordRes.data
        const musicRes = await getMusic()
        music.value = musicRes.data
    } catch (error) {
        console.error('请求失败', error)
    }
})
onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
})
</script>

<template>
    <div class="header-container">
        <Nav class="nav-wrapper" :class="{ 'nav-active': isScroll }"></Nav>
        <div class="title-wrapper">
            <div class="title">下雨天的个人博客</div>
            <div class="intro">一个既不热爱技术也不热爱生活的个人博客</div>
            <div class="word">随机一言：{{ word }}</div>
            <div class="music" v-if="music?.image">
                <img :src="music.image" alt="专辑封面" class="music-cover">
                <div class="music-info">
                    <div class="music-name">{{ music.name }}</div>
                    <div class="music-artist">{{ music.singer }}</div>
                    <audio :src="music.url" controls class="audio-player" preload="metadata">
                        您的浏览器不支持音频播放
                    </audio>
                </div>
            </div>
        </div>
    </div>
</template>

<style scoped>
.header-container {
    position: relative;
    width: 100%;
    height: 700px;
    background-image: url(../assets/bg.jpg);
    background-size: cover;
    background-position: center;
    overflow: hidden;
}

.nav-wrapper {
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 100;
    transform: translateY(-100%);
    opacity: 0;
    transition: all 0.3s ease;
}

.nav-active {
    transform: translateY(0);
    opacity: 1;
}

.title-wrapper {
    position: relative;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    width: 100%;
    height: 100%;
    color: #fff;
}

.title {
    font-size: 36px;
    font-weight: bold;
    overflow: hidden;
    width: 0;
    height: 50px;
    white-space: nowrap;
    border-right: 3px solid #ffffff;
    animation:
        write 3.5s steps(8) forwards,
        blink 0.7s infinite;
}

.intro {
    font-size: 24px;
    margin: 20px;
}

/* 每日鸡汤样式 */
.word {
    margin-top: 30px;
    padding: 15px 25px;
    max-width: 800px;
    width: 90%;
    text-align: center;
    background-color: rgba(0, 0, 0, 0.3);
    border-radius: 12px;
    backdrop-filter: blur(5px);
    border: 1px solid rgba(255, 255, 255, 0.1);
    /* 文字样式 */
    font-size: 18px;
    line-height: 1.5;
    color: #fff;
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
    /* 动画效果 */
    opacity: 0;
    transform: translateY(20px);
    animation: fadeIn 1s ease-out 3.5s forwards;
}

/* 引用符号装饰 */
.word::before,
.word::after {
    font-size: 40px;
    color: rgba(255, 255, 255, 0.6);
}

.word::before {
    content: "“";
    margin-right: 10px;
    vertical-align: -0.5em;
}

.word::after {
    content: "”";
    margin-left: 10px;
    vertical-align: -0.5em;
}

/* 淡入动画 */
@keyframes fadeIn {
    to {
        opacity: 1;
        transform: translateY(0);
    }
}


@keyframes write {
    0% {
        width: 0;
    }

    100% {
        width: 290px;
    }
}

@keyframes blink {

    0%,
    100% {
        border-color: #ffffff;
    }

    50% {
        border-color: transparent;
    }
}

.music {
    margin-top: 30px;
    display: flex;
    align-items: center;
    gap: 15px;
    padding: 15px;
    background-color: rgba(0, 0, 0, 0.2);
    border-radius: 10px;
    width: 90%;
    max-width: 600px;
}

.music-cover {
    width: 100px;
    height: 100px;
    border-radius: 8px;
    object-fit: cover;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.music-info {
    flex: 1;
    color: #fff;
}

.music-name {
    font-size: 18px;
    margin-bottom: 8px;
}

.music-artist {
    font-size: 14px;
    opacity: 0.8;
    margin-bottom: 12px;
}

.audio-player {
    width: 100%;
    background-color: rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    padding: 5px;
}
</style>