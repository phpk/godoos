<template>
    <div class="date-time-container" @click="showDatePicker">
        <div class="date">{{ currentDate }}</div>
        <div class="time-weekday">
            <div class="time">{{ currentTime }}</div>
            <div class="weekday">{{ currentWeekday }}</div>
        </div>
    </div>
    <Calendar v-if="isShowTimeBox" class="date-pop" />
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'

const currentTime = ref('')
const currentDate = ref('')
const currentWeekday = ref('')
let timeInterval: any
const isShowTimeBox = ref(false)
const showDatePicker = () => {
    isShowTimeBox.value = !isShowTimeBox.value
}
const updateDateTime = () => {
    const now = new Date()
    currentTime.value = now.toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit', second: '2-digit', hour12: false })
    currentDate.value = now.toLocaleDateString('zh-CN', { year: 'numeric', month: 'numeric', day: 'numeric' })
    currentWeekday.value = now.toLocaleDateString('zh-CN', { weekday: 'long' })
}

onMounted(() => {
    timeInterval = setInterval(updateDateTime, 1000)
    updateDateTime() // 初始化时间
})

onUnmounted(() => {
    clearInterval(timeInterval)
})
</script>

<style scoped>
.date-time-container {
    display: flex;
    align-items: center;
    gap: 12px;
}

.date {
    font-size: 14px;
}

.time-weekday {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
}

.time {
    font-size: 14px;
}

.weekday {
    font-size: 12px;
    color: #666;
}
.date-pop {
  position: absolute;
  bottom: 48px;
  right: 0;
  width: 450px;
  height: 465px;
  background-color: #fff; /* 更接近Win10的背景颜色 */
  border: 1px solid #e5e5e5; /* 边框颜色调整 */
  user-select: none;
  box-sizing: border-box;
  z-index: 100;
}
</style>