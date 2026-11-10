import axios from 'axios'
import type { AxiosRequestConfig } from 'axios'
import { API_BASE_URL } from '../config'
import { adminHeaders, requireAdminLogin } from './adminAuth'

const request = axios.create({
    baseURL: API_BASE_URL,
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
})

request.interceptors.request.use(
    (config) => {
        // Public third-party APIs (music/quotes) must never receive the session.
        if (!/^(?:[a-z][a-z0-9+.-]*:)?\/\//i.test(config.url || '')) {
            Object.entries(adminHeaders()).forEach(([key, value]) => config.headers.set(key, value))
        }
        return config
    },
    (error) => Promise.reject(error)
)

request.interceptors.response.use(
    (response) => {
        const res = response.data
        if (res.code !== 200 && res.code) {
            const msg = res.msg || res.message || '请求失败'
            return Promise.reject(new Error(msg))
        }
        return res
    },
    (error) => {
        let message = ''
        if (error.response) {
            const status = error.response.status
            const backendMessage = error.response.data?.msg
            const defaults: Record<number, string> = {
                400: '请求参数无效', 401: '登录已过期，请重新登录',
                404: '请求的内容不存在', 409: '当前操作存在冲突', 413: '提交内容过大',
                429: '请求过于频繁，请稍后重试', 500: '服务暂时不可用，请稍后重试',
                502: 'AI 服务暂时不可用', 503: '服务暂时不可用，请稍后重试', 504: '请求超时，请稍后重试',
            }
            message = typeof backendMessage === 'string' ? backendMessage : (defaults[status] || `请求错误（${status}）`)
            if (status === 401 && !/^(?:[a-z][a-z0-9+.-]*:)?\/\//i.test(error.config?.url || '')) requireAdminLogin()
        } else if (error.code === 'ECONNABORTED' || error.code === 'ETIMEDOUT') {
            message = '请求超时，请稍后重试'
        } else if (error.request) {
            message = '网络连接失败，请检查网络或后端服务'
        } else {
            message = error.message
        }
        console.error('请求错误:', message)
        return Promise.reject(new Error(message))
    }
)

const requestMethod = (url: string, method: string, data = {}, config: AxiosRequestConfig = {}) => {
    const params = method.toLowerCase() === 'get' ? data : {}
    const requestData = method.toLowerCase() !== 'get' ? data : {}

    return request({
        url,
        method,
        params,
        data: requestData,
        ...config
    })
}

export const http = {
    get: (url: string, params = {}, config: AxiosRequestConfig = {}) => {
        return requestMethod(url, 'get', params, config)
    },
    post: (url: string, data = {}, config: AxiosRequestConfig = {}) => {
        return requestMethod(url, 'post', data, config)
    },
    put: (url: string, data = {}, config: AxiosRequestConfig = {}) => {
        return requestMethod(url, 'put', data, config)
    },
    delete: (url: string, params = {}, config: AxiosRequestConfig = {}) => {
        return requestMethod(url, 'delete', params, config)
    },
    upload: (url: string, formData: FormData, config: AxiosRequestConfig = {}) => {
        return request.post(url, formData, {
            headers: { 'Content-Type': 'multipart/form-data' },
            ...config
        })
    }
}

export default request
