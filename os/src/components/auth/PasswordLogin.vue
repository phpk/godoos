<script setup lang="ts">
import { useLoginStore } from "@/stores/login";
import { ElMessage } from "element-plus";
import { reactive, ref } from "vue";

const store = useLoginStore();
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
			min: 3,
			max: 20,
			message: "长度在 3 到 20 个字符",
			trigger: "blur",
		},
	],
});

const formRef = ref();

const submitForm = () => {
	formRef.value.validate((valid: boolean) => {
		if (valid) {
			showCaptcha.value = true;
			// store.onLogin({
			// 	loginType: "password",
			// 	params: store.loginForm,
			// });
		} else {
			ElMessage.error("请检查输入内容");
			return false;
		}
	});
};
const showCaptcha = ref(false); // 控制是否显示验证码
const captchaPassed = ref(false); // 验证码是否通过
const handleCaptchaConfirm = (success: boolean) => {
	if (success) {
		captchaPassed.value = true;
		store.onLogin({
			loginType: "password",
			params: store.loginForm,
		});
		showCaptcha.value = false; // 隐藏验证码
	} else {
		ElMessage.error("验证码验证失败，请重试");
		showCaptcha.value = false;
	}
};
const closeCaptcha = () => {
	showCaptcha.value = false;
};
</script>

<template>
	<el-form ref="formRef" label-position="left" :model="store.loginForm" :rules="rules" label-width="0px">
		<el-form-item prop="username">
			<el-input style="border-radius: 50px; height: 45px" v-model="store.loginForm.username" size="large"
				placeholder="请输入用户名" autofocus prefix-icon="UserFilled"></el-input>
		</el-form-item>
		<el-form-item prop="password">
			<el-input style="border-radius: 50px; height: 45px" v-model="store.loginForm.password" size="large"
				type="password" placeholder="请输入登录密码" show-password prefix-icon="Key"
				@keyup.enter="submitForm"></el-input>
		</el-form-item>
		<el-form-item class="remember-me">
			<el-checkbox v-model:checked="store.loginForm.rememberMe" size="large">记住密码</el-checkbox>
		</el-form-item>
		<el-form-item class="button-center">
			<el-button style="height: 45px" class="login-button" type="primary" size="large"
				@click="submitForm">登录</el-button>
		</el-form-item>
	</el-form>
	<teleport to="body">
		<div v-if="showCaptcha" class="custom-captcha-dialog">
			<div class="overlay" @click.self="closeCaptcha"></div>
			<div class="dialog-content">
				<SlideCaptcha ref="captcha" @onSuccess="handleCaptchaConfirm" />
			</div>
		</div>
	</teleport>
</template>

<style scoped>
.button-center {
	width: 100%;

	.login-button {
		width: 100%;
		border-radius: 50px;
	}
}

.remember-me {
	display: flex;
	justify-content: flex-start;
	margin-left: 5px;
}

:deep(.el-input__wrapper) {
	border-radius: 50px;
}

.custom-captcha-dialog {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	z-index: 99;
	display: flex;
	align-items: center;
	justify-content: center;
}

.overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100vw;
	height: 100vh;
	background-color: rgba(0, 0, 0, 0.5);
	z-index: 9998;
}

.dialog-content {
	padding: 20px;
	z-index: 9998;
	width: 350px;
}
</style>
