<template>
  <div class="lockscreen" :class="lockClassName">
    <el-card class="login-box" shadow="never">
      <div class="avatar-container">
        <el-avatar size="large">
          <img src="/logo.png" alt="Logo" />
        </el-avatar>
      </div>
      <el-form v-if="!isRegisterMode" label-position="left" label-width="0px">
        <el-form-item>
          <el-input
            v-model="userName"
            placeholder="请输入用户名"
            autofocus
            prefix-icon="UserFilled"
          ></el-input>
        </el-form-item>
        <el-form-item v-if="!sys._options.noPassword">
          <el-input
            v-model="userPassword"
            type="password"
            placeholder="请输入登录密码"
            show-password
            prefix-icon="Key"
            @keyup.enter="onLogin"
          ></el-input>
        </el-form-item>
        <el-button type="primary" @click="onLogin">登录</el-button>
        <div class="actions" v-if="config.userType === 'member'">
          <a href="#" @click.prevent="toggleRegister">注册新用户</a>
          <a href="#" @click.prevent="toggleUserSwitch">切换角色</a>
        </div>
      </el-form>
      <el-form v-else label-position="left" label-width="0px" :model="regForm" ref="regFormRef" :rules="rules">
        <el-form-item prop="username">
          <el-input
            v-model="regForm.username"
            placeholder="请输入用户名"
            prefix-icon="UserFilled"
          ></el-input>
        </el-form-item>
        <el-form-item prop="nickname">
          <el-input
            v-model="regForm.nickname"
            placeholder="请输入真实姓名"
            prefix-icon="Avatar"
          ></el-input>
        </el-form-item>
        <el-form-item prop="email">
          <el-input
            v-model="regForm.email"
            placeholder="请输入邮箱"
            prefix-icon="Message"
          ></el-input>
        </el-form-item>
        <el-form-item prop="phone">
          <el-input
            v-model="regForm.phone"
            placeholder="请输入手机号"
            prefix-icon="Iphone"
          ></el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input
            v-model="regForm.password"
            type="password"
            placeholder="请输入密码"
            show-password
            prefix-icon="Key"
          ></el-input>
        </el-form-item>
        <el-form-item prop="confirmPassword">
          <el-input
            v-model="regForm.confirmPassword"
            type="password"
            placeholder="请再次输入密码"
            show-password
            prefix-icon="Lock"
          ></el-input>
        </el-form-item>
        <el-button type="primary" @click="onRegister">注册</el-button>
        <div class="actions">
          <a href="#" @click.prevent="toggleRegister">返回登录</a>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { useSystem } from '@/system';
import { getSystemConfig, setSystemConfig } from '@/system/config';
import { notifyError } from '@/util/msg';
import { RestartApp } from '@/util/goutil';

const sys = useSystem();
const loginCallback = sys._options.loginCallback;
const config = getSystemConfig();
const lockClassName = ref('screen-show');
const isRegisterMode = ref(false);

function loginSuccess() {
  lockClassName.value = 'screen-hidean';
  setTimeout(() => {
    lockClassName.value = 'screen-hide';
  }, 500);
}

const userName = ref('');
const userPassword = ref('');

onMounted(() => {
  if (config.userType === 'person') {
    userName.value = sys._options.login?.username || 'admin';
    userPassword.value = sys._options.login?.password || '';
  } else {
    userName.value = config.userInfo.username;
    userPassword.value = config.userInfo.password;
  }
});

async function onLogin() {
  localStorage.removeItem("godoosClientId");
  if (loginCallback) {
    const res = await loginCallback(userName.value, userPassword.value);
    if (res) {
      loginSuccess();
    }
  }
}

function toggleRegister() {
  isRegisterMode.value = !isRegisterMode.value;
}

function toggleUserSwitch() {
  config.userType = 'person'
  setSystemConfig(config);
  RestartApp();
}
const regForm = ref({
  username: '',
  password: '',
  confirmPassword: '',
  email: '',
  phone: '',
  nickname: '',
})
const regFormRef:any = ref(null);

const rules = {
  username: [
    { required: true, message: '用户名不能为空', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度应在3到20个字符之间', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '密码不能为空', trigger: 'blur' },
    { min: 6, message: '密码长度不能小于6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请再次输入密码', trigger: 'blur' },
    { validator: (rule:any, value:any, callback:any) => {
      console.log(rule)
      if (value === '') {
        callback(new Error('请再次输入密码'));
      } else if (value !== regForm.value.password) {
        callback(new Error('两次输入的密码不一致'));
      } else {
        callback();
      }
    }, trigger: 'blur' }
  ],
  email: [
    { required: true, message: '邮箱不能为空', trigger: 'blur' },
    { type: 'email', message: '请输入有效的邮箱地址', trigger: ['blur', 'change'] }
  ],
  phone: [
    { required: true, message: '手机号不能为空', trigger: 'blur' },
    { pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '昵称不能为空', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度应在2到20个字符之间', trigger: 'blur' }
  ]
};

async function onRegister() {
  try {
    await regFormRef.value.validate();
    const save = toRaw(regForm.value);
    const userInfo = config.userInfo;
    const comp = await fetch(userInfo.url + '/member/register', {
      method: 'POST',
      body: JSON.stringify(save),
    });
    if(!comp.ok){
      notifyError('网络错误，注册失败');
      return
    }
    const res = await comp.json();
    if(res.success){
      notifyError('注册成功');
      toggleRegister();
    }else{
      notifyError(res.message);
      return
    }
  } catch (error) {
    console.error(error);
  }

}

</script>

<style scoped lang="scss">
.lockscreen {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  z-index: 201;
  display: flex;
  justify-content: center;
  align-items: center;
  overflow: hidden;
  color: #fff;
  background-color: rgba(25, 28, 34, 0.78);
  backdrop-filter: blur(7px);

  .login-box{
    width: 300px;
    padding: 20px;
    text-align: center;
    background: #ffffff;
    border-radius: 10px;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    border: 1px solid #e0e0e0;

    .avatar-container {
      margin-bottom: 20px;
    }

    .el-input {
      width: 100%;
      margin-bottom: 10px;
      background: #f9f9f9;
      border-radius: 4px;
    }

    .el-button {
      width: 100%;
      margin-top: 10px;
      background: #409eff;
      color: #ffffff;
      border: none;
      border-radius: 4px;
      transition: background 0.3s ease;
    }

    .el-button:hover {
      background: #66b1ff;
    }

    .tip {
      padding: 4px 0;
      font-size: 12px;
      color: red;
      height: 30px;
    }

    .actions {
      margin-top: 10px;
      display: flex;
      justify-content: space-between;

      a {
        color: #409eff;
        text-decoration: none;
        cursor: pointer;

        &:hover {
          text-decoration: underline;
        }
      }
    }
  }

  .screen-hidean {
    animation: outan 0.5s forwards;
  }

  .screen-hide {
    display: none;
  }
}

@keyframes outan {
  0% {
    opacity: 1;
  }

  30% {
    opacity: 0;
  }

  100% {
    transform: translateY(-100%);
    opacity: 0;
  }
}
</style>