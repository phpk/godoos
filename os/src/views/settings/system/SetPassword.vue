<template>
	<el-form
		:model="formData"
		:rules="rule"
		ref="pwdRef"
		label-position="top"
	>
		<el-row justify="end">
			<el-form-item>
				<el-button
					type="primary"
					icon="Plus"
					circle
					@click="closeDialog(true)"
				/>
			</el-form-item>
		</el-row>
		<el-form-item>
			<div class="file-pwd-list-box">
				<div class="file-pwd-list">
					<div
						class="pwd-box"
						v-for="item in paginatedPwdList"
						:key="item.id"
					>
						<p>{{ item.pwdName }}</p>
						<el-button
							type="danger"
							icon="Delete"
							circle
							@click="operationPwd('del', item)"
						/>
						<el-tag
							type="primary"
							v-if="item.isDefault == 1"
							>default</el-tag
						>
					</div>
				</div>
				<el-pagination
					background
					layout="prev, pager, next"
					v-model:current-page="filePwdStore.page.current"
					v-model:page-size="filePwdStore.page.size"
					:total="filePwdStore.page.total"
					@current-change="handleCurrentChange"
				/>
			</div>
		</el-form-item>
		<el-dialog
			v-model="dialogShow"
			title="添加密码"
			width="400px"
		>
			<span>
				<el-form-item
					label="密码提示"
					prop="pwdName"
				>
					<el-input v-model="formData.pwdName" />
				</el-form-item>
				<el-form-item
					label="密码"
					prop="pwd"
				>
					<el-input
						v-model="formData.pwd"
						show-password
					/>
				</el-form-item>
				<el-form-item label="是否为默认密码">
					<el-switch
						v-model="formData.isDefault"
						:active-value="1"
						:inactive-value="0"
					/>
				</el-form-item>
				<el-form-item>
					<el-button
						type="primary"
						@click="operationPwd('add')"
						style="margin: 0 auto"
					>
						确认
					</el-button>
				</el-form-item>
			</span>
		</el-dialog>
	</el-form>
</template>

<script lang="ts" setup>
	import { computed, reactive, ref, onMounted } from "vue";
	import { successMsg } from "@/utils/msg";
	const dialogShow = ref(false);
	const formData = reactive({
		pwdName: "",
		pwd: "",
		isDefault: 0,
	});
	const defaultChoose = ref(false);

	// 模拟数据
	const filePwdStore = reactive({
		pwdList: [
			{ id: 1, pwdName: "密码1", isDefault: 1 },
			{ id: 2, pwdName: "密码2", isDefault: 0 },
			{ id: 3, pwdName: "密码3", isDefault: 0 },
			{ id: 4, pwdName: "密码4", isDefault: 0 },
			{ id: 5, pwdName: "密码5", isDefault: 1 },
			{ id: 6, pwdName: "密码6", isDefault: 0 },
			{ id: 7, pwdName: "密码7", isDefault: 0 },
			{ id: 8, pwdName: "密码8", isDefault: 0 },
			{ id: 9, pwdName: "密码9", isDefault: 0 },
			{ id: 10, pwdName: "密码10", isDefault: 0 },
		],
		page: {
			current: 1,
			size: 8,
			total: 20,
		},
		hasDefaultPwd: true,
		addPwd: async (pwd: any) => {
			console.log("添加密码", pwd);
		},
		delPwd: async (id: number) => {
			console.log("删除密码", id);
		},
		getPage: async () => {
			console.log("获取页面数据");
		},
		pageChange: async (val: number) => {
			console.log("页面改变", val);
		},
		setDefaultPwd: async () => {
			console.log("设置默认密码");
		},
	});

	const rule = {
		pwdName: [
			{ required: true, message: "密码提示不能为空", trigger: "blur" },
			{
				min: 2,
				max: 10,
				message: "昵称长度应该在2到10位",
				trigger: "blur",
			},
		],
		pwd: [
			{ required: true, message: "密码不能为空", trigger: "blur" },
			{
				min: 6,
				max: 10,
				message: "密码长度应该在6到10位",
				trigger: "blur",
			},
		],
	};
	const pwdRef: any = ref(null);

	const paginatedPwdList = computed(() => {
		const start = (filePwdStore.page.current - 1) * filePwdStore.page.size;
		const end = start + filePwdStore.page.size;
		return filePwdStore.pwdList.slice(start, end);
	});

	function closeDialog(val: boolean) {
		dialogShow.value = val;
	}

	async function operationPwd(type: string, item?: any) {
		if (type === "add") {
			if (formData.isDefault === 1) {
				showPwdDialog("setDefault");
			} else {
				addPwd();
			}
		} else {
			if (item.isDefault === 1) {
				showPwdDialog("del", item);
			} else {
				delPwd(item);
			}
		}
	}

	function showPwdDialog(type: string, item?: any) {
		const content = type === "setDefault" ? "设置默认密码成功" : "删除成功";
		successMsg(content);
		if (type === "setDefault") {
			addPwd();
		} else {
			delPwd(item);
		}
	}

	async function addPwd() {
		await pwdRef.value.validate();
		closeDialog(false);
		const temp = {
			id: Date.now(),
			...formData,
		};
		if (formData.isDefault == 1) {
			filePwdStore.hasDefaultPwd
				? await filePwdStore.setDefaultPwd()
				: "";
		}
		await filePwdStore.addPwd(temp);
		filePwdStore.pwdList.push(temp);
		successMsg("添加成功");
		await initData();
	}

	async function initData() {
		await filePwdStore.getPage();
		formData.pwdName = "";
		formData.pwd = "";
		if (!filePwdStore.hasDefaultPwd) {
			formData.isDefault = 1;
			defaultChoose.value = true;
			return;
		}
		formData.isDefault = 0;
		defaultChoose.value = false;
		filePwdStore.page.total = filePwdStore.pwdList.length;
	}

	async function handleCurrentChange(val: number) {
		filePwdStore.page.current = val;
		await filePwdStore.pageChange(val);
		await initData();
	}

	async function delPwd(item: any) {
		if (item.isDefault == 1) {
			successMsg("默认密码不能删除");
			return;
		}
		await filePwdStore.delPwd(item.id);
		filePwdStore.pwdList = filePwdStore.pwdList.filter(
			(pwd) => pwd.id !== item.id
		);
		successMsg("删除成功");
		await initData();
	}

	onMounted(async () => {
		await initData();
		console.log("数据：", filePwdStore.pwdList, filePwdStore.page);
	});
</script>

<style scoped lang="scss">
	.file-pwd-list-box {
		width: 100%;
		height: 100%;
		.file-pwd-list {
			display: flex;
			flex-wrap: wrap;
			justify-content: space-around;
			height: auto;
			box-sizing: border-box;
			margin: 10px auto;

			.pwd-box {
				width: 45%;
				margin: 10px;
				padding: 15px;
				border: 1px solid #e0e0e0;
				border-radius: 8px;
				background-color: #f9f9f9;
				box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
				box-sizing: border-box;
				transition: transform 0.2s;

				&:hover {
					transform: translateY(-5px);
				}

				p {
					display: inline-block;
					height: 38px;
					line-height: 38px;
					margin: 0;
				}
				.el-tag {
					float: right;
					margin: 5px 20px;
				}
				.el-button {
					float: right;
				}
			}
		}
		.el-pagination {
			display: flex;
			justify-content: center;
		}
	}

	.el-form {
		:deep(.el-form-item__label) {
			min-width: 80px;
		}
	}
</style>
