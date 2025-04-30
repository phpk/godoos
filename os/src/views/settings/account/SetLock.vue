<template>
  <div class="lock-screen">
    <el-form :model="form" class="lock-form" label-position="top">
      <el-form-item label="超时锁屏">
        <el-select v-model="form.timeout" placeholder="请选择超时时间">
          <el-option v-for="item in locktype" :key="item.value" :label="item.label" :value="item.value"></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="密码" v-if="form.timeout > 0">
        <el-input type="password" v-model="form.password" placeholder="请输入密码"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitForm">设置</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useSettingsStore } from "@/stores/settings";
import { successMsg } from "@/utils/msg";
import { md5 } from 'js-md5';
const store = useSettingsStore();
const form = ref({...store.config.lock});
form.value.password = ''

const locktype = ref([
  {
    label: '永不',
    value: 0
  },
  {
    label: '30秒',
    value: 30
  },
  {
    label: '1分钟',
    value: 60
  },
  {
    label: '5分钟',
    value: 300
  },
  {
    label: '15分钟',
    value: 900
  }
])

function submitForm() {
  if(form.value.password != '') form.value.password = md5(form.value.password)
  form.value.activeTime = new Date().getTime() + form.value.timeout*1000;
  store.setConfig('lock',form.value);
  successMsg('设置成功');
}
</script>

<style scoped>
.lock-screen {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}

.lock-form {
  width: 300px;
}

.el-button {
  width: 100%;
}
</style>
