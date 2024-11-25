<template>
  <el-button type="primary" :icon="Plus" circle  @click="addPwd"/>
  <div class="file-pwd-list-box">
    <div class="file-pwd-list">
      <div class="pwd-box">
        <p>密码提示</p>
        <el-button type="danger" :icon="Delete" circle  @click="addPwd"/>
        <el-tag type="primary">default</el-tag>
      </div>
    </div>
  </div>
  <el-dialog v-model="dialogShow" title="添加密码" width="350px">
  <span>
    <el-form>
      <el-form-item label="密码提示">
        <el-input v-model="formData.pwdName"/>
      </el-form-item>
      <el-form-item label="密码">
        <el-input v-model="formData.pwd" show-password/>
      </el-form-item>
      <el-form-item label="是否为默认密码">
        <el-switch  v-model="formData.isDefault"/>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="dialogShow = false">
          确认
        </el-button>
      </el-form-item>
    </el-form>
  </span>
  </el-dialog>
</template>

<script lang="ts" setup>
import { Plus, Delete } from '@element-plus/icons-vue'
import { reactive, ref } from "vue";
const dialogShow = ref(false)
const formData = reactive({
  pwdName: '',
  pwd: '',
  isDefault: ''
})
function addPwd() {
  dialogShow.value = true
}
	// import { t } from "@/system";
	// import {
	// 	fetchGet,
	// 	getApiUrl,
	// 	getSystemConfig,
	// 	setSystemKey,
	// } from "@/system/config";
	// const config = getSystemConfig();
	// const isSetPwd = ref(false);
  // async function setFilePwd() {
  //   isSetPwd.value = !isSetPwd.value
  //   const params = {
  //     isPwd: isSetPwd.value ? 1 : 0
  //   }
	// 	await fetchGet(`${getApiUrl()}/file/changeispwd?ispwd=${params.isPwd}`);
	// 	setSystemKey("file", params);

  // }
	// onMounted(() => {
	// 	isSetPwd.value = config.file.isPwd ? true : false;
	// });
</script>

<style scoped lang="scss">
@import "./setStyle.css";
.file-pwd-list-box {
  width: 100%;
  height: 100%;
  .file-pwd-list {
    width: 100%;
    height: 90%;
    padding: 10px;
    box-sizing: border-box;
    background-color: rgb(248, 247, 247);
    overflow-y: scroll;

    .pwd-box {
      width: 95%;
      height: 60px;
      padding: 10px 30px;
      margin: 10px auto;
      border: 1px solid black;
      background-color: white;
      box-sizing: border-box;

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
}
	/* .file-pwd-box {
		padding-top: 20px;
	}
	.setting-item {
		display: flex;
		align-items: center;
	} */
</style>
