<template>
    <div class="bottom-bar">
        <MobileApp v-for="item in bottomBarList" :key="item.name" :item="item"></MobileApp>
    </div>
</template>

<script setup lang="ts">
import { ref,watch } from 'vue';

const props = defineProps(['list']);
const bottomBarList = ref<any[]>([]);

watch(props.list,() => {
  console.log(props.list, '~~~~~~~');

  if (Array.isArray(props.list) && props.list.length > 0) {
    props.list.forEach((item: any) => {
      console.log(item.name.includes('computer'), 'item.name=======');
      if (item.name.includes('computer') || item.name.includes('localchat')) {
        bottomBarList.value.push(item);
      }
    });
    console.log(bottomBarList.value, 'bottomBarList============================================');
  } else {
    console.log('props.list is empty or not an array');
  }
},{
    immediate:true
});
</script>

<style lang="scss" scoped>
.bottom-bar{
    display: flex;
    justify-content: space-evenly;
    position: absolute;
    width: vw(340);
    left:50%;
    transform: translateX(-50%);
    bottom: vh(15);
    height: vh(80);
    border-radius: vw(50);
    color: #e5e3e3d5;
    background-color: rgba(255,255,255,.2);
    backdrop-filter: blur(15px);
}
</style>