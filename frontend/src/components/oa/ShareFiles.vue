<template>
    <el-form :model="form" label-width="auto" style="max-width: 560px;margin-top:20px;padding: 20px;">
        <el-form-item label="分享给">
            <el-select v-model="form.users" filterable multiple clearable collapse-tags placeholder="选择人员"
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
            <el-switch v-model="form.canEditor" />
        </el-form-item>
        <div class="btn-group">
            <el-button type="primary" @click="onSubmit">发布分享</el-button>
        </div>
    </el-form>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useSystem } from '@/system';
const sys = useSystem()
const userInfo: any = sys.getConfig('userInfo')
const userList = ref(userInfo.user_shares)
const checkAll = ref(false)
const form: any = ref({
    users: [],
    canEditor: false,
})

const handleCheckAll = (val: any) => {
    if (val) {
        form.value.users = userList.value.map((d: any) => d.value)
    } else {
        form.value.users.value = []
    }
}
const checkUsers = (val: any) => {
    const res:any = []
    val.forEach((item: any) => {
        if(item){
            res.push(item)
        }
    })
    form.value.users = res
}
const onSubmit = () => {
    console.log('submit!')
}
</script>
<style scoped>
.btn-group{
    display: flex;
    justify-content: center;
}
</style>