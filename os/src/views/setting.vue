<template>
	<div class="settings">
		<div class="sidebar">
			<div class="profile">
				<el-avatar src="" alt="Avatar" class="avatar" />
				<div class="username">用户名</div>
			</div>
			<el-input placeholder="Find a setting" class="search" />
			<el-menu class="menu" default-active="1" @select="handleMenuSelect">
				<el-menu-item class="menu-item" v-for="item in settingStore.settingList" :index="item.key">
					<icon :name="item.icon" size="18" />
					{{ item.title }}
				</el-menu-item>
			</el-menu>
		</div>
		<div class="content">
			<div v-if="!showContent">
				<h1 class="content-title">{{ currentSetting.title }}</h1>
				<el-menu class="content-menu">
					<el-menu-item @click="navigateTo(child)" class="content-menu-item"
						v-for="child in currentSetting.children" :key="child.key">
						<div class="content-menu-item-left">
							<icon :name="child.icon" size="18" />
							{{ child.title }}
						</div>
						<div class="content-menu-item-right">
							<el-icon :size="16">
								<ArrowRight />
							</el-icon>
						</div>
					</el-menu-item>
				</el-menu>
			</div>
			<template v-else>
				<div class="setting-item-title">
					<el-icon @click="showContent = false">
						<ArrowLeft />
					</el-icon>
					<span style="display: flex; align-items: center; gap: 8px;">
						<icon :name="currentChildren.icon" size="18" />
						{{ currentChildren.title }}
					</span>
				</div>
				<component :is="currentContent" />
			</template>
		</div>
	</div>
</template>

<script setup lang="ts">
import { useSettingsStore } from "@/stores/settings";
import { ref, defineAsyncComponent } from "vue";

const settingStore = useSettingsStore();
const currentSetting: any = ref(settingStore.settingList[0]);
const currentChildren = ref<any>(null);
const showContent = ref(false);
const currentContent = ref(null);

function handleMenuSelect(key: string) {
	currentSetting.value = settingStore.settingList.find(
		(item) => item.key === key
	);
	showContent.value = false;;
}

const navigateTo = async (item: any) => {
	showContent.value = true;
	currentChildren.value = item;
	currentContent.value = defineAsyncComponent(() =>
		import(`/src/views/settings/${item.content}.vue`).then((module) => module.default)
	);
};
</script>
<style lang="scss" scoped>
@use "@/styles/setting.scss";
</style>