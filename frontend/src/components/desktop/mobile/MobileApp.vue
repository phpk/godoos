<template>
    <div class="file-item">
        <FileIcon :file="item" class="icon"</FileIcon>
        <span class="title">{{ getName(item) }}</span>
    </div>
</template>

<script setup lang="ts">
import {
	basename,
} from "@/system/index.ts";
import { dealSystemName, t } from "@/i18n";
defineProps(['item'])
  function getName(item: any) {
	const name = dealSystemName(basename(item.path));
	// console.log(name)
	// console.log(item.path)
	if (name.endsWith(".exe")) {
		return t(name.replace(".exe", ""));
	} else {
		return name;
	}
}
</script>

<style lang="scss" scoped>
.file-item{
    height: vh(70);
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    .icon{
        width: vw(40);
        height: vh(40);
    }
    .title{
        font-size: vw(12);
        text-align: center;
    }
}
</style>