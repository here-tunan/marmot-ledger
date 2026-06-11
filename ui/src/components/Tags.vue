<template>
	<div class="tags" v-if="tags.show">
		<ul>
			<li
				class="tags-li"
				v-for="(item, index) in tags.list"
				:class="{ active: isActive(item.path) }"
				:key="item.path"
			>
				<router-link :to="item.path" class="tags-li-title">{{ t(item.titleKey) }}</router-link>
				<button class="tag-close" :aria-label="t('common.actions.delete')" @click="closeTags(index)">
					<el-icon><Close/></el-icon>
				</button>
			</li>
		</ul>
	</div>
</template>

<script setup lang="ts">
import { useTagsStore } from '@/stores/tags';
import {Close} from "@element-plus/icons-vue";
import { useI18n } from 'vue-i18n';
import { onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const isActive = (path: string) => {
  return path === route.fullPath;
};

const tags = useTagsStore();
// 关闭单个标签
const closeTags = (index: number) => {
  const delItem = tags.list[index];
  tags.delTagsItem(index);
  const item = tags.list[index] ? tags.list[index] : tags.list[index - 1];
  if (item) {
    delItem.path === route.fullPath && router.push(item.path);
  } else {
    router.push('/');
  }
};

// 设置标签
const setTags = (route: any) => {
  const isExist = tags.list.some(item => {
    return item.path === route.fullPath;
  });
  if (!isExist) {
    if (tags.list.length >= 8) tags.delTagsItem(0);
    tags.setTagsItem({
      name: route.name,
      titleKey: route.meta.titleKey || 'routes.dashboard',
      path: route.fullPath
    });
  }
};

setTags(route);

onBeforeRouteUpdate(to => {
  setTags(to);
});

</script>

<style scoped>
.tags {
	display: flex;
	align-items: center;
	height: 44px;
	background: rgba(247, 245, 239, 0.92);
	border-bottom: 1px solid rgba(100, 116, 139, 0.12);
	padding: 0 24px;
}

.tags ul {
	display: flex;
	align-items: center;
	flex: 1;
	margin: 0;
	padding: 0;
	overflow-x: auto;
	overflow-y: hidden;
}

.tags-li {
	display: flex;
	align-items: center;
	margin-right: 8px;
	border-radius: 999px;
	font-size: 13px;
	overflow: hidden;
	cursor: pointer;
	height: 30px;
	background: rgba(255, 255, 255, 0.72);
	padding: 0 6px 0 14px;
	color: #64748b;
	transition-property: transform, background-color, color, box-shadow;
	transition-duration: 160ms;
	box-shadow: inset 0 0 0 1px rgba(100, 116, 139, 0.12);
}

.tags-li:not(.active):hover {
	background: #ffffff;
	color: #1e293b;
	transform: translateY(-1px);
}

.tags-li.active {
	background: #1f2933;
	color: rgba(255,255,255,0.9);
	box-shadow: 0 8px 18px rgba(31, 41, 51, 0.16);
}

.tags-li-title {
	max-width: 120px;
	overflow: hidden;
	white-space: nowrap;
	text-overflow: ellipsis;
	margin-right: 6px;
	color: inherit;
	text-decoration: none;
	font-weight: 600;
}

.tag-close {
	display: inline-flex;
	align-items: center;
	justify-content: center;
	width: 22px;
	height: 22px;
	border: 0;
	border-radius: 999px;
	background: transparent;
	color: inherit;
	cursor: pointer;
	opacity: 0.72;
	transition-property: transform, opacity, background-color;
	transition-duration: 140ms;
}

.tag-close:hover {
	opacity: 1;
	background: rgba(255, 255, 255, 0.16);
}

.tag-close:active {
	transform: scale(0.9);
}

@media (max-width: 768px) {
	.tags {
		height: 36px;
		padding: 0 15px;
	}

	.tags-li {
		height: 26px;
		font-size: 12px;
		padding: 0 4px 0 10px;
		margin-right: 6px;
	}

	.tags-li-title {
		max-width: 70px;
	}
}

@media (prefers-reduced-motion: reduce) {
	.tags-li,
	.tag-close {
		transition: none;
	}
}
</style>
