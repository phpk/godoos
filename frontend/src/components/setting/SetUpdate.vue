<template>
  <div class="container">
    <div class="nav">
      <ul>
        <li
          v-for="(item, index) in items"
          :key="index"
          @click="selectItem(index)"
          :class="{ active: index === activeIndex }"
        >
          {{ item }}
        </li>
      </ul>
    </div>
    <div class="setting">
      <!-- <div v-if="0 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">
            {{ t('system.backup.and.import') }}
          </h1>
        </div>
        <div class="setting-item">
          <label>
            {{ t('export.system.status') }}
          </label>
          <WinButton @click="handleClick(0)">
            {{ t('export') }}
          </WinButton>
        </div>
        <div class="setting-item">
          <label>
            {{ t('import.status.file') }}
          </label>
          <textarea v-model="inputConfig" type="text"></textarea>
        </div>
        <div class="setting-item">
          <label></label>
          <WinButton @click="handleClick(1)">
            {{ t('import') }}
          </WinButton>
        </div>
      </div> -->
      <div v-if="0 === activeIndex">
        <div class="setting-item">
          <h1 class="setting-title">
            {{ t("system.version") }}
          </h1>
        </div>
        <div class="setting-item">
          <label>
            {{ t("version") }}
          </label>
          <span>{{ version }}</span>
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
import { ref } from "vue";
import { t } from "@/i18n";
import { getSystemKey } from "@/system/config";
// const system = useSystem();
const items = [
  t("version"), // '版本',
];
const version = getSystemKey("version");
const activeIndex = ref(0);
// const inputConfig = ref('');
const selectItem = (index: number) => {
  activeIndex.value = index;
};

</script>

<style scoped>
@import "./setStyle.css";
</style>
