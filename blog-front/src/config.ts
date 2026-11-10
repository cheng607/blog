/**
 * API 基础地址配置
 * - 开发环境默认使用 Vite 代理 (/api/)
 * - 生产环境必须通过环境变量 VITE_API_BASE_URL 配置（推荐使用域名）
 * - 不再硬编码 IP 地址
 */
export const API_BASE_URL =
	import.meta.env.VITE_API_BASE_URL ?? (import.meta.env.DEV ? '/api/' : '')
