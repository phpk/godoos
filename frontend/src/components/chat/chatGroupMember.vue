<template>
	<div class="group-member-container">
		<!-- 群成员列表 -->
		<div
			class="member-list"
			v-for="item in groupMembers"
			:key="item.userId"
		>
			<!-- 头像和昵称 -->
			<div class="member-item">
				<div class="avatar-container">
					<el-avatar
						style="width: 100%; height: 100%;"
						:src="item.avatar"
					/>
				</div>
				<span>{{ item.nickname }}</span>
			</div>
		</div>
	</div>
</template>
<script setup lang="ts">
	import { useChatStore } from "@/stores/chat";
	const store = useChatStore();

	// 使用模拟数据

	// 使用模拟数据，定义40条
	const groupMembers = ref(
		Array.from({ length: 40 }, (v, i) => ({
			userId: `${i + 1}`,
			avatar: `path/to/avatar${i + 1}.jpg`,
			nickname: `用户1111111111111111${i + 1}`,
		}))
	);
</script>
<style scoped>
	.group-member-container {
		margin: 10px 0;
		width: 100%;
		height: 100%;
	}

	.member-item {
		display: flex;
		align-items: center;
		margin-bottom: 5px;
	}

	.avatar-container {
		width: 25px;
		height: 25px;
		margin-right: 5px;
	}

	.member-item span {
		width: 90px;
		margin-left: 10px;
    font-size: 12px;
    color: #86909C;
		overflow: hidden; /* 防止内容溢出 */
		white-space: nowrap; /* 保持文本在一行 */
		text-overflow: ellipsis; /* 文本溢出显示省略号 */
		max-width: 150px; /* 最大宽度，根据需要调整 */
	}
</style>
