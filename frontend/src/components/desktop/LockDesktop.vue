<template>
	<div
		class="lockscreen"
		:class="lockClassName"
	>
		<el-card
			class="login-box"
			shadow="never"
		>
			<div class="avatar-container">
				<el-avatar size="large">
					<img
						src="/logo.png"
						alt="Logo"
					/>
				</el-avatar>
			</div>
			<el-form
				v-if="!isRegisterMode"
				label-position="left"
				label-width="0px"
			>
				<el-form-item>
					<el-input
						v-model="userName"
						placeholder="请输入用户名"
						autofocus
						prefix-icon="UserFilled"
					></el-input>
				</el-form-item>
				<el-form-item v-if="!sys._options.noPassword">
					<el-input
						v-model="userPassword"
						type="password"
						placeholder="请输入登录密码"
						show-password
						prefix-icon="Key"
						@keyup.enter="onLogin"
					></el-input>
				</el-form-item>
				<el-button
					type="primary"
					@click="onLogin"
					>登录</el-button
				>
				<div class="divider">
					<span>第三方登录</span>
				</div>
				<div class="third-party-login">
					<img
						v-for="platform in thirdPartyPlatforms"
						:key="platform.name"
						:src="platform.icon"
						:alt="platform.name"
						style="
							width: 25px;
							height: 25px;
							cursor: pointer;
							color: #409eff;
						"
						@click="onThirdPartyLogin(platform.name)"
					/>
				</div>
				<div
					class="actions"
					v-if="config.userType === 'member'"
				>
					<a
						href="#"
						@click.prevent="toggleRegister"
						>注册新用户</a
					>
					<a
						href="#"
						v-if="config.userType !== 'member'"
						@click.prevent="toggleUserSwitch"
						>切换角色</a
					>
				</div>
			</el-form>
			<el-form
				v-else
				label-position="left"
				label-width="0px"
				:model="regForm"
				ref="regFormRef"
				:rules="rules"
			>
				<el-form-item prop="username">
					<el-input
						v-model="regForm.username"
						placeholder="请输入用户名"
						prefix-icon="UserFilled"
					></el-input>
				</el-form-item>
				<el-form-item prop="nickname">
					<el-input
						v-model="regForm.nickname"
						placeholder="请输入真实姓名"
						prefix-icon="Avatar"
					></el-input>
				</el-form-item>
				<el-form-item prop="email">
					<el-input
						v-model="regForm.email"
						placeholder="请输入邮箱"
						prefix-icon="Message"
					></el-input>
				</el-form-item>
				<el-form-item prop="phone">
					<el-input
						v-model="regForm.phone"
						placeholder="请输入手机号"
						prefix-icon="Iphone"
					></el-input>
				</el-form-item>
				<el-form-item prop="password">
					<el-input
						v-model="regForm.password"
						type="password"
						placeholder="请输入密码"
						show-password
						prefix-icon="Key"
					></el-input>
				</el-form-item>
				<el-form-item prop="confirmPassword">
					<el-input
						v-model="regForm.confirmPassword"
						type="password"
						placeholder="请再次输入密码"
						show-password
						prefix-icon="Lock"
					></el-input>
				</el-form-item>
				<el-button
					type="primary"
					@click="onRegister"
					>注册</el-button
				>
				<div class="actions">
					<a
						href="#"
						@click.prevent="toggleRegister"
						>返回登录</a
					>
				</div>
			</el-form>
		</el-card>
	</div>
</template>

