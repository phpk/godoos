<template>
  <div class="appicon-group">
    <template v-for="node in treeRoot">
      <AppIcon
        :key="node.id"
        v-if="node.windowInfo.isCreated && !node.windowInfo.skipTaskbar"
        :window-node="node"
      />
    </template>
  </div>
</template>
<script lang="ts" setup>
import { useSystem } from "@/system";
const rootState = useSystem()._rootState;
const treeRoot = rootState.windowOrder;
</script>
<style lang="scss" scoped>
.appicon-group {
  display: flex;
  flex-direction: row;
  flex-wrap: nowrap;
  flex-grow: 1;
  overflow-x: auto;
  overflow: hidden;
  scrollbar-width: none;
  z-index: 41;
  background-color: inherit;

  &::-webkit-scrollbar {
    display: none;
  }
}
</style>
