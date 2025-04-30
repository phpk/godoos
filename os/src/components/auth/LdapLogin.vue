<script setup lang="ts">
	import { useLoginStore } from "@/stores/login";
	import { ElMessage } from "element-plus";
	import { reactive, ref } from "vue";

	const store = useLoginStore();

	const loginForm = ref({
		username: "",
		password: "",
	});

	const rules = reactive({
		username: [
			{ required: true, message: "请输入用户名", trigger: "blur" },
			{
				min: 3,
				max: 20,
				message: "长度在 3 到 20 个字符",
				trigger: "blur",
			},
		],
		password: [
			{ required: true, message: "请输入密码", trigger: "blur" },
			{
				min: 6,
				max: 20,
				message: "长度在 6 到 20 个字符",
				trigger: "blur",
			},
		],
	});

	const formRef = ref();

	const submitForm = () => {
		formRef.value.validate((valid: boolean) => {
			if (valid) {
				store.onLogin({ loginType: "ldap", params: loginForm.value });
			} else {
				ElMessage.error("请检查输入内容");
				return false;
			}
		});
	};
</script>

<template>
	<el-form
		ref="formRef"
		label-position="left"
		:model="loginForm"
		:rules="rules"
		label-width="0px"
	>
		<el-form-item prop="username">
			<el-input
				v-model="loginForm.username"
				size="large"
				placeholder="请输入用户名"
				autofocus
				prefix-icon="UserFilled"
			></el-input>
		</el-form-item>
		<el-form-item prop="password">
			<el-input
				v-model="loginForm.password"
				size="large"
				type="password"
				placeholder="请输入登录密码"
				show-password
				prefix-icon="Key"
				@keyup.enter="submitForm"
			></el-input>
		</el-form-item>
		<el-form-item class="button-center">
			<el-button
				class="login-button"
				type="primary"
				@click="submitForm"
				size="large"
				>登录</el-button
			>
		</el-form-item>
	</el-form>
</template>

<style scoped>
	:deep(.el-input__wrapper) {
		border-radius: 50px;
		height: 45px;
	}
	.login-button {
		height: 45px;
		width: 100%;
		border-radius: 50px;
	}
</style>
