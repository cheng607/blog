<script setup lang="ts">
import {  ref, onMounted } from 'vue'
// 接收来自父组件的 TOC 数据：数组项包含 id,text,level
const props = defineProps({
  toc: {
    type: Array<{ id: string; text: string; level: number }>,
    default: () => []
  }
})

const activeId = ref('')

// 监听滚动以高亮当前章节（可选）
const onScroll = () => {
  if (!props.toc || props.toc.length === 0) return
  let current = ''
  for (const item of props.toc) {
    const el = document.getElementById(item.id)
    if (!el) continue
    const rect = el.getBoundingClientRect()
    if (rect.top <= 120) {
      current = item.id
    }
  }
  activeId.value = current
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})

// 点击目录项时平滑滚动（带顶部偏移）
const TOP_OFFSET = 100 // 与 onScroll 中使用的阈值一致，保证标题不会紧贴顶部
const scrollTo = (id: string) => {
  const el = document.getElementById(id)
  if (!el) return
  const rect = el.getBoundingClientRect()
  const target = window.scrollY + rect.top - TOP_OFFSET
  window.scrollTo({ top: Math.max(0, target), behavior: 'smooth' })
}
</script>

<template>
  <aside class="sidebar-toc">
    <div class="toc-box" v-if="props.toc && props.toc.length">
      <div class="toc-title">目录</div>
      <ul class="toc-list">
        <li v-for="item in props.toc" :key="item.id" :class="['toc-item', 'level-' + item.level, { active: activeId === item.id }]">
          <a @click.prevent="scrollTo(item.id)">{{ item.text }}</a>
        </li>
      </ul>
    </div>
    <div class="toc-empty" v-else>本文暂无目录</div>
  </aside>
</template>

<style scoped>
.sidebar-toc {
    position: sticky;
    top: 80px;
}

.toc-box {
    background-color: #fff;
    padding: 15px;
    border-radius: 6px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.04);
}
:root.dark .toc-box{
    background-color:#1A1A1A;
    border: 1px solid #fff;
}
.toc-title {
    font-weight: 600;
    margin-bottom: 8px;
}

.toc-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 6px;
}

.toc-item a {
    color: #333;
    text-decoration: none;
    cursor: pointer;
}
:root.dark .toc-item a{
    color: #fff;
}
.toc-item.level-1 {
    padding-left: 0;
}

.toc-item.level-2 {
    padding-left: 10px;
}

.toc-item.level-3 {
    padding-left: 20px;
}

.toc-item.level-4 {
    padding-left: 30px;
}

.toc-item.level-5 {
    padding-left: 40px;
}

.toc-item.level-6 {
    padding-left: 50px;
}

.toc-item.active a,
:root.dark .toc-item.active a{
    color: #49B1F5;
}

.toc-empty {
    color: #888;
    font-size: 13px;
}
</style>
