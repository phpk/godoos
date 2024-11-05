<template>
	<div
		class="chatbox-main"
		v-if="store.targetChatId"
	
  
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
					>{{ store.targetGroupInfo.displayName }}{{ 1 }}</span
				>

				<div @click="store.drawerVisible = true">
					<el-icon><Tools /></el-icon>
					<!-- 抽屉 -->
					<el-drawer
						v-model="store.drawerVisible"
						direction="ltr"
						title="群设置"
						:with-header="true"
					>
						<div
							style="
								display: flex;
								flex-direction: column;
								justify-content: space-between;
								height: 100%;
							"
						>
							<div>
								<!-- 群名 -->
								<div class="group-name">
									{{ store.targetGroupInfo.displayName }}
								</div>

								<!-- 群成员 -->
								<div class="group-member">
									<div><span>群成员</span></div>
									<div class="group-member-list">
										<button
											@click="store.addMember"
											class="group-member-add"
										>
											<el-icon><Plus /></el-icon>
										</button>
										<div
											style="
												display: flex;
												flex-direction: row;
											"
											v-for="member in store.groupMemberList"
										>
											<el-avatar :src="member.avatar" />
										</div>
									</div>
								</div>
							</div>

							<!-- 退出按钮 -->
							<div class="group-exit">
								<el-button
									style="
										background-color: #0078d4;
										color: #fff;
									"
									@click="
										store.quitGroup(
											store.targetGroupInfo.group_id
										)
									"
									>退出群聊</el-button
								>
							</div>
						</div>
					</el-drawer>
				</div>

				<!-- 邀请好友对话框 -->
				<el-dialog
					v-model="store.inviteFriendDialogVisible"
					title="邀请好友"
					width="80%"
					style="height: 550px"
					align-center
				>
					<div>
						<el-transfer
							v-model="value"
							:data="data"
						/>
					</div>

					<template #footer>
						<span class="dialog-footer">
							<el-button
								style="background-color: #0078d4; color: #fff"
								@click="store.inviteFriendDialogVisible = false"
								>取消</el-button
							>
							<el-button
								style="background-color: #0078d4; color: #fff"
								>确定</el-button
							>
						</span>
					</template>
				</el-dialog>
			</div>
		</el-header>

		<!--聊天主体区-->
		<el-main class="msg-main">
			<el-scrollbar ref="store.scrollbarRef">
				<div ref="store.innerRef">
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
		class="no-message-container"
		v-else
	>
		<el-icon
			:size="180"
			color="#0078d7"
		>
			<ChatDotSquare />
		</el-icon>
		<p>欢迎使用GodoOS</p>
	</div>
</template>

<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
	import { useChooseStore } from "@/stores/choose";
	const store: any = useChatStore();
	const choose = useChooseStore();
	const imgExt = ["png", "jpg", "jpeg", "gif", "bmp", "webp", "svg"];
	const choosetype = ref("");

	// 监听store.drawerVisible
	watch(
		() => store.drawerVisible,
		(newVal, _) => {
			if (newVal) {
				store.getGroupMember(store.targetGroupInfo.group_id);
			}
		}
	);

	const generateData = () => {
		return store.allUserList.map((user: any) => ({
			key: user.id,
			label: user.nickname,
			avatar: user.avatar, // 添加头像数据
		}));
	};

	const data = ref(generateData());
	const value = ref([]);

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
