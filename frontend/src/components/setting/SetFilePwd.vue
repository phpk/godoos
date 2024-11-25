<template>
	<div class="file-pwd-box">
		<!-- <div v-if="setPwd">
			<div class="setting-item">
				<label>文件密码</label>
				<el-input
					v-model="filePwd"
					placeholder="请设置6-10位的密码"
					type="password"
					show-password
				/>
			</div>
			<div class="setting-item">
				<label></label>
				<el-button
					@click="toSetFilePwd"
					type="primary"
					>{{ t("setFilePwd") }}</el-button
				>
				<el-button
					@click="clearPwd"
					type="primary"
					>取消文件加密</el-button
				>
			</div>
		</div> -->
		<div
			class="setting-item"
		>
			<label></label>
			<el-button
				@click="setFilePwd"
				type="primary"
				>{{ isSetPwd ? t("cancleFilePwd") : t("setFilePwd") }}</el-button
			>
		</div>
	</div>
</template>

<script lang="ts" setup>
	import { t } from "@/system";
	import {
		fetchGet,
		getApiUrl,
		getSystemConfig,
		setSystemKey,
	} from "@/system/config";
	import { onMounted, ref } from "vue";
	const config = getSystemConfig();
	const isSetPwd = ref(false);
  async function setFilePwd() {
    isSetPwd.value = !isSetPwd.value
    const params = {
      isPwd: isSetPwd.value ? 1 : 0
    }
		await fetchGet(`${getApiUrl()}/file/changeispwd?ispwd=${params.isPwd}`);
		setSystemKey("file", params);

  }
	onMounted(() => {
		isSetPwd.value = config.file.isPwd ? true : false;
	});
</script>

<style scoped>
	@import "./setStyle.css";
	.file-pwd-box {
		padding-top: 20px;
	}
	.setting-item {
		display: flex;
		align-items: center;
	}
</style>
