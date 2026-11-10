<script setup lang="ts">
import { ElButton } from 'element-plus'
import 'element-plus/es/components/button/style/css'
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { getSiteSummary } from '../api/site'
import { resolveMediaUrl } from '../utils/media'
import defaultAvatar from '../assets/loading.gif'
interface user {
    username: string,
    avatar: string,
    signature: string
}
const ArticleCount = ref<number>(0)
const HollowCount = ref<number>(0)
const TypeCount = ref<number>(0)
const UserData = ref<user>()
const error = ref('')
const loadSummary = async () => {
    error.value = ''
    try {
        const res = await getSiteSummary()
        UserData.value = res.data?.profile ?? undefined
        ArticleCount.value = res.data?.articleCount ?? 0
        HollowCount.value = res.data?.hollowCount ?? 0
        TypeCount.value = res.data?.tagCount ?? 0
    } catch (e) { error.value = e instanceof Error ? e.message : '加载统计失败' }
}
onMounted(loadSummary)
const router = useRouter()
const handleLink = () => {
    router.push('/about')
}

const avatarSrc = () => {
    const url = resolveMediaUrl(UserData.value?.avatar)
    return url || defaultAvatar
}
</script>

<template>
    <div class="mine-container">
        <img :src="avatarSrc()" class="avatar" alt="头像">
        <div class="name">{{ UserData?.username || '博主尚未设置资料' }}</div>
        <div class="note">个性签名：{{ UserData?.signature }}</div>
        <div v-if="error" role="alert">{{ error }} <el-button text @click="loadSummary">重试</el-button></div>
        <div v-else class="summary">
            <div class="item">
                <div class="target">文章</div>
                <div class="num">{{ ArticleCount }}</div>
            </div>
            <div class="item">
                <div class="target">树洞</div>
                <div class="num">{{ HollowCount }}</div>
            </div>
            <div class="item">
                <div class="target">分类</div>
                <div class="num">{{ TypeCount }}</div>
            </div>
        </div>
        <el-button type="primary" class="link" @click="handleLink">查看个人</el-button>
        <div class="ports">
            <a href="https://github.com/cheng607" target="_blank">
                <img src="../assets/github.png">
            </a>
            <a href="https://gitee.com/cheng-zhengren" target="_blank">
                <img src="../assets/gitee.png">
            </a>
            <a href="https://xiaohome.oss-cn-wuhan-lr.aliyuncs.com/qrcode.png" target="_blank">
                <img src="../assets/qq.png">
            </a>
        </div>
    </div>
</template>

<style scoped>
.mine-container {
    width: 100%;
    height: 420px;
    display: flex;
    flex-direction: column;
    align-items: center;
    box-shadow: 0 3px 8px rgba(0, 0, 0, 0.2);
    padding: 15px;
    border-radius: 10px;
    color: #444;
    background-color: #fff;
}

/* 移动端紧凑优化 */
@media (max-width: 768px) {
  .mine-container {
    height: auto !important;
    padding: 16px 14px !important;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.08) !important;
  }

  .avatar {
    width: 72px !important;
    height: 72px !important;
    margin-bottom: 8px !important;
  }

  .name {
    font-size: 17px !important;
    margin-bottom: 4px !important;
  }

  .note {
    font-size: 13px !important;
    margin-bottom: 10px !important;
    line-height: 1.4 !important;
  }

  .summary {
    margin-bottom: 12px !important;
    gap: 16px !important;
  }

  .item .target {
    font-size: 12px !important;
  }

  .item .num {
    font-size: 16px !important;
  }

  .link {
    height: 36px !important;
    font-size: 14px !important;
    margin-bottom: 12px !important;
  }

  .ports {
    gap: 16px !important;
  }

  .ports img {
    width: 26px !important;
    height: 26px !important;
  }
}

.mine-container:hover {
    scale: 1.01;
}

.avatar {
    width: 150px;
    height: 150px;
    border-radius: 50%;
    border: 1px solid #333;
}

.name {
    font-size: 24px;
    font-weight: bold;
}

.note {
    text-align: center;
    margin: 10px 0;
}

.summary {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-around;
    text-align: center;
    margin: 10px 0;
}

.link {
    padding: 20px 70px;
    margin: 10px 0;
}

.ports {
    display: flex;
    align-items: center;
    gap: 10px;
}

.ports img {
    width: 30px;
    height: 30px;
    padding: 10px;
    border-radius: 50%;
    background-color: #fff;
}

.ports img:hover {
    transform: scale(1.2);
}

:root.dark .mine-container {
    color: #f5f5f5;
    background-color: #1A1A1A;
    border: 1px solid #fff;
}

:root.dark .avatar {
    border: 1px solid #fff;
}
</style>