<script lang="ts" setup>
	import { useLoginStore } from "@/stores/login";
	import { useSystem } from "@/system";
	import { getSystemConfig, setSystemConfig } from "@/system/config";
	import router from "@/system/router";
	import { RestartApp } from "@/util/goutil";
	import { notifyError } from "@/util/msg";
	import { onMounted, ref, watchEffect } from "vue";
	import { useRoute } from "vue-router";
	const route = useRoute();
	const store = useLoginStore();
	const sys = useSystem();
	const loginCallback = sys._options.loginCallback;
	const config = getSystemConfig();
	const lockClassName = ref("screen-show");
	const isRegisterMode = ref(false);

	const thirdPartyPlatforms = [
		// {
		// 	name: "wechat",
		// 	icon: new URL("@/assets/login/wechat.png", import.meta.url).href,
		// },
		// {
		// 	name: "qq",
		// 	icon: new URL("@/assets/login/qq.png", import.meta.url).href,
		// },
		// {
		// 	name: "sina",
		// 	icon: new URL("@/assets/login/sina.png", import.meta.url).href,
		// },
		{
			name: "github",
			icon: new URL("@/assets/login/github.png", import.meta.url).href,
		},
		{
			name: "gitee",
			icon: new URL("@/assets/login/gitee.png", import.meta.url).href,
		},
	];

	function loginSuccess() {
		lockClassName.value = "screen-hidean";
		setTimeout(() => {
			lockClassName.value = "screen-hide";
		}, 500);
	}

	const userName = ref("");
	const userPassword = ref("");

	onMounted(() => {
		if (config.userType === "person") {
			userName.value = sys._options.login?.username || "admin";
			userPassword.value = sys._options.login?.password || "";
		} else {
			userName.value = config.userInfo.username;
			userPassword.value = config.userInfo.password;
		}
	});

	// 使用 watchEffect 监听 code 参数
	watchEffect(() => {
		const code = route.query.code;
		if (code) {
			console.log(code, "---");
			onLogin();
		}
	});

	const onLogin = async () => {
		localStorage.removeItem("godoosClientId");
		if (loginCallback) {
			const platform = store.ThirdPartyPlatform;
			const code = router.currentRoute.value.query.code as string;

			// 使用映射对象，将平台名称映射到相应的参数名称
			const platformCodeMap: Record<string, string> = {
				github: "github_code",
				gitee: "gitee_code",
				wechat: "wechat_code",
				qq: "qq_code",
			};

			// 获取对应的参数名称
			const codeParam = platform ? platformCodeMap[platform] : null;

			let res;
			if (codeParam) {
				// 第三方登录统一调用
				const returnedState = router.currentRoute.value.query
					.state as string;

				if (returnedState !== store.State) {
					notifyError("登录失败，请重试");
					return;
				}

				res = await loginCallback(userName.value, userPassword.value, {
					[codeParam]: code,
				});
			} else {
				// 普通登录
				res = await loginCallback(userName.value, userPassword.value);
			}

			if (res) {
				store.ThirdPartyPlatform = null;
				loginSuccess();
			} else {
				notifyError("登录失败，请重试");
			}
		}
	};

	const toggleRegister = () => {
		isRegisterMode.value = !isRegisterMode.value;
	};

	const toggleUserSwitch = () => {
		config.userType = "person";
		setSystemConfig(config);
		RestartApp();
	};

	const regForm = ref({
		username: "",
		password: "",
		confirmPassword: "",
		email: "",
		phone: "",
		nickname: "",
	});
	const regFormRef: any = ref(null);

	const rules = {
		username: [
			{ required: true, message: "用户名不能为空", trigger: "blur" },
			{
				min: 3,
				max: 20,
				message: "用户名长度应在3到20个字符之间",
				trigger: "blur",
			},
		],
		password: [
			{ required: true, message: "密码不能为空", trigger: "blur" },
			{ min: 6, message: "密码长度不能小于6位", trigger: "blur" },
		],
		confirmPassword: [
			{ required: true, message: "请再次输入密码", trigger: "blur" },
			{
				validator: (rule: any, value: any, callback: any) => {
					console.log(rule);
					if (value === "") {
						callback(new Error("请再次输入密码"));
					} else if (value !== regForm.value.password) {
						callback(new Error("两次输入的密码不一致"));
					} else {
						callback();
					}
				},
				trigger: "blur",
			},
		],
		email: [
			{ required: true, message: "邮箱不能为空", trigger: "blur" },
			{
				type: "email",
				message: "请输入有效的邮箱地址",
				trigger: ["blur", "change"],
			},
		],
		phone: [
			{ required: true, message: "手机号不能为空", trigger: "blur" },
			{
				pattern: /^1[3-9]\d{9}$/,
				message: "请输入有效的手机号",
				trigger: "blur",
			},
		],
		nickname: [
			{ required: true, message: "昵称不能为空", trigger: "blur" },
			{
				min: 2,
				max: 20,
				message: "昵称长度应在2到20个字符之间",
				trigger: "blur",
			},
		],
	};

	const onRegister = async () => {
		try {
			await regFormRef.value.validate();
			const save = toRaw(regForm.value);
			const userInfo = config.userInfo;
			const comp = await fetch(userInfo.url + "/member/register", {
				method: "POST",
				body: JSON.stringify(save),
			});
			if (!comp.ok) {
				notifyError("网络错误，注册失败");
				return;
			}
			const res = await comp.json();
			if (res.success) {
				notifyError("注册成功");
				toggleRegister();
			} else {
				notifyError(res.message);
				return;
			}
		} catch (error) {
			console.error(error);
		}
	};

	const onThirdPartyLogin = async (platform: string) => {
		let loginFunction: (() => Promise<boolean>) | undefined;
		// 当前选择的第三方登录方式
		store.ThirdPartyPlatform = platform;
		switch (platform) {
			case "github":
				// 跳转到 GitHub 授权页面
				loginFunction = async function () {
					return await authWithGithub();
				};
				break;
			case "wechat":
				loginFunction = async function () {
					return await authWithWechat();
				};
				break;
			case "qq":
				loginFunction = async function () {
					return await authWithQQ();
				};
				break;
			case "sina":
				loginFunction = async function () {
					return await authWithSina();
				};
				break;
			case "gitee":
				loginFunction = async function () {
					return await authWithGitee();
				};
				break;
			default:
				notifyError("不支持的第三方登录平台");
				return;
		}

		if (loginFunction) {
			const success = await loginFunction();
			if (success) {
				loginSuccess();
			}
		}
	};

	const authWithGithub = async (): Promise<boolean> => {
		// 传递state用于防止CSRF攻击,使用时间戳加随机字符串
		const state = Date.now() + Math.random().toString(36).substring(2, 15);
		store.State = state;
		// 获取当前页面url当做回调参数
		const currentUrl = window.location.href;
		const url = config.userInfo.url + "/github/authorize?state=" + state;
		const res: any = await fetch(url, {
			method: "POST",
			body: JSON.stringify({
				state: state,
				redirect_url: currentUrl,
			}),
		});
		if (!res.ok) {
			return false;
		}
		const data = await res.json();
		if (data && data.data && data.data.url) {
			console.log(data.data.url, "---");
			// 使用正则表达式检查URL格式
			const urlPattern = /client_id=[^&]+/;
			if (urlPattern.test(data.data.url)) {
				window.location.href = data.data.url;
				return true;
			} else {
				notifyError("请先在系统配置中设置github登陆配置");
				return false;
			}
		} else {
			notifyError("获取授权URL失败");
			return false;
		}
	};

	const authWithWechat = (): boolean | PromiseLike<boolean> => {
		throw new Error("Function not implemented.");
	};

	const authWithQQ = (): boolean | PromiseLike<boolean> => {
		throw new Error("Function not implemented.");
	};

	const authWithSina = (): boolean | PromiseLike<boolean> => {
		throw new Error("Function not implemented.");
	};

	const authWithGitee = async (): Promise<boolean> => {
		// 传递state用于防止CSRF攻击,使用时间戳加随机字符串
		const state = Date.now() + Math.random().toString(36).substring(2, 15);
		store.State = state;
		// 获取当前页面url当做回调参数
		const currentUrl = window.location.href;
		const url = config.userInfo.url + "/gitee/authorize?state=" + state;
		const res: any = await fetch(url, {
			method: "POST",
			body: JSON.stringify({
				state: state,
				redirect_url: currentUrl,
			}),
		});
		if (!res.ok) {
			return false;
		}
		const data = await res.json();
		if (data && data.data && data.data.url) {
			// 使用正则表达式检查URL格式
			const urlPattern = /client_id=[^&]+/;
			if (urlPattern.test(data.data.url)) {
				window.location.href = data.data.url;
				return true;
			} else {
				notifyError("请先在系统配置中设置gitee登陆配置");
				return false;
			}
		} else {
			notifyError("获取授权URL失败");
			return false;
		}
	};
