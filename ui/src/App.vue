<template>
  <div :class="globalMode.mode">
    <el-config-provider :locale="elementLocale">
      <router-view />
    </el-config-provider>
  </div>
</template>
<script setup>
import { computed, watch } from 'vue';
import zhCn from 'element-plus/es/locale/lang/zh-cn';
import en from 'element-plus/es/locale/lang/en';
import dayjs from 'dayjs';
import 'dayjs/locale/zh-cn';
import 'dayjs/locale/en';
import './assets/css/main.css';
import './assets/css/color-dark.css';
import './assets/css/color-light.css';
import {useConfigStore} from "@/stores/config";

const globalMode = useConfigStore()

const elementLocale = computed(() => globalMode.locale === 'en-US' ? en : zhCn)

watch(
  () => globalMode.locale,
  (locale) => {
    dayjs.locale(locale === 'en-US' ? 'en' : 'zh-cn')
  },
  { immediate: true }
)
</script>
