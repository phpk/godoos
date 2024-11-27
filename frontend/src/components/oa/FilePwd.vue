<template>
	<el-form
		label-width="auto"
		style="max-width: 560px; margin-top: 20px; padding: 20px"
	>
		<el-form-item label="文件密码">
			<el-input
				v-model="filePwd"
				type="password"
				show-password
			/>
		</el-form-item>
		<div class="btn-group">
			<el-button
				type="primary"
				@click="setFilePwd"
				>提交</el-button
			>
		</div>
	</el-form>
</template>

<script lang="ts" setup>
import { BrowserWindow, useSystem } from "@/system";
import { notifyError, notifySuccess } from "@/util/msg";
// import { md5 } from "js-md5";
import { ref } from "vue";
import { getSystemConfig, setSystemKey } from "@/system/config";
const window: BrowserWindow | undefined = inject("browserWindow");
const filePwd = ref("");
const sys = useSystem();
async function setFilePwd() {
  if (
    filePwd.value !== "" &&
    filePwd.value.length >= 6 &&
    filePwd.value.length <= 10
  ) {
    const path = window?.config.path || "";
    const header = {
      pwd: filePwd.value
      // pwd: getSystemConfig().userType == 'person' ? md5(filePwd.value) : filePwd.value
    };
    const file = await sys.fs.readFile(path);
    if (file === false) return;
    const res = await sys.fs.writeFile(path, file, header);
    //console.log('res:', res);
    
    if (res && res.code == 0) {
      notifySuccess("文件密码设置成功");
      localStorageFilePwd(path, filePwd.value)
    } else {
      notifyError("文件密码设置失败");
    }
    //console.log("路径：", res, path);
  }
}
// 开源版存储文件密码
function localStorageFilePwd (path:string, pwd: string) {
  if (getSystemConfig().file.isPwd && getSystemConfig().userType == 'person') {
    let fileInputPwd = getSystemConfig().fileInputPwd
    const pos = fileInputPwd.findIndex((item: any) => item.path == path)
    if (pos !== -1) {
      fileInputPwd[pos].pwd = pwd
    } else {
      fileInputPwd.push({
        path: path,
        pwd: pwd
      })
    }
    setSystemKey('fileInputPwd', fileInputPwd)
  }
}
</script>
<style scoped>
	.btn-group {
		display: flex;
		justify-content: center;
	}
</style>
