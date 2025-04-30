<template>
	<div class="container">
		<div class="setting">
			<el-form
				class="setting-item"
				:model="settingStore.config.userInfo"
				label-position="top"
				label-width="100px"
			>
				<el-form-item label="服务器地址">
					<el-input
						v-model="settingStore.config.userInfo.url"
						placeholder="网址或域名，例子：https://godoos.com 不要加斜杠"
						class="input"
					/>
				</el-form-item>

				<el-form-item label="用户名">
					<el-input
						v-model="settingStore.config.userInfo.username"
						placeholder="登录用户名"
						class="input"
					/>
				</el-form-item>

				<el-form-item label="密码">
					<el-input
						v-model="settingStore.config.userInfo.password"
						type="password"
						placeholder="登录密码"
						class="input"
					/>
				</el-form-item>

				<el-form-item>
					<el-button
						type="primary"
						class="submit-button"
						>提交</el-button
					>
				</el-form-item>
			</el-form>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import { useSettingsStore } from "@/stores/settings";
	const settingStore = useSettingsStore();
</script>

<style scoped>
	.setting-item-title {
		display: flex;
		align-items: center;
		gap: 8px;
		padding-bottom: 20px;
		.el-icon {
			cursor: pointer;
			width: 20px;
			height: 20px;
		}
		.el-icon:hover {
			color: #409eff;
		}
	}

	.container {
		display: flex;
		justify-content: center;
		align-items: center;
	}

	.setting {
		width: 100%;
		max-width: 500px;
		border-radius: 8px;
	}

	.setting-item {
		margin-bottom: 20px;
		width: 22rem;
		margin: 0 auto;
	}

	.input {
		width: 100%;
		border-radius: 4px;
	}

	.submit-button {
		width: 100%;
		background-color: #409eff;
		color: #fff;
		border-radius: 4px;
		padding: 10px 0;
		font-size: 16px;
		font-weight: bold;
	}

	.submit-button:hover {
		background-color: #66b1ff;
	}
</style>
