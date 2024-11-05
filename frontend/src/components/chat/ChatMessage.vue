<template>
	<div
		v-if="store.chatHistory && store.chatHistory.length > 0"
		v-for="item in store.chatHistory"
		:key="item.chatId"
	>
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
							v-else-if="item.content_type == 'file'"
							:class="['chat-item-file', 'file-me']"
						>
							<!-- 文件图标 -->
							<div class="chat-item-file-icon">
								<el-icon
									size="30"
									color="#303133"
									><Document
								/></el-icon>
							</div>
							<!-- 文件名 -->
							<div class="chat-item-file-name">
								{{ item.message }}
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
								:src="item.message"
							/>
						</div>
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
							v-else-if="item.content_type == 'file'"
							:class="['chat-item-file', 'file-other']"
						>
							<!-- 文件图标 -->
							<div class="chat-item-file-icon">
								<el-icon
									size="30"
									color="#303133"
									><Document
								/></el-icon>
							</div>
							<!-- 文件名 -->
							<div class="chat-item-file-name">
								{{ item.message }}
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
								:src="item.message"
							/>
						</div>
					</div>
				</el-col>
				<el-col :span="8" />
			</el-row>
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
	const store = useChatStore();
</script>

<style scoped>
	.bubble-me,
	.bubble-other {
		display: flex;
		align-items: center;
		background-color: #ffffff;
		float: left;
		border-radius: 4px;
		margin-left: 5px;
		margin-bottom: 20px;
	}

	.bubble-me {
		background-color: #95ec69;
		float: right;
		margin-right: 5px;
		margin-left: 0;
	}

	.chat-name-me,
	.chat-name-other {
		font-size: 14px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		color: #b2b2b2;
		margin-bottom: 2px;
		margin-left: 2px;
	}

	.chat-name-me {
		text-align: right; /* 右对齐昵称 */
		margin-right: 2px;
	}

	.chat-font {
		font-size: 15px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		margin: 10px;
	}

	.chat-item-file {
		border-radius: 4px;
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin: 10px;
	}

	.file-me {
		flex-direction: row; /* 图标在左边 */
	}

	.file-other {
		flex-direction: row-reverse; /* 图标在右边 */
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
		width: 180px;
		height: 180px;
		margin: 10px;
	}
</style>
