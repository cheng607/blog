<script setup lang="ts">

import { ref, onMounted } from 'vue'
interface linkItem {
    linkName: string,
    linkUrl: string
}
interface messageItem {
    messageName: string,
    messageContent: string
}
const linkList = ref<linkItem[]>()
const messageList = ref<messageItem[]>()
onMounted(() => {
    linkList.value = [
        { linkName: 'MDN', linkUrl: 'https://developer.mozilla.org/zh-CN/' },
        { linkName: 'React', linkUrl: 'https://zh-hans.react.dev/' },
        { linkName: 'Vue', linkUrl: 'https://cn.vuejs.org/' },
        { linkName: 'CSDN', linkUrl: 'https://www.csdn.net/' },
        { linkName: '稀土掘金', linkUrl: 'https://juejin.cn/' },
        { linkName: '菜鸟教程', linkUrl: 'https://www.runoob.com/' },
    ]
    messageList.value = [
        { messageName: '电话', messageContent: '18934689565' },
        { messageName: '邮箱', messageContent: 'czrxiayutian@163.com' }
    ]
})

</script>

<template>
    <div class="footer-container">
        <div class="footer">
            <!-- 友情链接 -->
            <div class="link">
                <div class="title">友情链接</div>
                <div class="list">
                    <a v-for="link in linkList" :href="link.linkUrl" target="_blank">
                        {{ link.linkName }}
                    </a>
                </div>
            </div>
            <!-- 社媒 -->
            <div class="social">
                <div class="soc-box">
                    <img src="../assets/wechat.png" class="soc-icon">
                    <div class="soc-title">微信号</div>
                </div>
                <div class="soc-box">
                    <img src="../assets/qqIcon.png" class="soc-icon">
                    <div class="soc-title">QQ号</div>
                </div>
            </div>
            <!-- 联系方式 -->
            <div class="message">
                <div class="title">联系方式</div>
                <div class="message-box">
                    <div class="message-item" v-for="msg in messageList">
                        {{ msg.messageName }}:{{ msg.messageContent }}
                    </div>
                </div>
            </div>
        </div>
        <div class="beian-wrap">
            <a href="https://beian.mps.gov.cn/#/query/webSearch?code=42122402000254" target="_blank"
                rel="nofollow noopener noreferrer" class="jj-link">
                <img src="//lf-web-assets.juejin.cn/obj/juejin-web/xitu_juejin_web/img/police.d0289dc.png" alt="公安备案图标"
                    style="vertical-align:middle; width:16px; height:16px; margin-right:4px;">
                鄂公网安备42122402000254号
            </a>
            <span class="copyright">©2025 下雨天 | 本站原创内容转载请注明出处</span>
        </div>
    </div>
</template>

<style scoped>
.footer-container {
    width: 100%;
    height: 250px;
    padding: 50px 0;
    background-color: #fff;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: space-between;
}

:root.dark .footer-container {
    background-color: #333;
}

.footer {
    width: 100%;
    display: flex;
    align-items: center;
    justify-content: space-around;
}

.title {
    font-size: 24px;
    text-align: center;
    margin: 15px 0;
}

.link .list {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
    gap: 8px;
}

.link .list a {
    text-align: center;
    text-decoration: none;
    color: #333;
    border: 1px solid #333;
    border-radius: 16px;
    padding: 4px 10px;
    font-size: 13px;
    white-space: nowrap;
}

:root.dark .link .list a {
    color: #fff;
    border: 1px solid #fff;
}

.link .list a:hover {
    color: skyblue;
    border: 1px solid skyblue;
}

.social {
    display: flex;
    gap: 20px;
}

.social .soc-icon {
    width: 150px;
    height: 150px;
    border: 1px solid #333;
    border-radius: 20px;
}

.soc-box {
    text-align: center;
    margin: 10px 0;
}

.message-box {
    display: flex;
    flex-direction: column;
    gap: 10px;
}

.beian-wrap {
    margin-top: 20px;
    font-size: 14px;
    display: flex;
    align-items: center;
    gap: 8px;
}

.jj-link {
    cursor: pointer;
    color: #666;
    text-decoration: none;
    -webkit-tap-highlight-color: transparent;
    user-select: none;
}

:root.dark .jj-link {
    color: #ddd;
}

.copyright {
    color: #666;
}

:root.dark .copyright {
    color: #ddd;
}

/* 移动端 Footer 优化 - 减少冗余 */
@media (max-width: 768px) {
  .footer-container {
    height: auto !important;
    padding: 24px 12px 20px !important;
  }

  .footer {
    flex-direction: column !important;
    gap: 16px !important;
    align-items: stretch !important;
  }

  /* 友情链接 - flex 自然换行，更整齐 */
  .link .list {
    gap: 6px !important;
  }

  .link .list a {
    font-size: 12.5px !important;
    padding: 3px 9px !important;
    border-radius: 14px !important;
  }

  .title {
    font-size: 16px !important;
    margin: 4px 0 !important;
  }

  /* 二维码大幅缩小 */
  .social {
    justify-content: center !important;
    gap: 24px !important;
  }

  .social .soc-icon {
    width: 90px !important;
    height: 90px !important;
    border-radius: 12px !important;
  }

  .soc-box {
    margin: 0 !important;
  }

  .soc-title {
    font-size: 12px !important;
    margin-top: 4px !important;
  }

  /* 联系方式 */
  .message-box {
    flex-direction: row !important;
    justify-content: center !important;
    gap: 16px !important;
    font-size: 13px !important;
  }

  .beian-wrap {
    flex-direction: column !important;
    font-size: 12px !important;
    margin-top: 12px !important;
    gap: 4px !important;
    text-align: center;
  }
}
</style>
