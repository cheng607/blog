export interface ArticleDraft {
    title: string
    content: string
    tagIds: number[]
    cover: File | null
    coverUrl: string
    updatedAt: number
    coverKey?: string
}
const KEY = 'blog-article-draft-v1'
let queue: Promise<unknown> = Promise.resolve()
const fingerprint = (file: File | null) => file ? `${file.name}:${file.size}:${file.lastModified}` : ''
function valid(value: unknown): value is ArticleDraft {
    const v = value as ArticleDraft | null
    return !!v && typeof v.title === 'string' && typeof v.content === 'string'
        && Array.isArray(v.tagIds) && v.tagIds.every(id => Number.isInteger(id) && id > 0)
        && typeof v.updatedAt === 'number' && typeof v.coverUrl === 'string'
}
async function database(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open('blog-drafts', 1)
        request.onupgradeneeded = () => request.result.createObjectStore('drafts')
        request.onsuccess = () => resolve(request.result)
        request.onerror = () => reject(request.error)
        request.onblocked = () => reject(new Error('草稿存储被其他页面占用'))
    })
}
async function stored(write?: ArticleDraft | null): Promise<ArticleDraft | null> {
    const db = await database()
    return new Promise((resolve, reject) => {
        const transaction = db.transaction('drafts', write === undefined ? 'readonly' : 'readwrite')
        const store = transaction.objectStore('drafts')
        const request = write === undefined ? store.get(KEY) : write === null ? store.delete(KEY) : store.put(write, KEY)
        transaction.oncomplete = () => { db.close(); resolve(write === undefined ? request.result ?? null : write) }
        transaction.onerror = () => { db.close(); reject(transaction.error) }
        transaction.onabort = () => { db.close(); reject(transaction.error || new Error('草稿保存中断')) }
    })
}
// Synchronous text backup also runs before reload; cover files are saved in IndexedDB.
export function backupDraft(draft: ArticleDraft): boolean {
    try {
        localStorage.setItem(KEY, JSON.stringify({ ...draft, cover: null, coverKey: fingerprint(draft.cover) }))
        return true
    } catch { return false }
}
export function saveDraft(draft: ArticleDraft): Promise<void> {
    const snapshot = { ...draft, tagIds: [...draft.tagIds], coverKey: fingerprint(draft.cover) }
    const textSaved = backupDraft(snapshot)
    const pending = queue.then(async () => {
        try { await stored(snapshot) }
        catch { if (!textSaved || snapshot.cover) throw new Error('草稿未完整保存，请勿关闭页面；请检查浏览器存储空间') }
    })
    queue = pending.catch(() => {})
    return pending
}
export async function loadDraft(): Promise<{ draft: ArticleDraft | null; missingCover: boolean }> {
    let text: ArticleDraft | null = null
    try { const value = JSON.parse(localStorage.getItem(KEY) || 'null'); if (valid(value)) text = value } catch {}
    let file: ArticleDraft | null = null
    try { const value = await stored(); if (valid(value)) file = value } catch {}
    const draft = text && (!file || text.updatedAt >= file.updatedAt) ? text : file
    if (!draft) {
        // Preserve the body cached by the previous editor version on upgrade.
        try {
            const content = localStorage.getItem('article-editor')
            if (content) return { draft: { title: '', content, tagIds: [], cover: null, coverUrl: '', updatedAt: Date.now() }, missingCover: false }
        } catch {}
        return { draft: null, missingCover: false }
    }
    const cover = file?.cover && fingerprint(file.cover) === draft.coverKey ? file.cover : null
    return { draft: { ...draft, cover }, missingCover: !!draft.coverKey && !cover && !draft.coverUrl }
}
export function clearDraft(): Promise<void> {
    // A newer tombstone prevents an old IndexedDB record being restored if deletion fails.
    const empty: ArticleDraft = { title: '', content: '', tagIds: [], cover: null, coverUrl: '', updatedAt: Date.now() }
    backupDraft(empty)
    try { localStorage.removeItem('article-editor') } catch {}
    const pending = queue.then(async () => { await stored(null) })
    queue = pending.catch(() => {})
    return pending
}
