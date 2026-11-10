import { API_BASE_URL } from '../config'

/** 将后端返回的相对路径转为可访问的完整 URL */
export function resolveMediaUrl(url?: string | null): string {
  if (!url) return ''
  url = url.trim().replace(/^\/+((?:https?):\/\/)/i, '$1')
  if (!url) return ''
  if (/^https?:\/\//i.test(url)) return url
  const base = API_BASE_URL.replace(/\/$/, '')
  const path = url.startsWith('/') ? url : `/${url}`
  return `${base}${path}`
}
