<template>
    <div class="file-pwd-box">
        <div class="setting-item" >
            <label>文件密码</label>
            <el-input v-model="filePwd" placeholder="请输入文件加密密码" />
        </div>
        <div class="setting-item">
            <label></label>
            <el-button @click="toSetFilePwd" type="primary">{{ t("setFilePwd") }}</el-button>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import { t } from "@/system";
const filePwd = ref('')
const toSetFilePwd = ()=> {}
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
