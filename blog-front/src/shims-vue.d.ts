// src/shims-vue.d.ts
import type { ComponentCustomProperties } from 'vue'
import { formatTime } from './utils/format' // 导入 formatTime 函数（确保路径正确）

// 声明 .vue 文件模块类型（原有内容保留）
declare module '*.vue' {
    import type { DefineComponent } from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}

// 扩展 Vue 全局属性类型
declare module 'vue' {
    // 给组件实例添加 $formatTime 类型声明
    interface ComponentCustomProperties {
        $formatTime: typeof formatTime
    }
}