<script setup>
	import { useChatStore } from "@/stores/chat";
	import {
		ElAvatar,
		ElCol,
		ElCollapse,
		ElCollapseItem,
		ElRow,
	} from "element-plus";
	const store = useChatStore();
</script>

<template>
	<el-collapse v-model="store.activeNames">
		<el-collapse-item name="1">
			<template #title>
				<span
					v-if="store.userList.length > 0"
					class="title"
					>同事（{{ store.userList.length }}）</span
				>
				<span
					v-else
					class="title"
					>同事</span
				>
			</template>
			<div v-if="store.userList.length > 0">
				<div
					v-for="item in store.userList"
					:key="item.id"
				>
					<div
						class="list-item"
						@click="store.changeChatList(item.id)"
						:style="{
							backgroundColor:
								item.id === store.targetUserId ? '#bae7ff' : '',
						}"
					>
						<el-row>
							<el-col :span="6">
								<el-avatar
									shape="square"
									:size="40"
									class="avatar"
									:src="item.avatar"
								/>
							</el-col>
							<el-col :span="18">
								<el-row>
									<el-col :span="18">
										<div class="previewName">
											{{ item.nickname }}
											<el-icon
												class="online-icon"
												v-if="item.isOnline"
											>
												<CircleCheckFilled />
											</el-icon>
										</div>
									</el-col>
									<el-col :span="6">
										<!-- 空白占位 -->
									</el-col>
								</el-row>
								<el-row>
									<div class="previewIP">
										{{ item.ip }}
									</div>
								</el-row>
							</el-col>
						</el-row>
					</div>
				</div>
			</div>
			<div v-else>
				<p class="no-data">暂无数据</p>
			</div>
		</el-collapse-item>
		<el-collapse-item name="2">
			<template #title>
				<span
					v-if="store.groupList.length > 0"
					class="title"
					>部门（{{ store.groupList.length }}）</span
				>
				<span
					v-else
					class="title"
					>部门</span
				>
			</template>
			<div v-if="store.groupList.length > 0">
				<div
					v-for="group in store.groupList"
					:key="group.id"
				>
					<div
						class="list-item"
						@click="store.changeGroupList(group)"
						:style="{
							backgroundColor:
								group.id === store.targetGroupId
									? '#C4C4C4'
									: '',
						}"
					>
						<el-row>
							<el-col :span="6">
								<el-avatar
									shape="square"
									:size="40"
									class="avatar"
									:src="group.avatar"
								/>
							</el-col>
							<el-col :span="18">
								<el-row>
									<el-col :span="18">
										<div class="previewName">
											{{ group.name }}
										</div>
									</el-col>
									<el-col :span="6">
										<div class="previewTime">
											{{ group.previewTimeFormat }}
										</div>
									</el-col>
								</el-row>
								<el-row>
									<div
										v-if="group.previewType === 0"
										class="previewChat"
									>
										{{ group.previewMessage }}
									</div>
									<div
										v-if="group.previewType === 1"
										class="previewChat"
									>
										{{
											group.userId === id
												? "你"
												: '"' + group.name + '"'
										}}撤回了一条消息
									</div>
								</el-row>
							</el-col>
						</el-row>
					</div>
				</div>
			</div>
			<div v-else>
				<p class="no-data">暂无数据</p>
			</div>
		</el-collapse-item>
	</el-collapse>
</template>

<style scoped>
	.title {
		padding-left: 10px;
	}
	.list-item {
		width: 94%;
		height: 60px;
		margin: 0 auto;
		border-radius: 4px;
		transition: all 0.5s;
		margin-bottom: 5px;
	}

	.list-item:hover {
		background-color: #bae7ff;
	}

	.avatar {
		margin: 10px;
	}

	.previewName {
		margin-left: 10px;
		font-size: 14px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		overflow: hidden; /* 隐藏超出部分 */
		text-overflow: ellipsis; /* 显示为省略号 */
		white-space: nowrap; /* 不换行 */
		max-width: 100%; /* 最大宽度 */
		display: flex; /* 使用 Flexbox */
		align-items: center; /* 垂直居中 */
	}

	/* 为了使内部的 el-row 垂直居中，我们也需要设置其父级元素 */
	.el-row {
		display: flex;
		align-items: center;
	}

	.previewIP {
		margin-left: 10px;
		font-size: 12px;
		font-family: Arial, sans-serif;
		color: #999999;
		display: flex; /* 使用 Flexbox */
		align-items: center; /* 垂直居中 */
	}

	.online-icon {
		font-size: 16px; /* 调整图标大小 */
		color: green; /* 在线状态颜色 */
		margin-left: 5px; /* 与用户名之间的间距 */
	}

	.previewChat {
		margin-left: 10px;
		font-size: 12px; /* 调整字体大小 */
		font-family: Arial, sans-serif;
		color: #999999;
		overflow: hidden; /* 隐藏超出部分 */
		text-overflow: ellipsis; /* 显示为省略号 */
		white-space: nowrap; /* 不换行 */
		max-width: 100%; /* 最大宽度 */
	}

	.previewTime {
		float: right;
		margin-top: 10px;
		margin-right: 10px;
		font-size: 12px;
		font-family: Arial, sans-serif;
		color: #999999;
	}

	.no-data {
		text-align: center;
		color: #999999;
	}
</style>
