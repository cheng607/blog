import { createRouter, createWebHashHistory } from 'vue-router';
import Main from '../views/Main.vue';
import Edit from '../views/Edit.vue';
import About from '../views/About.vue';
import Categories from '../views/Categories.vue';
import CategoryDetail from '../views/CategoryDetail.vue';
import Article from '../views/Article.vue'
import Hollow from '../views/Hollow.vue'
const routes = [
    { path: '/', component: Main },
    { path: '/edit', component: Edit },
    { path: '/about', component: About },
    { path: '/categories', component: Categories },
    { path: '/categoryDetail', component: CategoryDetail },
    { path: '/article', component: Article },
    { path: '/hollow', component: Hollow },
];

const router = createRouter({
    history: createWebHashHistory(),
    routes,
});
export default router;