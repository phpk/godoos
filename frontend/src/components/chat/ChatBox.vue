<template>
	<!-- 邀请群聊对话框 -->
	<el-dialog
		v-model="store.inviteFriendDialogVisible"
		title="邀请群聊"
		width="80%"
		style="height: 450px"
		align-center
	>
		<div>
			<el-transfer
				:titles="['可选项', '已选项']"
				filterable
				filter-placeholder="搜索用户名"
				style="height: 300px"
				v-model="users"
				:data="data"
				:props="{ key: 'key', label: 'label', avatar: 'avatar' }"
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
			<el-button
				style="
					background-color: #0078d4;
					color: #fff;
					position: absolute;
					bottom: 10px;
					right: 120px;
				"
				@click="store.inviteFriendDialogVisible = false"
				>取消</el-button
			>
			<el-button
				@click="
					store.inviteFriend(store.targetGroupInfo.group_id, users)
				"
				class="invite-group-button"
				style=""
				>确定</el-button
			>
		</div>
	</el-dialog>
	<div
		class="chatbox-main"
		v-if="store.targetChatId"
	>
		<!--聊天顶部区-->
		<el-header class="chat-header">
			<div class="header-title">
				<span
					v-if="store.targetUserInfo.displayName"
					class="header-title-name"
					>{{ store.targetUserInfo.displayName }}</span
				>
				<span
					v-else-if="store.targetGroupInfo.displayName"
					class="header-title-name"
					>{{ store.targetGroupInfo.displayName }}</span
				>

				<div
					v-if="
						store.targetGroupInfo &&
						Object.keys(store.targetGroupInfo).length > 0
					"
				>
					<div style="display: flex; gap: 10px">
						<el-dropdown
							placement="bottom"
							style="border: none"
						>
							<el-icon
								style="
									cursor: pointer;
									color: black;
									font-size: 15px;
								"
								><More
							/></el-icon>
							<template #dropdown>
								<el-dropdown-menu>
									<el-dropdown-item
										@click="openInviteGroupDialog()"
										>邀请群聊</el-dropdown-item
									>
									<el-dropdown-item
										@click="
											store.quitGroup(
												store.targetGroupInfo.group_id
											)
										"
										>退出群聊</el-dropdown-item
									>
								</el-dropdown-menu>
							</template>
						</el-dropdown>
						<el-icon
							style="
								cursor: pointer;
								color: black;
								font-size: 15px;
							"
							@click="openDrawer()"
							><Tools
						/></el-icon>
					</div>
					<!-- 抽屉 -->
					<el-drawer
						v-model="store.drawerVisible"
						direction="ltr"
						title="群成员"
						:with-header="true"
					>
						<div
							style="
								display: flex;
								height: 100%;
								padding-top: 40px;
								flex-direction: column;
								justify-content: space-between;
							"
						>
							<div
								class="group-member-container"
								style="
									display: flex;
									flex-wrap: wrap;
									max-width: 248px; /* 每行4个用户，每个60px + 2px间距 */
								"
							>
								<div
									class="group-member"
									v-for="member in store.groupMembers"
									:key="member.id"
									style="
										display: flex;
										flex-direction: column;
										align-items: center;
										justify-content: center;
										width: 50px;
										height: 70px;
									"
								>
									<el-avatar
										size="40"
										src=""
									/>
									<div
										style="
											max-width: 46px;
											height: 20px;
											font-size: 10px;
											text-align: center;
											white-space: nowrap;
											overflow: hidden;
											text-overflow: ellipsis;
											display: inline-block;
											line-height: 20px;
										"
									>
										{{ member.nickname }}
									</div>
								</div>
							</div>
						</div>
					</el-drawer>
				</div>

				<div
					v-else
					style="
						height: 50px;
						display: flex;
						align-items: center;
						justify-content: center;
					"
				>
					<el-dropdown>
						<el-icon
							style="
								cursor: pointer;
								color: black;
								font-size: 15px;
							"
							><Tools
						/></el-icon>
						<template #dropdown>
							<el-dropdown-menu>
								<el-dropdown-item>删除好友</el-dropdown-item>
							</el-dropdown-menu>
						</template>
					</el-dropdown>
				</div>
			</div>
		</el-header>

		<!--聊天主体区-->
		<el-main class="msg-main">
			<el-scrollbar ref="scrollbarRef">
				<div ref="innerRef">
					<ChatMessage />
				</div>
			</el-scrollbar>
		</el-main>

		<!--聊天输入区和发送按钮等-->
		<el-footer
			class="msg-footer"
			style="
				display: flex;
				flex-direction: column;
				justify-content: center;
				align-items: start;
			"
		>
			<!-- 输入主体 -->
			<div
				class="input-main"
				style="width: 100%; height: 40px; background-color: #ffffff"
			>
				<el-input
					size="large"
					style="width: 100%; height: 100%"
					placeholder="请输入内容"
					@keyup.enter.exact="store.sendMessage('text')"
					v-model="store.message"
				>
					<template #suffix>
						<el-icon
							:size="20"
							class="input-option-icon"
							@click="selectImg"
						>
							<Picture />
						</el-icon>
						<el-icon
							:size="20"
							class="input-option-icon"
							@click="selectFile"
						>
							<Link />
						</el-icon>
						<el-icon
							:size="20"
							class="input-option-icon"
						>
							<Delete />
						</el-icon>
						<el-icon
							@click="store.sendMessage('text')"
							:size="22"
							class="input-button"
						>
							<Promotion />
						</el-icon>
					</template>
				</el-input>
			</div>
		</el-footer>
	</div>

	<div
		style="
			width: 100%;
			height: 100%;
			display: flex;
			flex-direction: column;
			align-items: center;
			justify-content: center;
		"
		v-else
	>
		<el-icon
			:size="180"
			color="#0078d7"
		>
			<ChatDotSquare />
		</el-icon>
		<p
			style="
				font-size: 18px;
				font-weight: 600;
				margin-top: 10px;
				color: #333;
			"
		>
			欢迎使用GodoOS
		</p>
	</div>
