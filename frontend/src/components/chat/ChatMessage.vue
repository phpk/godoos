<template>
	<div
		v-if="store.chatHistory && store.chatHistory.length > 0"
		v-for="item in store.chatHistory"
		:key="item.time"
	>
		<div
			v-if="item.userId == store.userInfo.id"
			class="chat-item"
		>
			<el-row>
				<el-col :span="8" />
				<el-col :span="14">
					<el-row>
						<el-col :span="24">
							<div class="chat-name-me">
								{{ item.userInfo.nickname }}
							</div>
						</el-col>
					</el-row>
					<div
						class="bubble-me"
						@contextmenu.prevent="
							store.showContextMenu($event, item.userId)
						"
					>
						<div class="chat-font">
							{{ item.message }}
						</div>
					</div>
					<!-- 时间显示在消息框外 -->
					<div class="chat-time">{{ formatTime(item.time) }}</div>
				</el-col>
				<el-col :span="2">
					<div class="chat-avatar">
						<el-avatar
							shape="square"
							style="margin: 0; float: left"
							:size="32"
							class="userAvatar"
							:src="item.userInfo.avatar"
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
							:src="item.userInfo.avatar"
						/>
					</div>
				</el-col>
				<el-col :span="14">
					<el-row>
						<el-col :span="24">
							<div class="chat-name-other">
								{{ item.userInfo.nickname }}
							</div>
						</el-col>
					</el-row>
					<div class="bubble-other">
						<div class="chat-font">
							{{ item.message }}
						</div>
					</div>
					<!-- 时间显示在消息框外 -->
					<div class="chat-time">{{ formatTime(item.time) }}</div>
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

	const formatTime = (timestamp: number) => {
		const date = new Date(timestamp);
		const options: Intl.DateTimeFormatOptions = {
			hour: "numeric",
			minute: "numeric",
			hour12: false,
		};
		return date.toLocaleString("default", options);
	};
</script>

<style scoped>
	.bubble-me {
		background-color: #95ec69;
		float: right;
		border-radius: 4px;
		margin-right: 5px;
		margin-top: 5px;
		padding: 5px;
	}

	.bubble-other {
		background-color: #ffffff;
		float: left;
		border-radius: 4px;
		margin-left: 5px;
		margin-top: 5px;
		padding: 5px;
	}

	.chat-name-me,
	.chat-name-other {
		font-size: 14px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		color: #b2b2b2;
	}

	.chat-font {
		margin: 8px;
		font-size: 15px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
	}

	.chat-time {
		font-size: 12px;
		color: #999;
		height: 50px;
		display: flex;
		align-items: self-end;
		justify-content: start;
		padding-left: 10px;
		text-align: left;
	}
</style>
