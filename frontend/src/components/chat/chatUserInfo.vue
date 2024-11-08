<template>
	<div
		v-if="store.targetChatId"
		class="chat-user-info"
	>
		<!-- 用户信息区域 -->
		<div
			v-if="
				store.targetUserInfo &&
				Object.keys(store.targetUserInfo).length > 0
			"
			class="user-content"
		>
			<div class="user-details">
				<h2>{{ store.targetUserInfo.displayName }}</h2>
				<p>工号：{{ targetUserInfo.jobNumber }}</p>
				<p>岗位：{{ targetUserInfo.desc }}</p>
				<p>邮箱：{{ targetUserInfo.email }}</p>
				<p>电话：{{ targetUserInfo.phone }}</p>
				<p>入职日期：{{ targetUserInfo.hiredDate }}</p>
			</div>
			<div class="avatar">
				<el-avatar
					style="width: 80px; height: 80px"
					:src="targetUserInfo.avatar"
					alt="avatar"
				/>
			</div>
		</div>

		<!-- 群信息区域 -->
		<div
			v-else-if="
				store.targetGroupInfo &&
				Object.keys(store.targetGroupInfo).length > 0
			"
			class="group-content"
		>
			<div class="group-details">
				<h2>{{ store.targetGroupInfo.displayName }}</h2>
				<p>群ID：{{ store.targetGroupInfo.chatId }}</p>
				<p>群人数：{{ targetGroupInfo.memberCount }}人</p>
				<p>创建时间：{{ targetGroupInfo.createdAt }}</p>
			</div>
			<div class="group-avatar">
				<div class="avatar-container">
					<el-avatar
						style="width: 80px; height: 80px"
						:src="targetGroupInfo.avatar"
						alt="group-avatar"
					/>
				</div>
			</div>
		</div>

		<!-- 分割线 -->
		<div class="divider"></div>

		<!-- 发送按钮 -->
		<div class="send-button-container">
			<el-button
				type="primary"
				@click="
					sendMessage(
						store.targetGroupInfo?.chatId ||
							store.targetUserInfo.chatId,
						store.targetGroupInfo
							? 'group'
							: store.targetUserInfo.type
					)
				"
			>
				发送消息
			</el-button>
		</div>
	</div>
	<div v-else>1</div>
</template>

<script lang="ts" setup>
	import { useChatStore } from "@/stores/chat";
	const store = useChatStore();

	// 模拟用户信息
	const targetUserInfo = {
		avatar: "./logo.png",
		jobNumber: 12345678,
		desc: "测试岗位",
		email: "12345678@qq.com",
		phone: "12345678910",
		hiredDate: "2024-01-01",
	};

	// 模拟群信息
	const targetGroupInfo = {
		avatar: "./logo.png",
		chatId: "1234567890",
		memberCount: 100,
		createdAt: "2024-01-01",
	};

	const sendMessage = (chatId: string, type: string) => {
		store.currentNavId = 0;
		store.getGroupMemberList(chatId);
		store.addChatListAndGetChatHistory(chatId, type);
	};
</script>

<style scoped>
	.chat-user-info {
		display: flex;
		flex-direction: column;
		width: 100%;
		padding: 20px;
	}

	.user-content,
	.group-content {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.user-details,
	.group-details {
		flex: 1;
		display: flex;
		flex-direction: column;
	}

	.user-details h2,
	.group-details h2 {
		margin: 0;
		font-size: 1.5rem;
	}

	.user-details p,
	.group-details p {
		margin: 5px 0;
		color: #666;
	}

	.avatar,
	.group-avatar {
		width: 80px;
		height: 100%;
		object-fit: cover;
		margin-left: 20px;
	}

	.divider {
		width: 100%;
		height: 1px;
		background-color: #e0e0e0;
		margin: 20px 0;
	}

	.send-button-container {
		display: flex;
		justify-content: center;
		margin-top: 20px;
	}

	.el-button {
		background-color: #0d42d2;
		color: #fff;
	}
	.el-button:hover {
		background-color: #4080ff;
		color: #fff;
	}
</style>
