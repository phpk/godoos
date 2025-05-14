<template>
  <div class="start-menu" v-if="desktopStore.isStartMenuOpen"  @mouseleave="desktopStore.isStartMenuOpen = false">
    <div class="start-menu-item">
      <div class="start-menu-left">
        <div class="start-menu-left-inner">
          <div
            class="start-menu-button glowing"
            @click.stop="handleLogOut"
          >
            <div class="start-menu-button-img">
              <el-icon><SwitchButton /></el-icon>
            </div>
            <div class="start-menu-button-title">
              退出
            </div>
          </div>
          <div
            class="start-menu-button glowing"
            @click.stop="handleRecover"
          >
            <div class="start-menu-button-img">
              <el-icon><Compass /></el-icon>
            </div>
            <div class="start-menu-button-title">
              重置
            </div>
          </div>
          <div
            class="start-menu-button glowing"
            @click.stop="handleSetting"
          >
            <div class="start-menu-button-img">
              <el-icon><Setting /></el-icon>
            </div>
            <div class="start-menu-button-title">
              设置
            </div>
          </div>
        </div>
      </div>
    </div>
    <div class="start-menu-item">
      <div class="start-menu-right-group scroll-bar">
        <div
          @click.stop="handle(item)"
          class="start-menu-right-item glowing"
          :style="{
            animationDelay: `${Math.floor(index / 4) * 0.02}s`,
            animationDuration: `${Math.floor(index / 4) * 0.04 + 0.1}s`,
          }"
          v-for="(item, index) in desktopStore.menuList"
          :key="index"
        >
          <icon class="start-menu-right-item-img" :name="dealIcon(item)" :size="32"/>
          <span class="start-menu-right-item-title">{{ t(item.title) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLoginStore } from '@/stores/login'
import { useDesktopStore } from '@/stores/desktop'

import {dealIcon} from '@/utils/icon'
import { clear } from '@/api/files'
import { t } from '@/i18n'
import { useRouter } from 'vue-router'
import { confirmMsg } from '@/utils/msg'

const router = useRouter()
const loginStore = useLoginStore();
const desktopStore = useDesktopStore();

async function handleRecover() {
  desktopStore.isStartMenuOpen = false
  const res = await confirmMsg('确定要重置桌面吗？这将删除系统盘内所有文件和文件夹。', '提示')
  if (res) {
    await clear()
    //await desktopStore.initDesktop()
    localStorage.clear()
    window.indexedDB.deleteDatabase("GodoDatabase");
    window.location.reload()
  }
}
function handleLogOut() {
  desktopStore.isStartMenuOpen = false
  loginStore.loginOut()
}
function handleSetting() {
  desktopStore.isStartMenuOpen = false
  router.push('/setting')
}

function handle(item: any) {
  desktopStore.isStartMenuOpen = false
  router.push('/'+item.title)
}

</script>

<style lang="scss" scoped>
@use '@/styles/startmenu.scss';
</style>