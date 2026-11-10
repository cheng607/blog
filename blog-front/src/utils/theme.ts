export type ThemeMode = 'light' | 'dark'

const STORAGE_KEY = 'app-theme'

export function getSavedTheme(): ThemeMode | null {
    const v = localStorage.getItem(STORAGE_KEY)
    return v === 'light' || v === 'dark' ? v : null
}

export function getSystemTheme(): ThemeMode {
    return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

export function resolveTheme(): ThemeMode {
    return getSavedTheme() ?? getSystemTheme()
}

export function applyTheme(theme: ThemeMode) {
    const root = document.documentElement
    root.classList.remove('light', 'dark')
    root.classList.add(theme)
    root.style.colorScheme = theme
    localStorage.setItem(STORAGE_KEY, theme)
}

export function initTheme() {
    applyTheme(resolveTheme())
}

export function watchSystemTheme(onChange: (theme: ThemeMode) => void) {
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    const handler = () => {
        if (!getSavedTheme()) {
            const next = mq.matches ? 'dark' : 'light'
            applyTheme(next)
            onChange(next)
        }
    }
    mq.addEventListener('change', handler)
    return () => mq.removeEventListener('change', handler)
}
