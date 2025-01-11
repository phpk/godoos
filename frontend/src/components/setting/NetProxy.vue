<script setup lang="ts">
	import { useProxyStore } from "@/stores/proxy";
	import { notifyError, notifySuccess } from "@/util/msg";
	import { Plus } from "@element-plus/icons-vue";
	import { ref } from "vue";

	const proxyStore = useProxyStore();

	const proxyDialogShow = ref(false);
	const isEditing = ref(false);
	const pwdRef = ref<any>(null);
	const internalPortDialogVisible = ref(false);

	// 定义表单验证规则
	const proxyRules = {
		name: [{ required: true, message: "请输入代理名称", trigger: "blur" }],
		port: [
			{ required: true, message: "请输入端口", trigger: "blur" },
			{ type: "number", message: "端口必须是数字", trigger: "blur" },
		],
		domain: [{ required: true, message: "请输入域名", trigger: "blur" }],
		// 其他字段的验证规则...
	};

	const addProxy = () => {
		if (pwdRef.value.validate()) {
			// const isNameDuplicate = proxyStore.proxies.some(
			// 	(p) => p.name === proxyStore.proxyData.name
			// );
			// if (isNameDuplicate) {
			// 	notifyError("代理名称已存在");
			// 	return;
			// }
			// proxyStore.addProxy({ ...proxyStore.proxyData });
			// proxyDialogShow.value = false;
			// proxyStore.resetProxyData();
			// 调用 store 中的 createFrpcConfig 方法
			proxyStore
				.createFrpcConfig()
				.then(() => {
					proxyDialogShow.value = false;
					proxyStore.resetProxyData();
					notifySuccess("代理配置已成功创建");
				})
				.catch((error) => {
					notifyError(`创建代理配置失败: ${error.message}`);
				});
		}
	};

	const updateProxy = () => {
		if (pwdRef.value.validate()) {
			const index = proxyStore.proxies.findIndex(
				(p) => p.id === proxyStore.proxyData.id
			);
			if (index !== -1) {
				const isNameDuplicate = proxyStore.proxies.some(
					(p, i) =>
						p.name === proxyStore.proxyData.name && i !== index
				);
				if (isNameDuplicate) {
					notifyError("代理名称已存在");
					return;
				}
				proxyStore.updateProxy({ ...proxyStore.proxyData });
				proxyDialogShow.value = false;
				proxyStore.resetProxyData();
				isEditing.value = false;
			}
		}
	};

	const saveProxy = () => {
		pwdRef.value.validate((valid: boolean) => {
			if (valid) {
				if (isEditing.value) {
					updateProxy();
				} else {
					addProxy();
				}
			} else {
				console.log("表单验证失败");
			}
		});
	};

	const openInternalPortDialog = () => {
		internalPortDialogVisible.value = true;
	};

	const selectPort = (port: number) => {
		console.log(`Selected port: ${port}`);
		internalPortDialogVisible.value = false;
	};

	const proxyTypes = ref([
		"http",
		"https",
		"tcp",
		"udp",
		"stcp",
		"xtcp",
		"sudp",
	]);

	const stcpModels = ref([
		{ label: "访问者", value: "visitors" },
		{ label: "被访问者", value: "visited" },
	]);

	const handleSelectFile = (type: number, ext: string[]) => {
		ipcRenderer.invoke("file.selectFile", ext).then((r) => {
			switch (type) {
				case 1:
					proxyStore.proxyData.https2httpCaFile = r[0];
					break;
				case 2:
					proxyStore.proxyData.https2httpKeyFile = r[0];
					break;
			}
		});
	};
</script>

