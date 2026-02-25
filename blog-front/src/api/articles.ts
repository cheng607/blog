import { http } from '../utils/request'
import { type Article } from '../types/articles'
export const getList = () => {
    return http.get('/article/list')
}
export const getArticle = (id: number) => {
    return http.get(`/article/${id}`)
}

export const addArticle = (article: Article) => {
    return http.post('/article', article)
}

export const addViewsCount = (id: number) => {
    return http.post(`/article/view/${id}`)
}