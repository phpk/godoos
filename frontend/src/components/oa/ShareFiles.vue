<template>
    <el-form :model="form" label-width="auto" style="max-width: 560px;margin-top:20px;padding: 20px;">
        <el-form-item label="分享给">
            <el-select v-model="form.receverid" filterable multiple clearable collapse-tags placeholder="选择人员"
                popper-class="custom-header" :max-collapse-tags="1" value-key="id" style="width: 240px" @change="checkUsers">
                <template #header>
                    <el-checkbox v-model="checkAll" @change="handleCheckAll">
                        全选
                    </el-checkbox>
                </template>
                <el-option v-for="item in userList" :key="item.id" :label="item.nickname" :value="item.id" />
            </el-select>
        </el-form-item>
        <el-form-item label="编辑权限">
            <el-switch v-model="form.iswrite" active-value="1" inactive-value="0" />
        </el-form-item>
        <div class="btn-group">
            <el-button type="primary" @click="onSubmit">发布分享</el-button>
        </div>
    </el-form>
</template>

<script lang="ts" setup>
import { ref, inject } from 'vue';
import { useSystem, BrowserWindow } from '@/system';
import { getSystemConfig, fetchPost } from "@/system/config";

const window: BrowserWindow | undefined = inject("browserWindow");
const sys = useSystem()
const userInfo: any = sys.getConfig('userInfo')
const userList = ref(userInfo.user_shares)
const checkAll = ref(false)
const form: any = ref({
    senderid: '',
    receverid: [],
    path: '',
    iswrite: '0'
})
const config = ref(getSystemConfig())

const handleCheckAll = (val: any) => {
    if (val) {
        form.value.receverid = userList.value.map((d: any) => d.value)
    } else {
        form.value.receverid.value = []
    }
}
const checkUsers = (val: any) => {
    const res:any = []
    val.forEach((item: any) => {
        if(item){
            res.push(item)
        }
    })
    form.value.receverid = res
}
const onSubmit = async () => {
    const apiUrl = config.value.userInfo.url + '/files/share'
    form.value.senderid = config.value.userInfo.id
    form.value.path = window?.config.path || ''
    const temp = {...form.value}
    temp.senderid = temp.senderid.toString()
    temp.receverid = temp.receverid.map((item:any) => item.toString())
    await fetchPost(apiUrl, new URLSearchParams(temp))
}
</script>
<style scoped>
.btn-group{
    display: flex;
    justify-content: center;
}
</style>