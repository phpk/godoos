
<template>
    <div class="custom-captcha-dialog">
			<div class="overlay" @click.self="close"></div>
			<div class="dialog-content">
      <gocaptcha-slide
          :data="handler.data"
          :scope="0"
          :events="{
          close: close,
          refresh: handler.refreshEvent,
          confirm: handler.confirmEvent,
        }"
          ref="domRef"
      />
    </div>
  </div>
  </template>
  
  <script setup lang="ts">
  import {useHandler} from '@/api/captcha'
  import {ref, watch} from 'vue'
  
  const emit = defineEmits(['onSuccess','onCancel'])
  
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
    emit("onCancel",true)
    state.popoverVisible = false
  }
  
  defineExpose({open, close})
  
  </script>
  <style>
.custom-captcha-dialog {
	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
  overflow: hidden;

	z-index: 99;
	display: flex;
	align-items: center;
	justify-content: center;
}

.overlay {
	position: fixed;
	top: 0;
	left: 0;
	width: 100vw;
	height: 100vh;
  overflow: hidden;
	background-color: rgba(0, 0, 0, 0.5);
	z-index: 9998;
}

.dialog-content {
	padding: 20px;
	z-index: 9998;
	width: 350px;
}
</style>