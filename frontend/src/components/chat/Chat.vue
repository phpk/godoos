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
	// 声明 users 时指定类型为 any[]
const users = ref<any[]>([]);

	watchEffect(() => {
		if (store.allUserList.length > 0) {
			data.value = generateData();
		}
	});

  function toggleSelectItem(item: any) {
	const index = users.value.indexOf(item.key);
	if (index === -1) {
		users.value.push(item.key);
	} else {
		users.value.splice(index, 1);
	}
}

function removeItem(userId:string) {
	users.value = users.value.filter((user) => user !== userId);
}
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
		align-center
	>
		<div class="dialog-body">
			<!-- 添加输入部门名的输入框 -->
			<el-input
				v-model="store.departmentName"
				placeholder="请输入群聊名称"
			></el-input>
			<div class="transfer-container">
				<!-- 自定义穿梭框组件 -->
				<div class="transfer-box">
					<div class="list-box">
						<h3>可选项</h3>
						<ul>
							<li
								v-for="item in data"
								:key="item.key"
								@click="toggleSelectItem(item)"
								:class="{ selected: users.includes(item.key) }"
							>
								<el-avatar
									:size="10"
									:src="item.avatar"
									class="avatar"
								/>
								<input
									type="checkbox"
									:checked="users.includes(item.key)"
								/>
								<span>{{ item.label }}</span>
								<span
									v-if="users.includes(item.key)"
									class="remove-icon"
									@click.stop="removeItem(item.key)"
									>✖</span
								>
							</li>
						</ul>
					</div>
					<div class="list-box">
						<h3>已选项</h3>
						<ul>
							<li
								v-for="user in users"
								:key="user"
							>
								<span>{{
									data.find((item: any) => item.key === user)?.label
								}}</span>
								<span
									class="remove-icon"
									@click="removeItem(user)"
									>✖</span
								>
							</li>
						</ul>
					</div>
				</div>
			</div>
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

	.transfer-container >>> .el-transfer-panel {
		width: 300px;
	}
	.el-transfer {
		display: flex;
		align-items: center;
		justify-content: center;
	}

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

  .transfer-box {
	display: flex;
	gap: 0;
}

.list-box {
	width: 50%;
	border: 1px solid #ccc;
	border-radius: 5px;
	padding: 10px;
	box-shadow: 0px 4px 8px rgba(0, 0, 0, 0.1);
}

.list-box ul {
	list-style: none;
	padding: 0;
	margin: 0;
}

.list-box li {
	display: flex;
	align-items: center;
	padding: 5px 0;
	cursor: pointer;
	transition: background-color 0.2s;
}
.list-box li:hover {
	background-color: #f0f8ff;
}

.list-box li.selected {
	background-color: #e6f7ff;
}

input[type="checkbox"] {
	margin-right: 8px;
}

.remove-icon {
	margin-left: auto;
	color: #d32f2f;
	cursor: pointer;
	font-size: 14px;
}

.remove-icon:hover {
	color: #ff5a5a;
}
</style>
