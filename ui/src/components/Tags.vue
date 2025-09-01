<template>
	<div class="tags" v-if="tags.show">
		<ul>
			<li
				class="tags-li"
				v-for="(item, index) in tags.list"
				:class="{ active: isActive(item.path) }"
				:key="index"
			>
				<router-link :to="item.path" class="tags-li-title">{{ item.title }}</router-link>
				<el-icon @click="closeTags(index)"><Close/></el-icon>
			</li>
		</ul>
	</div>
</template>

<script setup lang="ts">
import { useTagsStore } from '@/stores/tags';
import {Close} from "@element-plus/icons-vue";

import { onBeforeRouteUpdate, useRoute, useRouter } from 'vue-router';

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
      // 根据路由中的meta中的title字段设置标签名称
      title: route.meta.title,
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
	height: 48px;
	background: rgba(248, 250, 252, 0.95);
	backdrop-filter: blur(10px);
	border-bottom: 1px solid #e2e8f0;
	transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	padding: 0 24px;
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
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
	border-radius: 16px;
	font-size: 13px;
	overflow: hidden;
	cursor: pointer;
	height: 32px;
	border: 1px solid #e2e8f0;
	background: #ffffff;
	padding: 0 10px 0 14px;
	color: #64748b;
	transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
	box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
}

.tags-li:not(.active):hover {
	background: #f8fafc;
	border-color: #cbd5e1;
	color: #475569;
	transform: translateY(-1px);
	box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.tags-li.active {
	background: linear-gradient(135deg, #3b82f6 0%, #1d4ed8 100%);
	border-color: #3b82f6;
	color: #ffffff;
	box-shadow: 0 4px 12px rgba(59, 130, 246, 0.3);
}

.tags-li-title {
	max-width: 120px;
	overflow: hidden;
	white-space: nowrap;
	text-overflow: ellipsis;
	margin-right: 8px;
	color: inherit;
	text-decoration: none;
	font-weight: 500;
}

.tags-li .el-icon {
	font-size: 14px;
	opacity: 0.7;
	transition: all 0.2s ease;
}

.tags-li .el-icon:hover {
	opacity: 1;
	transform: scale(1.1);
}

.tags-li.active .el-icon {
	opacity: 0.9;
}




/* 响应式设计 */
@media (max-width: 768px) {
	.tags {
		height: 35px;
		padding: 0 15px;
	}
	
	.tags-li {
		height: 24px;
		font-size: 12px;
		padding: 0 6px 0 12px;
		margin-right: 6px;
	}
	
	.tags-li-title {
		max-width: 60px;
		margin-right: 6px;
	}
	
}

@media (max-width: 480px) {
	.tags {
		padding: 0 10px;
	}
	
	
	.tags ul {
		padding-right: 10px;
	}
}
</style>
