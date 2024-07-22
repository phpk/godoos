<template>
  <div class="desk-group">
    <FileList :on-chosen="props.onChosen" :on-open="openapp" :file-list="appList"></FileList>
  </div>
</template>
<script lang="ts" setup>
import { mountEvent } from '@/system/event';
import { useSystem } from '@/system/index.ts';
import { useAppOpen } from '@/hook/useAppOpen';
import { onMounted } from 'vue';

const { openapp, appList } = useAppOpen('apps');
const props = defineProps({
  onChosen: {
    type: Function,
    required: true,
  },
});
onMounted(() => {
  mountEvent('file.props.edit', async () => {
    useSystem().initAppList();
  });
});
</script>
<style lang="scss" scoped>
.desk-group {
  display: flex;
  flex-direction: column;
  flex-wrap: wrap;
  height: 100%;
}
</style>
