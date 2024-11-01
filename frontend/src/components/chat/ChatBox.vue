<template>
	<div
		class="chatbox-main"
		v-if="store.targetChatId"
	>
		<!--聊天顶部区-->
		<el-header class="chat-header">
			<div class="header-title">
				<span v-if="store.targetUserInfo.displayName">{{
					store.targetUserInfo.displayName
				}}</span>
				<span v-else-if="store.targetGroupInfo.displayName">{{
					store.targetGroupInfo.displayName
				}}</span>
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
		<el-footer class="msg-footer">
			<!--聊天输入选项-->
			<div class="input-option">
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
			</div>
			<!--聊天输入区-->
			<div class="input-textarea">
				<el-input
					type="textarea"
					maxlength="1000"
					resize="none"
					class="textarea"
					@keyup.enter.exact="store.sendMessage('text')"
					v-model="store.message"
				/>
			</div>
			<!--聊天发送按钮-->
			<el-tooltip
				placement="top"
				content="按enter键发送，按ctrl+enter键换行"
			>
				<el-icon
					@click="store.sendMessage('text')"
					:size="22"
					class="input-button"
				>
					<Promotion />
				</el-icon>
			</el-tooltip>
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
	const choosetype = ref("image");

	function selectImg() {
		choosetype.value = "image";
		choose.select("选择图片", imgExt);
	}
	function selectFile() {
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
				store.sendMessage(choosetype.value);
			}
		},
		{ deep: true } // 添加deep: true以深度监听数组或对象内部的变化
	);
</script>

<style scoped>
	.chatbox-main {
		width: 100%;
		height: 100%;
	}

	.chat-header {
		width: 100%;
		/* 占据整个宽度 */
		height: 50px;
		line-height: 50px;
		padding: 0;
		-webkit-app-region: drag;
	}

	header-title {
		font-size: 20px;
		text-align: left;
		margin-left: 15px;
	}

	.msg-main {
		width: 100%;
		/* 占据整个宽度 */
		height: calc(70% - 50px);
		padding: 0;
		border-top-width: 1px;
		border-bottom-width: 1px;
		border-left-width: 0;
		border-right-width: 0;
		border-color: #d6d6d6;
		border-style: solid;
	}

	.msg-footer {
		width: 100%;
		/* 占据整个宽度 */
		height: 30%;
		padding: 0;
	}

	.input-option {
		height: 20px;
		padding: 5px;
	}

	.input-option-icon {
		color: #494949;
		margin-left: 20px;
		margin-top: 5px;
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
		position: absolute;
		bottom: 12px;
		right: 12px;
		width: 30px;
		/* 缩小宽度 */
		height: 30px;
		/* 减小高度 */
		border-radius: 50%;
		/* 较小的圆角 */
		background-color: #e8f0fe;
		/* 浅蓝色，符合Win11的轻量风格 */
		color: #0078d4;
		/* 使用Win11的强调色作为文字颜色 */
		font-weight: bold;
		border: 1px solid #b3d4fc;
		/* 添加边框，保持简洁风格 */
		box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
		/* 轻微阴影 */
		transition: all 0.2s ease;
		/* 快速过渡效果 */
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

	.no-message-container {
		height: 100%;
		margin: 120px auto;
		text-align: center;
		justify-content: center;
	}
</style>
