<template>
	<div>
		<el-row style="margin-bottom: 20px;" justify="end">
			<el-button
				type="primary"
				icon="Plus"
				circle
				@click="addProxy"
			/>
		</el-row>
		<el-table
			:data="proxies"
			style="width: 98%; border: none"
		>
			<el-table-column
				prop="proxyType"
				label="代理类型"
			/>
			<el-table-column
				prop="port"
				label="本地端口"
			/>
			<!-- <el-table-column prop="domain" label="代理域名" /> -->

			<el-table-column label="状态">
				<template #default="scope">
					<!-- <el-switch v-model="scope.row.status" active-color="#ff4949" inactive-color="#13ce66" @change="SetStatus(scope.row.id)"></el-switch> -->
					<el-button
						size="small"
						@click="SetStatus(scope.row.id)"
						>{{ scope.row.status ? "启用" : "禁用" }}</el-button
					>
				</template>
			</el-table-column>
			<el-table-column label="操作">
				<template #default="scope">
					<el-button
						size="small"
						circle
						icon="Edit"
						@click="editProxy(scope.row)"
					></el-button>
					<el-button
						size="small"
						circle
						icon="Delete"
						@click="DeleteProxy(scope.row.id)"
					></el-button>
				</template>
			</el-table-column>
		</el-table>
		<el-pagination
			v-if="page.total > page.size"
			layout="prev, pager, next"
			:total="page.total"
			:page-size="page.size"
			v-model:current-page="page.current"
			@current-change="changePage"
		/>
		<el-dialog
			v-model="proxyDialogShow"
			:title="isEditing ? '编辑代理' : '添加代理'"
			width="400px"
		>
			<span>
				<el-form
					:model="proxyData"
					:rules="proxyRules"
					ref="pwdRef"
				>
					<el-form-item
						label="代理类型"
						prop="type"
					>
						<el-select
							v-model="proxyData.proxyType"
							placeholder="代理类型"
						>
							<el-option
								v-for="type in types"
								:key="type.value"
								:label="type.label"
								:value="type.value"
							/>
						</el-select>
					</el-form-item>
					<el-form-item
						label="本地端口"
						prop="port"
					>
						<el-input v-model="proxyData.port" />
					</el-form-item>
					<el-form-item
						label="状态"
						prop="status"
					>
						<el-switch
							v-model="proxyData.status"
							active-color="#13ce66"
							inactive-color="#ff4949"
							active-text="启用"
							inactive-text="禁用"
						/>
					</el-form-item>
					<div v-if="proxyData.proxyType === 'http'">
						<el-form-item
							label="代理域名"
							prop="domain"
						>
							<el-input v-model="proxyData.domain" />
						</el-form-item>
						<!-- <el-form-item
							label="代理端口"
							prop="listenPort"
						>
							<el-input v-model="proxyData.listenPort" />
						</el-form-item> -->
					</div>
					<el-form-item
						label="文件路径"
						prop="path"
						v-if="proxyData.proxyType === 'file'"
					>
						<el-input
							v-model="proxyData.path"
							@click="selectFile()"
						/>
					</el-form-item>
					<el-form-item
						label="转发IP+端口"
						prop="domain"
						v-if="proxyData.proxyType === 'udp'"
					>
						<el-input v-model="proxyData.domain" />
					</el-form-item>
					<el-form-item>
						<el-button
							type="primary"
							@click="saveProxy"
							style="margin: 0 auto"
						>
							确认
						</el-button>
					</el-form-item>
				</el-form>
			</span>
		</el-dialog>
	</div>
</template>

<script setup>
import { ref, reactive } from 'vue';

const proxies = ref([]); // 模拟代理数据
const page = reactive({
	total: 0,
	size: 10,
	current: 1
});
const proxyDialogShow = ref(false);
const isEditing = ref(false);
const proxyData = reactive({
	proxyType: '',
	port: '',
	status: false,
	domain: '',
	path: ''
});
const types = [
	{ label: 'HTTP', value: 'http' },
	{ label: 'File', value: 'file' },
	{ label: 'UDP', value: 'udp' }
];
const proxyRules = {}; // 模拟验证规则

function addProxy() {
	// 模拟添加代理的方法
	proxyDialogShow.value = true;
	isEditing.value = false;
	Object.assign(proxyData, {
		proxyType: '',
		port: '',
		status: false,
		domain: '',
		path: ''
	});
}

function editProxy(proxy) {
	// 模拟编辑代理的方法
	proxyDialogShow.value = true;
	isEditing.value = true;
	Object.assign(proxyData, proxy);
}

function DeleteProxy(id) {
	// 模拟删除代理的方法
	proxies.value = proxies.value.filter(proxy => proxy.id !== id);
}

function SetStatus(id) {
	// 模拟设置状态的方法
	const proxy = proxies.value.find(proxy => proxy.id === id);
	if (proxy) {
		proxy.status = !proxy.status;
	}
}

function saveProxy() {
	// 模拟保存代理的方法
	if (isEditing.value) {
		const index = proxies.value.findIndex(proxy => proxy.id === proxyData.id);
		if (index !== -1) {
			proxies.value.splice(index, 1, { ...proxyData });
		}
	} else {
		proxies.value.push({ ...proxyData, id: Date.now() });
	}
	proxyDialogShow.value = false;
}

function changePage(pageNumber) {
	page.current = pageNumber;
}

function selectFile() {

	console.log('选择文件');
}
</script>

<style scoped>
</style>
