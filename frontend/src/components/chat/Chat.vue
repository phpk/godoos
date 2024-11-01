<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
	import { getWorkflowUrl } from "@/system/config";
	import { Search } from "@element-plus/icons-vue";
	const store = useChatStore();
	const workUrl = getWorkflowUrl();
	onMounted(() => {
		store.initChat();
	});

	// 将用户列表转换为 el-transfer 组件所需的数据格式
	const generateData = () => {
		return store.allUserList.map((user: any) => ({
			key: user.id,
			label: user.nickname,
			avatar: user.avatar, // 添加头像数据
		}));
	};

	const data = ref(generateData());
	const users = ref([]);

	watchEffect(() => {
		if (store.allUserList.length > 0) {
			data.value = generateData();
		}
	});
</script>
<template>
	<el-container class="container">
		<!--菜单-->
		<el-aside class="menu">
			<chat-menu />
		</el-aside>
		<el-container
			class="side"
			v-if="store.currentNavId < 3"
		>
			<!--搜索栏-->
			<el-header
				class="search"
				v-if="store.currentNavId < 2"
			>
				<el-input
					placeholder="搜索"
					:prefix-icon="Search"
					class="search-input"
					v-model="store.search"
				/>
				<!-- 邀请群聊 -->
				<button
					class="inviteGroupChats"
					@click="store.setGroupChatInvitedDialogVisible(true)"
				>
					<el-icon><Plus /></el-icon>
				</button>
			</el-header>
			<!--好友列表-->
			<el-main class="list">
				<el-scrollbar>
					<chat-msg-list v-if="store.currentNavId == 0" />
					<chat-user-list v-if="store.currentNavId == 1" />
				</el-scrollbar>
			</el-main>
		</el-container>
		<el-container class="chat-box">
			<chat-box v-if="store.currentNavId < 1" />
			<chat-user-info v-if="store.currentNavId == 1"></chat-user-info>
		</el-container>
		<el-container
			class="chat-setting"
			v-if="store.currentNavId == 2"
		>
			<iframe
				class="workflow"
				:src="workUrl"
			></iframe>
		</el-container>
		<el-container
			class="chat-setting"
			v-if="store.currentNavId == 5"
		>
			<ChatUserSetting />
		</el-container>
	</el-container>

	<!-- 邀请群聊弹窗 -->
	<el-dialog
		v-model="store.groupChatInvitedDialogVisible"
		title="发起群聊"
		width="80%"
	>
		<div class="dialog-body">
			<!-- 添加输入部门名的输入框 -->
			<el-input
				v-model="store.departmentName"
				placeholder="请输入群聊名称"
			></el-input>
			<el-transfer
				v-model="users"
				:data="data"
			>
				<!-- 义穿梭框 -->
				<template #default="{ option }">
					<div class="user-item">
						<el-avatar
							:size="20"
							:src="option.avatar"
							class="avatar"
						/>
						<span>{{ option.label }}</span>
					</div>
				</template>
			</el-transfer>
		</div>
		<template #footer>
			<span class="dialog-footer">
				<el-button @click="store.groupChatInvitedDialogVisible = false"
					>取消</el-button
				>
				<el-button
					type="primary"
					@click="store.createGroupChat(users)"
					>确定</el-button
				>
			</span>
		</template>
	</el-dialog>
</template>

<style scoped>
	.container {
		display: flex;
		height: 100%;
		width: 100%;
		overflow-y: hidden;
		overflow-x: hidden;
	}

	.menu {
		width: 55px;

		background-color: #f0f0f0;
		overflow-y: hidden;
		overflow-x: hidden;
		-webkit-app-region: drag;
	}

	.side {
		flex: 1;
		/* 占据剩余宽度 */
		max-height: max-content;
		border-right: 1px solid #edebeb;
		overflow-y: hidden;
		overflow-x: hidden;
	}

	.search {
		display: flex;
		align-items: center;
		justify-content: space-evenly;
		width: 90%;
		/* 占据整个宽度 */
		height: 50px;
		padding: 0;
		-webkit-app-region: drag;
	}
	.inviteGroupChats {
		width: 40px;
		height: 30px;
		border: none;
		border-radius: 4px;
		background-color: #f0f0f0;
	}

	.user-item {
		width: 100%;
		height: 30px;
		display: flex;
		align-items: center;
	}

	.search-input {
		width: calc(100% - 20px);
		/* 减去左右边距 */
		margin: 10px;
		height: 32px;
		-webkit-app-region: no-drag;
		--el-input-placeholder-color: #bfbfbf !important;
		--el-input-icon-color: #bfbfbf !important;
	}

	.list {
		width: 100%;
		/* 占据整个宽度 */
		padding: 0;
		overflow-y: hidden;
		overflow-x: hidden;
	}
	.dialog-body {
		width: 100%;
	}
	.el-transfer {
		display: flex;
		align-items: center;
		justify-content: center;
	}

	/* .el-transfer >>> .el-transfer-panel >>> .el-transfer__buttons {
		width: 50px;
	} */
	/* .el-transfer__button {
    width: 20px;;
	} */

	.chat-box {
		flex: 3;
		/* 占据剩余宽度的三倍 */
		max-height: max-content;
		background-color: #f5f5f5;
	}

	.chat-setting {
		width: calc(100% - 65px);
		/* 占据整个宽度 */
		height: 100%;
		overflow: hidden;
	}
	.workflow {
		width: 100%;
		height: 100%;
		object-fit: contain;
		border: none;
	}

	.no-message-container {
		height: 100%;
		margin: 120px auto;
		text-align: center;
		font-size: 14px;
		justify-content: center;
	}
</style>
