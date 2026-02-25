<script setup lang="ts">
import Footer from './components/Footer.vue';
import Header from './components/Header.vue';
import Agent from './components/Agent.vue';
import { onMounted, onUnmounted } from 'vue';
// 拖尾特效相关
let lastTime = 0;
const TRAIL_INTERVAL = 20;
const trailColors = [
  'rgba(73, 177, 245, 0.8)',
  'rgba(127, 200, 248, 0.6)',
  'rgba(76, 50, 120, 0.7)',
  'rgba(255, 255, 255, 0.5)'
];

const createTrail = (x: number, y: number) => {
  const container = document.getElementById('trail-container');
  if (!container) return;

  // 创建拖尾粒子
  const trail = document.createElement('div');
  const size = Math.random() * 10 + 5; // 5-15px的随机大小
  const color = trailColors[Math.floor(Math.random() * trailColors.length)];

  // 基础样式
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
  `;

  container.appendChild(trail);

  // 动画效果 - 结合现有项目中的动画风格
  const duration = Math.random() * 1000 + 5000; // 1s-5s动画时长
  trail.animate([
    { opacity: 1, transform: `translate(-50%, -50%) scale(1)` },
    { opacity: 0, transform: `translate(-50%, -50%) scale(0.2) translateY(-${Math.random() * 30 + 20}px)` }
  ], {
    duration,
    easing: 'cubic-bezier(0.2, 0.8, 0.2, 1)'
  });

  setTimeout(() => {
    trail.remove();
  }, duration);
};

// 鼠标移动监听
const handleMouseMove = (e: MouseEvent) => {
  const now = Date.now();
  if (now - lastTime > TRAIL_INTERVAL) {
    createTrail(e.clientX, e.clientY);
    lastTime = now;
  }
};

onMounted(() => {
  window.addEventListener('mousemove', handleMouseMove);
});

onUnmounted(() => {
  window.removeEventListener('mousemove', handleMouseMove);
});
</script>

<template>
  <div class="main">
    <Header></Header>
    <router-view></router-view>
    <Footer></Footer>
    <div id="trail-container"></div>
    <Agent class="agent"></Agent>
    <el-backtop :right="100" :bottom="100" class="backtop-box">
      <img src="./assets/backTop.png" alt="返回顶部">
    </el-backtop>
  </div>
</template>

<style scoped>
.app-container {
  width: 100%;
  min-height: 100vh;
  margin: 0 auto;
}

.backtop-box {
  padding: 5px;
  box-shadow: 0 3px 10px rgba(0, 0, 0, 0.2);
}

.agent {
  position: fixed;
  left: 30px;
  bottom: 100px;
  z-index: 99;
  transition: transform 0.3s ease;
}

.agent.hidden {
  transform: translateX(-120%);
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
