import test from 'node:test'
import assert from 'node:assert/strict'
import { readFileSync } from 'node:fs'
import vm from 'node:vm'
import ts from 'typescript'

const source = readFileSync(new URL('../src/utils/media.ts', import.meta.url), 'utf8')
const code = ts.transpileModule(source, { compilerOptions: { module: ts.ModuleKind.CommonJS } }).outputText
const exports = {}
vm.runInNewContext(code, { exports, require: () => ({ API_BASE_URL: '/api/' }) })
test('media URLs tolerate whitespace and a leading slash before an absolute address', () => {
  assert.equal(exports.resolveMediaUrl(' /https://example.com/avatar.jpg '), 'https://example.com/avatar.jpg')
  assert.equal(exports.resolveMediaUrl('https://example.com/avatar.jpg'), 'https://example.com/avatar.jpg')
  assert.equal(exports.resolveMediaUrl('/upload/images/a.png'), '/api/upload/images/a.png')
  assert.equal(exports.resolveMediaUrl('  '), '')
})
