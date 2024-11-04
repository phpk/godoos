<template>
	<div
		v-if="store.targetChatId == null"
		class="no-message-container"
	>
		<el-icon
			:size="180"
			color="#0078d7"
		>
			<ChatDotSquare />
		</el-icon>
		<p>欢迎使用GodoOS</p>
	</div>
	<div
		v-else
		class="user-info-container"
	>
		<div class="user-details">
			<p>昵称: {{ store.targetUserInfo.displayName }}</p>
			<p>邮箱: {{ store.targetUserInfo.email }}</p>
			<p>电话: {{ store.targetUserInfo.phone }}</p>
			<p>描述: {{ store.targetUserInfo.desc }}</p>
			<p>工号: {{ store.targetUserInfo.jobNumber }}</p>
			<p>工作地点: {{ store.targetUserInfo.workPlace }}</p>
			<p>入职日期: {{ store.targetUserInfo.hiredDate }}</p>
		</div>
		<el-button
			type="primary"
			@click="sendMessage(store.targetChatId)"
			>发送消息</el-button
		>
	</div>
</template>

<script lang="ts" setup>
	import { useChatStore } from "@/stores/chat";

	const store = useChatStore();

	const sendMessage = (chatId: string) => {
		store.currentNavId = 0;
		store.addChatListAndGetChatHistory(chatId);
	};
</script>

<style scoped>
	.no-message-container {
		height: 100%;
		margin: 0px auto;
		text-align: center;
		font-size: 14px;
	}

	.user-info-container {
		height: 100%;
		margin: 0px auto;
		text-align: center;
		font-size: 14px;
		display: flex;
		flex-direction: column;
		align-items: center;
		justify-content: center;
	}

	.user-details {
		margin-bottom: 20px;
		text-align: left;
		max-width: 300px;
	}

	.user-details p {
		margin: 5px 0;
	}

	.el-button {
		width: 100%;
		max-width: 300px;
	}
</style>
