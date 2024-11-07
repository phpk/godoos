<template>
	<div
		v-if="store.targetChatId == null"
		style="height: 280px; width: 100%"
		class="no-message-container"
	>
		<el-icon
			:size="180"
			color="#0078d7"
		>
			<ChatDotSquare />
		</el-icon>
		<p class="welcome-text">欢迎使用GodoOS</p>
	</div>

	<div
		v-else-if="
			store.targetUserInfo && Object.keys(store.targetUserInfo).length > 0
		"
		class="user-info-container"
	>
		<el-avatar
			class="user-avatar"
			:size="80"
			:src="store.targetUserInfo.avatarUrl || ''"
		/>
		<h2 class="user-name">{{ store.targetUserInfo.displayName }}</h2>
		<div class="user-details">
			<p><strong>邮箱:</strong> {{ store.targetUserInfo.email }}</p>
			<p><strong>电话:</strong> {{ store.targetUserInfo.phone }}</p>
			<p><strong>描述:</strong> {{ store.targetUserInfo.desc }}</p>
			<p><strong>工号:</strong> {{ store.targetUserInfo.jobNumber }}</p>
			<p>
				<strong>工作地点:</strong> {{ store.targetUserInfo.workPlace }}
			</p>
			<p>
				<strong>入职日期:</strong> {{ store.targetUserInfo.hiredDate }}
			</p>
		</div>
		<el-button
			class="send-message-btn"
			@click="sendMessage(store.targetChatId, 'user')"
		>
			发送消息
		</el-button>
	</div>
	<!-- 群聊信息 -->
	<div
		v-else
		class="group-info-container"
	>
		<el-avatar
			class="group-avatar"
			:size="80"
		><span style="font-size: 40px;">群</span></el-avatar>

		<h2 class="group-name">{{ store.targetGroupInfo.displayName }}</h2>

		<el-button
			class="send-message-btn"
			@click="sendMessage(store.targetChatId, 'group')"
		>
			发送消息
		</el-button>
	</div>
</template>

<script lang="ts" setup>
	import { useChatStore } from "@/stores/chat";

	const store = useChatStore();

	const sendMessage = (chatId: string, type: string) => {
		store.currentNavId = 0;
		console.log(chatId, type);
		store.addChatListAndGetChatHistory(chatId, type);
	};
</script>

<style scoped>
	/* 无消息提示样式 */
	.no-message-container {
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		text-align: center;
		color: #4a4a4a;
	}

	.welcome-text {
		font-size: 18px;
		font-weight: 600;
		margin-top: 10px;
		color: #333;
	}

	/* 用户信息容器 */
	.user-info-container {
		display: flex;
		width: 100%;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		padding: 20px;
		text-align: center;
	}

	/* 用户头像 */
	.user-avatar {
		margin-bottom: 15px;
		border-radius: 50%;
	}

	/* 用户名标题 */
	.user-name {
		font-size: 22px;
		font-weight: 700;
		color: #333;
		margin-bottom: 15px;
	}

	/* 用户详细信息 */
	.user-details {
		text-align: left;
		background-color: #f8f9fa;
		padding: 20px;
		border-radius: 10px;
		box-shadow: 0px 1px 1px 1px rgba(0, 0, 0, 0.1);
		max-width: 200px;
		width: 100%;
		margin-bottom: 30px;
	}

	.user-details p {
		margin: 8px 0;
		color: #333;
	}

	.user-details p strong {
		color: #333;
	}

	/* 发送消息按钮 */
	.send-message-btn {
		width: 100%;
		max-width: 240px;
		font-size: 16px;
		background-color: #0078d4;
		color: #fff;
		border-radius: 8px;
		box-shadow: 0 4px 8px rgba(0, 120, 212, 0.3);
		font-weight: 600;
		transition: background-color 0.3s;
	}

	.send-message-btn:hover {
		background-color: #005a9e;
	}

	/* 群聊信息容器 */
	.group-info-container {
		display: flex;
		width: 100%;
		flex-direction: column;
		align-items: center;
		justify-content: center;
		height: 100%;
		padding: 20px;
		text-align: center;
	}
</style>
