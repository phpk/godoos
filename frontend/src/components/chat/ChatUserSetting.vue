<script setup>
import { t } from '@/i18n';
import { ref, onMounted } from "vue";
import { useChatStore } from "@/stores/chat";
import { notifyError, notifySuccess } from "@/util/msg";
import { fetchGet, getSystemConfig, setSystemKey } from "@/system/config";
const store = useChatStore()
const form = ref({
    nickname: '',
})
const userInfo = getSystemConfig().userInfo
const dialogShow = ref(false)
let imgList = reactive([])
let pos = ref(0)
const showChange = (val) =>{
    dialogShow.value = val
}
const toChangeHead = async ()=>{
    let res = await fetchGet(`${userInfo.url}/files/saveavatar?id=${userInfo.id}&name=${imgList[pos.value].name}`)
    console.log('submit!:',res)
    if (!res.ok) {
        notifyError("头像换取失败")
    } else {
        res = await res.json()
        notifySuccess(res.message)
        userInfo.avatar = imgList[pos.value].name
        store.userInfo.avatar = userInfo.url + '/upload/avatar/' + userInfo.avatar
        setSystemKey('userInfo', userInfo)
    }
    showChange(false)
}
// 获取头像列表
const getHeadList = async () => {
    let res = await fetchGet (`${userInfo.url}/files/avatarlist`)
    // const apiUrl = getSystemConfig().apiUrl + '/upload/avatar/';
    const apiUrl = userInfo.url + '/upload/avatar/';
    if (res.ok) {
        res = await res.json()
        for (const item of res.data) {
            imgList.push({
                url: apiUrl + item,
                name: item
            })
        }
    }
}
// 换头像
const chooseImg = (index) => {
    pos.value = index
}
const onSubmit = async () => {
    console.log('提交');
}
onMounted(()=>{
    getHeadList()
    store.userInfo.avatar = userInfo.url + '/upload/avatar/' + userInfo.avatar
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
            <span
                v-for="(item, index) in imgList"
                :key="item.name"
                class="img-box"
            >
                <el-avatar  
                    :size="80" 
                    :src="item.url"
                    @click="chooseImg(index)"
                    :class="{'is-active': pos == index}"
                >
                </el-avatar>
                <el-icon v-show="pos == index"><Select /></el-icon>
            </span>
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
.img-box {
    position: relative;
    .el-avatar {
        margin: 10px;
    }
    .is-active {
        width: 90px;
        height: 90px;
        border: 2px solid green;
        box-sizing: border-box;
    }
    .el-icon {
        position: absolute;
        bottom: 0;
        left: 50%;
        width: 25px;
        height: 25px;
        border-radius: 50%;
        color: white;
        background-color: green;
        transform: translate(-50%,0)
    }
}

</style>
