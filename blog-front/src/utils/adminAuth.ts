import axios from 'axios'
import { ref } from 'vue'
import { API_BASE_URL } from '../config'

const AUTH_KEY = 'blog-admin-session'
// Remove the password cached by older versions.
sessionStorage.removeItem('blog-admin-auth')
export const reauthRequired = ref(false)
const authClient = axios.create({ baseURL: API_BASE_URL, timeout: 10000 })

export function getAdminToken(): string {
    return sessionStorage.getItem(AUTH_KEY) || ''
}
export function isAdminAuthed(): boolean {
    return !!getAdminToken()
}
export function clearAdminAuth() {
    sessionStorage.removeItem(AUTH_KEY)
}
export function requireAdminLogin() {
    clearAdminAuth()
    reauthRequired.value = true
}
export function adminHeaders(): Record<string, string> {
    const token = getAdminToken()
    return token ? { Authorization: `Bearer ${token}` } : {}
}
export async function loginAdmin(password: string) {
    try {
        const res = await authClient.post('/auth/login', { password })
        const token = res.data?.data?.token
        if (res.data?.code !== 200 || typeof token !== 'string' || !token) {
            throw new Error('登录响应无效')
        }
        sessionStorage.setItem(AUTH_KEY, token)
        reauthRequired.value = false
    } catch (error) {
        if (axios.isAxiosError(error)) {
            throw new Error(error.response?.data?.msg || '登录失败，请检查网络后重试')
        }
        throw error
    }
}
export async function checkAdminSession(): Promise<boolean> {
    if (!getAdminToken()) return false
    try {
        await authClient.get('/auth/session', { headers: adminHeaders() })
        return true
    } catch (error) {
        if (axios.isAxiosError(error) && error.response?.status === 401) clearAdminAuth()
        return false
    }
}
export async function logoutAdmin() {
    // Do not claim logout succeeded when the server could not revoke the session.
    await authClient.post('/auth/logout', {}, { headers: adminHeaders() })
    clearAdminAuth()
    reauthRequired.value = false
}
