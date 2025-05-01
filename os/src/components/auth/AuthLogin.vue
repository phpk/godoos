<script setup lang="ts">
import { useLoginStore } from "@/stores/login";
import { onMounted, ref } from "vue";
import { List, MapLocation } from "@element-plus/icons-vue";
import LdapLogin from "./LdapLogin.vue";
const store = useLoginStore();

onMounted(() => {
	store.initThirdPartyLogin();
});
// 是否锁定
const isLock = ref(false);
</script>
<template>
	<div class="login-container">
		<SlideBackground />
		<div class="lockscreen">

			<el-card v-if="!isLock" class="login-box" shadow="never">
				<!-- <div class="absolute top-2 right-2 group flex items-center gap-2">
					<el-button type="primary" :icon="MapLocation" circle />
				</div> -->
				<!-- 顶部欢迎词和 logo -->
				<div class="header">
					<el-avatar size="large" class="logo">
						<img src="/logo.png" alt="Logo" />
					</el-avatar>
					<h3>欢迎使用 GodoOS</h3>
				</div>

				<!-- 登录、注册 -->
				<div class="login-register" v-if="
					store.thirdPartyLoginMethod !== 'dingding' &&
					store.thirdPartyLoginMethod !== 'qyweixin'
				">
					<!-- 滑块 -->
					<div class="slider">
						<span @click="store.thirdPartyLoginMethod = 'password'" :class="{
							active:
								store.thirdPartyLoginMethod !== 'register' && store.thirdPartyLoginMethod !== 'setup'
						}">登录</span>
						<span @click="store.thirdPartyLoginMethod = 'register'" :class="{
							active:
								store.thirdPartyLoginMethod === 'register',
						}">注册</span>
						<span @click="store.thirdPartyLoginMethod = 'setup'" :class="{
							active:
								store.thirdPartyLoginMethod === 'setup',
						}">设置</span>
					</div>
				</div>
				<!-- 登录功能 -->
				<PasswordLogin v-if="store.thirdPartyLoginMethod === 'password'" />
				<SetUp v-if="store.thirdPartyLoginMethod === 'setup'" />
				<PhoneLogin v-if="store.thirdPartyLoginMethod === 'phone'" />
				<EmailLogin v-if="store.thirdPartyLoginMethod === 'email'" />
				<LdapLogin v-if="store.thirdPartyLoginMethod === 'ldap'" />
				<div class="qr-code dingding" v-if="store.thirdPartyLoginMethod === 'dingding'">
					<div class="qr-code-container" id="dd-qr-code"></div>
				</div>
				<div class="qr-code-qyweixin" v-if="store.thirdPartyLoginMethod === 'qyweixin'">
					<div id="qywechat-qr-code"></div>
				</div>
				<UserRegister class="register-form" v-if="store.thirdPartyLoginMethod === 'register'" />

				<template v-else-if="store.thirdPartyLoginMethod !== 'register'">
					<div class="divider" v-if="store.thirdpartyList.length > 0">
						<span>第三方登录</span>
					</div>
					<div class="third-party-login">
						<el-button v-for="platform in store.thirdpartyList" :key="platform.name"
							class="third-party-login-button" @click="store.onThirdPartyLogin(platform.name)" :style="{
								backgroundColor:
									store.thirdPartyLoginMethod ===
										platform.name
										? '#f2ecec'
										: 'transparent',
							}">
							<img class="third-party-login-icon" :src="platform.icon" />
						</el-button>
					</div>
				</template>
			</el-card>
			<UnLock v-else />
		</div>
	</div>
</template>
<style lang="scss" scoped>
@use "@/styles/login.scss";
</style>
