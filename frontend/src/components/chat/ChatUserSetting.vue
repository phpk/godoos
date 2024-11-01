<script setup>
import { t } from '@/i18n';
import { ref, onMounted } from "vue";
import { useChatStore } from "@/stores/chat";
import { fetchGet, getSystemConfig } from "@/system/config";
const store = useChatStore()
const form = ref({
    nickname: '',
})
const userInfo = getSystemConfig().userInfo
const dialogShow = ref(false)
let imgList = reactive([])
const showChange = (val) =>{
    dialogShow.value = val
}
const toChangeHead = ()=>{
    console.log('换头像');
}
// 获取头像列表
const getHeadList = async() => {
    const res = await fetchGet(`${userInfo.url}/files/avatarlist`)
    if (res.ok) {
        imgList = await res.json()
        imgList = imgList.data.map(item => {
            item.isChoose = false
            return item
        })
        imgList[0].isChoose = true
    }
    console.log('imgList:' , imgList);
}
let pos = ref(0)
// 换头像
const chooseImg = (index) => {
    imgList.forEach(item => item.isChoose = false)
    imgList[index].isChoose = true
    pos.value = index
}
const onSubmit = () => {
    console.log('submit!')
}
onMounted(()=>{
    getHeadList()
})
</script>
<template>
    <div>
    <el-form :model="form" label-width="160px" style="padding: 30px;">
        <el-form-item label="头像">
            <el-avatar :size="90" :src="store.userInfo.avatar" @click="showChange(true)"/>
        </el-form-item>
        <el-form-item label="昵称">
            <el-input v-model="store.userInfo.nickname" />
        </el-form-item>
        <el-form-item label="手机号">
            <el-input v-model="store.userInfo.phone" />
        </el-form-item>
        <el-form-item label="工号">
            <el-input v-model="store.userInfo.job_number" />
        </el-form-item>
        <el-form-item label="自我介绍">
            <el-input v-model="store.userInfo.desc" />
        </el-form-item>
        <el-form-item>
            <el-button type="primary" @click="onSubmit">保存</el-button>
        </el-form-item>
    </el-form>
    <el-dialog v-model="dialogShow" title="修改头像" width="500px" draggable>
        <div>
            <el-avatar 
                v-for="(item, index) in imgList" 
                :key="item.name" 
                :size="90" 
                :src="item.url"
                :class="{'is-active': item.isChoose}"
                @click="chooseImg(index)"/>
        </div>
        <template #footer>
        <div class="dialog-footer">
            <el-button @click="showChange(false)">{{  t("cancel") }}</el-button>
            <el-button type="primary" @click="toChangeHead">
            {{  t("confirm") }}
            </el-button>
        </div>
        </template>
    </el-dialog>

</div>
</template>

<style scoped lang="scss">
.el-avatar {
    margin: 10px;
}
.is-active {
    width: 100px;
    height: 100px;
    border: 2px solid red;
}
</style>
