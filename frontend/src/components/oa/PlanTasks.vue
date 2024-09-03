<template>
    <el-form :model="form" label-width="auto" style="max-width: 560px;margin-top:20px;padding: 20px;">
        <el-form-item label="任务名称">
            <el-input v-model="form.name" />
        </el-form-item>
        <el-form-item label="任务跨度">
            <el-select v-model="form.duration">
                <el-option label="按天" value="1" />
                <el-option label="按周" value="7" />
                <el-option label="按月" value="30" />
                <el-option label="按季" value="90" />
                <el-option label="按年" value="365" />
            </el-select>
        </el-form-item>
        <el-form-item label="任务模板">
            <el-input v-model="form.tpl" />
        </el-form-item>
        <el-form-item label="选择人员">
            <el-select v-model="form.users" filterable multiple clearable collapse-tags placeholder="选择人员"
                popper-class="custom-header" :max-collapse-tags="1" value-key="id" style="width: 240px" @change="checkUsers">
                <template #header>
                    <el-checkbox v-model="checkAll" @change="handleCheckAll">
                        全选
                    </el-checkbox>
                </template>
                <el-option v-for="item in userList" :key="item.id" :label="item.nickname" :value="item.id" />
            </el-select></el-form-item>
        <el-form-item label="任务描述">
            <el-input v-model="form.desc" type="textarea" />
        </el-form-item>
        <div class="btn-group">
            <el-button type="primary" @click="onSubmit">创建任务</el-button>
        </div>
    </el-form>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { useSystem } from '@/system';
const sys = useSystem()
const userInfo: any = sys.getConfig('userInfo')
const userList = ref(userInfo.user_auths)
const checkAll = ref(false)
const form: any = ref({
    name: '',
    duration: '1',
    tpl: '',
    after: '1',
    users: [],
    desc: '',
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