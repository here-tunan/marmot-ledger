<!--首页的Header控件-->
<template>
  <div class="header">
    <!-- 左侧区域 -->
    <div class="header-left">
      <!--  折叠展开按钮  -->
      <div class="collapse-btn" @click="collapseChange">
        <el-icon v-if="sidebar.collapse">
          <Expand/>
        </el-icon>
        <el-icon v-else>
          <Fold/>
        </el-icon>
      </div>

      <div class="logo">This is my life!</div>
    </div>

    <!-- 右侧区域 -->
    <div class="header-right">
      <div class="header-user-con">
        <!-- 用户头像 -->
        <el-avatar class="user-avatar" :size="30" :src="avatarImg"/>

        <!-- 用户名下拉菜单 -->
        <el-dropdown class="user-name" trigger="click" @command="handleCommand" placement="bottom-end">
					<span class="el-dropdown-link">
						{{ username }}
						<el-icon class="el-icon--right">
							<arrow-down/>
						</el-icon>
					</span>
          <template #dropdown>
            <el-dropdown-menu>
              <a href="https://gitee.com/yaodao666/vue-my-life" target="_blank">
                <el-dropdown-item>项目仓库</el-dropdown-item>
              </a>
              <el-dropdown-item command="user">个人中心</el-dropdown-item>
              <el-dropdown-item divided command="loginOut">退出登录</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </div>
  </div>
</template>

<script setup>
import {useSidebarStore} from "@/stores/sidebar";
import {ArrowDown, Expand, Fold} from "@element-plus/icons-vue";
import {computed, onMounted} from "vue";
import {useRouter} from "vue-router";
import {getLoginUserInfo} from "@/api/user/user";
import {ElMessage} from "element-plus";
import {useUserStore} from "@/stores/user";

const sidebar = useSidebarStore()

const avatarImg = computed(() => {
  return useUserStore().avatar
})

const username = computed(() => {
  return useUserStore().username
})

// 用户名下拉菜单选择事件
const router = useRouter();

onMounted(() => {
  if (document.body.clientWidth < 1000) {
    collapseChange();
  }
});

onMounted(() => {
  getLoginUserInfo().then(
      (res) => {
        if (res.success) {
          useUserStore().account = res.data.account
          useUserStore().username = res.data.name
          useUserStore().desc = res.data.desc
          useUserStore().avatar = res.data.avatar
          localStorage.setItem("avatar", res.data.avatar)
        } else {
          ElMessage.error("获取用户信息失败")
        }
      }
  )
})


const handleCommand = (command) => {
  if (command === 'loginOut') {
    localStorage.removeItem('token');
    ElMessage.success('成功退出，请重新登录')
    router.push('/login');
  } else if (command === 'user') {
    router.push('/user');
  }
};

const collapseChange = () => {
  sidebar.handleCollapse();
};
</script>

<style scoped>
.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-sizing: border-box;
  width: 100%;
  height: 70px;
  padding: 0 24px;
  font-size: 22px;
  color: #fff;
  background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
  box-shadow: 0 4px 20px rgba(59, 130, 246, 0.15);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.collapse-btn {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 50px;
  padding: 0 15px;
  cursor: pointer;
  transition: all 0.3s ease;
  border-radius: 8px;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: scale(1.05);
}

.header .logo {
  font-weight: 700;
  font-size: 24px;
  background: linear-gradient(45deg, #ffffff, #dbeafe);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin-left: 16px;
  letter-spacing: -0.5px;
}

.header-left {
  display: flex;
  align-items: center;
}

.header-right {
  display: flex;
  align-items: center;
}

.header-user-con {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-avatar {
  border: 2px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.user-avatar:hover {
  border-color: rgba(255, 255, 255, 0.8);
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.2);
}

.user-name {
  margin-left: 10px;
}

.el-dropdown-link {
  color: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  padding: 8px 12px;
  border-radius: 6px;
  transition: all 0.3s ease;
  font-weight: 500;
}

.el-dropdown-link:hover {
  background: rgba(255, 255, 255, 0.1);
}

:deep(.el-dropdown-menu__item) {
  text-align: center;
  transition: all 0.2s ease;
}

:deep(.el-dropdown-menu__item:hover) {
  background: var(--primary-color, #409eff);
  color: #fff;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .header {
    height: 60px;
    font-size: 18px;
    padding: 0 15px;
  }
  
  .header .logo {
    font-size: 18px;
    margin-left: 10px;
  }
  
  .collapse-btn {
    height: 40px;
    padding: 0 10px;
  }
  
  .header-user-con {
    gap: 12px;
  }
  
  .user-avatar {
    width: 28px !important;
    height: 28px !important;
  }
}

@media (max-width: 480px) {
  .header {
    padding: 0 10px;
  }
  
  .header .logo {
    font-size: 16px;
    margin-left: 8px;
  }
  
  .header-user-con {
    gap: 8px;
  }
}
</style>