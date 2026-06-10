import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Dashboard from "@/views/Dashboard.vue";
import Login from "@/views/login/Login.vue"
import User from "@/views/User.vue";

const router = createRouter({
    history: createWebHashHistory(),
    routes: [
        {
            path: '/',
            redirect: '/dashboard',
        },
        {
            path: '/login',
            name: 'login',
            component: Login,
        },
        {
            path: '/',
            name: 'home',
            component: HomeView,
            children: [
                {
                    path: '/dashboard',
                    name: 'dashboard',
                    meta: {
                        title: '系统首页',
                    },
                    component: Dashboard,
                },
                {
                    path: '/user',
                    name: 'user',
                    component: User,
                    meta: {
                        title: "用户中心"
                    }
                },
            ]
        }
    ]
})

export default router
