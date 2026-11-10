export const IMAGE_ACCEPT = '.jpg,.jpeg,.png'
export const MAX_IMAGE_BYTES = 10 * 1024 * 1024

export function validateImageFile(file: Pick<File, 'name' | 'size' | 'type'>): string | null {
    if (!/\.(jpe?g|png)$/i.test(file.name) || (file.type && !['image/jpeg', 'image/png'].includes(file.type))) {
        return '仅支持 JPG、JPEG、PNG 图片'
    }
    if (!file.size) return '图片不能为空'
    if (file.size > MAX_IMAGE_BYTES) return '图片大小不能超过10MB'
    return null
}
