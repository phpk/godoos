<template>
	<el-drawer
		v-model="store.groupMemberDrawerVisible"
		direction="rtl"
		title="群成员"
		:with-header="true"
		:show-close="true"
		size="250px"
	>
		<!-- 群成员列表 -->
		<div class="group-member-container">
			<div
				class="member-list"
				v-for="item in store.groupMembers"
				:key="item.id"
			>
				<div class="member-item">
					<!-- 头像和昵称 -->
					<div class="avatar-container">
						<el-avatar
							style="width: 100%; height: 100%"
							:src="item.avatar"
						/>
					</div>
					<span>{{ item.nickname }}</span>
				</div>
			</div>
		</div>
	</el-drawer>
</template>
<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
	const store = useChatStore();

	// 监听群成员抽屉的可见性，当其变为true时获取群成员列表
	watch(
		() => store.groupMemberDrawerVisible,
		(newVal) => {
			if (newVal) {
				store.getGroupMemberList(store.targetGroupInfo.group_id);
			}
		}
	);

	// 监听抽屉状态变化，当打开时获取群成员列表
	watch(
		() => store.drawerVisible,
		(newVal) => {
			if (newVal) {
				store.getGroupMemberList(store.targetGroupInfo.group_id);
			}
		}
	);
</script>
<style scoped>
	:deep(.el-drawer__header) {
		background-color: red;
		padding: 0px 20px;
		height: 50px;
		color: #000000;
		margin-bottom: 0px;
	}
	:deep(.el-drawer__title) {
		font-size: 20px;
	}

	.group-member-container {
		width: 100%;
		height: 100%;
		overflow-y: auto; /* 添加滚动条 */
	}

	.member-item {
		display: flex;
		align-items: center;
		margin-bottom: 5px;
	}

	.avatar-container {
		width: 30px;
		height: 30px;
		margin-right: 10px;
	}

	.member-item span {
		width: 150px;
		margin-left: 10px;
		font-size: 12px;
		color: #86909c;
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
	}
</style>
