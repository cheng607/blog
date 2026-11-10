<script setup lang="ts">
import { ElIcon } from 'element-plus'
import 'element-plus/es/components/icon/style/css'
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { ArrowDown, ArrowRight } from '@element-plus/icons-vue'

export interface TocItem {
    id: string
    text: string
    level: number
}

export interface TocNode extends TocItem {
    children: TocNode[]
}

const props = defineProps({
    toc: {
        type: Array as () => TocItem[],
        default: () => [],
    },
})

const activeId = ref('')
const expandedIds = ref<Set<string>>(new Set())

function buildTocTree(items: TocItem[]): TocNode[] {
    const root: TocNode[] = []
    const stack: TocNode[] = []
    for (const item of items) {
        const node: TocNode = { ...item, children: [] }
        while (stack.length > 0) {
            const top = stack[stack.length - 1]
            if (!top || top.level < item.level) break
            stack.pop()
        }
        const parent = stack[stack.length - 1]
        if (!parent) {
            root.push(node)
        } else {
            parent.children.push(node)
        }
        stack.push(node)
    }
    return root
}

const tocTree = computed(() => buildTocTree(props.toc || []))

function collectAncestors(nodes: TocNode[], targetId: string, path: string[] = []): string[] | null {
    for (const n of nodes) {
        const next = [...path, n.id]
        if (n.id === targetId) return path
        if (n.children.length) {
            const found = collectAncestors(n.children, targetId, next)
            if (found) return found
        }
    }
    return null
}

/** 默认展开到 h2，避免目录过长时一打开就占满屏 */
function initExpanded() {
    const set = new Set<string>()
    const walk = (nodes: TocNode[]) => {
        for (const n of nodes) {
            if (n.children.length > 0 && n.level < 3) {
                set.add(n.id)
            }
            walk(n.children)
        }
    }
    walk(tocTree.value)
    expandedIds.value = set
}

function toggleExpand(id: string) {
    const next = new Set(expandedIds.value)
    if (next.has(id)) next.delete(id)
    else next.add(id)
    expandedIds.value = next
}

function isExpanded(id: string) {
    return expandedIds.value.has(id)
}

interface FlatTocRow {
    node: TocNode
    depth: number
    hasChildren: boolean
}

function flattenVisible(nodes: TocNode[], depth = 0, out: FlatTocRow[] = []): FlatTocRow[] {
    for (const n of nodes) {
        const hasChildren = n.children.length > 0
        out.push({ node: n, depth, hasChildren })
        if (hasChildren && isExpanded(n.id)) {
            flattenVisible(n.children, depth + 1, out)
        }
    }
    return out
}

const visibleRows = computed(() => flattenVisible(tocTree.value))

const TOP_OFFSET = 100

const onScroll = () => {
    if (!props.toc?.length) return
    let current = ''
    for (const item of props.toc) {
        const el = document.getElementById(item.id)
        if (!el) continue
        if (el.getBoundingClientRect().top <= TOP_OFFSET) {
            current = item.id
        }
    }
    activeId.value = current
}

const scrollTo = (id: string) => {
    const el = document.getElementById(id)
    if (!el) return
    const rect = el.getBoundingClientRect()
    const target = window.scrollY + rect.top - TOP_OFFSET
    window.scrollTo({ top: Math.max(0, target), behavior: 'smooth' })
}

watch(
    () => props.toc,
    () => {
        initExpanded()
        onScroll()
    },
    { deep: true }
)

watch(activeId, (id) => {
    if (!id) return
    const ancestors = collectAncestors(tocTree.value, id)
    if (!ancestors?.length) return
    const next = new Set(expandedIds.value)
    ancestors.forEach((aid) => next.add(aid))
    expandedIds.value = next
})

onMounted(() => {
    initExpanded()
    window.addEventListener('scroll', onScroll, { passive: true })
    onScroll()
})

onUnmounted(() => {
    window.removeEventListener('scroll', onScroll)
})
</script>

<template>
    <aside class="sidebar-toc">
        <div class="toc-box" v-if="props.toc && props.toc.length">
            <div class="toc-header">
                <span class="toc-title">目录</span>
            </div>
            <div class="toc-list-wrap">
                <ul class="toc-list">
                    <li
                        v-for="row in visibleRows"
                        :key="row.node.id"
                        :class="[
                            'toc-item',
                            'level-' + row.node.level,
                            { active: activeId === row.node.id },
                        ]"
                        :style="{ paddingLeft: `${10 + row.depth * 14}px` }"
                    >
                        <button
                            v-if="row.hasChildren"
                            type="button"
                            class="toc-toggle"
                            :aria-label="isExpanded(row.node.id) ? '收起' : '展开'"
                            @click.stop="toggleExpand(row.node.id)"
                        >
                            <el-icon :size="16">
                                <ArrowDown v-if="isExpanded(row.node.id)" />
                                <ArrowRight v-else />
                            </el-icon>
                        </button>
                        <span v-else class="toc-toggle-placeholder" />
                        <a
                            class="toc-link"
                            :title="row.node.text"
                            @click.prevent="scrollTo(row.node.id)"
                        >{{ row.node.text }}</a>
                    </li>
                </ul>
            </div>
        </div>
        <div class="toc-empty" v-else>本文暂无目录</div>
    </aside>
