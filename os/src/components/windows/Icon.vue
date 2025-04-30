<template>
    <svg aria-hidden="true" :style="{ width: size + 'px', height: size + 'px' }">
        <use :xlink:href="iconName"></use>
    </svg>
</template>

<script setup lang="ts">
import '@/utils/iconfont.js'
import { computed } from 'vue';

const props = defineProps<{ name: string, size: number }>()
const iconName = computed(() => `#icon-${props.name}`)
</script>
