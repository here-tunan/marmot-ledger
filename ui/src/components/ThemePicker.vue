<template>
  <div class="theme-picker">
    <!-- 明暗模式切换 -->
    <el-tooltip content="切换明暗模式" placement="bottom">
      <div class="mode-toggle" @click="toggleMode">
        <el-icon class="mode-icon" :class="{ 'dark': config.isDark }">
          <Sunny v-if="!config.isDark" />
          <Moon v-else />
        </el-icon>
      </div>
    </el-tooltip>

    <!-- 主题颜色选择 -->
    <el-dropdown trigger="click" @command="handleThemeChange">
      <el-tooltip content="选择主题颜色" placement="bottom">
        <div class="theme-toggle">
          <div 
            class="current-theme" 
            :style="{ background: config.currentTheme.gradient }"
          >
            <el-icon class="theme-icon">
              <Brush />
            </el-icon>
          </div>
        </div>
      </el-tooltip>
      
      <template #dropdown>
        <el-dropdown-menu class="theme-dropdown">
          <div class="theme-dropdown-header">
            <span>选择主题颜色</span>
          </div>
          <div class="theme-colors">
            <div 
              v-for="(theme, key) in THEMES" 
              :key="key"
              class="theme-color-item"
              :class="{ 'active': config.theme === key }"
              @click.stop="handleThemeChange(key)"
            >
              <div 
                class="color-preview" 
                :style="{ background: theme.gradient }"
              >
                <el-icon v-if="config.theme === key" class="check-icon">
                  <Check />
                </el-icon>
              </div>
              <span class="color-name">{{ theme.name }}</span>
            </div>
          </div>
        </el-dropdown-menu>
      </template>
    </el-dropdown>
  </div>
</template>

<script setup>
import { useConfigStore, THEMES } from '@/stores/config'
import { Sunny, Moon, Brush, Check } from '@element-plus/icons-vue'
import { onMounted } from 'vue'

const config = useConfigStore()

const toggleMode = () => {
  config.toggle()
}

const handleThemeChange = (themeName) => {
  config.setTheme(themeName)
}

onMounted(() => {
  config.initTheme()
})
</script>

<style scoped>
.theme-picker {
  display: flex;
  align-items: center;
  gap: 12px;
}

.mode-toggle,
.theme-toggle {
  width: 36px;
  height: 36px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.mode-toggle:hover,
.theme-toggle:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

.mode-icon,
.theme-icon {
  font-size: 18px;
  color: white;
  transition: all 0.3s ease;
}

.mode-icon.dark {
  color: #ffd93d;
}

.current-theme {
  width: 24px;
  height: 24px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.2);
}

.theme-dropdown {
  padding: 0;
  min-width: 200px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.theme-dropdown-header {
  padding: 16px 20px;
  background: var(--gradient-bg, linear-gradient(135deg, #667eea 0%, #764ba2 100%));
  color: white;
  font-weight: 600;
  font-size: 14px;
}

.theme-colors {
  padding: 16px;
  background: white;
}

.theme-color-item {
  display: flex;
  align-items: center;
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  gap: 12px;
}

.theme-color-item:hover {
  background: #f5f7fa;
}

.theme-color-item.active {
  background: #e8f4fd;
}

.color-preview {
  width: 32px;
  height: 32px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.15);
  position: relative;
}

.check-icon {
  color: white;
  font-size: 16px;
  font-weight: bold;
}

.color-name {
  font-size: 14px;
  color: #333;
  font-weight: 500;
}

.theme-color-item.active .color-name {
  color: var(--primary-color, #409eff);
  font-weight: 600;
}

/* 暗黑模式适配 */
:deep(.dark) .theme-colors {
  background: #2d2d2d;
}

:deep(.dark) .theme-color-item:hover {
  background: #404040;
}

:deep(.dark) .theme-color-item.active {
  background: #1a365d;
}

:deep(.dark) .color-name {
  color: #e2e8f0;
}
</style>