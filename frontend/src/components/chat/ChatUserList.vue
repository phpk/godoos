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
	const handleNodeClick = (data) => {
		if (data.isUser) {
			store.getSessionInfo(data.id, "user");
		}
	};

	const defaultProps = {
		children: "children",
		label: "label",
	};

	// 将数据转换成树形结构
	function transformData(data) {
		return data.map((dept) => ({
			label: dept.dept_name,
			children: transformUsers(dept.users).concat(
				transformSubDepts(dept.sub_depts)
			),
		}));
	}

	function transformUsers(users) {
		if (!users) return [];
		console.log(users);
		return users.map((user) => ({
			label: `${user.nick_name}`,
			id: user.user_id, // 添加用户ID
			isUser: true,
			children: [],
		}));
	}

	function transformSubDepts(sub_depts) {
		if (!sub_depts) return [];
		return transformData(sub_depts);
	}

	const data = transformData(store.departmentList);
	console.log(data);
</script>

<template>
	<el-collapse v-model="store.activeNames">
		<el-collapse-item name="1">
			<template #title>
				<span
					v-if="store.onlineUserList.length > 0"
					class="title"
					>在线（{{ store.onlineUserList.length }}）</span
				>
				<span
					v-else
					class="title"
					>在线</span
				>
			</template>
			<div v-if="store.onlineUserList.length > 0">
				<div
					v-for="item in store.onlineUserList"
					:key="item.id"
				>
					<div
						class="list-item"
						@click="store.getSessionInfo(item.chatId, 'user')"
						:style="{
							backgroundColor:
								item.id === store.targetChatId ? '#F5F5F5' : '',
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
								<el-row
									style="display: flex; align-items: center"
								>
									<el-col>
										<div class="previewName">
											{{ item.nickname }}
										</div>
									</el-col>
									<el-col :span="6">
										<!-- 空白占位 -->
									</el-col>
								</el-row>
								<el-row>
									<div class="previewIP">
										{{ item.login_ip }}
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
					v-if="store.departmentList.length > 0"
					class="title"
				>
					部门（{{ store.departmentList.length }}）
				</span>
				<span
					v-else
					class="title"
					>部门</span
				>
			</template>
			<div
				class="tree-container"
				v-if="data.length > 0"
			>
				<el-tree
					:data="data"
					node-key="dept_id"
					:props="{ label: 'label', children: 'children' }"
					@node-click="handleNodeClick"
					:default-expand-all="false"
				>
				</el-tree>
			</div>

			<div v-else>
				<p class="no-data">暂无数据</p>
			</div>
		</el-collapse-item>
		<el-collapse-item name="3">
			<template #title>
				<span
					v-if="store.groupList.length > 0"
					class="title"
					>群组（{{ store.groups.length }}）</span
				>
				<span
					v-else
					class="title"
					>群组</span
				>
			</template>
			<div v-if="store.groupList.length > 0">
				<div
					v-for="group in store.groups"
					:key="group.id"
				>
					<div
						class="list-item"
						@click="store.getSessionInfo(group.id, 'group')"
						:style="{
							backgroundColor:
								group.id === store.targetChatId
									? '#F5F5F5'
									: '',
						}"
					>
						<el-row>
							<el-col :span="6">
								<el-avatar
									shape="square"
									:size="40"
									class="avatar"
									>群</el-avatar
								>
							</el-col>
							<el-col :span="18">
								<div class="previewName">{{ group.name }}</div>
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
		width: 95%;
		height: 60px;
		margin: 0 auto;
		border-radius: 4px;
		transition: all 0.5s;
		margin-bottom: 5px;
		background-color: #fff;
	}

	.list-item:hover {
		background-color: #e8f3ff;
	}

	.avatar {
		margin: 10px;
	}

	.previewName {
		font-weight: 400;
		margin-left: 10px;
		font-size: 14px;
		font-family: Arial, sans-serif;
		line-height: 1.5;
		color: #000000;
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

	.no-data {
		text-align: center;
		color: #999999;
	}
</style>
