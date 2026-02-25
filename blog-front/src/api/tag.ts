import { http } from "../utils/request";

export const getTag = () => {
    return http.get('/tag')
}
export const getTagByid = (tagId: number) => {
    return http.get(`/tag/${tagId}`)
}