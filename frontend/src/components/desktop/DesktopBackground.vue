<template>
  <div class="background" draggable="false" :style="{ backgroundColor: background }">
    <template v-if="backgroundType === 'image'">
      <Transition name="fade">
        <img
          v-show="loaded"
          draggable="false"
          @load="imgload"
          class="background_load"
          :src="background"
        />
      </Transition>
    </template>
  </div>
</template>
<script lang="ts" setup>
import { onMounted, ref, watch } from "vue";
import { useSystem } from "@/system";
const rootState = useSystem()._rootState;
const backgroundType = ref("color");
const background:any = ref("#3A98CE");
const loaded = ref(false);

function imgload() {
  loaded.value = true;
}
onMounted(() => {
  refershBack(rootState.options.background);
});
watch(rootState.options, (nv) => {
  refershBack(nv.background);
});
function refershBack(val: string | undefined) {
  background.value = val || "#3A98CE";
  //console.log(background.value)
  if (background.value || background.value.startsWith("/image/")) {
    backgroundType.value = "image";
  } else {
    backgroundType.value = "color";
  }
}
</script>
<style lang="scss" scoped>
.background {
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.704);
  user-select: none;
  .background_unload {
    width: 100%;
    height: 100%;
  }

  .background_load {
    width: 100%;
    height: 100%;
    object-fit: cover;
    user-select: none;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: all 0.5s ease;
}

.fade-enter-to,
.fade-leave-from {
  opacity: 0.8;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
