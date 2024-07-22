<template>
  <div class="screen" @contextmenu.prevent ref="screenref" :style="rootState?.options?.rootStyle">
    <template v-if="rootState.state == SystemStateEnum.close">
      <CloseDesktop></CloseDesktop>
    </template>
    <template v-else-if="rootState.state == SystemStateEnum.opening">
      <OpeningDesktop></OpeningDesktop>
    </template>
    <template v-else-if="rootState.state == SystemStateEnum.open || rootState.state == SystemStateEnum.lock">
      <Transition name="moveup">
        <div class="login" v-if="rootState.state == SystemStateEnum.lock">
          <LockDesktop> </LockDesktop>
        </div>
      </Transition>
      <Transition name="fadeout">
        <DesktopBackground v-if="rootState.state == SystemStateEnum.lock" class="mask"></DesktopBackground>
      </Transition>
      <Desktop v-if="rootState.state == SystemStateEnum.open"></Desktop>
    </template>
  </div>
</template>

<script lang="ts" setup>

import { SystemStateEnum } from '@/system/type/enum';
import { useSystem } from '@/system';
import { onMounted, ref } from 'vue';
import { RootState } from '@/system/root';
const screenref = ref();
defineProps<{
  rootState: RootState;
}>();
onMounted(() => {
  useSystem().rootRef = screenref.value;
});
</script>
<style lang="scss" scoped>
@import '@/assets/root.scss';

.screen {
  position: relative;
  width: 100vw;
  height: 100vh;
  overflow: hidden;
}
</style>
