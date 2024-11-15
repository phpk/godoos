<template>
  <div class="version">
    <div class="version-item">
      <div class="version-title">Godo OS</div>
      <div class="version-content">Version: {{ plusVersion }}</div>
      <div class="version-content">
        Home:<a href="https://godoos.com" target="_blank">godoos.com</a>
      </div>
    </div>
  </div>
</template>
<script setup lang="ts">
import { getSystemKey } from "@/system/config";
const plusVersion = getSystemKey("version");
</script>
<style scoped lang="scss">
$primary-color: #0078d4; // Win11的主题蓝色
$neutral-light: #f3f2f1; // 浅中性色
$neutral-dark: #1e1e1e; // 深中性色
$shadow-depth-1: 0 1px 2px 0 rgba(0, 0, 0, 0.2), 0 1px 4px 0 rgba(0, 0, 0, 0.16); // 浅阴影
$shadow-depth-2: 0 2px 4px 0 rgba(0, 0, 0, 0.24), 0 1px 8px 0 rgba(0, 0, 0, 0.16); // 深阴影
$acrylic-opacity: 0.3; // 亚克力效果的透明度

.version {
  width: 100%;
  height: 100%;
  background-color: $neutral-light;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  border-radius: 12px;
  box-shadow: $shadow-depth-2;
  backdrop-filter: blur(10px) brightness(0.95); // 亚克力模糊效果
  user-select: none;
}

.version-item {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  margin: 10px;
  padding: 30px 50px;
  border-radius: 8px;
  background-color: rgba(white, 0.9);
  box-shadow: $shadow-depth-1;
  transition: all 0.2s;

  &:hover {
    box-shadow: $shadow-depth-2;
    transform: translateY(-2px);
  }
}

.version-title {
  font-size: 24px;
  font-weight: bold;
  color: $neutral-dark;
  text-shadow: none;
}

.version-content {
  font-size: 16px;
  margin-top: 10px;
  color: $neutral-dark;
  text-shadow: none;
}
</style>