<template>
	<div>
		<el-row justify="end">
			<el-button
				type="primary"
				:icon="Plus"
				circle
				@click="proxyDialogShow = true"
			/>
		</el-row>
		<el-table
			:data="proxyStore.proxies"
			style="width: 98%; border: none"
		>
			<el-table-column
				prop="name"
				label="名称"
				width="100"
			/>
			<el-table-column
				prop="type"
				label="类型"
				width="80"
			/>
			<el-table-column
				prop="port"
				label="本地端口"
				width="80"
			/>
			<el-table-column
				prop="domain"
				label="代理域名"
				width="180"
			/>
			<el-table-column label="操作">
				<template #default="scope">
					<el-button
						size="small"
						@click="editProxy(scope.row)"
						>编辑</el-button
					>
					<el-button
						size="small"
						type="danger"
						@click="deleteProxy(scope.row.id)"
						>删除</el-button
					>
				</template>
			</el-table-column>
		</el-table>
		<el-pagination
			v-if="proxyStore.totalPages > 1"
			layout="prev, pager, next"
			:total="proxyStore.proxies.length"
			:page-size="proxyStore.pageSize"
			v-model:current-page="proxyStore.currentPage"
			@next-click="proxyStore.nextPage"
			@prev-click="proxyStore.prevPage"
		/>

		<!-- 代理配置抽屉 -->
		<el-drawer
			v-model="proxyDialogShow"
			:title="isEditing ? '编辑代理' : '添加代理'"
			direction="rtl"
			size="60%"
			@close="proxyStore.resetProxyData"
		>
			<el-form
				:model="proxyStore.proxyData"
				:rules="proxyRules"
				ref="pwdRef"
				label-position="top"
			>
				<!-- 代理类型选择 -->
				<el-form-item
					label="代理类型："
					prop="type"
				>
					<el-radio-group v-model="proxyStore.proxyData.type">
						<el-radio-button
							v-for="type in proxyTypes"
							:key="type"
							:value="type"
							>{{ type }}</el-radio-button
						>
					</el-radio-group>
				</el-form-item>

				<!-- HTTP/HTTPS模式 -->
				<template
					v-if="
						proxyStore.proxyData.type === 'http' ||
						proxyStore.proxyData.type === 'https' ||
						proxyStore.proxyData.type === 'tcp' ||
						proxyStore.proxyData.type === 'udp' ||
						proxyStore.proxyData.type === 'stcp' ||
						proxyStore.proxyData.type === 'xtcp' ||
						proxyStore.proxyData.type === 'sudp'
					"
				>
					<el-form-item
						label="代理名称："
						prop="name"
					>
						<el-input
							v-model="proxyStore.proxyData.name"
							placeholder="代理名称"
						/>
					</el-form-item>
					<el-row
						v-if="
							proxyStore.proxyData.type === 'http' ||
							proxyStore.proxyData.type === 'https' ||
							proxyStore.proxyData.type === 'tcp' ||
							proxyStore.proxyData.type === 'udp'
						"
						:gutter="20"
					>
						<el-col :span="12">
							<el-form-item
								label="内网地址："
								prop="serverAddr"
							>
								<el-input
									v-model="proxyStore.proxyData.serverAddr"
									placeholder="内网地址"
								/>
							</el-form-item>
						</el-col>
						<el-col :span="8">
							<el-form-item
								label="端口地址："
								prop="port"
							>
								<el-input-number
									v-model="proxyStore.proxyData.port"
									:min="1"
									:max="65535"
								/>
							</el-form-item>
						</el-col>
					</el-row>
					<el-form-item
						v-if="
							proxyStore.proxyData.type === 'http' ||
							proxyStore.proxyData.type === 'https'
						"
						label="子域名："
						prop="domain"
					>
						<el-input
							v-model="proxyStore.proxyData.domain"
							placeholder="subdomain"
						/>
					</el-form-item>
					<el-form-item
						v-if="
							proxyStore.proxyData.type === 'http' ||
							proxyStore.proxyData.type === 'https'
						"
						label="自定义域名："
						prop="customDomain"
					>
						<el-row
							v-for="(domain, index) in proxyStore.customDomains"
							:key="index"
							:gutter="24"
						>
							<el-col :span="12">
								<el-input
									v-model="proxyStore.customDomains[index]"
									placeholder="example.com"
								/>
							</el-col>
							<el-col :span="5">
								<el-button
									type="primary"
									icon="Plus"
									style="width: 80px"
									@click="proxyStore.addCustomDomain"
									>添加</el-button
								>
							</el-col>
							<el-col :span="5">
								<el-button
									type="primary"
									icon="Plus"
									style="width: 80px"
									@click="
										proxyStore.removeCustomDomain(index)
									"
									>删除</el-button
								>
							</el-col>
						</el-row>
					</el-form-item>
					<el-form-item
						v-if="proxyStore.proxyData.type === 'http'"
						label="HTTP基本认证："
						prop="httpAuth"
					>
						<el-switch v-model="proxyStore.proxyData.httpAuth" />
					</el-form-item>
					<el-form-item
						v-if="proxyStore.proxyData.httpAuth"
						label="认证用户名："
						prop="authUsername"
					>
						<el-input
							v-model="proxyStore.proxyData.authUsername"
							placeholder="username"
						/>
					</el-form-item>
					<el-form-item
						v-if="proxyStore.proxyData.httpAuth"
						label="认证密码："
						prop="authPassword"
					>
						<el-input
							v-model="proxyStore.proxyData.authPassword"
							type="password"
							placeholder="password"
						/>
					</el-form-item>
					<el-form-item
						v-if="proxyStore.proxyData.type === 'https'"
						label="证书文件："
						prop="https2httpCaFile"
					>
						<el-input
							v-model="proxyStore.proxyData.https2httpCaFile"
							placeholder="点击选择证书文件"
							readonly
							@click="handleSelectFile(1, ['crt', 'pem'])"
						/>
					</el-form-item>
					<el-form-item
						v-if="proxyStore.proxyData.type === 'https'"
						label="密钥文件："
						prop="https2httpKeyFile"
					>
						<el-input
							v-model="proxyStore.proxyData.https2httpKeyFile"
							placeholder="点击选择密钥文件"
							readonly
							@click="handleSelectFile(2, ['key'])"
						/>
					</el-form-item>
				</template>

				<el-form-item
					v-if="
						proxyStore.proxyData.type === 'tcp' ||
						proxyStore.proxyData.type === 'udp'
					"
					label="外网端口："
					prop="remotePort"
				>
					<el-input-number
						v-model="proxyStore.proxyData.remotePort"
						:min="1"
						:max="65535"
					/>
				</el-form-item>

				<!-- STCP/XTCP/SUDP模式 -->
				<template
					v-if="
						proxyStore.proxyData.type === 'stcp' ||
						proxyStore.proxyData.type === 'xtcp' ||
						proxyStore.proxyData.type === 'sudp'
					"
				>
					<el-row :gutter="22">
						<el-col :span="14">
							<el-form-item
								label="STCP模式："
								prop="stcpModel"
							>
								<el-radio-group
									v-model="proxyStore.proxyData.stcpModel"
								>
									<el-radio
										v-for="model in stcpModels"
										:key="model.value"
										:value="model.value"
										>{{ model.label }}</el-radio
									>
								</el-radio-group>
							</el-form-item>
						</el-col>
						<el-col :span="10">
							<el-form-item
								label="共享密钥："
								prop="secretKey"
							>
								<el-input
									v-model="proxyStore.proxyData.secretKey"
									type="password"
									placeholder="密钥"
								/>
							</el-form-item>
						</el-col>
					</el-row>

					<!-- 被访问者代理名称 -->
					<el-form-item
						v-if="
							proxyStore.proxyData.type === 'stcp' ||
							proxyStore.proxyData.type === 'xtcp' ||
							proxyStore.proxyData.type === 'sudp'
						"
						label="被访问者代理名称："
						prop="visitedName"
					>
						<el-input
							v-model="proxyStore.proxyData.visitedName"
							placeholder="被访问者代理名称"
						/>
					</el-form-item>

					<template
						v-if="
							proxyStore.proxyData.type === 'stcp' ||
							proxyStore.proxyData.type === 'xtcp' ||
							proxyStore.proxyData.type === 'sudp'
						"
					>
						<el-row :gutter="20">
							<el-col :span="10">
								<el-form-item
									label="绑定地址："
									prop="bindAddr"
								>
									<el-input
										v-model="proxyStore.proxyData.bindAddr"
										placeholder="127.0.0.1"
									/>
								</el-form-item>
							</el-col>
							<el-col :span="10">
								<el-form-item
									label="绑定端口："
									prop="bindPort"
								>
									<el-input-number
										v-model="proxyStore.proxyData.bindPort"
										:min="1"
										:max="65535"
									/>
								</el-form-item>
							</el-col>
						</el-row>
					</template>
					<template v-if="proxyStore.proxyData.type === 'xtcp'">
						<el-row :gutter="20">
							<el-col :span="10">
								<el-form-item
									label="回退代理名称："
									prop="fallbackTo"
								>
									<el-input
										v-model="
											proxyStore.proxyData.fallbackTo
										"
										placeholder="回退代理名称"
									/>
								</el-form-item>
							</el-col>
							<el-col :span="10">
								<el-form-item
									label="回退超时毫秒："
									prop="fallbackTimeoutMs"
								>
									<el-input-number
										v-model="
											proxyStore.proxyData
												.fallbackTimeoutMs
										"
										:min="0"
									/>
								</el-form-item>
							</el-col>
						</el-row>
						<!-- 保持隧道开启 -->
						<el-form-item
							label="保持隧道开启："
							prop="keepAlive"
						>
							<el-switch
								v-model="proxyStore.proxyData.keepAlive"
							/>
						</el-form-item>
					</template>
				</template>

				<!-- 保存和取消按钮 -->
				<el-row justify="start">
					<el-button
						type="primary"
						@click="saveProxy"
						style="width: 100px"
						>保存</el-button
					>
					<el-button
						type="primary"
						style="width: 100px"
						@click="proxyDialogShow = false"
						>取消</el-button
					>
				</el-row>
			</el-form>
		</el-drawer>
	</div>
</template>

<style scoped>
	:deep(.el-drawer) {
		width: 60% !important;
	}
</style>
