<template>

	<div class="table-container">
		<!-- tab -->
		<el-tabs type="border-card" v-model="activeName">
			<!-- 用户 -->
			<el-tab-pane label="编辑用户" name="first">
				<div class="users-container">
					<div v-for="(user, index) in editUsers.users" :key="index" class="user-tag">
						<el-avatar :src="user.avatar ? user.avatar : ''" :style="user.avatar
								? {}
								: {
									backgroundColor: getRandomColor(
										user.nickname[0]
									),
								}
							">
							{{ user.avatar ? "" : user.nickname[0] }}
						</el-avatar>
						<div class="user-info">
							<div class="username">{{ user.nickname }}</div>
						</div>
					</div>
				</div>
				<!-- 分页 -->
				<!-- <div style="display: flex; justify-content: center">
            <el-pagination
              background
              layout="prev, pager, next"
              :total="props.editUsers?.total"
            />
          </div> -->
			</el-tab-pane>
			<!-- 历史 -->
			<el-tab-pane label="历史" name="second">
				<el-table :data="editHistory.list">
					<el-table-column property="created_at" label="时间" width="150">
						<template #default="scope">
							{{ formatTimestamp(scope.row.created_at) }}
						</template>
					</el-table-column>
					<el-table-column property="nickname" label="操作人" width="150"></el-table-column>
					<el-table-column label="操作">
						<template #default="scope">
							<el-button @click="handleView(scope.row)" type="text" size="small">查看</el-button>
							<el-button @click="handleRecover(scope.row)" type="text" size="small">还原</el-button>
						</template>
					</el-table-column>
				</el-table>
			</el-tab-pane>
		</el-tabs>
	</div>

</template>

<script setup lang="ts">
import {
	getEditHistory,
	getShareUserList,
	restoreEditData,
} from "@/api/share";
import { useFileSystemStore } from "@/stores/filesystem";
import { errMsg, successMsg } from "@/utils/msg";
import {  onMounted, ref } from "vue";

const store = useFileSystemStore();
const props = defineProps({
	path: {
		type: String,
		default: "",
	},
	truePath: {
		type: String,
		default: "",
	},
});
// 定义响应式变量
const editUsers: any = ref([]);
const editHistory: any = ref([]);

const activeName = ref("first");

const handleView = async (row: any) => {
	const regex = /(\.log)(.*)/;
	const filePach = row.file_path.match(regex)[0];
	store.openFile(filePach);
};

const handleRecover = async (row: any) => {
	const res = await restoreEditData(row.id);
	if (!res.success) {
		errMsg("还原失败");
		return
	}
	store.currentShareFile = { path: props.truePath };
	// 关闭协同对话框
	successMsg("还原成功")
};

// 加载用户数据
const loadEditUsers = async () => {
	if (!props.truePath) return;
	const res = await getShareUserList(props.truePath);
	editUsers.value = res;
};

// 加载历史记录数据
const loadEditHistoryList = async () => {
	if (!props.truePath) return;
	const res = await getEditHistory(props.truePath, "", "");
	editHistory.value = res;
};

onMounted(() => {
	loadEditUsers();
	loadEditHistoryList();
});

const formatTimestamp = (timestamp: number): string => {
	const date = new Date(timestamp * 1000);
	return date.toLocaleString();
};

const getRandomColor = (char: string): string => {
	const colors = [
		"#f44336",
		"#e91e63",
		"#9c27b0",
		"#673ab7",
		"#3f51b5",
		"#2196f3",
		"#03a9f4",
		"#00bcd4",
		"#009688",
		"#4caf50",
		"#8bc34a",
		"#cddc39",
		"#ffeb3b",
		"#ffc107",
		"#ff9800",
		"#ff5722",
	];
	const hash = char.charCodeAt(0);
	return colors[hash % colors.length];
};



</script>

<style scoped>
.table-container {
	max-height: 350px;
	overflow-y: scroll;
}

.table-container::-webkit-scrollbar {
	display: none;
}

.table-container {
	-ms-overflow-style: none;
	scrollbar-width: none;
}

.users-container {
	display: flex;
	gap: 20px;
	height: 200px;
	flex-wrap: wrap;
}

.user-tag {
	display: flex;
	flex-direction: column;
	align-items: center;
	gap: 5px;
}

.user-info {
	display: flex;
	flex-direction: column;
	align-items: center;
}

.username {
	font-weight: bold;
	color: #333;
	text-align: center;
}
</style>
