<template>
    <el-form label-width="auto" style="max-width: 560px;margin-top:20px;padding: 20px;">
        <el-form-item label="文件密码">
            <el-input v-model="filePwd" type="password" show-password/>
        </el-form-item>
        <div class="btn-group">
            <el-button type="primary" @click="setFilePwd">提交</el-button>
        </div>
    </el-form>
</template>

<script lang="ts" setup>
import { fetchPost } from '@/system/config';
import { ref } from 'vue' 
import { getSystemConfig } from "@/system/config";
import { BrowserWindow } from '@/system';
const window: BrowserWindow | undefined = inject("browserWindow");
const filePwd = ref('')
const userInfo = getSystemConfig().userInfo
async function setFilePwd() {
    if (filePwd.value !== '' && filePwd.value.length >= 6 && filePwd.value.length <= 10) {
        const path = window?.config.path || ''
        const res = await fetchPost(`${userInfo.url}`, filePwd.value)
        console.log('路径：', res, path);
    }
}
</script>
<style scoped>
.btn-group{
    display: flex;
    justify-content: center;
}
</style>