</template>

<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
	import { useChooseStore } from "@/stores/choose";
	const store: any = useChatStore();
	const choose = useChooseStore();
	const imgExt = ["png", "jpg", "jpeg", "gif", "bmp", "webp", "svg"];
	const choosetype = ref("");
	const scrollbarRef = ref(null);
	const innerRef = ref(null);

	const openInviteGroupDialog = () => {
		store.inviteFriendDialogVisible = true;
		store.getInviteUserList();
	};

	function scrollToBottom() {
		store.setScrollToBottom(innerRef, scrollbarRef);
	}

	const generateData = () => {
		return store.inviteUserList.map((user: any) => ({
			key: user.id,
			label: user.nickname,
			avatar: user.avatar, // 添加头像数据
		}));
	};

	const data = ref([]);

	watchEffect(() => {
		if (store.inviteUserList && store.inviteUserList.length > 0) {
			data.value = generateData();
		}
	});
	// 声明 users 时指定类型为 any[]
	const users = ref<any[]>([]);

	watchEffect(() => {
		if (store.allUserList.length > 0) {
			data.value = generateData();
		}
	});

	// 监听store中的messageSendStatus.value = true，调用scrollToBottom
	watch(
		() => store.messageSendStatus,
		(newVal, _) => {
			console.log("messageSendStatus 变化了:", newVal);
			if (newVal) {
				console.log("messageSendStatus 变化了:", newVal);
				scrollToBottom();
			}
		}
	);

	// 监听store中的messageReceiveStatus，调用scrollToBottom
	watch(
		() => store.messageReceiveStatus,
		(newVal, _) => {
			if (newVal) {
				scrollToBottom();
			}
		}
	);

	function openDrawer() {
		store.drawerVisible = true;
		store.getGroupMemberList(store.targetGroupInfo.group_id);
	}

	// 监听store.drawerVisible
	watch(
		() => store.drawerVisible,
		(newVal, _) => {
			if (newVal) {
				store.getGroupMember(store.targetGroupInfo.group_id);
			}
		}
	);

	function selectImg() {
		choosetype.value = "image";
		console.log(choosetype.value);
		choose.select("选择图片", imgExt);
	}
	function selectFile() {
		console.log("selectFile");
		choosetype.value = "applyfile";
		choose.select("选择文件", "*");
	}

	watch(
		() => choose.path,
		(newVal, _) => {
			console.log("choose.path 变化了:", newVal);
			const paths = toRaw(newVal);
			if (paths.length > 0) {
				store.sendInfo = paths;
				choose.path = [];
				console.log(choosetype.value);
				store.sendMessage(choosetype.value);
			}
		},
		{ deep: true } // 添加deep: true以深度监听数组或对象内部的变化
	);
</script>

