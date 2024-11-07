<template>
	<div
		v-if="store.chatHistory && store.chatHistory.length > 0"
		v-for="item in store.chatHistory"
		style="margin-top: 10px"
		:key="item.chatId"
	>
		<div
			v-if="item.content_type == 'invite_group_message'"
			class="system-message"
		>
			{{ item.message }}
			<!-- <span v-if="item.isMe">{{ item.message }}</span>
			<span v-else>{{ item.message }}</span> -->
		</div>

		<div v-else>
			<div
				v-if="item.isMe"
				class="chat-item"
			>
				<el-row>
					<el-col :span="8" />
					<el-col :span="14">
						<el-row>
							<el-col :span="24">
								<div class="chat-name-me">
									{{ store.userInfo.nickname }}
								</div>
							</el-col>
						</el-row>
						<div
							class="bubble-me"
							@contextmenu.prevent="
								store.showContextMenu($event, item.chatId)
							"
						>
							<!-- 文本消息展示 -->
							<div
								v-if="item.content_type == 'text'"
								class="chat-font"
							>
								{{ item.message }}
							</div>

							<!-- 文件消息展示 -->
							<div
								@click="sys.openFile(item.file_path)"
								v-else-if="item.content_type == 'file'"
								:class="['chat-item-file', 'file-me']"
							>
								<div class="chat-item-file-icon">
									<el-icon
										size="30"
										color="#303133"
										><Document
									/></el-icon>
								</div>
								<div class="chat-item-file-name">
									{{ item.file_name }}
								</div>
							</div>

							<!-- 图片消息展示 -->
							<div
								v-else-if="item.content_type == 'image'"
								class=""
							>
								<el-image
									fit="cover"
									loading="lazy"
									:preview-src-list="[item.message]"
									:src="item.message"
								/>
							</div>
						</div>
						<!-- 发送时间展示，在消息框外部的下方 -->
						<div class="chat-time-me-outer">
							{{ item.previewTimeFormat }}
						</div>
					</el-col>
					<el-col :span="2">
						<div class="chat-avatar">
							<el-avatar
								shape="square"
								style="margin: 0; float: left"
								:size="32"
								class="userAvatar"
								:src="item.avatar"
							/>
						</div>
					</el-col>
				</el-row>
			</div>
			<div
				v-else
				class="chat-item"
			>
				<el-row>
					<el-col :span="2">
						<div class="chat-avatar">
							<el-avatar
								shape="square"
								style="margin: 0; float: right"
								:size="32"
								class="userAvatar"
								:src="item.avatar"
							/>
						</div>
					</el-col>
					<el-col :span="14">
						<el-row>
							<el-col :span="24">
								<div class="chat-name-other">
									{{ item.displayName }}
								</div>
							</el-col>
						</el-row>
						<div class="bubble-other">
							<!-- 文本消息展示 -->
							<div
								v-if="item.content_type == 'text'"
								class="chat-font"
							>
								{{ item.message }}
							</div>

							<!-- 文件消息展示 -->
							<div
								@click="sys.openFile(item.message)"
								v-else-if="item.content_type == 'file'"
								:class="['chat-item-file', 'file-other']"
							>
								<div class="chat-item-file-icon">
									<el-icon
										size="30"
										color="#303133"
										><Document
									/></el-icon>
								</div>
								<div class="chat-item-file-name">
									{{ item.file_name }}
								</div>
							</div>

							<!-- 图片消息展示 -->
							<div
								v-else-if="item.content_type == 'image'"
								class="chat-item-image"
							>
								<el-image
									fit="cover"
									loading="lazy"
									:preview-src-list="[item.message]"
									:src="item.message"
								/>
							</div>
						</div>
						<!-- 发送时间展示，在消息框外部的下方 -->
						<div class="chat-time-other-outer">
							{{ item.previewTimeFormat }}
						</div>
					</el-col>
					<el-col :span="8" />
				</el-row>
			</div>
		</div>
	</div>
	<!-- 悬浮菜单 -->
	<div
		class="context-menu"
		v-if="store.contextMenu.visible"
		:style="{
			top: `${store.contextMenu.y}px`,
			left: `${store.contextMenu.x}px`,
		}"
	>
		<div
			v-for="contextItem in store.contextMenu.list"
			:key="contextItem.id"
			class="context-menu-item"
		>
			<div
				class="context-menu-item-font"
				@click="store.handleContextMenu()"
			>
				{{ contextItem.label }}
			</div>
		</div>
	</div>
</template>

<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
import { System } from "@/system";
	const store = useChatStore();
	const sys: any = inject<System>("system");
</script>

<style scoped>
	.system-message {
		background-color: #e8e8e8; /* 设置背景颜色 */
		color: #333; /* 设置文字颜色 */
		font-size: 12px; /* 设置文字大小 */
		padding: 2px 20px; /* 设置内边距 */
		font-family: Arial, sans-serif;
		border-radius: 10px; /* 设置边角为圆角 */
		text-align: center; /* 文本居中显示 */
		margin: 10px auto; /* 上下外边距为10px，左右自动（居中） */
		max-width: 60%; /* 最大宽度为80% */
	}

	.bubble-me,
	.bubble-other {
		display: flex;
		flex-direction: column;
		background-color: #ffffff;
		float: left;
		margin-left: 5px;
		margin-bottom: 10px;
	}

	.bubble-me {
		background-color: #d6e4f6;
		float: right;
		border-radius: 12px 0px 12px 12px;
		margin-right: 5px;
		margin-left: 0;
	}

	.bubble-other {
		border-radius: 0 12px 12px 12px;
		background-color: #e8eaed;
	}

	.chat-name-me,
	.chat-name-other {
		font-size: 14px;
		font-family: Arial, sans-serif;
		color: #b2b2b2;
		margin-bottom: 2px;
		margin-left: 2px;
	}

	.chat-name-me {
		text-align: right;
		margin-right: 5px;
	}

	.chat-name-other {
		text-align: left;
		margin-left: 5px;
	}

	.chat-font {
		font-size: 15px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		margin: 10px;
		word-break: break-all;
	}

	.chat-item-file {
		border-radius: 4px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin: 10px;
	}

	.file-me {
		flex-direction: row;
	}

	.file-other {
		flex-direction: row-reverse;
	}

	.chat-item-file-icon {
		width: 40px;
	}

	.chat-item-file-name {
		font-size: 14px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		color: #409eff;
	}

	.chat-item-image {
		margin: 10px;
	}

	.chat-time-me-outer,
	.chat-time-other-outer {
		font-size: 12px;
		color: #999999;
		margin-top: 25px;
		text-align: center;
	}

	.chat-time-me-outer {
		display: flex;
		padding-right: 5px;
		justify-content: end;
	}

	.chat-time-other-outer {
		display: flex;
		padding-left: 5px;
		justify-content: start;
	}
</style>
