// 该文件必须挂载到main.js中才能生效

import router from './router'
import {isNeedRefresh, refreshToken, removeToken, validToken2} from './api/auth/auth'
import {getCurrentUserRole} from './api/user/user'
import {ElMessage} from 'element-plus'
import {t} from './i18n'

const publicPages = ['/login', '/register']

async function ensureAdmin() {
    const res = await getCurrentUserRole()
    return Boolean(res.success && res.data?.role === 'admin')
}

router.beforeEach((to, from, next) => {
    if (publicPages.includes(to.path)) {
        next()
        return
    }

    let token = localStorage.getItem("token")

    if (token == null || token === '') {
        ElMessage.warning(t('auth.notLoggedIn'))
        next('/login')
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

        const continueNavigation = () => {
            if (to.path === '/login') {
                ElMessage.success(t('auth.alreadyLoggedIn'))
                next("/")
                return
            }

            if (to.meta?.requiresAdmin) {
                ensureAdmin().then(isAdmin => {
                    if (isAdmin) {
                        next()
                    } else {
                        ElMessage.warning(t('auth.adminRequired'))
                        next('/dashboard')
                    }
                }).catch(() => {
                    ElMessage.warning(t('auth.adminRequired'))
                    next('/dashboard')
                })
                return
            }

            next()
        }

        if (isNeedRefresh()) {
            refreshToken().then(res => {
                if (res.success) {
                    localStorage.setItem("token", res.data.token)
                    let expiredSeconds = res.data.expiredTime
                    let now = new Date()
                    now.setTime(now.getTime() + expiredSeconds * 1000)
                    localStorage.setItem("validTime", now.getTime())
                    continueNavigation()
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
            continueNavigation()
        }
    });
})
