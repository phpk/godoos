<template>
    <el-tabs v-model="tabsStore.activeId" type="card" class="demo-tabs" closable @tab-remove="tabsStore.removeTab">
        <el-tab-pane v-for="item in tabsStore.tabsList" :key="item.id" :label="item.title" :name="item.id">
            <!-- <template #label>
                <span class="custom-tabs-label">
                    <icon :name="item.icon" size="20" />
                    <span>{{ item.title }}</span>
                </span>
            </template> -->
            <!-- <div class="window-content">
                <IframeApp v-if="typeof item.component === 'string'" :win="item" />
                <component v-else :is="item.component" :win="item" />
            </div> -->
        </el-tab-pane>
    </el-tabs>
</template>
<script setup lang="ts">
import { useTabsStore } from '@/stores/tabs';
import { onMounted } from 'vue';
const tabsStore = useTabsStore();
const props = defineProps({
    win: {
        type: Object,
        required: true,
    },
});
const win = props.win;
onMounted(() => {
    tabsStore.addTab(win);
});
</script>
<style scoped>
::deep(.el-tabs__header){
    margin: 0px;
    padding:0px;
}
</style>