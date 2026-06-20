<template>
  <div class="icon-color-picker">
    <!-- 统一的两行网格布局 -->
    <div class="picker-row">
      <!-- 图标选择 -->
      <div class="picker-col">
        <div class="picker-label">
          <span>{{ iconLabel }}</span>
          <span class="selected-preview" v-if="selectedIcon">{{ selectedIcon }}</span>
        </div>
        <div class="picker-grid icon-grid">
          <button
            v-for="icon in iconList"
            :key="icon"
            type="button"
            :class="['grid-btn icon-btn', { active: selectedIcon === icon }]"
            @click="selectIcon(icon)"
          >
            {{ icon }}
          </button>
        </div>
        <div class="picker-input">
          <el-input
            v-model="customIcon"
            placeholder="自定义"
            maxlength="8"
            size="small"
            clearable
            @change="handleCustomIconChange"
          />
        </div>
      </div>

      <!-- 颜色选择 -->
      <div class="picker-col" v-if="showColor">
        <div class="picker-label">
          <span>{{ colorLabel }}</span>
          <span class="color-preview" v-if="selectedColor" :style="{ background: selectedColor }"></span>
        </div>
        <div class="picker-grid color-grid">
          <button
            v-for="color in colorList"
            :key="color"
            type="button"
            :class="['grid-btn color-btn', { active: selectedColor === color }]"
            :style="{ background: color }"
            @click="selectColor(color)"
          ></button>
        </div>
        <div class="picker-input">
          <el-input
            v-model="customColor"
            placeholder="自定义"
            maxlength="9"
            size="small"
            clearable
            @change="handleCustomColorChange"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  iconValue: { type: String, default: '' },
  colorValue: { type: String, default: '' },
  showColor: { type: Boolean, default: true },
  iconLabel: { type: String, default: '图标' },
  colorLabel: { type: String, default: '颜色' }
})

const emit = defineEmits(['update:iconValue', 'update:colorValue'])

const iconList = [
  '💰', '💸', '🏠', '🚗', '🍔', '🛍️', '✈️', '🎮',
  '📚', '🎁', '💼', '🎓'
]

const colorList = [
  '#ef4444', '#f97316', '#f59e0b', '#84cc16',
  '#22c55e', '#3b82f6', '#8b5cf6', '#ec4899',
  '#2f7d5c', '#64748b', '#0ea5e9', '#facc15'
]

const selectedIcon = computed({
  get: () => props.iconValue,
  set: (val) => emit('update:iconValue', val)
})

const selectedColor = computed({
  get: () => props.colorValue,
  set: (val) => emit('update:colorValue', val)
})

const customIcon = ref('')
const customColor = ref('')

watch(() => props.iconValue, (val) => {
  if (!iconList.includes(val)) {
    customIcon.value = val
  }
}, { immediate: true })

watch(() => props.colorValue, (val) => {
  customColor.value = val || ''
}, { immediate: true })

function selectIcon(icon) {
  customIcon.value = ''
  selectedIcon.value = icon
}

function selectColor(color) {
  customColor.value = color
  selectedColor.value = color
}

function handleCustomIconChange() {
  if (customIcon.value) {
    selectedIcon.value = customIcon.value
  }
}

function handleCustomColorChange() {
  if (customColor.value) {
    selectedColor.value = customColor.value
  }
}
</script>

<style scoped>
.icon-color-picker {
  width: 100%;
}

/* 一行两列布局 */
.picker-row {
  display: flex;
  gap: 16px;
  align-items: stretch;
}

.picker-col {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.picker-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  font-weight: 600;
  color: #64748b;
  height: 24px;
}

.selected-preview,
.color-preview {
  margin-left: auto;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.selected-preview {
  background: #f8faf7;
  font-size: 14px;
}

.color-preview {
  border: 2px solid #e2e8f0;
}

/* 统一的网格容器 - 4列3行，两边完全对称 */
.picker-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 8px;
  padding: 12px;
  background: #f8fafc;
  border-radius: 10px;
}

/* 统一的按钮尺寸 */
.grid-btn {
  width: 100%;
  aspect-ratio: 1;
  border: 2px solid transparent;
  border-radius: 8px;
  background: #fff;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.grid-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.1);
}

.grid-btn.active {
  border-color: #2f7d5c;
  background: #dcfce7;
  transform: scale(1.05);
  box-shadow: 0 2px 6px rgba(47, 125, 92, 0.15);
}

.icon-btn {
  font-size: 16px;
}

.color-btn {
  border-radius: 50%;
  box-shadow: inset 0 1px 3px rgba(0, 0, 0, 0.1);
}


/* 输入框统一高度 */
.picker-input {
  height: 32px;
}

.picker-input :deep(.el-input__wrapper) {
  box-shadow: 0 0 0 1px #e2e8f0 inset;
}
</style>