<style scoped>
	.invite-group-button {
		background-color: #0078d4;
		color: #fff;
		position: absolute;
		bottom: 10px;
		right: 50px;
	}

	.infinite-list {
		display: flex;
		flex-direction: column;
		align-items: start;
		padding: 0;
		width: 100%;
		margin: 0;
	}
	.infinite-list-item {
		width: 100%;
		display: flex;
		justify-content: center;
		align-items: center;
		margin-bottom: 10px; /* 每个成员之间的间距 */
	}

	:deep(.el-drawer__header) {
		color: #000000;
	}

	:deep(.el-drawer__body) {
		margin-top: -40px; /* 根据需要调整这个值 */
	}

	.group-name {
		font-size: 16px;
		border-top: 1px solid #edebeb;
		border-bottom: 1px solid #edebeb;
		display: flex;
		align-items: center;
		justify-content: start;
		height: 50px;
	}

	.group-member-list {
		display: flex;
		gap: 5px;
		width: 100%;
		align-items: center;
		height: 50px;
		overflow: hidden;
		padding-bottom: 10px;
		border-bottom: 1px solid #edebeb;
	}

	.group-member-add {
		background-color: #0078d4;
		border-radius: 50px;
		width: 40px;
		height: 40px;
		color: #fff;
		border: none;
	}

	.group-member {
		display: flex;
		flex-direction: column;
		align-items: start;
		justify-content: start;
	}

	.el-transfer {
		display: flex;
		flex-direction: row; /* 将布局方向设置为横向 */
		width: 550px; /* 让穿梭框占满宽度 */
	}

	.el-transfer-panel {
		width: 300px !important;
		height: 530px !important;
	}

	.el-transfer-panel__body {
		height: 450px !important;
	}

	.el-checkbox__label {
		margin-left: 10px !important;
	}

	.el-transfer-panel
		.el-transfer-panel__header
		.el-checkbox
		.el-checkbox__label
		span {
		left: 150px;
		right: 0px;
	}

	.el-transfer__buttons {
		display: flex;
		flex-direction: column;
		align-items: center !important;
		/* 水平居中对齐 */
		justify-content: center !important;
		gap: 10px;
		/* 子元素之间的间距 */
		padding: 0 15px;
	}

	.el-transfer__buttons .el-button {
		min-width: 35px !important;
		text-align: center;
		margin-left: 0 !important;
	}

	.el-transfer__buttons .el-transfer__button {
		width: 35px;
		height: 35px;
		border-radius: 50%;
	}

	.chatbox-main {
		width: 100%;
		height: 100%;
	}

	.chat-header {
		width: 100%;
		/* 占据整个宽度 */
		height: 49px;
		line-height: 50px;
	}

	.header-title {
		display: flex;
		align-items: center;
		justify-content: space-between;
	}

	.header-title-name {
		height: 50px;
		line-height: 50px;
		font-size: 20px;
	}

	.msg-main {
		width: 100%;
		/* 占据整个宽度 */
		height: 75%;
		padding: 0;
		border-top: 1px solid #edebeb;
	}

	.msg-footer {
		width: 100%;
		/* 占据整个宽度 */
		height: calc(100% - 75% - 49px);
		border: none;
	}

	.input-option {
		height: 20px;
		padding: 5px;
	}

	.input-option-icon {
		margin-left: 10px;
		color: #656a72;
		cursor: pointer;
	}

	.input-textarea {
		height: calc(100% - 50px);
		width: calc(100% - 20px);
		/* 减去左右边距 */
		margin: 10px;
	}

	.textarea {
		font-size: 16px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		width: 100%;
		height: 100%;
		overflow-y: hidden;
		overflow-x: hidden;
		--el-input-border-radius: 0;
		--el-input-border-color: transparent;
		--el-input-hover-border-color: transparent;
		--el-input-clear-hover-color: transparent;
		--el-input-focus-border-color: transparent;
	}

	.input-button {
		margin-left: 10px;
		color: #2a6bf2;
	}

	.input-button:hover {
		background-color: #d1e4ff;
		/* 悬浮时颜色略深，保持浅色调 */
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
		/* 稍微增强阴影 */
	}

	.input-button:active {
		background-color: #b3d4fc;
		/* 按下时颜色更深，但依然保持清新 */
		box-shadow: 0 1px 2px rgba(0, 0, 0.1);
		/* 回复初始阴影 */
		transform: translateY(1px);
		/* 微小下移，模拟按下 */
	}
</style>
