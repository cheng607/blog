import { readFileSync } from 'node:fs'
import { gzipSync } from 'node:zlib'
const manifest = JSON.parse(readFileSync('dist/.vite/manifest.json', 'utf8'))
const seen = new Set()
function visit(key) {
  if (seen.has(key)) return
  seen.add(key)
  for (const dependency of manifest[key].imports || []) visit(dependency)
}
visit('index.html')
visit('src/views/Main.vue')
let raw = 0, gzip = 0
for (const key of seen) {
  const bytes = readFileSync('dist/' + manifest[key].file)
  raw += bytes.length
  gzip += gzipSync(bytes).length
}
console.log(JSON.stringify({ initialJsBytes: raw, initialJsGzipBytes: gzip, chunks: [...seen] }, null, 2))
