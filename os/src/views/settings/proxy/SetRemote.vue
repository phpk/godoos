<template>
	<div>
		<el-row
			style="margin-bottom: 20px"
			justify="end"
		>
			<el-button
				type="primary"
				icon="VideoPause"
				circle
			/>
			<el-button
				type="primary"
				icon="VideoPlay"
				circle
			/>
			<el-button
				type="primary"
				icon="Setting"
				circle
			/>
			<el-button
				@click="addProxy"
				type="primary"
				icon="Plus"
				circle
			/>
		</el-row>
		<el-table
			:data="proxyStore.proxies"
			style="width: 100%; border: none"
		>
			<el-table-column
				prop="name"
				label="名称"
			/>
			<el-table-column
				prop="type"
				label="类型"
			/>
			<el-table-column
				prop="localPort"
				label="本地端口"
			/>
			<el-table-column
				prop="localIp"
				label="本地Ip"
			/>
			<el-table-column
				fixed="right"
				label="操作"
			>
				<template #default="{ row, $index }">
					<el-row
						:gutter="24"
						justify="start"
					>
						<el-col :span="10">
							<el-button
								size="small"
								@click="editProxy(row)"
								>编辑</el-button
							>
						</el-col>
						<el-col :span="10">
							<el-button
								size="small"
								type="danger"
								@click="deleteProxy($index)"
								>删除</el-button
							>
						</el-col>
					</el-row>
				</template>
			</el-table-column>
		</el-table>
		<el-pagination layout="prev, pager, next" />
		<el-drawer
			:title="proxyStore.isEditor ? '编辑代理' : '添加代理'"
			direction="rtl"
			size="80%"
			v-model="proxyStore.addShow"
		>
			<FrpcAdd />
		</el-drawer>
	</div>
</template>

<script setup lang="ts">
	import { useProxyStore } from "@/stores/proxy";
	import FrpcAdd from "./FrpcAdd.vue";
	import { successMsg } from "@/utils/msg";
	const proxyStore = useProxyStore();

	const addProxy = () => {
		proxyStore.addShow = true;
		proxyStore.resetProxyData();
		proxyStore.isEditor = false;
	};

	const editProxy = (proxy: any) => {
		proxyStore.proxyData = { ...proxy };
		proxyStore.isEditor = true;
		proxyStore.addShow = true;
	};

	const deleteProxy = (index: any) => {
		proxyStore.proxies.splice(index, 1);
		proxyStore.addShow = false;
		successMsg("删除成功");
	};
</script>

<style scoped>
	:deep(.el-drawer) {
		width: 60% !important;
	}
</style>
