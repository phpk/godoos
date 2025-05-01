<script lang="ts" setup>
import { reactive, ref } from "vue";
const setForm = ref({
	userRole: "person",
	storeType: "local",
	netUrl: "",
});
const rules = reactive({
	username: [
		{ required: true, message: "请输入用户名", trigger: "blur" },
		{
			min: 3,
			max: 20,
			message: "长度在 3 到 20 个字符",
			trigger: "blur",
		},
	],
	password: [
		{ required: true, message: "请输入密码", trigger: "blur" },
		{
			min: 6,
			max: 20,
			message: "长度在 6 到 20 个字符",
			trigger: "blur",
		},
	],
});
const setSave = () => {
};

</script>
<template>
	<el-form :model="setForm" :rules="rules" ref="emailLoginFormRef" label-position="left" label-width="0px">
		<el-form-item label-position="right">
			<el-button-group>
				<el-button :type="setForm.userRole == 'person' ? 'primary' : ''" @click="setForm.userRole = 'person'" icon="ArrowLeft" round>个人用户</el-button>
				<el-button  :type="setForm.userRole == 'member' ? 'primary' : ''" @click="setForm.userRole = 'member'" icon="ArrowRight" round>企业用户</el-button>
			</el-button-group>
		</el-form-item>
		<el-form-item label-position="right" v-if="setForm.userRole === 'person'">
			<el-button-group>
				<el-button :type="setForm.storeType == 'local' ? 'primary' : ''" @click="setForm.storeType = 'local'" icon="ArrowLeft" round>本地服务</el-button>
				<el-button :type="setForm.storeType == 'net' ? 'primary' : ''" @click="setForm.storeType = 'net'" icon="ArrowRight" round>远程服务</el-button>
			</el-button-group>
		</el-form-item>
		<!-- <el-form-item label-position="right">
			<el-radio-group v-model="setForm.userRole" aria-label="label position">
				<el-radio-button value="person" round>个人用户</el-radio-button>
				<el-radio-button value="member">企业用户</el-radio-button>
			</el-radio-group>
		</el-form-item>
		<el-form-item label-position="right" v-if="setForm.userRole === 'person'">
			<el-radio-group v-model="setForm.storeType" aria-label="label position">
				<el-radio-button value="local">本地存储</el-radio-button>
				<el-radio-button value="brower">浏览器存储</el-radio-button>
				<el-radio-button value="net">远程存储</el-radio-button>
			</el-radio-group>
		</el-form-item> -->
		<el-form-item prop="neturl" v-if="setForm.storeType === 'net' || setForm.userRole === 'member'">
			<el-input v-model="setForm.netUrl" size="large" placeholder="请输入远程地址" autofocus
				prefix-icon="Link"></el-input>
		</el-form-item>
		<el-form-item class="button-center">
			<el-button class="login-button" type="primary" size="large" @click="setSave">保存</el-button>
		</el-form-item>
	</el-form>
</template>
<style scoped>
.button-center {
	width: 100%;

	.login-button {
		width: 100%;
		border-radius: 50px;
	}
}
:deep(.el-button-group){
	margin: auto;
}
:deep(.el-button){
	height: 45px;
}
:deep(.el-radio-group) {
	margin-left: 30px;
}

:deep(.el-input__wrapper) {
	border-radius: 50px;
	height: 45px;
}
</style>