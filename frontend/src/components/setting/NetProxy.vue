<script setup lang="ts">
import { useProxyStore } from "@/stores/proxy";
import { notifyError, notifySuccess } from "@/util/msg";
import { Plus, Setting, VideoPlay, VideoPause, Money } from "@element-plus/icons-vue";
import { onMounted } from "vue";
const proxyStore = useProxyStore();
const {
	fetchProxies,
	fetchProxy,
	deleteProxyById,
} = proxyStore;


const editProxyBefore = async (id: number) => {
	if (!id) {
		notifyError("请选择要编辑的代理");
		return
	}
	//console.log(id)
	await fetchProxy(id);
	proxyStore.addShow = true;
	proxyStore.isEditor = true;
};

const deleteProxy = async (id: number) => {
	if (!id) {
		notifyError("请选择要删除的代理");
		return
	}
	const res = await deleteProxyById(id);
	if (res) {
		notifySuccess("删除成功");
	} else {
		notifyError("删除失败");
	}
};
const addProxy = () => {
	proxyStore.addShow = true;
	proxyStore.resetProxyData();
	proxyStore.isEditor = false;
};
onMounted(async () => {
	await fetchProxies();
});
</script>

<template>
	<div>
		<el-row justify="end">
			<el-button :type="proxyStore.status ? 'primary' : 'success'" :icon="VideoPause" circle @click="proxyStore.stopFrpc"  />
			<el-button :type="proxyStore.status ? 'success' : 'primary'" :icon="VideoPlay" circle @click="proxyStore.startFrpc" />
			
			<el-button type="primary" :icon="Setting" circle @click="proxyStore.settingShow = true" />
			<el-button type="primary" :icon="Plus" circle @click="addProxy" />
		</el-row>
		<el-table :data="proxyStore.proxies" style="width: 98%; border: none">
			<el-table-column prop="name" label="名称" width="100" />
			<el-table-column prop="type" label="类型" width="80" />
			<el-table-column prop="localPort" label="本地端口" width="80" />
			<el-table-column prop="localIp" label="本地Ip" />
			<el-table-column label="操作">
				<template #default="scope">
					<el-row :gutter="24" justify="start">
						<el-col :span="10">
							<el-button size="small" @click="editProxyBefore(scope.row.id)">编辑</el-button>
						</el-col>
						<el-col :span="10">
							<el-button size="small" type="danger" @click="deleteProxy(scope.row.id)">删除</el-button>
						</el-col>
					</el-row>
				</template>
			</el-table-column>
		</el-table>
		<el-pagination v-if="proxyStore.page.total > proxyStore.page.size" layout="prev, pager, next"
			:page-size="proxyStore.page.size" v-model:current-page="proxyStore.page.current"
			:total="proxyStore.page.total" @current-change="proxyStore.pageChange" />
		<el-drawer v-model="proxyStore.settingShow" title="编辑配置" direction="rtl" size="80%">
			<FrpcConfig />
		</el-drawer>
		<!-- 代理配置抽屉 -->
		<el-drawer v-model="proxyStore.addShow" :title="proxyStore.isEditor ? '编辑代理' : '添加代理'" direction="rtl"
			size="80%" @close="proxyStore.resetProxyData">
			<FrpcEdit />
		</el-drawer>
	</div>
</template>

<style scoped>
:deep(.el-drawer) {
	width: 60% !important;
}
</style>
