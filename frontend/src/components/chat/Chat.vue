<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
	import { getWorkflowUrl } from "@/system/config";
	import { Search } from "@element-plus/icons-vue";
	const store = useChatStore();
	const workUrl = getWorkflowUrl();
	onMounted(() => {
		store.initChat();
	});

	const generateData = () => {
		return store.allUserList.map((user: any) => ({
			key: user.id,
			label: user.nickname,
			avatar: user.avatar, // 添加头像数据
		}));
	};

	const data = ref(generateData());
	// 声明 users 时指定类型为 any[]
	const users = ref<any[]>([]);

	watchEffect(() => {
		if (store.allUserList.length > 0) {
			data.value = generateData();
		}
	});

	// 监听搜索输入并更新 searchList
	watch(
		() => store.searchInput,
		(newSearchInput) => {
			if (newSearchInput === "") {
				store.searchList = []; // 当搜索输入为空时，清空搜索列表
			} else {
				store.searchList = store.chatList.filter(
					(user: { displayName: string | string[] }) =>
						user.displayName.includes(newSearchInput)
				);
			}
		}
	);
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
					v-model="store.searchInput"
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
		width="80%"
		title="创建群聊"
		style="height: 550px"
		align-center
	>
		<div class="dialog-body">
			<!-- 添加输入部门名的输入框 -->
			<div>
				<el-form label-position="top">
					<el-form-item label="群聊名称:">
						<el-input
							maxlength="8"
							show-word-limit
							style="width: 240px; height: 30px"
							class="department-name"
							v-model="store.departmentName"
							placeholder="请输入群聊名称"
						></el-input>
					</el-form-item>
				</el-form>
			</div>

			<!-- 使用 Element 的 el-transfer 组件替换自定义穿梭框 -->
			<el-transfer
				v-model="users"
				:data="data"
				:titles="['可选项', '已选项']"
				filterable
				style="height: 250px"
				filter-placeholder="搜索用户名"
				:props="{ key: 'key', label: 'label', avatar: 'avatar' }"
				class="transfer-container"
			>
				<!-- 自定义穿梭框列表项模板 -->
				<template #default="{ option }">
					<el-avatar
						:src="option.avatar"
						size="small"
						style="margin-right: 5px"
					/>
					<span>{{ option.label }}</span>
				</template>
			</el-transfer>
		</div>

		<template #footer>
			<span class="dialog-footer">
				<el-button
					style="background-color: #0078d4; color: #fff"
					@click="store.groupChatInvitedDialogVisible = false"
					>取消</el-button
				>
				<el-button
					style="background-color: #0078d4; color: #fff"
					@click="store.createGroupChat(users)"
					>确定</el-button
				>
			</span>
		</template>
	</el-dialog>
</template>

<style scoped>
	:deep(.el-transfer) {
		display: flex;
		flex-direction: row; /* 将布局方向设置为横向 */
		width: 550px; /* 让穿梭框占满宽度 */
		justify-content: center;
		align-items: center;
		.el-transfer__buttons {
			display: flex;
			justify-content: space-evenly;
			align-items: center;
			width: 120px;
			padding: 0 0 !important;
			.el-button {
				width: 40px !important;
				height: 30px !important;
				background-color: #0078d4;
				color: #fff;
			}
		}
	}

	::v-deep .el-dialog .el-dialog__header {
		border: none !important;
	}

	::v-deep .el-transfer-panel {
		width: 270px !important; /* 设置每个穿梭框面板的宽度 */
		height: 250px !important;
	}

	::v-deep .el-transfer-panel__body {
		height: 200px !important;
	}

	::v-deep .el-checkbox__label {
		margin-left: 10px !important;
	}

	.container {
		display: flex;
		height: 100%;
		width: 100%;
		overflow-y: hidden;
		overflow-x: hidden;
		border-top: 1px solid #e8e8e8;
	}

	.menu {
		width: 55px;
		background-color: #ffffff;
		overflow-y: hidden;
		overflow-x: hidden;
		-webkit-app-region: drag;
	}

	.side {
		flex: 1;
		/* 占据剩余宽度 */
    max-width: 200px;
    min-width: 200px;
		max-height: max-content;
		border-right: 1px solid #e8e8e8;
		overflow-y: hidden;
		overflow-x: hidden;
	}

	.search {
		display: flex;
		align-items: center;
		justify-content: space-evenly;
		width: 100%;
		height: 50px;
		padding: 0;
		padding-right: 10px;
		-webkit-app-region: drag;
		border-bottom: 1px solid #e8e8e8;
		border-left: 1px solid #e8e8e8;
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
		margin: 10px;
		height: 32px;
		-webkit-app-region: no-drag;
		--el-input-placeholder-color: #bfbfbf !important;
		--el-input-icon-color: #bfbfbf !important;
	}

	.list {
		width: 100%;
		height: 100%;
		padding: 0;
		overflow-y: hidden;
		overflow-x: hidden;
		border-left: 1px solid #e8e8e8;
	}
	.dialog-body {
		width: 100%;
		height: 350px;
	}

	.chat-box {
		flex: 3;
		max-height: max-content;
		background-color: #f5f5f5;
	}

	.chat-setting {
		width: calc(100% - 65px);
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

	.department-name {
		margin: 10px 0;
	}
</style>
