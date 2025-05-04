<script setup lang="ts">
import { useLoginStore } from "@/stores/login";
import { ref, Ref } from "vue";
const store = useLoginStore();
interface RegisterInfo {
	username: string;
	nickname: string;
	password: string;
	email: string;
	phone: string;
	// third_user_id: string;
	// union_id: string;
	// platform: string;
	confirmPassword: string;
}
// 注册信息
const registerInfo: Ref<RegisterInfo> = ref({
	username: "",
	nickname: "",
	password: "",
	email: "",
	phone: "",
	// third_user_id: "",
	// union_id: "",
	// platform: "",
	confirmPassword: "",
});
const rules = {
	username: [
		{ required: true, message: "用户名不能为空", trigger: "blur" },
		{
			min: 3,
			max: 20,
			message: "用户名长度应在3到20个字符之间",
			trigger: "blur",
		},
		{
			pattern: /^[a-zA-Z0-9]+$/,
			message: "用户名只能包含英文和数字",
			trigger: "blur",
		},
	],
	password: [
		{ required: true, message: "密码不能为空", trigger: "blur" },
		{ min: 3, message: "密码长度不能小于3位", trigger: "blur" },
	],
	confirmPassword: [
		{ required: true, message: "请再次输入密码", trigger: "blur" },
		{
			validator: (rule: any, value: any, callback: any) => {
				console.log(rule);
				if (value === "") {
					callback(new Error("请再次输入密码"));
				} else if (value !== registerInfo.value.password) {
					callback(new Error("两次输入的密码不一致"));
				} else {
					callback();
				}
			},
			trigger: "blur",
		},
	],
	// email: [
	// 	{ required: true, message: "邮箱不能为空", trigger: "blur" },
	// 	{
	// 		type: "email",
	// 		message: "请输入有效的邮箱地址",
	// 		trigger: ["blur", "change"],
	// 	},
	// ],
	// phone: [
	// 	{ required: true, message: "手机号不能为空", trigger: "blur" },
	// 	{
	// 		pattern: /^1[3-9]\d{9}$/,
	// 		message: "请输入有效的手机号",
	// 		trigger: "blur",
	// 	},
	// ],
	// nickname: [
	// 	{ required: true, message: "昵称不能为空", trigger: "blur" },
	// 	{
	// 		min: 2,
	// 		max: 20,
	// 		message: "昵称长度应在2到20个字符之间",
	// 		trigger: "blur",
	// 	},
	// ],
};
</script>
<template>
	<el-form label-position="left" label-width="0px" :model="registerInfo" ref="regFormRef" :rules="rules">
		<el-form-item prop="username">
			<el-input v-model="registerInfo.username" size="large" placeholder="请输入用户名"
				prefix-icon="UserFilled"></el-input>
		</el-form-item>
		<!-- <el-form-item prop="nickname">
			<el-input
				v-model="registerInfo.nickname"
				size="large"
				placeholder="请输入真实姓名"
				prefix-icon="Avatar"
			></el-input>
		</el-form-item>
		<el-form-item prop="email">
			<el-input
				v-model="registerInfo.email"
				size="large"
				placeholder="请输入邮箱"
				prefix-icon="Message"
			></el-input>
		</el-form-item>
		<el-form-item prop="phone">
			<el-input
				v-model="registerInfo.phone"
				size="large"
				placeholder="请输入手机号"
				prefix-icon="Iphone"
			></el-input>
		</el-form-item> -->
		<el-form-item prop="password">
			<el-input v-model="registerInfo.password" size="large" type="password" placeholder="请输入密码" show-password
				prefix-icon="Key"></el-input>
		</el-form-item>
		<el-form-item prop="confirmPassword">
			<el-input v-model="registerInfo.confirmPassword" size="large" type="password" placeholder="请再次输入密码"
				show-password prefix-icon="Lock"></el-input>
		</el-form-item>

		<el-form-item class="button-center">
			<el-button class="login-button" type="primary" size="large" @click="store.onRegister(registerInfo)">注册</el-button>
		</el-form-item>
	</el-form>
</template>
<style scoped>
.button-center {
	width: 100%;

	.login-button {
		width: 100%;
		height: 45px;
		border-radius: 50px;
	}
}

:deep(.el-input__wrapper) {
	border-radius: 50px;
	height: 45px;
}


</style>
