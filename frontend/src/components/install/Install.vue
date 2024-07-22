<template>
  <div v-cloak class="container">
    <div class="title-bar">
      <h1>安装系统</h1>
    </div>
    <el-row class="row-bg" justify="space-around" v-if="currentUserType == ''">
      <el-card
        hover
        v-for="(item, key) in userTypeList"
        :key="key"
        class="box-card"
        style="width: 30%"
        @click="changeView(item.value, item.title)"
      >
        <template #header>
          <div class="card-header">
            <span>{{ item.title }}</span>
          </div>
        </template>
        <div class="text item">{{ item.desc }}</div>
        <template #footer>
          <el-text>选择</el-text>
        </template>
      </el-card>
    </el-row>
    <div class="formbox" v-if="currentUserType != ''">
      <el-row>
        <el-button circle icon="ArrowLeftBold" @click="currentUserType = ''"></el-button>
        <el-text style="margin-left: 5px">{{ currentTitle }}</el-text>
      </el-row>
      <div v-if="currentUserType == 'person'">
        <InstallPerson />
      </div>
      <div v-if="currentUserType == 'member'">
        <InstallMember />
      </div>
      <div v-if="currentUserType == 'compony'">
        <InstallCompony />
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
const currentUserType = ref("");
const currentTitle = ref("");
const userTypeList = [
  {
    title: "个人用户",
    value: "person",
    desc: "个人用户是指独立的用户，可以创建自己的OS，并管理自己的文件，无交互功能。",
  },
  {
    title: "企业用户",
    value: "member",
    desc: "企业用户是指企业员工，可以创建企业员工的个人OS，能和其他员工交互。",
  },
  {
    title: "企业管理员",
    value: "compony",
    desc: "企业管理员是指企业管理员，可以创建企业员工的帐户，并管理维护企业的各种信息。",
  },
];
function changeView(value: string, title: string) {
  currentUserType.value = value;
  currentTitle.value = title;
}
</script>
<style scoped>
[v-cloak] {
  display: none;
}
/* :root {
  --win11-bg-color: #222;
  --win11-accent-color: #00bcf2;
  --win11-accent-hover: #0095dd;
  --win11-text-color: #eee;
} */
.container {
  margin: 5%;
  border-radius: 12px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  background-color: #fff;
  padding-bottom: 20px;
}
.formbox {
  padding: 0 5%;
}
/* 添加标题栏样式 */
.title-bar {
  /* 保持宽度 100%，并确保在容器内水平居中 */
  width: 100%;
  background-color: #00bcf2;
  padding: 10px 0px;
  display: flex;
  justify-content: center;
  align-items: center;
  border-radius: 12px 12px 0 0;

  /* 添加底部外边距以创建与下方卡片的间距 */
  margin-bottom: 20px;
}

.title-bar h1 {
  color: #eee;
  font-size: 24px;
  margin: 0;
}
/* 使 el-card 的 footer 内容居中 */
::v-deep .el-card__footer {
  display: flex;
  justify-content: center;
  align-items: center;
}
.subBtn {
  display: flex;
  justify-content: center;
  align-items: center;
}
.subBtn .el-button {
  margin: 30px auto;
}
.userForm {
  max-width: 600px;
  margin-top: 15px;
  padding: 0 50px;
}
</style>
