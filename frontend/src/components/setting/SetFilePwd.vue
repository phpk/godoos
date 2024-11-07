<template>
    <div class="file-pwd-box">
        <div v-if="setPwd">
            <div class="setting-item" >
                <label>文件密码</label>
                <el-input v-model="filePwd" placeholder="请设置6-10位的密码" type="password" show-password/>
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
import { ref, onMounted } from "vue";
import { t } from "@/system";
import { fetchGet, getApiUrl, setSystemKey, getSystemConfig } from "@/system/config";
import { notifySuccess, notifyError } from "@/util/msg";
const filePwd = ref('')
const setPwd = ref(false)
const config = getSystemConfig()
const params = {
    isPwd: 1,
    pwd: '',
    salt: getSystemConfig().file.salt
}
// 设置文件密码
async function toSetFilePwd() {
    if (filePwd.value.length < 6 || filePwd.value.length > 10) {
        notifyError("密码长度应该在6-10位之间")
        return
    }
    params.pwd = md5(filePwd.value)
    params.isPwd = filePwd.value === '' ? 0 : 1
    const url = getApiUrl() + '/file/setfilepwd'
    const header = {
        'Salt': params.salt ? params.salt : 'vIf_wIUedciAd0nTm6qjJA==',
        'FilePwd': params.pwd
    }
    await fetchGet(`${getApiUrl()}/file/changeispwd?ispwd=${params.isPwd}`)
    const res = await fetchGet(url, header)
    if (res.ok){
        notifySuccess("设置文件密码成功");
    } else {
        params.isPwd = 0
        params.pwd = ''
        notifyError("设置文件密码失败")
    }
    setSystemKey('file',params)
}
async function clearPwd() {
    setPwd.value = false
    filePwd.value = ''
    params.isPwd = 0
    await fetchGet(`${getApiUrl()}/file/changeispwd?ispwd=0`)
    setSystemKey('file',params)
}
onMounted(()=>{
    params.isPwd = config.file.isPwd
    setPwd.value = params.isPwd ? true : false
})
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
