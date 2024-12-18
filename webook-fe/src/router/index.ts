// src/router/index.ts

import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import login from '../views/login/login.vue'
import register from '../views/register/register.vue'
const routes: Array<RouteRecordRaw> = [
    {
        path: '/',
        name: 'Home',
        redirect:'/login'
    },
    {
        path: '/login',
        name: '登录界面',
        component: login,
    },
    {
        path: '/register',
        name: '注册界面',
        component: register,
    },

];

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes,
});

// 在这里添加路由的导航守卫
router.beforeEach((to, from, next) => {
    console.log('Navigating to:', to.path);
    next();
});

export default router;