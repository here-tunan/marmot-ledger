import {createRouter, createWebHashHistory} from 'vue-router'
import HomeView from '../views/HomeView.vue'
import Dashboard from "@/views/Dashboard.vue";
import Login from "@/views/login/Login.vue"
import User from "@/views/User.vue";
import Accounts from "@/views/Accounts.vue";
import Buckets from "@/views/Buckets.vue";
import Record from "@/views/Record.vue";
import Records from "@/views/Records.vue";
import Categories from "@/views/Categories.vue";
import Family from "@/views/Family.vue";

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
                        titleKey: 'routes.dashboard',
                    },
                    component: Dashboard,
                },
                {
                    path: '/accounts',
                    name: 'accounts',
                    meta: {
                        titleKey: 'routes.accounts',
                    },
                    component: Accounts,
                },
                {
                    path: '/categories',
                    name: 'categories',
                    meta: {
                        titleKey: 'routes.categories',
                    },
                    component: Categories,
                },
                {
                    path: '/buckets',
                    name: 'buckets',
                    meta: {
                        titleKey: 'routes.buckets',
                    },
                    component: Buckets,
                },
                {
                    path: '/record',
                    name: 'record',
                    meta: {
                        titleKey: 'routes.record',
                    },
                    component: Record,
                },
                {
                    path: '/records',
                    name: 'records',
                    meta: {
                        titleKey: 'routes.records',
                    },
                    component: Records,
                },
                {
                    path: '/family',
                    name: 'family',
                    meta: {
                        titleKey: 'routes.family',
                    },
                    component: Family,
                },
                {
                    path: '/user',
                    name: 'user',
                    component: User,
                    meta: {
                        titleKey: "routes.userCenter"
                    }
                },
            ]
        }
    ]
})

export default router
