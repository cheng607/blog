import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'
import * as vue from 'vue'

function load(file, names, deps = {}) {
    let source = readFileSync(new URL(file, import.meta.url), 'utf8')
    if (file.endsWith('.vue')) source = source.match(/<script setup lang="ts">([\s\S]*?)<\/script>/)[1]
    if (names.length) source += `\nexport { ${names.join(',')} }`
    const mounts = [], stops = [], exports = {}
    const vueMock = { ...vue, onMounted: fn => mounts.push(fn), onUnmounted: fn => stops.push(fn) }
    const output = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 } }).outputText
    vm.runInNewContext(output, { exports, require: name => {
        if (name === 'vue') return vueMock
        if (name.endsWith('/style/css') || name.endsWith('.css')) return {}
        if (name === 'element-plus') return { ElLoading: { directive: {} }, ...deps[name] }
        if (!(name in deps)) throw new Error(`Missing mock ${name}`)
        return deps[name]
    }, setTimeout: (fn, ms) => { const timer = setTimeout(fn, ms); timer.unref(); return timer }, clearTimeout, window: { addEventListener() {}, removeEventListener() {} }, console: { error() {} }, URL: { createObjectURL: () => 'blob:cover', revokeObjectURL() {} } })
    return { ...exports, mount: async () => { for (const fn of mounts) await fn() }, stop: () => stops.forEach(fn => fn()) }
}
const images = load('../src/utils/imageUpload.ts', [])
const deferred = () => { let resolve, reject; const promise = new Promise((a,b) => { resolve=a; reject=b }); return { promise, resolve, reject } }
const flush = async () => { await vue.nextTick(); await Promise.resolve(); await vue.nextTick() }

async function editor(overrides = {}) {
    const published = []
    let text = 'Before\n\n&lt;example&gt;\n\nAfter', options
    const e = load('../src/views/Edit.vue', ['title','tagIds','profile','tags','newTagName','vditorRef','coverFile','publishing','handleSaveProfile','handleCreateTag','handlePublish'], {
        vditor: { default: class {
            constructor(el, opts) { options = opts }
            getValue() { return text }
            setValue(value) { text=value }
            disabled() {} enable() {} clearCache() {} destroy() {}
        } },
        'vditor/dist/index.css': {},
        '@element-plus/icons-vue': {},
        'element-plus': { ElMessage: { success() {}, error() {}, warning() {} } },
        'vue-router': { useRouter: () => ({ push() {} }), onBeforeRouteLeave() {} },
        '../api/articles': { addArticle: async payload => { published.push(payload); return { code:200 } } },
        '../api/tag': { getTag: async () => ({ data: null }), createTag: async name => ({ data:{ id:9, name } }) },
        '../api/user': { getUser: async () => ({ data:null }), saveProfile: async profile => ({ data: { ...profile, id:7 } }) },
        '../api/upload': { uploadImage: async () => ({ data:{ imageUrl:'/upload/images/test.png' } }) },
        '../utils/imageUpload': images,
        '../utils/draft': { loadDraft: async () => ({ draft: null, missingCover: false }), saveDraft: async () => {}, backupDraft() {}, clearDraft: async () => {} },
        '../utils/media': { resolveMediaUrl: url => url || '' },
        '../utils/adminAuth': { logoutAdmin() {}, adminHeaders: () => ({}), requireAdminLogin() {} },
        '../utils/vditorUpload': { VDITOR_UPLOAD_URL:'/api/upload/image', formatVditorUploadResponse() {} },
        ...overrides,
    })
    e.vditorRef.value = { addEventListener: () => assert.fail('must not intercept native paste') }
    await e.mount()
    e.profile.value.username='博主'
    await e.handleSaveProfile()
    e.newTagName.value='技术'
    await e.handleCreateTag()
    e.title.value='First article'
    return { e, published, options }
}

test('fresh blog can create profile/tag and publish using returned IDs; preserves Markdown entities', async () => {
    const { e, published, options } = await editor()
    assert.equal(e.profile.value.id, 7)
    assert.equal(e.tags.value.length, 1)
    assert.equal(options.upload.accept, images.IMAGE_ACCEPT)
    await e.handlePublish()
    assert.equal(published.length, 1)
    assert.equal(published[0].userId, 7)
    assert.deepEqual(Array.from(published[0].tagIds), [9])
    assert.match(published[0].content, /&lt;example&gt;/)
    assert.equal(e.title.value, '')
    assert.equal(e.publishing.value, false)
})

