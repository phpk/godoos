<template>
	<el-form
		label-width="auto"
		style="max-width: 560px; margin-top: 20px; padding: 20px"
	>
		<el-form-item label="文件密码">
			<el-input
				v-model="filePwd"
				type="password"
				show-password
			/>
		</el-form-item>
		<div class="btn-group">
			<el-button
				type="primary"
				@click="setFilePwd"
				>提交</el-button
			>
		</div>
	</el-form>
</template>

<script lang="ts" setup>
	import { BrowserWindow, useSystem } from "@/system";
	import { notifyError, notifySuccess } from "@/util/msg";
	import { ref } from "vue";
	const window: BrowserWindow | undefined = inject("browserWindow");
	const filePwd = ref("");
	const sys = useSystem();
	async function setFilePwd() {
		if (
			filePwd.value !== "" &&
			filePwd.value.length >= 6 &&
			filePwd.value.length <= 10
		) {
			const path = window?.config.path || "";
			const header = {
				pwd: filePwd.value,
			};
			const file = await sys.fs.readFile(path);
			if (file === false) return;
			const res = await sys.fs.writeFile(path, file, header);
			if (res && res.success) {
				notifySuccess("文件密码设置成功");
			} else {
				notifyError("文件密码设置失败");
			}
			//console.log("路径：", res, path);
		}
	}
</script>
<style scoped>
	.btn-group {
		display: flex;
		justify-content: center;
	}
</style>
