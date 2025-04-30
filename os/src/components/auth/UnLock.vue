<template>
	<div class="unlock-container">
		<div class="user-info">
			<el-avatar class="avatar">
				<el-icon style="font-size: 25px">
					<User />
				</el-icon>
			</el-avatar>
			<span class="username">{{ displayName }}</span>
		</div>
		<el-input
			v-model="password"
			placeholder="输入密码"
			show-password
			class="password-input"
			clearable
			@keyup.enter="unlock"
		>
			<template #prefix>
				<el-icon>
					<Lock />
				</el-icon>
			</template>
		</el-input>
		<div class="gohome">
			<el-icon><ArrowRightBold /></el-icon>
		</div>
	</div>
</template>

<script setup lang="ts">
	import { ElMessage } from "element-plus";
	import { ref } from "vue";

	const username = ref("");
	const password = ref("");
	const displayName = "admin";

	function unlock() {
		if (
			username.value === "your_username" &&
			password.value === "your_password"
		) {
			ElMessage.success("解锁成功");
		} else {
			ElMessage.error("用户名或密码错误");
		}
	}
</script>

<style scoped>
	.unlock-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100vh;
		gap: 20px;
		padding: 20px;
		width: 100%;
		opacity: 0;
		transition: opacity 0.5s ease-in-out, background 0.5s ease-in-out,
			box-shadow 0.5s ease-in-out;
		animation: fadeIn 1s forwards;
	}

	@keyframes fadeIn {
		to {
			opacity: 1;
		}
	}

	.unlock-container:focus-within {
		background: rgba(255, 255, 255, 0.1); /* 半透明白色背景 */
		box-shadow: 0 15px 30px rgba(0, 0, 0, 0.5);
		backdrop-filter: blur(10px); /* 磨砂玻璃效果 */
		-webkit-backdrop-filter: blur(10px); /* 兼容Safari */
	}

	.user-info {
		display: flex;
		flex-direction: column;
		align-items: center;
	}

	.avatar {
		width: 60px;
		height: 60px;
		border-radius: 50%;
		box-shadow: 0 4px 8px rgba(0, 0, 0, 0.2);
		transition: transform 0.3s;
	}

	.avatar:hover {
		transform: scale(1.1);
	}

	.username {
		font-size: 20px;
		font-weight: bold;
		margin-top: 10px;
		color: #333;
	}

	.password-input {
		width: 250px;
		overflow: hidden;
		box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
		transition: box-shadow 0.3s;
	}

	.password-input:focus {
		box-shadow: 0 6px 12px rgba(0, 0, 0, 0.2);
	}
	.gohome {
		display: flex;
		justify-content: center;
		align-items: center;
		width: 30px;
		height: 30px;
		border-radius: 50px;
	}
	.gohome:hover {
		background-color: rgba(0, 0, 0, 0.1); 
	}
</style>