test('publish lock covers cover upload and is released on failure without clearing draft', async () => {
    const upload = deferred()
    let uploads = 0
    const { e, published } = await editor({ '../api/upload': { uploadImage: () => { uploads++; return upload.promise } } })
    e.coverFile.value = { name:'cover.png' }
    const first = e.handlePublish()
    await e.handlePublish()
    assert.equal(uploads, 1)
    assert.equal(e.publishing.value, true)
    upload.reject(new Error('offline'))
    await first
    assert.equal(e.publishing.value, false)
    assert.equal(e.title.value, 'First article')
    assert.equal(e.coverFile.value.name, 'cover.png')
    assert.equal(published.length, 0)
})

test('repeated publish clicks during article save create only one request', async () => {
    const save = deferred()
    let calls = 0
    const { e } = await editor({ '../api/articles': { addArticle: () => { calls++; return save.promise } } })
    const first = e.handlePublish()
    await e.handlePublish()
    assert.equal(calls, 1)
    save.resolve({ code:200 })
    await first
    assert.equal(e.publishing.value, false)
})

test('category query changes ignore stale responses and clear missing categories', async () => {
    const route = vue.reactive({ query:{ typeId:'1' } })
    const requests = new Map()
    const e = load('../src/views/CategoryDetail.vue', ['typeData','error','loading'], {
        'vue-router': { useRoute: () => route },
        '../components/List.vue': {},
        '../api/tag': { getTagByid: id => { const d=deferred(); requests.set(id,d); return d.promise } },
        'element-plus': { ElMessage:{ error() {} } },
    })
    route.query.typeId='2'; await flush()
    requests.get(2).resolve({ data:{ id:2, name:'B' } }); await flush()
    requests.get(1).resolve({ data:{ id:1, name:'A' } }); await flush()
    assert.equal(e.typeData.value.name, 'B')
    route.query.typeId='3'; await flush()
    assert.equal(e.typeData.value, undefined)
    requests.get(3).resolve({ data:null }); await flush()
    assert.equal(e.error.value, '分类不存在')
    assert.equal(e.loading.value, false)
    e.stop()
})

test('image picker accepts only PNG/JPEG within 10 MB', () => {
    assert.equal(images.validateImageFile({ name:'photo.PNG',type:'image/png',size:100 }), null)
    assert.ok(images.validateImageFile({ name:'photo.webp',type:'image/webp',size:100 }))
    assert.ok(images.validateImageFile({ name:'photo.png',type:'text/html',size:100 }))
    assert.ok(images.validateImageFile({ name:'photo.jpg',type:'image/jpeg',size:0 }))
    assert.ok(images.validateImageFile({ name:'photo.jpg',type:'image/jpeg',size:10*1024*1024+1 }))
})

test('empty homepage data produces zero counts and an empty category list', async () => {
    const empty = async () => ({ data:null })
    const mine = load('../src/components/Mine.vue', ['UserData','ArticleCount','HollowCount','TypeCount'], {
        'vue-router': { useRouter: () => ({}) },
        '../api/site': { getSiteSummary:empty },
        '../utils/media':{ resolveMediaUrl: () => '' }, '../assets/loading.gif':{},
    })
    await mine.mount()
    assert.equal(mine.UserData.value, undefined)
    assert.equal(mine.ArticleCount.value, 0)
    assert.equal(mine.HollowCount.value, 0)
    assert.equal(mine.TypeCount.value, 0)
    const tags = load('../src/components/Type.vue', ['displayTypeData'], {
        'vue-router': { useRouter: () => ({}) }, '../api/tag':{ getTag:empty },
    })
    await tags.mount()
    assert.equal(tags.displayTypeData.value.length, 0)
})

test('article list ignores previous category response after fast navigation', async () => {
    const route = vue.reactive({ query:{ typeId:'1' } })
    const requests = new Map()
    const list = load('../src/components/List.vue', ['listData','isLoading'], {
        'vue-router': { useRoute: () => route, useRouter: () => ({}) },
        '../api/articles': { getArticlesByTag: id => { const d=deferred(); requests.set(id,d); return d.promise } },
        '../utils/media':{}, '../utils/markdown':{}, '../assets/loading.gif':{},
        'element-plus': { ElMessage:{ error() {} } },
    })
    route.query.typeId='2'; await flush()
    requests.get(2).resolve({ data:{ items:[{ id:2 }], total:1, page:1 } }); await flush()
    requests.get(1).reject(new Error('old request failed')); await flush()
    assert.equal(list.listData.value[0].id, 2)
    assert.equal(list.isLoading.value, false)
    list.stop()
})
