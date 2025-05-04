<script setup lang="ts">
import { getEmailCode } from "@/api/net/auth"; // 假设你有一个发送邮箱验证码的API
import { useLoginStore } from "@/stores/login";
import { ElMessage } from "element-plus";
import { ref } from "vue";

const store = useLoginStore();

const emailLoginRules = {
	email: [
		{
			required: true,
			message: "邮箱不能为空",
			trigger: ["blur"],
		},
		{
			type: "email" as const, // 显式声明为 RuleType
			message: "请输入有效的邮箱地址",
			trigger: ["blur", "change"],
		},
	],
	code: [
		{
			required: true,
			message: "验证码不能为空",
			trigger: ["blur"],
		},
		{
			pattern: /^\d{4,6}$/,
			message: "验证码必须为4到6位数字",
			trigger: ["blur"],
		},
	],
};

const emailLoginFormRef = ref();
const emailForm = ref({
	email: "",
	code: "",
});

const isSendCodeButtonDisabled = ref(false);
const sendCodeButtonText = ref("发送验证码");

function getSms() {
	if (!emailForm.value.email) {
		ElMessage.error("请输入有效的邮箱地址");
		return;
	}
	getEmailCode(emailForm.value.email).then((res) => {
		if (res.success) {
			ElMessage.success("验证码发送成功");
			startCountdown();
		} else {
			ElMessage.error(res.message);
		}
	});
}

const startCountdown = () => {
	let countdown = 60; // 倒计时60秒
	isSendCodeButtonDisabled.value = true;
	sendCodeButtonText.value = `${countdown}秒后重试`;

	const interval = setInterval(() => {
		countdown--;
		if (countdown > 0) {
			sendCodeButtonText.value = `${countdown}秒后重试`;
		} else {
			clearInterval(interval);
			sendCodeButtonText.value = "发送验证码";
			isSendCodeButtonDisabled.value = false;
		}
	}, 1000);
};

function toEmailLogin() {
	emailLoginFormRef.value.validate((valid: boolean) => {
		if (valid) {
			store.onLogin({ loginType: "email", params: emailForm.value });
		} else {
			ElMessage.error("请检查输入内容");
			return false;
		}
	});
}
</script>

<template>
	<el-form :model="emailForm" :rules="emailLoginRules" ref="emailLoginFormRef" label-position="left"
		label-width="0px">
		<el-form-item prop="email">
			<el-input v-model="emailForm.email" size="large" placeholder="请输入邮箱" prefix-icon="Message"></el-input>
		</el-form-item>
		<el-form-item class="codeitem" prop="code">
			<el-input v-model="emailForm.code" size="large" placeholder="请输入验证码" prefix-icon="Key">
			</el-input>

			<button class="send-code-btn" @click.stop="getSms" :disabled="isSendCodeButtonDisabled">
				{{ sendCodeButtonText }}
			</button>

		</el-form-item>
		<el-form-item class="button-center">
			<el-button class="login-button" type="primary" size="large" @click="toEmailLogin">登录</el-button>
		</el-form-item>
	</el-form>
</template>

<style scoped>
.button-center {
	display: flex;
	justify-content: center;
	text-align: center;
	width: 100%;

	.login-button {
		width: 100%;
		border-radius: 50px;
		height: 45px;
	}
}

.send-code-btn {
	border: none;
	height: 45px;
	background: #409eff;
	color: #fff;
	width: 100px;
	border-radius: 50px;
	line-height: 36px;
	padding: 2px 6px;
	cursor: pointer;
	transition: background 0.3s ease;

	&:disabled {
		background: #c0c4cc;
		/* 灰色背景 */
		cursor: not-allowed;
		/* 禁用时的鼠标样式 */
	}
}

:deep(.el-input__wrapper) {
	border-radius: 50px;
	height: 45px;
}

.codeitem {
	display: flex;
	align-items: center;
}

.codeitem .el-input {
	flex: 1;
	margin-right: 10px;
}
</style>
