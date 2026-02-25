export interface Article {
    title: string;
    content: string;
    tagIds: number[];
    // coverImage?: string;
    id?: number;
    userId?: number;
}

export interface ApiResponse<T = any> {
    code: number
    data: T
    message?: string
}