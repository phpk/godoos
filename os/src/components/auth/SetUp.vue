<script lang="ts" setup>
import { ref,toRaw } from "vue";
import { errMsg, successMsg } from "@/utils/msg";
import { useSettingsStore } from "@/stores/settings";
const settingsStore = useSettingsStore();
const setForm = ref({
	userRole: "person",
	storeType: "local",
	netUrl: "",
});
const validateURL = (value: string) => {
  const urlPattern = /^(https?:\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})(:[0-9]+)?(\/[\w \.-]*)*\/?$/;
  return urlPattern.test(value)
};

const setSave = () => {
	const data = toRaw(setForm.value);
	//console.log("保存数据",data);
	if((data.storeType === 'net' || data.userRole === 'member') && !validateURL(data.netUrl)){
		errMsg("请输入正确的URL地址！");
		return;
	}
	settingsStore.setConfig('system',data)
	successMsg("设置成功！");
};

</script>
<template>
	<el-form :model="setForm" ref="formRef" label-position="left" label-width="0px">
		<!-- <el-form-item label-position="right">
			<el-button-group>
				<el-button :type="setForm.userRole == 'person' ? 'primary' : ''" @click="setForm.userRole = 'person'" icon="UserFilled" round>个人用户</el-button>
				<el-button  :type="setForm.userRole == 'member' ? 'primary' : ''" @click="setForm.userRole = 'member'" icon="Football" round>企业用户</el-button>
			</el-button-group>
		</el-form-item> -->
		<el-form-item label-position="right" v-if="setForm.userRole === 'person'">
			<el-button-group>
				<el-button :type="setForm.storeType == 'local' ? 'primary' : ''" @click="setForm.storeType = 'local'" icon="LocationInformation" round>本地服务</el-button>
				<el-button :type="setForm.storeType == 'net' ? 'primary' : ''" @click="setForm.storeType = 'net'" icon="Promotion" round>远程服务</el-button>
			</el-button-group>
		</el-form-item>
		<el-form-item prop="netUrl" v-if="setForm.storeType === 'net' || setForm.userRole === 'member'">
			<el-input v-model="setForm.netUrl" size="large" placeholder="请输入远程地址" autofocus prefix-icon="Link"></el-input>
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