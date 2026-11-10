import { http } from '../utils/request'
import { type Article } from '../types/articles'

export const getList = (page = 1, pageSize = 5) => {
    return http.get('/article/list', { page, pageSize })
}

export const getArticlesByTag = (tagId: number, page = 1, pageSize = 5) => {
    return http.get(`/article/tag/${tagId}`, { page, pageSize })
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
