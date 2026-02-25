<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router';
import { getUser } from '../api/user';
import { getList } from '../api/articles';
import { getTag } from '../api/tag'
import { getHollow } from '../api/hollow';
interface user {
    username: string,
    avatar: string,
    signature: string
}
const ArticleCount = ref<number>(0)
const HollowCount = ref<number>(0)
const TypeCount = ref<number>(0)
const UserData = ref<user>()
onMounted(async () => {
    try {
        const res = await getUser()
        UserData.value = res.data[0]
        const ArticleRes = await getList()
        ArticleCount.value = ArticleRes.data.length
        const HollowRes=await getHollow()
        HollowCount.value=HollowRes.data.length
        const TagRes = await getTag()
        TypeCount.value = TagRes.data.length
    } catch (error) {
        console.error('用户数据错误', error)
    }
})
const router = useRouter()
const handleLink = () => {
    router.push('/about')
}
</script>
 
<template>
    <div class="mine-container">
        <img :src="UserData?.avatar || '../assets/loading.gif'" class="avatar">
        <div class="name">{{ UserData?.username }}</div>
        <div class="note">个性签名：{{ UserData?.signature }}</div>
        <div class="summary">
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
