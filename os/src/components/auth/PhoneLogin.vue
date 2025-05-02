<script setup lang="ts">
	import { getSmsCode } from "@/api/net/auth";
	import { useLoginStore } from "@/stores/login";
	import { ElMessage } from "element-plus";
	import { ref } from "vue";
	const store = useLoginStore();
	const phoneLoginRules = {
		phone: [
			{ required: true, message: "手机号不能为空", trigger: "blur" },
			{
				pattern: /^\d{11}$/,
				message: "手机号格式不正确，必须为11位有效数字",
				trigger: "blur",
			},
		],
		sms_code: [
			{ required: true, message: "验证码不能为空", trigger: "blur" },
			{
				pattern: /^\d{4,6}$/,
				message: "验证码必须为4到6位数字",
				trigger: "blur",
			},
		],
	};
	const phoneLoginFormRef = ref();
	const phoneForm = ref({
		phone: "",
		sms_code: "",
	});
	const isSendCodeButtonDisabled = ref(false);
	const sendCodeButtonText = ref("发送验证码");
	const phoneRegex = /^1[3-9]\d{9}$/;
	function getSms() {
		if (!phoneRegex.test(phoneForm.value.phone)) {
			ElMessage.error("请输入有效的手机号码");
			return;
		}
		getSmsCode(phoneForm.value.phone).then((res) => {
			if (res.success) {
				ElMessage.success("验证码发送成功");
				startCountdown();
			} else {
				ElMessage.error(res.message);
			}
		});
	}
	const startCountdown = () => {
		let countdown = 2;
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
	function toLogin() {
		phoneLoginFormRef.value.validate((valid: boolean) => {
			if (valid) {
				store.onLogin({
					loginType: "sms_code",
					params: phoneForm.value,
				});
			} else {
				ElMessage.error("请检查输入内容");
				return false;
			}
		});
	}
</script>
<template>
	<el-form
		:model="phoneForm"
		:rules="phoneLoginRules"
		ref="phoneLoginFormRef"
		label-position="left"
		label-width="0px"
	>
		<el-form-item prop="phone">
			<el-input
				style="height: 45px"
				v-model="phoneForm.phone"
				size="large"
				placeholder="请输入手机号"
				prefix-icon="Message"
			></el-input>
		</el-form-item>
		<el-form-item
			class="codeitem"
			prop="code"
		>
			<el-input
				style="height: 45px"
				v-model="phoneForm.sms_code"
				size="large"
				placeholder="请输入验证码"
				prefix-icon="Key"
			></el-input>

			<button
				class="send-code-btn"
				@click.stop="getSms"
				:disabled="isSendCodeButtonDisabled"
			>
				{{ sendCodeButtonText }}
			</button>
		</el-form-item>
		<el-form-item class="button-center">
			<el-button
				style="height: 45px"
				class="login-button"
				type="primary"
				size="large"
				@click="toLogin"
				>登录</el-button
			>
		</el-form-item>
	</el-form>
</template>
<style scoped lang="scss">
	.button-center {
		width: 100%;
		.login-button {
			width: 100%;
			border-radius: 50px;
		}
	}
	.code-input-container {
		display: flex;
		height: 45px;
		justify-content: space-between;
		margin-bottom: 16px;
	}

	.code-input {
		flex: 1;
		margin-right: 16px;
	}

	.send-button-container {
		flex: 1;
		margin-right: 0;
	}

	.login-button {
		width: 100%;
	}

	.send-code-btn {
		border: none;
		height: 45px;
		width: 100px;
		background: #409eff;
		color: #fff;
		border-radius: 50px;
		line-height: 36px;
		padding: 2px 6px;
		cursor: pointer;
		transition: background 0.3s ease;

		&:disabled {
			background: #c0c4cc; // 灰色背景
			cursor: not-allowed; // 禁用时的鼠标样式
		}
	}
	:deep(.el-input__wrapper) {
		border-radius: 50px;
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