</script>

<style scoped lang="scss">
	.lockscreen {
		position: absolute;
		top: 0;
		right: 0;
		bottom: 0;
		left: 0;
		z-index: 201;
		display: flex;
		justify-content: center;
		align-items: center;
		overflow: hidden;
		color: #fff;
		background-color: rgba(25, 28, 34, 0.78);
		backdrop-filter: blur(7px);

		.login-box {
			width: 300px;
			padding: 20px;
			text-align: center;
			background: #ffffff;
			border-radius: 10px;
			box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
			border: 1px solid #e0e0e0;

			.avatar-container {
				margin-bottom: 20px;
			}

			.el-input {
				width: 100%;
				margin-bottom: 10px;
				background: #f9f9f9;
				border-radius: 4px;
			}

			.el-button {
				width: 100%;
				margin-top: 10px;
				background: #409eff;
				color: #ffffff;
				border: none;
				border-radius: 4px;
				transition: background 0.3s ease;
			}

			.el-button:hover {
				background: #66b1ff;
			}

			.tip {
				padding: 4px 0;
				font-size: 12px;
				color: red;
				height: 30px;
			}

			.actions {
				margin-top: 10px;
				display: flex;
				justify-content: space-between;

				a {
					color: #409eff;
					text-decoration: none;
					cursor: pointer;

					&:hover {
						text-decoration: underline;
					}
				}
			}
		}

		.screen-hidean {
			animation: outan 0.5s forwards;
		}

		.screen-hide {
			display: none;
		}

		.third-party-login {
			margin-top: 10px;
			display: flex;
			justify-content: center;
			gap: 15px;

			el-icon {
				font-size: 24px;
				cursor: pointer;
				transition: color 0.3s ease;

				&:hover {
					color: #409eff;
				}
			}
		}

		.divider {
			display: flex;
			align-items: center;
			text-align: center;
			margin: 20px 0;
			color: #999;
			font-size: 14px;

			&::before,
			&::after {
				content: "";
				flex: 1;
				border-bottom: 1px solid #ddd;
			}

			&::before {
				margin-right: 0.25em;
			}

			&::after {
				margin-left: 0.25em;
			}
		}
	}

	@keyframes outan {
		0% {
			opacity: 1;
		}

		30% {
			opacity: 0;
		}

		100% {
			transform: translateY(-100%);
			opacity: 0;
		}
	}
</style>