</template>

<style scoped>
.sidebar-toc {
    position: sticky;
    top: 80px;
    max-height: calc(100vh - 100px);
    display: flex;
    flex-direction: column;
}

.toc-box {
    background-color: #fff;
    padding: 12px 12px 10px;
    border-radius: 8px;
    box-shadow: 0 1px 4px rgba(0, 0, 0, 0.08);
    display: flex;
    flex-direction: column;
    min-height: 0;
    max-height: inherit;
}

:root.dark .toc-box {
    background-color: #1a1a1a;
    border: 1px solid #444;
}

.toc-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 8px;
    margin-bottom: 8px;
    flex-shrink: 0;
}

.toc-title {
    font-weight: 600;
    font-size: 17px;
    color: var(--text-color, #333);
}

.toc-list-wrap {
    flex: 1;
    min-height: 0;
    max-height: min(70vh, calc(100vh - 180px));
    overflow-y: auto;
    overflow-x: hidden;
    margin-right: -4px;
    padding-right: 4px;
}

.toc-list-wrap::-webkit-scrollbar {
    width: 6px;
}

.toc-list-wrap::-webkit-scrollbar-thumb {
    background: rgba(73, 177, 245, 0.35);
    border-radius: 3px;
}

.toc-list-wrap::-webkit-scrollbar-track {
    background: transparent;
}

.toc-list {
    list-style: none;
    padding: 0;
    margin: 0;
    display: flex;
    flex-direction: column;
    gap: 2px;
}

.toc-item {
    display: flex;
    align-items: flex-start;
    gap: 4px;
    line-height: 1.45;
    border-radius: 4px;
    transition: background-color 0.15s;
}

.toc-item:hover {
    background-color: rgba(73, 177, 245, 0.08);
}

:root.dark .toc-item:hover {
    background-color: rgba(127, 200, 248, 0.12);
}

.toc-toggle {
    flex-shrink: 0;
    width: 20px;
    height: 22px;
    padding: 0;
    border: none;
    background: none;
    cursor: pointer;
    color: #888;
    display: flex;
    align-items: center;
    justify-content: center;
}

.toc-toggle:hover {
    color: #49b1f5;
}

.toc-toggle-placeholder {
    flex-shrink: 0;
    width: 20px;
}

.toc-link {
    flex: 1;
    min-width: 0;
    color: #333;
    text-decoration: none;
    cursor: pointer;
    font-size: 15px;
    line-height: 1.5;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    word-break: break-word;
}

.toc-item.level-1 .toc-link {
    font-size: 16px;
    font-weight: 500;
}

.toc-item.level-2 .toc-link {
    font-size: 15px;
}

.toc-item.level-3 .toc-link,
.toc-item.level-4 .toc-link,
.toc-item.level-5 .toc-link,
.toc-item.level-6 .toc-link {
    font-size: 14px;
}

:root.dark .toc-link {
    color: #e8e8e8;
}

.toc-item.active .toc-link {
    color: #49b1f5;
    font-weight: 600;
}

:root.dark .toc-item.active .toc-link {
    color: #7fc8f8;
}

.toc-empty {
    color: #888;
    font-size: 13px;
    padding: 12px;
}

/* 抽屉内或非 sticky 场景 */
:global(.toc-drawer) .sidebar-toc,
:global(.el-drawer__body) .sidebar-toc {
    position: static;
    top: auto;
    max-height: none;
}

:global(.toc-drawer) .toc-list-wrap,
:global(.el-drawer__body) .toc-list-wrap {
    max-height: calc(100vh - 160px);
}

/* 目录抽屉移动端优化 */
@media (max-width: 768px) {
  :global(.toc-drawer) .toc-title {
    font-size: 16px !important;
  }

  :global(.toc-drawer) .toc-item {
    padding: 6px 16px !important;
  }

  :global(.toc-drawer) .toc-link {
    font-size: 15px !important;
    line-height: 1.5 !important;
  }

  /* 当前选中项左侧蓝色竖线高亮 */
  :global(.toc-drawer) .toc-item.active {
    position: relative;
  }

  :global(.toc-drawer) .toc-item.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 6px;
    bottom: 6px;
    width: 3px;
    background-color: #49b1f5;
    border-radius: 2px;
  }
}
</style>
