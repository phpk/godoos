<template>
  <div class="outer">
    <input class="win-input" v-model="type" />
    <WinButton @click="confirm">{{ t('confirm') }}</WinButton>
  </div>
</template>
<script setup lang="ts">
import { inject, ref } from 'vue';
// import { OsFileWithoutContent } from '@/system/core/FileSystem';
// import { useSystem } from '../system';
// import { BrowserWindow } from '../window/BrowserWindow';
import { emitEvent } from '../../system/event';
// import { basename, extname } from '../core/Path';
// import { t } from '@/i18n';
import { t,basename,extname,BrowserWindow,useSystem,OsFileWithoutContent } from '../../system';
const browserWindow: BrowserWindow = inject('browserWindow')!;
const fileBaseName = basename((browserWindow.config.content as OsFileWithoutContent).path);
const type = ref(extname(fileBaseName));

function confirm() {
  useSystem()
    ?.fs.rename(
      browserWindow.config.content.path,
      browserWindow.config.content.path.replace(
        fileBaseName,
        fileBaseName.replace(extname(fileBaseName), type.value)
      )
    )
    .then(() => {
      emitEvent('file.props.edit');
      browserWindow.close();
    });
}
</script>
<style lang="scss" scoped>
.outer {
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  height: 100%;

  .win-input {
    font-size: 20px;
    width: 200px;
    height: 40px;
    margin-bottom: 40px;
    outline: none;
    border: 1px solid black;

    &:focus {
      border: 1px solid var(--color-blue);
    }
  }
}
</style>
@/system/core/FileSystem
