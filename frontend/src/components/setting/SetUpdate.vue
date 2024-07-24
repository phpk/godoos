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
// import WinButton from '@/components/win/WinButton.vue';
import { ref } from "vue";
// import { useSystem } from '@/system';
// import { Dialog } from '@/dialog/Dialog';
import { t } from "@/i18n";
import { getSystemKey } from "@/system/config";
// const system = useSystem();
const items = [
  // t('backup'), // '备份',
  t("version"), // '版本',
];
const version = getSystemKey("version");
const activeIndex = ref(0);
// const inputConfig = ref('');
const selectItem = (index: number) => {
  activeIndex.value = index;
};

// async function handleClick(num: number) {
//   if (num === 0) {
//     //导出配置
//     const cfg = await system?.serializeState();
//     try {
//       await navigator.clipboard.writeText(cfg!);
//       Dialog.showMessageBox({
//         title: t('export.config'),
//         message: t('export.success.saved.to.clipboard'),
//         type: 'info',
//         buttons: [t('confirm')],
//       });
//     } catch (err) {
//       Dialog.showMessageBox({
//         title: t('export.config'),
//         message: t('export.failed'),
//         type: 'error',
//         buttons: [t('confirm')],
//       });
//     }
//   } else if (num === 1) {
//     // 导入配置
//     try {
//       const req = await Dialog.showMessageBox({
//         title: t('import.config'),
//         // message: '导入会覆盖现有的文件,是否继续?',
//         message: t('import.config.will.cover.existing.files.continue'),
//         type: 'warning',
//         buttons: [t('confirm'), t('cancel')],
//       });
//       if (req.response === 1) return;
//       await system?.deserializeState(inputConfig.value);
//       setTimeout(() => {
//         system?.reboot();
//       }, 10000);
//       await Dialog.showMessageBox({
//         title: t('import.success'),
//         message: t('import.success.reboot.soon'),
//         type: 'warning',
//         buttons: [t('confirm')],
//       });
//       system?.reboot();
//     } catch (err) {
//       Dialog.showMessageBox({
//         title: t('import.config'),
//         message: t('import.failed'),
//         type: 'error',
//         buttons: [t('confirm')],
//       });
//     }
//   }
// }
</script>

<style scoped>
@import "./setStyle.css";
</style>
