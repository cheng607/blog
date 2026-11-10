import { http } from '../utils/request'

export interface ChatResult {
    content: string
}

/** AI 超时预算：单模型30秒、后端总65秒、Nginx85秒、浏览器90秒。 */
export function getChatStatus() {
    return http.get('chat/status', {}, { timeout: 5000 })
}

export function chatWithAI(message: string) {
    return http.post('chat', { message }, { timeout: 90000 })
}
