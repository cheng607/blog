import { createRouter, createWebHistory } from 'vue-router'
import { ElMessage } from 'element-plus'
import 'element-plus/es/components/message/style/css'
const Main = () => import('../views/Main.vue')
const Edit = () => import('../views/Edit.vue')
const About = () => import('../views/About.vue')
const Categories = () => import('../views/Categories.vue')
const CategoryDetail = () => import('../views/CategoryDetail.vue')
const Article = () => import('../views/Article.vue')
const Hollow = () => import('../views/Hollow.vue')
const AdminLogin = () => import('../views/AdminLogin.vue')
import { isAdminAuthed, checkAdminSession } from '../utils/adminAuth'

const NotFound = () => import('../views/NotFound.vue')

const routes = [
    { path: '/:pathMatch(.*)*', name: 'notFound', component: NotFound, meta: { title: '页面不存在', public: true } },
    { path: '/', name: 'home', component: Main, meta: { title: '首页', public: true } },
    { path: '/about', name: 'about', component: About, meta: { title: '关于', public: true } },
    { path: '/categories', name: 'categories', component: Categories, meta: { title: '分类', public: true } },
    { path: '/categoryDetail', name: 'categoryDetail', component: CategoryDetail, meta: { title: '分类详情', public: true } },
    { path: '/article', name: 'article', component: Article, meta: { title: '文章', public: true } },
    { path: '/hollow', name: 'hollow', component: Hollow, meta: { title: '树洞', public: true } },
    { path: '/admin', name: 'admin', component: AdminLogin, meta: { title: '管理登录', public: true } },
    {
        path: '/admin/write',
        name: 'adminWrite',
        component: Edit,
        meta: { title: '写文章', requiresAdmin: true },
    },
    {
        path: '/edit',
        redirect: () =>
            isAdminAuthed()
                ? '/admin/write'
                : { path: '/admin', query: { redirect: '/admin/write' } },
    },
]

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
})

router.beforeEach(async (to) => {
    if (to.meta.requiresAdmin && !(await checkAdminSession())) {
        ElMessage.warning('请先验证管理密钥')
        return {
            path: '/admin',
            query: { redirect: to.fullPath },
        }
    }
})

router.afterEach((to) => {
    const title = (to.meta?.title as string) || '下雨天个人博客'
    document.title = `${title} | 下雨天个人博客`
})

export default router
