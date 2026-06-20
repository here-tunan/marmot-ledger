import axios from "axios";
import {ElMessage} from "element-plus";
import router from "@/router";
import {isNeedRefresh, refreshToken, removeToken} from "@/api/auth/auth";
import {useAuthStore} from "@/stores/auth";
import {t} from "@/i18n";

// create an axios instance
const service = axios.create({
    baseURL: import.meta.env.VITE_APP_URL,
    timeout: 50000,
    headers: {
        'Content-Type': 'application/json',
    },
})

// request interceptor
service.interceptors.request.use(
    config => {
        const token = localStorage.getItem('token');
        if (token) {
            config.headers['Authorization'] = token;

            if (!config.url.endsWith('/user/refreshToken') && isNeedRefresh()) {
                refreshToken().then(res => {
                    if (res.success) {
                        localStorage.setItem('token', res.data.token);
                        const expiredSeconds = res.data.expiredTime;
                        const now = new Date();
                        now.setTime(now.getTime() + expiredSeconds * 1000);
                        localStorage.setItem('validTime', now.getTime());

                        config.headers['Authorization'] = `${res.data.token}`;
                        useAuthStore().isRefreshing = false
                        return config;
                    } else {
                        removeToken()
                        router.push('/login');
                        useAuthStore().isRefreshing = false
                        return Promise.reject(new Error('Token refresh failed'));
                    }
                })
            }
        }
        return config
    },
    error => {
        return Promise.reject(error)
    }
)

// response interceptor
service.interceptors.response.use(
    response => {
        let url = response.config.url
        const data = response.data

        if (data.code !== '' && data.code !== 200) {
            if (data.code === 401) {
                if (!url.endsWith("/validToken") && !url.endsWith("/register")) {
                    if (!useAuthStore().isAxios401Failing) {
                        useAuthStore().isAxios401Failing = true
                        ElMessage.warning(t('auth.expired'))
                        router.push('/login')

                        setTimeout(() => {
                            useAuthStore().isAxios401Failing = false
                        }, 300)
                    }
                }
            }
        }

        return data
    },

    error => {
        console.log('axios instance response interceptors, error', error)
        ElMessage({
            message: t('messages.networkError'),
            type: 'error',
        })
        return Promise.reject(error)
    }
)

export default service