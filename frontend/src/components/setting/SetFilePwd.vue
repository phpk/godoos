<template>
    <div class="file-pwd-box">
        <div v-if="setPwd">
            <div class="setting-item" >
                <label>文件密码</label>
                <el-input v-model="filePwd" placeholder="请输入文件加密密码" type="password"/>
            </div>
            <div class="setting-item">
                <label></label>
                <el-button @click="toSetFilePwd" type="primary">{{ t("setFilePwd") }}</el-button>
                <el-button @click="clearPwd" type="primary">取消文件加密</el-button>
            </div>
        </div>
        <div v-else class="setting-item">
            <label></label>
            <el-button @click="setPwd = true" type="primary">设置文件密码</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { md5 } from "js-md5";
import { ref } from "vue";
import { t } from "@/system";
import { fetchGet, getSystemConfig } from "@/system/config";
import { notifySuccess } from "@/util/msg";
const filePwd = ref('')
const setPwd = ref(false)
async function toSetFilePwd() {
    console.log('pwd:', filePwd);
    const url = getSystemConfig().userInfo.url + '/file/setfilepwd'
    const header = {
        'Salt': 'vIf_wIUedciAd0nTm6qjJA==',
        'FilePwd': md5(filePwd.value)
    }
    const res = await fetchGet(url, header)
    if (res.ok){
        notifySuccess("保存成功111");
    }
}
function clearPwd() {
    setPwd.value = false
}
</script>

<style scoped>
@import "./setStyle.css";
.file-pwd-box {
    padding-top: 20px;
}
.setting-item {
  display: flex;
  align-items: center;
}
</style>
