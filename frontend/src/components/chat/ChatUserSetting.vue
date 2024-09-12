<template>
    <el-form :model="form" label-width="auto" style="padding: 30px;">
      <el-form-item label="昵称">
        <el-input v-model="form.nickname" />
      </el-form-item>
      <el-form-item>
      <el-button type="primary" @click="onSubmit">保存</el-button>
    </el-form-item>
    </el-form>
  </template>
  
  <script lang="ts" setup>
  import { reactive } from 'vue'
  
  // do not use same name with ref
  const form = reactive({
    nickname: '',
  })
  
  const onSubmit = () => {
    console.log('submit!')
  }
  </script>