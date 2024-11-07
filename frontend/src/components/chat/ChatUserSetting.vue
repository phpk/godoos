<script lang="ts" setup>
import { t } from '@/i18n';
import { ref, reactive, onMounted } from "vue";
import { useChatStore } from "@/stores/chat";
import { notifyError, notifySuccess } from "@/util/msg";
import { fetchGet, fetchPost, getSystemConfig, setSystemKey } from "@/system/config";
const store = useChatStore()
// const form = ref({
//     nickname: '',
// })
interface ImgItem {
    url: string;
    name: string;
}
const editType = ref(0)
const userInfo = getSystemConfig().userInfo
const userInfoForm: {[key: string]: any} = reactive({
    nickname: '',
    phone: '',
    email: '',
    desc: ''
})
const dialogShow = ref(false)
let imgList = reactive<ImgItem[]>([])
let pos = ref(0)
// 换头像
const showChange = (val: boolean) =>{
    dialogShow.value = val
}
// 获取头像列表
const getHeadList = async () => {
    let res = await fetchGet (`${userInfo.url}/files/avatarlist`)
    // const apiUrl = getSystemConfig().apiUrl + '/upload/avatar/';
    const apiUrl = userInfo.url + '/upload/avatar/';
    if (res.ok) {
        const response = await res.json()
        for (const item of response?.data) {
            imgList.push({
                url: apiUrl + item,
                name: item
            })
        }
    }
}
const chooseImg = (index: number) => {
    pos.value = index
}
const toChangeHead = async ()=>{
    let res = await fetchGet(`${userInfo.url}/files/saveavatar?id=${userInfo.id}&name=${imgList[pos.value]?.name}`)
    if (!res.ok) {
        notifyError("头像换取失败")
    } else {
        const response = await res.json()
        notifySuccess(response?.message || '头像换取成功')
        userInfo.avatar = imgList[pos.value]?.name
        store.userInfo.avatar = userInfo.url + '/upload/avatar/' + userInfo.avatar
        setSystemKey('userInfo', userInfo)
    }
    showChange(false)
}
//修改资料
const userRef:any = ref(null)
const userRule = {
    nickname: [
        { required: true, message: '昵称不能为空', trigger: 'blur'},
        { min: 2, max: 10, message: '昵称长度应该在2到10位', trigger: 'blur'}
    ],
    phone: [
        { required: true, message: '手机号不能为空', trigger: 'blur'},
        { pattern: /^1[3-9]\d{9}$/, message: '手机号格式不正确', trigger: 'blur'}
    ],
    email: [
        { required: true, message: '邮箱不能为空', trigger: 'blur'},
        { pattern: /^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$/, message: '邮箱格式不正确', trigger: 'blur'}
    ]
    // desc: [
    //     { required: true, message: '昵称不能为空', trigger: 'blur'},
    //     { min: 2, max: 10, message: '昵称长度应该在2到10位', trigger: 'blur'}
    // ],
}
const onSubmit = async () => {                                                                                                                                                                                                                                                                                                
    try {
        for (const key in userInfoForm) {
            userInfo[key] = userInfoForm[key]
        }
        const head = {
            'AuthorizationAdmin': userInfo.token
        }
        await userRef.value.validate()
        const res = await fetchPost(`${userInfo.url}/user/editdata?id=${userInfo.id}`, JSON.stringify(userInfoForm), head)
        const result = await res.json()
        if (result?.success) {
            notifySuccess(result?.message)
            setSystemKey('userInfo', userInfo)
        } else {
            notifyError(result?.message)
        }
    }catch (e) {
        console.log(e)
    }
}
//修改密码
const passwordForm = reactive({
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
})
const passwordRef:any = ref(null)
const rules = {
    oldPassword: [
        { required: true, message: '密码不能为空', trigger: 'blur' }
    ],
    newPassword: [
        { required: true, message: '密码不能为空', trigger: 'blur' },
        { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
    ],
    confirmPassword: [
        { required: true, message: '请再次输入密码', trigger: 'blur' },
        { validator: (rule: any, value:any, callback:any) => {
            console.log('rule:', rule);
            
        if (value === '') {
            callback(new Error('请再次输入密码'));
        } else if (value !== passwordForm.newPassword) {
            callback(new Error('两次输入的密码不一致'));
        } else {
            callback();
        }
        }, trigger: 'blur' }
    ],
}
const toChangePwd = async () => {
    try {
        await passwordRef.value.validate()
        const formData = new FormData()
        formData.append('oldpassword', passwordForm.oldPassword)
        formData.append('newpassword', passwordForm.newPassword)
        let res = await fetchPost(`${userInfo.url}/member/savepwd`, formData)
        const response = await res.json()
        if (response && response?.success) {
            notifySuccess(response?.message)
        } else {
            notifyError(response?.message)
        }
    } catch (e) {
        console.log(e)
    }
}
onMounted(()=>{
    getHeadList()
    for (const key in userInfoForm) {
        userInfoForm[key] = userInfo[key]
    }
    store.userInfo.avatar = userInfo.url + '/upload/avatar/' + userInfo.avatar
})
</script>
<template>
    <el-container>
        <el-aside>
            <div @click="editType = 0" :class="{'is-active': editType == 0}">修改资料</div>
            <div @click="editType = 1" :class="{'is-active': editType == 1}">修改密码</div>
        </el-aside>
        <el-main>  
            <el-form v-if="editType == 0" :model="userInfoForm" :rules="userRule"  ref="userRef" label-width="130px" style="padding: 30px;">
                <el-form-item label="头像">
                    <el-avatar :size="90" :src="store.userInfo.avatar" @click="showChange(true)"/>
                </el-form-item>
                <el-form-item label="昵称" prop="nickname">
                    <el-input v-model="userInfoForm.nickname" />
                </el-form-item>
                <el-form-item label="手机号" prop="phone">
                    <el-input v-model="userInfoForm.phone" />
                </el-form-item>
                <el-form-item label="邮箱" prop="email">
                    <el-input v-model="userInfoForm.email" />
                </el-form-item>
                <el-form-item label="自我介绍" prop="desc">
                    <el-input v-model="userInfoForm.desc" type="textarea" maxlength="100" show-word-limit/>
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="onSubmit">保存</el-button>
                </el-form-item>
            </el-form>
            <el-form v-else :model="passwordForm" :rules="rules" ref="passwordRef" label-width="120px" style="padding: 30px;">
                <el-form-item label="旧密码" prop="oldPassword">
                    <el-input 
                        v-model="passwordForm.oldPassword"
                        type="password" 
                        show-password
                    />
                </el-form-item>
                <el-form-item label="新密码" prop="newPassword">
                    <el-input  
                        v-model="passwordForm.newPassword"
                        type="password" 
                        show-password
                    />
                </el-form-item>
                <el-form-item label="再次输入密码" prop="confirmPassword">
                    <el-input  
                        v-model="passwordForm.confirmPassword" 
                        type="password" 
                        show-password
                    />
                </el-form-item>
                <el-form-item>
                    <el-button type="primary" @click="toChangePwd">保存</el-button>
                </el-form-item>
            </el-form>
        </el-main>
        <el-dialog v-model="dialogShow" title="修改头像" width="400px" draggable>
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
    </el-container>
</template>

<style scoped lang="scss">
:deep(.el-dialog) {
    .el-dialog__body {
        max-height: 300px;
        overflow-y: scroll
    }
}
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
.el-container {

    .el-aside {
        width: 20%;
        padding: 20px;
        background-color: #f9f9f9;
        box-sizing: border-box;

        div {
            width: 100%;
            min-width: 38px;
            height: 40px;
            font-size: 14px;
        }
        .is-active,
        div:hover {
            font-weight: bold;
            color: #16b777;
        }
        .el-input {
            width: 200px;
        }
    }
}
</style>
