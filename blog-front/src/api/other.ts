import { http } from "../utils/request"
// 随机一言
export const getWord = () => {
    return http.get('https://hmajax.itheima.net/api/ambition')
}
// 获取随机音乐
export const getMusic = () => {
    return http.get('https://v2.xxapi.cn/api/randomkuwo')
}