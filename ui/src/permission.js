// 该文件必须挂载到main.js中才能生效

import router from './router'
import {isNeedRefresh, refreshToken, removeToken, validToken2} from './api/auth/auth'
import {ElMessage} from 'element-plus'
import {t} from './i18n'

router.beforeEach((to, from, next) => {
    let token = localStorage.getItem("token")

    if (token == null || token === '') {
        if (to.path !== '/login') {
            ElMessage.warning(t('auth.notLoggedIn'))
            next('/login')
        } else {
            next()
        }
        return
    }

    validToken2(token).then(isValid => {
        if (!isValid) {
            if (to.path === '/login') {
                next()
                return;
            }
            ElMessage.warning(t('auth.invalid'))
            removeToken()
            next('/login')
            return;
        }

        if (isNeedRefresh()) {
            refreshToken().then(res => {
                if (res.success) {
                    localStorage.setItem("token", res.data.token)
                    let expiredSeconds = res.data.expiredTime
                    let now = new Date()
                    now.setTime(now.getTime() + expiredSeconds * 1000)
                    localStorage.setItem("validTime", now.getTime())
                    if (to.path === '/login') {
                        ElMessage.success(t('auth.alreadyLoggedIn'))
                        next("/")
                    } else {
                        next()
                    }
                } else {
                    removeToken()
                    if (to.path === '/login') {
                        next()
                    } else {
                        ElMessage.warning(t('auth.invalid'))
                        next('/login')
                    }
                }
            })
        } else {
            if (to.path === '/login') {
                ElMessage.success(t('auth.alreadyLoggedIn'))
                next("/")
            } else {
                next()
            }
        }
    });
})