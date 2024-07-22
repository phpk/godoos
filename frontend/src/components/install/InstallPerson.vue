<script setup lang="ts">
import { getSystemConfig, setSystemConfig } from "@/system/config";
import { OpenDirDialog } from "@/util/goutil";
import { Dialog } from "@/system";
const config = ref(getSystemConfig());
const storeList = [
  {
    title: "浏览器存储",
    value: "browser",
  },
  {
    title: "本地存储",
    value: "local",
  },
];
// if (config.value.isApp) {
//   storeList.push({
//     title: "本地存储",
//     value: "local",
//   });
// }
function selectFile() {
  OpenDirDialog().then((res: string) => {
    config.value.storePath = res;
  });
}
function handleSubmit() {
  const configData = toRaw(config.value);
  const postData: any = {
    name: "osInfo",
    type: "person",
  };
  if (configData.storeType === "local") {
    if (!configData.storePath || configData.storePath === "") {
      Dialog.showMessageBox({
        message: "存储地址不能为空",
        type: "error",
      });
      return;
    }
    postData.value = configData.storePath;
  }
  console.log(postData);
  const postUrl = config.value.apiUrl + "/system/setting";
  fetch(postUrl, {
    method: "POST",
    body: JSON.stringify(postData),
  })
    .then((res) => res.json())
    .then((res) => {
      console.log(res);
      if (res.code === 0) {
        setSystemConfig(configData);
        //localStorage.setItem("isGodoOSInstall", "true");
        Dialog.showMessageBox({
          message: "设置成功",
          type: "success",
        });
        location.reload();
      } else {
        Dialog.showMessageBox({
          message: res.message,
          type: "error",
        });
      }
    });
}
</script>
<template>
  <el-form label-width="auto" class="userForm">
    <el-form-item label="选择存储方式">
      <el-select v-model="config.storeType">
        <el-option
          v-for="(item, key) in storeList"
          :key="key"
          :label="item.title"
          :value="item.value"
        />
      </el-select>
    </el-form-item>
    <el-form-item label="存储路径" v-if="config.storeType === 'local'">
      <el-input v-model="config.storePath" @click="selectFile" />
    </el-form-item>
    <el-form-item class="subBtn">
      <el-button type="primary" @click="handleSubmit">确定</el-button>
    </el-form-item>
  </el-form>
</template>
