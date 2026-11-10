import { http } from '../utils/request'

export const getHollow = () => {
    return http.get('/hollow')
}

export const addHollow = (word: string, date: string) => {
    return http.post('/hollow', { word, date })
}
