import { http } from '../utils/request'
// 获取树洞列表
export const getHollow = () => {
    return http.get(`/hollow`)
}
// 添加树洞
export const addHollow = (word: string, date: string) => {
    console.log(word)
    console.log(date)
    return http.post('/hollow', { word, date })
}
