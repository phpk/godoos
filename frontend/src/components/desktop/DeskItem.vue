<template>
  <div class="desk-group">
    <FileList
      :on-chosen="props.onChosen"
      :on-open="openapp"
      :file-list="appList"
    ></FileList>
  </div>
</template>
<script lang="ts" setup>
import { mountEvent } from "@/system/event";
import { useSystem } from "@/system/index.ts";
import { useAppOpen } from "@/hook/useAppOpen";
import { onMounted } from "vue";

const { openapp, appList } = useAppOpen("apps");
const props = defineProps({
  onChosen: {
    type: Function,
    required: true,
  },
});
onMounted(() => {
  mountEvent("file.props.edit", async () => {
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
  // 应用镂空效果
  color: transparent; /* 文字颜色设为透明 */
  text-shadow: 0 0 0.5px white, 0 0 1px black, 0 0 2px rgba(0, 0, 0, 2); /* 多层阴影 */

  // 重置子元素的默认样式
  > * {
    color: inherit; /* 继承颜色 */
    text-shadow: inherit; /* 继承描边效果 */
    font-size: 0.8rem;
  }
}
</style>
