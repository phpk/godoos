
<template>
    
      <gocaptcha-slide
          :data="handler.data"
          :scope="0"
          :events="{
          close: handler.closeEvent,
          refresh: handler.refreshEvent,
          confirm: handler.confirmEvent,
        }"
          ref="domRef"
      />
  </template>
  
  <script setup lang="ts">
  import {useHandler} from '@/api/captcha'
  import {ref, watch} from 'vue'
  
  const emit = defineEmits(['onSuccess'])
  
  const domRef = ref(null)
  const handler = useHandler(domRef)
  const state = handler.state
  watch(() => state.type, (n) => {
        console.log(n)
        if (n === "success") {
          emit("onSuccess",true)
          state.type = "default"
        }
      }
  )
  const open = () => {
    state.popoverVisible = true
  }
  const close = () => {
    state.popoverVisible = false
  }
  
  defineExpose({open, close})
  
  </script>
  <style></style>