<template>
    <div class="container">
        <div class="nav">
            <ul>
                <li v-for="(item, index) in items" :key="index" @click="activeIndex = index"
                    :class="{ active: index === activeIndex }">
                    {{ item }}
                </li>
            </ul>
        </div>
        <div class="setting">
            <div v-if="0 === activeIndex">
                <LocalNas />
            </div>
            <div v-if="1 === activeIndex">
                <NasClient />
            </div>
            <div v-if="2 === activeIndex">
                <WebDavClient />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref } from "vue";
const items = ["服务端配置", "NAS客户端", "WebDav客户端"];
const activeIndex = ref(0);
</script>

<style scoped>
@import "./setStyle.css";
</style>