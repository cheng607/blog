import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'

function load(file, dependencies, globals = {}) {
    const source = readFileSync(new URL(file, import.meta.url), 'utf8')
    const { outputText } = ts.transpileModule(source, {
        compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 },
    })
    const exports = {}
    vm.runInNewContext(outputText, { exports, require: (name) => {
        assert.ok(name in dependencies, `Unexpected dependency: ${name}`)
        return dependencies[name]
    }, console: { error() {} }, ...globals })
    return exports
}

function authFixture() {
    const storage = new Map([['blog-admin-auth', 'legacy-password']])
    const calls = []
    let failLogout = false
    const client = {
        async post(path, body, options) {
            calls.push({ path, body, options })
            if (path === '/auth/logout' && failLogout) throw new Error('offline')
            if (path === '/auth/login' && body.password !== 'correct') {
                throw { isAxiosError: true, response: { status: 401, data: { msg: 'wrong password' } } }
            }
            return { data: { code: 200, data: { token: 'opaque-session' } } }
        },
        async get() { return { data: { code: 200 } } },
    }
    const auth = load('../src/utils/adminAuth.ts', {
        axios: { default: { create: () => client, isAxiosError: (e) => e.isAxiosError } },
        vue: { ref: (value) => ({ value }) },
        '../config': { API_BASE_URL: '/api/' },
    }, { sessionStorage: {
        getItem: (k) => storage.get(k) ?? null,
        setItem: (k, v) => storage.set(k, v),
        removeItem: (k) => storage.delete(k),
    } })
    return { auth, storage, calls, failLogout: () => { failLogout = true } }
}

test('login stores only the issued session; logout revokes it before clearing storage', async () => {
    const { auth, storage, calls } = authFixture()
    assert.equal(storage.has('blog-admin-auth'), false)
    await assert.rejects(auth.loginAdmin('wrong'), /wrong password/)
    assert.equal(auth.isAdminAuthed(), false)
    await auth.loginAdmin('correct')
    assert.equal(auth.getAdminToken(), 'opaque-session')
    assert.equal([...storage.values()].includes('correct'), false)
    assert.equal(auth.adminHeaders().Authorization, 'Bearer opaque-session')
    await auth.logoutAdmin()
    assert.equal(calls.at(-1).options.headers.Authorization, 'Bearer opaque-session')
    assert.equal(auth.isAdminAuthed(), false)
})

test('failed server logout does not pretend the session was revoked', async () => {
    const f = authFixture()
    await f.auth.loginAdmin('correct')
    f.failLogout()
    await assert.rejects(f.auth.logoutAdmin(), /offline/)
    assert.equal(f.auth.isAdminAuthed(), true)
})

test('API requests attach sessions, third-party requests do not; 401 opens reauth without navigation', async () => {
    const { auth } = authFixture()
    await auth.loginAdmin('correct')
    let beforeRequest, onError
    const client = { interceptors: {
        request: { use: (fn) => { beforeRequest = fn } },
        response: { use: (_, fn) => { onError = fn } },
    } }
    load('../src/utils/request.ts', {
        axios: { default: { create: () => client } },
        '../config': { API_BASE_URL: '/api/' },
        './adminAuth': auth,
    }, { window: { location: { assign: () => assert.fail('must retain current page') } } })
    for (const url of ['/article', '/hollow', '/upload/image']) {
        const headers = new Map()
        beforeRequest({ url, headers })
        assert.equal(headers.get('Authorization'), 'Bearer opaque-session')
    }
    for (const url of ['https://music.example/api', '//quotes.example/api']) {
        const headers = new Map()
        beforeRequest({ url, headers })
        assert.equal(headers.has('Authorization'), false)
    }
    await assert.rejects(onError({ response: { status: 401 }, config: { url: '/article' } }))
    assert.equal(auth.reauthRequired.value, true)
    assert.equal(auth.isAdminAuthed(), false)
    await auth.loginAdmin('correct')
    assert.equal(auth.reauthRequired.value, false)
    assert.equal(auth.adminHeaders().Authorization, 'Bearer opaque-session')
})

test('HTTP validation, quota and timeout messages remain actionable', async () => {
    let onError
    const client = { interceptors: {
        request: { use() {} }, response: { use: (_, fn) => { onError = fn } },
    } }
    load('../src/utils/request.ts', {
        axios: { default: { create: () => client } },
        '../config': { API_BASE_URL: '/api/' },
        './adminAuth': { adminHeaders: () => ({}), requireAdminLogin: () => assert.fail('non-auth error must not log out') },
    })
    for (const [status, msg] of [[400,'标题不能为空'],[404,'文章不存在'],[429,'今日 AI 额度已用完'],[504,'AI 回答超时']]) {
        await assert.rejects(onError({ response:{ status, data:{ msg } }, config:{ url:'chat' } }), error => error.message === msg)
    }
    await assert.rejects(onError({ code:'ECONNABORTED', request:{} }), /请求超时/)
})
