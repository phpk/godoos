<script setup>
	import { useChatStore } from "@/stores/chat";
	import { Setting as SettingIcon } from "@element-plus/icons-vue";

	const store = useChatStore();
</script>

<template>
	<!-- 头像 -->
	<el-row>
		<el-avatar
			shape="square"
			:size="40"
			class="userAvatar"
			:src="store.userInfo.avatar"
		/>
	</el-row>

	<!-- 菜单按钮 -->
	<el-row @click="store.setCurrentNavId(0)">
		<div class="menu-icon-box">
			<el-icon
				v-if="store.currentNavId === 0"
				class="menu-icon-on"
			>
				<Comment />
			</el-icon>
			<el-icon
				v-else
				class="menu-icon"
			>
				<Comment />
			</el-icon>
		</div>
	</el-row>

	<el-row @click="store.setCurrentNavId(1)">
		<div class="menu-icon-box">
			<el-icon
				v-if="store.currentNavId === 1"
				class="menu-icon-on"
			>
				<UserFilled />
			</el-icon>
			<el-icon
				v-else
				class="menu-icon"
			>
				<User />
			</el-icon>
		</div>
	</el-row>

	<el-row @click="store.setCurrentNavId(2)">
		<div class="menu-icon-box">
			<el-icon
				v-if="store.currentNavId === 2"
				class="menu-icon-on"
			>
				<Platform />
			</el-icon>
			<el-icon
				v-else
				class="menu-icon"
			>
				<Monitor />
			</el-icon>
		</div>
	</el-row>

	<!-- 设置按钮固定在底部 -->
	<el-row
		class="settings-icon-row"
		@click="store.setCurrentNavId(5)"
	>
		<div class="menu-icon-box">
			<el-icon
				v-if="store.currentNavId === 5"
				class="menu-icon-on"
			>
				<Tools />
			</el-icon>
			<el-icon
				v-else
				class="menu-icon"
			>
					<Tools />
			</el-icon>
		</div>
	</el-row>
</template>

<style scoped>
	.chat-user-info {
		position: relative;
		height: 100%;
	}

	.userAvatar {
		margin: 6px auto;
		-webkit-app-region: no-drag;
	}

	.menu-icon-box {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 40px;
		height: 40px;
		margin: 10px auto;
		transition: all 0.5s;
	}

	.menu-icon {
		margin: 0 auto;
		font-size: 25px;
		cursor: pointer;
		color: #86909C;
	}

	.menu-icon-on {
		margin: 0 auto;
		font-size: 25px;
		color: #0d42d2;
		cursor: pointer;
	}

	/* 设置按钮独立在底部 */
	.settings-icon-row {
		position: absolute;
		bottom: 20px;
		left: 0;
		right: 0;
		width: 55px;
		height: 55px;
	}
</style>
