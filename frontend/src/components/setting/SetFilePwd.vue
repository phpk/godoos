<template>
  <el-button type="primary" :icon="Plus" circle  @click="closeDialog(true)"/>
  <div class="file-pwd-list-box">
    <div class="file-pwd-list">
      <div class="pwd-box" v-for="item in filePwdStore.pwdList">
        <p>{{ item.pwdName }}</p>
        <el-button type="danger" :icon="Delete" circle  @click="operationPwd('del', item)"/>
        <el-tag type="primary" v-if="item.isDefault == 1">default</el-tag>
      </div>
    </div>
    <el-pagination 
      background 
      layout="prev, pager, next" 
      v-model:current-page="filePwdStore.page.current"
      v-model:page-size="filePwdStore.page.size"
      :total="filePwdStore.page.total"
      @current-change="handleCurrentChange"
    />
  </div>
  <el-dialog v-model="dialogShow" title="添加密码" width="400px">
    <span>
      <el-form :model="formData" :rules="rule" ref="pwdRef" >
        <el-form-item label="密码提示" prop="pwdName">
          <el-input v-model="formData.pwdName"/>
        </el-form-item>
        <el-form-item label="密码" prop="pwd">
          <el-input v-model="formData.pwd" show-password/>
        </el-form-item>
        <el-form-item label="是否为默认密码">
          <el-switch  
            v-model="formData.isDefault" 
            :active-value="1"
            :inactive-value="0"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="operationPwd('add')" style="margin: 0 auto;">
            确认
          </el-button>
        </el-form-item>
      </el-form>
    </span>
  </el-dialog>
</template>

<script lang="ts" setup>
import { Plus, Delete } from '@element-plus/icons-vue'
import { reactive, ref, onMounted } from "vue";
import { useFilePwdStore } from '@/stores/filePwd';
import { fetchGet, getApiUrl } from "@/system/config";
import { md5 } from 'js-md5';
import { ElMessageBox } from 'element-plus'
const dialogShow = ref(false)
const formData = reactive({
  pwdName: '',
  pwd: '',
  isDefault: 0
})
const defaultChoose = ref(false)
const filePwdStore: any = useFilePwdStore()

const rule = {
  pwdName: [
    { required: true, message: '密码提示不能为空', trigger: 'blur'},
    { min: 2, max: 10, message: '昵称长度应该在2到10位', trigger: 'blur'}
  ],
  pwd: [
    { required: true, message: '密码不能为空', trigger: 'blur'},
    { min: 6, max: 10, message: '密码长度应该在6到10位', trigger: 'blur'}
  ]
}
const pwdRef:any = ref(null)
function closeDialog (val: boolean) {
  dialogShow.value = val
}
async function operationPwd(type: string, item?: any) {
  if (type == 'add') {
    formData.isDefault == 1 ? showPwdDialog(type) : addPwd()
  } else {
    item.isDefault == 1 ? showPwdDialog(type, item) : delPwd(item)
  }
}
// 弹窗提示用户是否操作默认密码
function showPwdDialog(type: string, item?: any) {
  const content = type == 'add' ? '确定设置为默认密码吗？' : '确定删除默认密码吗？'
  ElMessageBox.alert(content, '提示', {
    confirmButtonText: '确定',
    callback: (action: any) => {
      if (type == 'add') {
        action !== 'cancle' ? addPwd() : ''
      }else {
        action !== 'cancle' ? delPwd(item) : ''
      }
    }
  })
}
async function addPwd() {
  await pwdRef.value.validate()
  closeDialog(false)
  const temp = { ...formData }
  if (formData.isDefault == 1) {
    //保证默认密码唯一性
    filePwdStore.hasDefaultPwd ? await filePwdStore.setDefaultPwd() : ''
    await fetchDefaultPwd(formData.pwd)
  }
  await filePwdStore.addPwd(temp)
  await initData()
}
// 将默认密码传递给后端
async function fetchDefaultPwd(pwd: string) {
  const head = { 
    pwd: pwd !== '' ? md5(pwd) : ''
  }
  await fetchGet(`${getApiUrl()}/file/setfilepwd`, head);
}
async function initData () {
  await filePwdStore.getPage()
  formData.pwdName = ''
  formData.pwd = ''
  if (!filePwdStore.hasDefaultPwd) {
    formData.isDefault = 1
    defaultChoose.value = true
    return
  }
  formData.isDefault = 0
  defaultChoose.value = false
}
async function handleCurrentChange(val: number) {
 await filePwdStore.pageChange(val)
}
async function delPwd(item: any) {
  if (item.isDefault == 1) {
    await fetchDefaultPwd('')
  }
  await filePwdStore.delPwd(item.id)
  await initData()
}
onMounted(async() => {
  await initData()
  console.log('数据：', filePwdStore.pwdList, filePwdStore.page);
})
</script>

<style scoped lang="scss">
@import "./setStyle.css";
.file-pwd-list-box {
  width: 100%;
  height: 100%;
  .file-pwd-list {
    width: 100%;
    height: 90%;
    padding: 10px;
    box-sizing: border-box;
    background-color: rgb(248, 247, 247);
    overflow-y: scroll;
    margin: 10px auto;

    .pwd-box {
      width: 95%;
      height: 60px;
      padding: 10px 30px;
      margin: 10px auto;
      border: 1px solid black;
      background-color: white;
      box-sizing: border-box;

      p {
        display: inline-block;
        height: 38px;
        line-height: 38px;
        margin: 0;
      }
      .el-tag {
        float: right;
        margin: 5px 20px;
      }
      .el-button {
        float: right;
      }
    }
  }
  .el-pagination {
    display: flex;
    justify-content: center;
  }
}
.el-form {
  :deep(.el-form-item__label) {
    min-width: 80px;
  }
}
</style>
