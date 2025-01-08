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
				v-if="ThirdPartyLoginMethod === 'login'"
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
						@keyup.enter="onLogin()"
					></el-input>
				</el-form-item>
				<el-button
					type="primary"
					@click="onLogin()"
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
				<div class="actions">
					<!-- <a
						href="#"
						@click.prevent="toggleRegister"
						>注册新用户</a
					> -->
					<a
						href="#"
						@click.prevent="toggleUserSwitch"
						>切换角色</a
					>
				</div>
			</el-form>
			<el-form
				v-else-if="ThirdPartyLoginMethod === 'register'"
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
			</el-form>

			<div v-else-if="ThirdPartyLoginMethod === 'dingding'">
				<el-row>
					<el-col :span="24">
						<div class="qr-code">
							<div id="dd-qr-code"></div>
						</div>
					</el-col>
				</el-row>
				<el-row>
					<el-col :span="24">
						<el-button @click="ThirdPartyLoginMethod = 'login'"
							>返回</el-button
						>
					</el-col>
				</el-row>
			</div>
			<!-- 企业微信 -->
			<div v-else-if="ThirdPartyLoginMethod === 'qywechat'">
				<el-row>
					<el-col :span="24">
						<div class="qr-code">
							<div id="qywechat-qr-code"></div>
						</div>
					</el-col>
				</el-row>
				<el-row>
					<el-col :span="24">
						<el-button @click="ThirdPartyLoginMethod = 'login'"
							>返回</el-button
						>
					</el-col>
				</el-row>
			</div>
			<!-- 手机登录 -->
			<div v-else-if="ThirdPartyLoginMethod === 'phone'">
				<el-row>
					<el-col :span="24">
						<el-form>
							<el-form-item>
								<el-input
									v-model="phoneForm.phone"
									placeholder="请输入手机号"
								/>
							</el-form-item>
							<!-- 验证码 -->
							<el-form-item>
								<el-row :gutter="10">
									<el-col :span="16">
										<el-input
											v-model="phoneForm.code"
											placeholder="请输入验证码"
										/>
									</el-col>
									<el-col :span="8">
										<button
											class="send-code-btn"
											@click.prevent="onSendCode"
										>
											发送验证码
										</button>
									</el-col>
								</el-row>
							</el-form-item>
						</el-form>
					</el-col>
				</el-row>
				<el-row>
					<el-col :span="24">
						<el-button @click="ThirdPartyLoginMethod = 'login'"
							>返回</el-button
						>
					</el-col>
				</el-row>
			</div>
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
	import QRCode from "qrcode";
	import { onMounted, ref, watchEffect } from "vue";
	import { useRoute } from "vue-router";
	const route = useRoute();
	const store = useLoginStore();
	const sys = useSystem();
	const loginCallback = sys._options.loginCallback;
	const config = getSystemConfig();
	const lockClassName = ref("screen-show");
	const ThirdPartyLoginMethod = ref("login");
	const phoneForm = ref({
		phone: "",
		code: "",
	});
	// 声明 DTFrameLogin 方法
	declare global {
		interface Window {
			DTFrameLogin: (
				frameParams: IDTLoginFrameParams,
				loginParams: IDTLoginLoginParams,
				successCbk: (result: IDTLoginSuccess) => void,
				errorCbk?: (error: any) => void
			) => void;
		}
	}
	// 包裹容器的尺寸与样式需要接入方自己使用css设置
	interface IDTLoginFrameParams {
		id: string; // 必传，包裹容器元素ID，不带'#'
		width?: number; // 选传，二维码iframe元素宽度，最小280，默认300
		height?: number; // 选传，二维码iframe元素高度，最小280，默认300
	}
	// 增加了isPre参数来设定运行环境
	interface IDTLoginLoginParams {
		redirect_uri: string; // 必传，注意url需要encode
		response_type: string; // 必传，值固定为code
		client_id: string; // 必传
		scope: string; // 必传，如果值为openid+corpid，则下面的org_type和corpId参数必传，否则无法成功登录
		prompt: string; // 必传，值为consent。
		state?: string; // 选传
		org_type?: string; // 选传，当scope值为openid+corpid时必传
		corpId?: string; // 选传，当scope值为openid+corpid时必传
		exclusiveLogin?: string; // 选传，如需生成专属组织专用二维码时，可指定为true，可以限制非组织帐号的扫码
		exclusiveCorpId?: string; // 选传，当exclusiveLogin为true时必传，指定专属组织的corpId
	}
	interface IDTLoginSuccess {
		redirectUrl: string; // 登录成功后的重定向地址，接入方可以直接使用该地址进行重定向
		authCode: string; // 登录成功后获取到的authCode，接入方可直接进行认证，无需跳转页面
		state?: string; // 登录成功后获取到的state
	}

	onMounted(() => {
		getThirdpartyList();
	});

	const getThirdpartyList = async () => {
		const result = await fetch(
			config.userInfo.url + "/user/thirdparty/list"
		);
		if (result.ok) {
			await result.json();
		}
	};

	const onDingDingScan = () => {
		ThirdPartyLoginMethod.value = "dingding";
		initDingDingScan();
	};

	const initDingDingScan = async () => {
		try {
			// 加载钉钉登录脚本
			await loadScript(
				"https://g.alicdn.com/dingding/h5-dingtalk-login/0.21.0/ddlogin.js"
			);

			const res = await fetch("http://192.168.1.10:8816/user/ding/conf");
			const data = await res.json();

			// 在这里可以调用DTFrameLogin或其他依赖于该脚本的方法
			window.DTFrameLogin(
				{
					id: "dd-qr-code",
					width: 200,
					height: 300,
				},
				{
					redirect_uri: encodeURIComponent(data.data.host),
					client_id: data.data.client_id,
					scope: "openid",
					response_type: "code",
					state: "xxxxxxxxx",
					prompt: "consent",
				},
				(loginResult: any) => {
					const { authCode } = loginResult;
					onLogin(authCode);
				},
				(errorMsg: any) => {
					console.log("二维码获取错误", errorMsg);
				}
			);
		} catch (error) {
			console.error(error);
		}
	};

	// 函数用于动态加载外部JS文件
	function loadScript(url: string): Promise<void> {
		return new Promise((resolve, reject) => {
			const script = document.createElement("script");
			script.src = url;
			script.onload = () => resolve();
			script.onerror = () =>
				reject(new Error(`Failed to load script ${url}`));
			document.head.appendChild(script);
		});
	}

	const onSendCode = () => {
		console.log("发送验证码");
	};

	const thirdPartyPlatforms = [
		{
			name: "qywechat",
			icon: new URL("@/assets/login/qywechat.png", import.meta.url).href,
		},
		{
			name: "dingding",
			icon: new URL("@/assets/login/dingding.png", import.meta.url).href,
		},
		{
			name: "phone",
			icon: new URL("@/assets/login/phone.png", import.meta.url).href,
		},
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
			console.log("code:", code);
			onLogin();
		}
	});

	async function onLogin(authCode?: string) {
		if (loginCallback) {
			const platform = localStorage.getItem("ThirdPartyPlatform");
			const code =
				authCode || (router.currentRoute.value.query.code as string);

			// 使用映射对象，将平台名称映射到相应的参数名称
			const platformCodeMap: Record<string, string> = {
				github: "github",
				gitee: "gitee",
				qywechat: "qywechat",
				qq: "qq",
				dingding: "dingtalk_scan",
			};

			// 获取对应的参数名称
			const codeParam = platform ? platformCodeMap[platform] : null;

			// 根据登录类型传递参数
			const login_type = codeParam ? codeParam : "password";

			const param = codeParam
				? { code: code }
				: { username: userName.value, password: userPassword.value };

			const res = await loginCallback(login_type, param);

			if (res) {
				localStorage.removeItem("ThirdPartyPlatform");
				loginSuccess();
			} else {
				notifyError("登录失败，请重试");
			}
		}
	}

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
				// toggleRegister();
			} else {
				notifyError(res.message);
				return;
			}
		} catch (error) {
			console.error(error);
		}

		// 		console.log("123");
		// const res = await fetch(config.userInfo.url + "/user/editdata", {
		// 	method: "POST",
		// 	body: JSON.stringify({
		// 		phone: userName.value,
		// 	}),
		// 	headers: {
		// 		ClientID: store.tempClientId,
		// 		Authorization: store.tempToken,
		// 	},
		// });
		// if (res.ok) {
		// 	notifySuccess("手机号设置成功");
		// } else {
		// 	notifyError("手机号设置失败");
		// }
	};

	const onThirdPartyLogin = async (platform: string) => {
		localStorage.setItem("ThirdPartyPlatform", platform);
		let loginFunction: (() => Promise<boolean>) | undefined;

		switch (platform) {
			case "github":
				loginFunction = async function () {
					return await authWithGithub();
				};
				break;
			case "gitee":
				loginFunction = async function () {
					return await authWithGitee();
				};
				break;
			case "dingding":
				loginFunction = async function () {
					return await authWithDingDing();
				};
				break;
			case "qywechat":
				loginFunction = async function () {
					return await authWithWechat();
				};
				break;
			case "phone":
				loginFunction = async function () {
					return await authWithPhone();
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

	const authWithPhone = async (): Promise<boolean> => {
		ThirdPartyLoginMethod.value = "phone";
		return true;
	};

	const authWithDingDing = async (): Promise<boolean> => {
		onDingDingScan();
		return true;
	};

	const authWithGithub = async (): Promise<boolean> => {
		// 传递state用于防止CSRF攻击,使用时间戳加随机字符串
		const state = Date.now() + Math.random().toString(36).substring(2, 15);
		store.State = state;
		// 获取当前页面url当做回调参数
		const currentUrl = window.location.href;
		const url =
			config.userInfo.url + "/user/github/authorize?state=" + state;
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

	const authWithWechat = async (): Promise<boolean> => {
		ThirdPartyLoginMethod.value = "qywechat";
		const res = await fetch(
			"http://server001.godoos.com/user/qyweixin/qrcode"
		);
		if (res.ok) {
			const data = await res.json();
			if (data.success && data.data.url) {
				// 生成二维码
				QRCode.toCanvas(
					document.getElementById("qywechat-qr-code"),
					data.data.url,
					function (error: any) {
						if (error) console.error(error);
						console.log("二维码生成成功");
					}
				);
			} else {
				notifyError("获取授权URL失败");
			}
		} else {
			notifyError("网络错误，无法获取二维码");
		}
		return true;
	};

	const authWithGitee = async (): Promise<boolean> => {
		// 传递state用于防止CSRF攻击,使用时间戳加随机字符串
		const state = Date.now() + Math.random().toString(36).substring(2, 15);
		store.State = state;
		// 获取当前页面url当做回调参数
		const currentUrl = window.location.href;
		const url =
			config.userInfo.url + "/user/gitee/authorize?state=" + state;
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
		// store.page = "phone";
		// return true;
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

	.qr-code {
		width: 200px;
		height: 200px;
		padding: 10px;
		margin: 0 auto;
		overflow: hidden;
		border: 1px solid rgb(203, 203, 203);
		border-radius: 10px;
		:deep(#dd-qr-code) {
			width: 200px;
			height: 200px;
			iframe {
				margin-top: -50px;
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

	.send-code-btn {
		border: none;
		height: 40px;
		background: #409eff;
		color: #fff;
		border-radius: 4px;
		line-height: 40px;
		cursor: pointer;
	}
</style>
