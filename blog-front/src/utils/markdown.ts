/** 列表摘要：去掉 Markdown 标记，避免整段源码展示在卡片上 */
export function stripMarkdownForPreview(text?: string, maxLen = 120): string {
  if (!text) return ''
  let plain = text
    .replace(/```[\s\S]*?```/g, ' ')
    .replace(/`[^`]+`/g, ' ')
    .replace(/!\[[^\]]*\]\([^)]*\)/g, ' ')
    .replace(/\[([^\]]+)\]\([^)]+\)/g, '$1')
    .replace(/^#{1,6}\s+/gm, '')
    .replace(/(\*\*|__|\*|_|~~)/g, '')
    .replace(/\n+/g, ' ')
    .trim()
  if (plain.length > maxLen) {
    plain = `${plain.slice(0, maxLen)}…`
  }
  return plain
}
