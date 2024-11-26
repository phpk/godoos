<template>
  <el-button type="primary" :icon="Plus" circle  @click="closeDialog(true)"/>
  <div class="file-pwd-list-box">
    <div class="file-pwd-list">
      <div class="pwd-box" v-for="item in filePwdStore.pwdList">
        <p>{{ item.pwdName }}</p>
        <el-button type="danger" :icon="Delete" circle  @click="addPwd"/>
        <el-tag type="primary" v-if="item.isDefault == 1">default</el-tag>
      </div>
    </div>
    <el-pagination 
      background 
      layout="prev, pager, next" 
      :total="filePwdStore.page.total"
      :size="filePwdStore.page.size"
      :current="filePwdStore.page.current"
    />
  </div>
  <el-dialog v-model="dialogShow" title="添加密码" width="400px">
    <span>
      <el-form>
        <el-form-item label="密码提示">
          <el-input v-model="formData.pwdName"/>
        </el-form-item>
        <el-form-item label="密码">
          <el-input v-model="formData.pwd" show-password/>
        </el-form-item>
        <el-form-item label="是否为默认密码">
          <el-switch  
            v-model="formData.isDefault" 
            active-value="1"
            inactive-value="0"/>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="addPwd" style="margin: 0 auto;">
            确认
          </el-button>
        </el-form-item>
      </el-form>
    </span>
  </el-dialog>
</template>

<script lang="ts" setup>
import { Plus, Delete } from '@element-plus/icons-vue'
import { reactive, ref, onMounted } from "vue";
import { useFilePwdStore } from '@/stores/filePwd';
const dialogShow = ref(false)
const formData = reactive({
  pwdName: '',
  pwd: '',
  isDefault: ''
})
const filePwdStore: any = useFilePwdStore()
function closeDialog (val: boolean) {
  dialogShow.value = val
}
async function addPwd() {
  closeDialog(false)
  const temp = { ...formData }
  await filePwdStore.addPwd(temp)
  await initData()
}
async function initData () {
  await filePwdStore.getPage()
}
onMounted(async() => {
  await initData()
  // console.log('数据：', filePwdStore.pwdList, filePwdStore.page.total);
})
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
    margin: 10px auto;

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
	/* .file-pwd-box {
		padding-top: 20px;
	}
	.setting-item {
		display: flex;
		align-items: center;
	} */
</style>
