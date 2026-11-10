import { resolveMediaUrl } from './media'
import { API_BASE_URL } from '../config'

/** 将后端 { code:200, data:{ imageUrl } } 转为 Vditor 要求的 { code:0, data:{ succMap } } */
export function formatVditorUploadResponse(files: File[], responseText: string): string {
    try {
        const res = JSON.parse(responseText) as {
            code?: number
            msg?: string
            data?: { imageUrl?: string }
        }
        if (res.code !== 200 || !res.data?.imageUrl) {
            return JSON.stringify({
                code: 1,
                msg: res.msg || '上传失败',
                data: { errFiles: files.map((f) => f.name), succMap: {} },
            })
        }
        const url = resolveMediaUrl(res.data.imageUrl)
        const name = files[0]?.name || 'image'
        return JSON.stringify({
            code: 0,
            msg: '',
            data: {
                errFiles: [],
                succMap: { [name]: url },
            },
        })
    } catch {
        return JSON.stringify({
            code: 1,
            msg: '响应解析失败',
            data: { errFiles: files.map((f) => f.name), succMap: {} },
        })
    }
}

export const VDITOR_UPLOAD_URL = `${API_BASE_URL.replace(/\/$/, '')}/upload/image`
