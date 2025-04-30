<template>
	<transition name="slide">
	<div class="lockscreen" v-if="settingsStore.isLockScreen">
		<SlideBackground />
		<div class="lock-card">
			<div class="avatar-container">
				<el-avatar size="large">
					<img src="/logo.png" alt="Logo" />
				</el-avatar>
			</div>
			<el-form>
				<el-form-item>
					<el-input ref="passwordInput" :prefix-icon="Lock" v-model="password" type="password"
						@keyup.enter="unLockScreen" placeholder="请输入锁屏密码">
						<template #suffix>
							<el-icon class="el-input__icon" @click="unLockScreen">
								<Promotion />
							</el-icon>
						</template>
					</el-input>

				</el-form-item>
			</el-form>
		</div>
	</div>
</transition>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted } from 'vue'
import { useSettingsStore } from "@/stores/settings";
import { errMsg, successMsg } from '@/utils/msg';
import { Lock, Promotion } from '@element-plus/icons-vue';

const password = ref('');
const passwordInput: any = ref(null);
const settingsStore = useSettingsStore();

const unLockScreen = () => {
	if (settingsStore.checkLockPassword(password.value)) {
		settingsStore.unLockScreen();
		successMsg('解锁成功');
	}
	else {
		errMsg('密码错误');
	}
};

onMounted(() => {
	nextTick(() => {
		if (passwordInput.value) {
			passwordInput.value.focus();
		}
	});
});
</script>

<style lang="scss" scoped>
.lockscreen {
	height: 100vh;
	width: 100vw;
	display: flex;
	justify-content: center;
	align-items: center;
	background-color: #121212; // 深色背景
}

.lock-card {
	border-radius: 12px; // 圆角更大
	display: flex;
	flex-direction: column;
	justify-content: center;
	align-items: center;
	width: 30vw; // 宽度增加
	height: 40vh; // 高度增加
	background-color: rgba(255, 255, 255, 0.5); // 半透明背景
	box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2); // 更深的阴影
	padding: 30px; // 内边距增加
	position: fixed;
	top: 20vh; // 调整位置
	left: 50%;
	transform: translateX(-50%);
	z-index: 1000;
	color: #ffffff; // 白色文字
	font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; // 使用 Segoe UI 字体
}

.avatar-container {
	margin-bottom: 30px; // 间距增加
	display: flex;
	justify-content: center;
	align-items: center;
}

.el-input {
	margin-top: 20px; // 间距调整
	margin-bottom: 20px; // 间距调整
	height: 40px; // 高度增加
	font-size: 16px; // 字体大小增加
	border: none; // 去掉边框

	&::placeholder {
		color: #ffffff; // 占位符颜色
		opacity: 0.5; // 占位符透明度
	}

	&:focus {
		background-color: rgba(255, 255, 255, 0.2); // 聚焦时背景颜色
	}
}
/* 定义进入和离开动画 */
.slide-enter-active,
.slide-leave-active {
	transition: all 1s ease;
}

.slide-enter-from,
.slide-leave-to {
	transform: translateY(-100%);
	opacity: 0;
}

</style>