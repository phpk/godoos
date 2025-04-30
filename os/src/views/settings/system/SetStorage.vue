<template>
	<div class="container">
		<div class="setting">
			<el-form
				:model="settingStore.config"
				label-position="top"
				class="form"
			>
				<el-form-item label="存储方式">
					<el-select v-model="settingStore.config.storeType">
						<el-option
							v-for="(item, key) in storeList"
							:key="key"
							:label="item.title"
							:value="item.value"
						/>
					</el-select>
				</el-form-item>

				<template v-if="settingStore.config.storeType === 'local'">
					<el-form-item label="存储地址">
						<el-input
							v-model="settingStore.config.storePath"
							placeholder="可为空，为空则取系统默认存储地址:/用户目录/.godoos/os"
						/>
					</el-form-item>
					<el-form-item label="自定义端口">
						<el-input
							v-model="settingStore.config.netPort"
							placeholder="可为空，为空则取系统默认56780"
						/>
					</el-form-item>
					<el-form-item label="自定义路径">
						<el-input
							v-model="settingStore.config.netPath"
							placeholder="自定义web访问路径，英文，不要加斜杠，可为空"
						/>
					</el-form-item>
				</template>

				<template v-if="settingStore.config.storeType === 'net'">
					<el-form-item label="服务器地址">
						<el-input
							v-model="settingStore.config.storenet.url"
							placeholder="可访问的地址，例如http://192.168.1.6:56780 不要加斜杠"
						/>
					</el-form-item>
					<el-form-item label="允许跨域">
						<el-switch
							v-model="settingStore.config.storenet.isCors"
							active-text="允许"
							inactive-text="不允许"
						/>
					</el-form-item>
				</template>

				<template v-if="settingStore.config.storeType === 'webdav'">
					<el-form-item label="服务器地址">
						<el-input
							v-model="settingStore.config.webdavClient.url"
							placeholder="https://godoos.com/webdav 不要加斜杠"
						/>
					</el-form-item>
					<el-form-item label="登陆用户名">
						<el-input
							v-model="settingStore.config.webdavClient.username"
						/>
					</el-form-item>
					<el-form-item label="登陆密码">
						<el-input
							v-model="settingStore.config.webdavClient.password"
							type="password"
						/>
					</el-form-item>
				</template>

				<el-form-item>
					<el-button type="primary">
						{{ t("confirm") }}
					</el-button>
				</el-form-item>
			</el-form>
		</div>
	</div>
</template>

<script setup>
	import { ref } from "vue";
	import { useSettingsStore } from "@/stores/settings";
	const settingStore = useSettingsStore();
	const activeIndex = ref(1);

	const storeList = ref([
		{ title: "本地", value: "local" },
		{ title: "网络", value: "net" },
		{ title: "WebDAV", value: "webdav" },
	]);

	const t = (key) => {
		const translations = {
			confirm: "确认",
		};
		return translations[key] || key;
	};
</script>

<style scoped>
	.container {
		display: flex;
		justify-content: center;
	}

	.setting {
		width: 100%;
		max-width: 22rem;
	}

	.form {
		width: 100%;
	}

	.el-button {
		width: 100%;
	}
</style>
