<template>
    <div class="icon-container">
        <Suspense>
            <icon :name="dealIcon(file)" :size="36" />
        </Suspense>
            <icon v-if="file.isPwd" class="lock-img" name="pwdbox" :size="12"></icon>
            <icon v-if="file.isFavorite" class="fav-img" name="xingxing" :size="12"></icon>
            <icon v-if="file.isShare" class="share-img" name="gallery" :size="12"></icon>
            <icon v-if="file.ext == 'ln'" class="ln-img" name="ink" :size="12"></icon>
        
    </div>
</template>
<script setup lang="ts">
import { dealIcon } from "@/utils/icon";
const props = defineProps({
    file: {
        type: Object,
        required: true
    },
})
const file = props.file;

</script>
<style lang="scss" scoped>
.icon-container {
    display: flex;
    justify-content: center;
    align-items: center;
    position: relative; // 添加相对定位

    svg {
        width: 100%;
        height: 100%;
        object-fit: contain;
    }

    .ln-img,
    .lock-img,
    .fav-img,
    .share-img {
        position: absolute; // 保留绝对定位，但调整位置
    }

    .ln-img {
        bottom: 0;
        right: 0;
    }

    .lock-img {
        bottom: 0;
        left: 0; // 调整位置
    }

    .fav-img {
        top: 0; // 调整位置
        left: 0;
    }

    .share-img {
        top: 0; // 调整位置
        right: 0; // 调整位置
    }
}
</style>