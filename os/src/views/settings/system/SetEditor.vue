<template>
	<div class="editor-settings">
		<el-form :model="settingStore.config" label-position="top" class="form-container">
			<el-form-item label="编辑器类型">
				<el-select v-model="settingStore.config.editorType">
					<el-option
						v-for="(item, key) in editorType"
						:key="key"
						:label="item.title"
						:value="item.value"
					/>
				</el-select>
			</el-form-item>
			<template v-if="settingStore.config.editorType === 'onlyoffice'">
				<el-form-item label="地址">
					<el-input
						v-model="settingStore.config.onlyoffice.url"
						placeholder="https://godoos.com/onlyoffice 不要加斜杠"
					/>
				</el-form-item>
				<!-- <el-form-item label="私钥">
					<el-input v-model="settingStore.config.onlyoffice.sceret" />
				</el-form-item> -->
			</template>
			<el-form-item>
				<el-button type="primary">
					确认
				</el-button>
			</el-form-item>
		</el-form>
	</div>
</template>

<script lang="ts" setup>
	import { useSettingsStore } from "@/stores/settings";

	const settingStore = useSettingsStore();

	const editorType = [
		{
			title: "系统默认",
			value: "local",
		},
		{
			title: "OnlyOffice",
			value: "onlyoffice",
		},
	];
</script>

<style scoped>
	.editor-settings {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 100%;
		height: 100%;
	}

	.form-container {
		width: 22rem;
	}

	.setting-item {
		display: flex;
		align-items: center;
		margin-bottom: 20px;
	}

	.setting-item label {
		width: 100px;
	}

	.el-select,
	.el-input {
		flex-grow: 1;
	}

	.el-button {
		width: 100%;
	}
</style>
