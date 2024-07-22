<template>
  <div class="taskbar">
    <StartMenu></StartMenu>
    <AppIconGroup></AppIconGroup>
    <StateIcon></StateIcon>
  </div>
</template>
<style lang="scss" scoped>
.taskbar {
  width: 100%;
  height: var(--task-bar-height);
  background-color: var(--color-task-bar);
  display: flex;
  position: relative;
}
</style>
