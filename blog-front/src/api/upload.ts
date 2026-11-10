import { http } from '../utils/request'

export const uploadImage = (file: File) => {
    const formData = new FormData()
    formData.append('file', file)
    return http.upload('/upload/image', formData)
}
