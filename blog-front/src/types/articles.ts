export interface Article {
    title: string
    content: string
    tagIds: number[]
    coverImage?: string
    id?: number
    userId?: number
}

export interface ApiResponse<T = unknown> {
    code: number
    data: T
    message?: string
    msg?: string
}
