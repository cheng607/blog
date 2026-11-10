import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'

function draftStore() {
  const values = new Map(), exports = {}
  const localStorage = { getItem: k => values.get(k) ?? null, setItem: (k,v) => values.set(k,v), removeItem: k => values.delete(k) }
  const source = readFileSync(new URL('../src/utils/draft.ts', import.meta.url), 'utf8')
  const code = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.CommonJS, target: ts.ScriptTarget.ES2022 } }).outputText
  vm.runInNewContext(code, { exports, localStorage, Date })
  return { ...exports, values }
}
const sample = () => ({ title:'草稿', content:'# 正文', tagIds:[1], cover:null, coverUrl:'', updatedAt:Date.now() })
test('text draft survives reload when IndexedDB is unavailable', async () => {
  const store = draftStore()
  await store.saveDraft(sample())
  const restored = await store.loadDraft()
  assert.equal(restored.draft.title, '草稿')
  assert.equal(restored.draft.content, '# 正文')
  assert.equal(restored.missingCover, false)
})
test('cover storage failure reports loss while keeping text recoverable', async () => {
  const store = draftStore()
  await assert.rejects(store.saveDraft({ ...sample(), cover:{ name:'cover.png', size:100, lastModified:1 } }), /草稿未完整保存/)
  const restored = await store.loadDraft()
  assert.equal(restored.draft.content, '# 正文')
  assert.equal(restored.missingCover, true)
})
test('clear tombstone prevents fallback draft and legacy body resurrection', async () => {
  const store = draftStore()
  store.values.set('article-editor', '旧正文')
  assert.equal((await store.loadDraft()).draft.content, '旧正文')
  await store.saveDraft(sample())
  await assert.rejects(store.clearDraft())
  assert.equal((await store.loadDraft()).draft.content, '')
  assert.equal(store.values.has('article-editor'), false)
})
test('malformed saved data is ignored', async () => {
  const store = draftStore()
  store.values.set('blog-article-draft-v1', '{broken')
  assert.equal((await store.loadDraft()).draft, null)
  store.values.set('blog-article-draft-v1', JSON.stringify({ ...sample(), tagIds:[-1] }))
  assert.equal((await store.loadDraft()).draft, null)
})
