import axios from 'axios'
import { useRouter } from 'vue-router'
const router = useRouter()
// 1. 创建Axios实例
const request = axios.create({
    baseURL: 'http://47.114.122.96:9090/',
    timeout: 10000,
    headers: {
        'Content-Type': 'application/json'
    }
})

// 2. 请求拦截器：发送请求前处理
request.interceptors.request.use(
    (config) => {
        return config
    },
    (error) => {
        return Promise.reject(error)
    }
)

// 3. 响应拦截器：接收响应后处理
request.interceptors.response.use(
    (response) => {
        const res = response.data
        if (res.code !== 200 && res.code) {
            return Promise.reject(new Error(res.message || 'Error'))
        }
        return res
    },
    (error) => {
        let message = ''
        console.log(error)
        if (error.response) {
            // 服务器返回错误
            switch (error.response.status) {
                case 401:
                    message = '登录已过期，请重新登录'
                    router.push('/')
                    break
                case 404:
                    message = '请求地址不存在'
                    break
                case 500:
                    message = '服务器内部错误'
                    break
                default:
                    message = `请求错误（${error.response.status}）`
            }
        } else if (error.request) {
            // 请求已发出但无响应（如网络错误）
            message = '网络连接失败，请检查网络'
        } else {
            // 其他错误
            message = error.message
        }

        // 使用message变量：打印错误信息（或根据需求展示给用户）
        console.error('请求错误:', message)
        // 也可以将message作为错误信息传递
        return Promise.reject(new Error(message))
    }
)
const requestMethod = (url: string, method: string, data = {}, config = {}) => {
    // 根据请求方法区分参数类型（GET用params，其他用data）
    const params = method.toLowerCase() === 'get' ? data : {}
    const requestData = method.toLowerCase() !== 'get' ? data : {}

    return request({
        url,
        method,
        params, // GET参数
        data: requestData, // POST/PUT参数
        ...config // 允许覆盖默认配置（如自定义headers）
    })
}

// 3. 暴露常用请求方法（简化调用）
export const http = {
    get: (url: string, params = {}, config = {}) => {
        return requestMethod(url, 'get', params, config)
    },
    post: (url: string, data = {}, config = {}) => {
        return requestMethod(url, 'post', data, config)
    },
    put: (url: string, data = {}, config = {}) => {
        return requestMethod(url, 'put', data, config)
    },
    delete: (url: string, params = {}, config = {}) => {
        return requestMethod(url, 'delete', params, config)
    }
}
export default request