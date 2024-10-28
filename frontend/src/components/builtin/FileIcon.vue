<template>
  <div class="icon">
    <Suspense>
      <FileIconImg v-if="isSvg === true" :file="file" :icon="icon" />
      <FileIconIs v-else :file="file" :icon="icon" />
    </Suspense>
    <div v-if="extname(file?.path || '') === '.ln' || file?.isShare === true" class="ln-img">
      <img :src="lnicon" alt="ln" />
    </div>
  </div>
</template>
<script setup lang="ts">
import lnicon from '@/assets/ln.png';
import { OsFileWithoutContent,extname } from '@/system';
// import { extname } from '../core/Path';
import { ref } from 'vue';
const props:any = defineProps<{
  file?: OsFileWithoutContent | null;
  icon?: string;
}>();
const isSvg = ref(true);

if(props.icon && props.icon.indexOf('.') !== -1){
  isSvg.value = false;
}
if(props.file && props.file.content) {
  if(typeof props.file.content === 'string') {
    const end = props.file.content.split("::").pop()
    if(end && end.indexOf('.') !== -1){
      isSvg.value = false;
    }
  }
  
}

</script>
<style lang="scss" scoped>
.icon {
  width: 100%;
  height: 100%;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;

  img {
    width: 100%;
    height: 100%;
    object-fit: contain;
  }
  .ln-img {
    position: absolute;
    bottom: 0;
    left: 0;
    width: 80%;
    height: 80%;
  }
}
</style>
