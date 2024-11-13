<template>
  <div class="editor-header">
    <div class="left">
      <div class="group-menu-item">
        <div
          class="menu-item"
          v-tooltip="'幻灯片放映'"
          @click="enterScreening()"
        >
          <IconPpt class="icon" />
        </div>

        <Popover trigger="click" center>
          <template #content>
            <PopoverMenuItem @click="enterScreeningFromStart()"
              >从头开始</PopoverMenuItem
            >
            <PopoverMenuItem @click="enterScreening()"
              >从当前页开始</PopoverMenuItem
            >
          </template>
          <div class="arrow-btn"><IconDown class="arrow" /></div>
        </Popover>
      </div>
      <div
        class="menu-item"
        v-tooltip="'清空'"
        @click="
          resetSlides();
          mainMenuVisible = false;
        "
      >
        <IconDelete class="icon" />
      </div>
      <div class="menu-item" v-tooltip="'导入pptist'">
        <FileInput
          accept=".pptist"
          @change="
            (files) => {
              importSpecificFile(files);
            }
          "
        >
          <IconImportPPtist class="icon" />
        </FileInput>
      </div>
      <div class="menu-item" v-tooltip="'导入pptx'">
        <FileInput
          accept="application/vnd.openxmlformats-officedocument.presentationml.presentation"
          @change="
            (files) => {
              importPPTXFile(files);
            }
          "
        >
          <IconUpload class="icon" />
        </FileInput>
      </div>
      <div class="menu-item" v-tooltip="'导出pptx'" @click="exportPPT()">
        <IconExport class="icon" />
      </div>
      <div
        class="menu-item"
        v-tooltip="'导出图片'"
        @click="setDialogForExport('image')"
      >
        <IconDownPic class="icon" />
      </div>
      <a
        class="github-link"
        v-tooltip="'Copyright © 2020-PRESENT pipipi-pikachu'"
        href="https://github.com/pipipi-pikachu/PPTist"
        target="_blank"
      >
        <div class="menu-item"><IconGithub class="icon" /></div>
      </a>
      <div
        class="menu-item"
        v-tooltip="'快捷键'"
        @click="
          mainMenuVisible = false;
          hotkeyDrawerVisible = true;
        "
      >
        <IconHelp class="icon" />
      </div>
    </div>

    <div class="right">
      <div class="title">
        <Input
          class="title-input"
          ref="titleInputRef"
          v-model:value="titleValue"
          @blur="handleUpdateTitle()"
          v-if="editingTitle"
        ></Input>
        <div class="title-text" @click="startEditTitle()" :title="title" v-else>
          {{ title }}
        </div>
      </div>
      <div
        class="menu-item"
        v-tooltip="'保存'"
        @click="saveData"
      >
        <IconSave class="icon" />
      </div>
    </div>

    <Drawer
      :width="320"
      v-model:visible="hotkeyDrawerVisible"
      placement="right"
    >
      <HotkeyDoc />
    </Drawer>

    <FullscreenSpin :loading="exporting" tip="正在导入..." />
  </div>
</template>

<script lang="ts" setup>
import { nextTick, ref } from 'vue'
import { storeToRefs } from 'pinia'
import { useMainStore, useSlidesStore } from '@/store'
import useScreening from '@/hooks/useScreening'
import useImport from '@/hooks/useImport'
import useExport from '@/hooks/useExport'
import useSlideHandler from '@/hooks/useSlideHandler'
import type { DialogForExportTypes } from '@/types/export'

import HotkeyDoc from './HotkeyDoc.vue'
import FileInput from '@/components/FileInput.vue'
import FullscreenSpin from '@/components/FullscreenSpin.vue'
import Drawer from '@/components/Drawer.vue'
import Input from '@/components/Input.vue'
import Popover from '@/components/Popover.vue'
import PopoverMenuItem from '@/components/PopoverMenuItem.vue'
import message from '@/utils/message'
import { md5 } from 'js-md5'

const mainStore = useMainStore()
const slidesStore = useSlidesStore()
const { title } = storeToRefs(slidesStore)
const { enterScreening, enterScreeningFromStart } = useScreening()
const { importSpecificFile, importPPTXFile, exporting, importData, isBase64, base64ToBuffer, arrayBufferToBase64 } = useImport()
const { exportPPT, exportBuffer } = useExport()
const { resetSlides } = useSlideHandler()

const mainMenuVisible = ref(false)
const hotkeyDrawerVisible = ref(false)
const editingTitle = ref(false)
const titleInputRef = ref<InstanceType<typeof Input>>()
const titleValue = ref('')

const startEditTitle = () => {
  titleValue.value = title.value
  editingTitle.value = true
  nextTick(() => titleInputRef.value?.focus())
}

const handleUpdateTitle = () => {
  slidesStore.setTitle(titleValue.value)
  editingTitle.value = false
}

const setDialogForExport = (type: DialogForExportTypes) => {
  mainStore.setDialogForExport(type)
  mainMenuVisible.value = false
}

const eventHandler = (e: any) => {
  const eventData = e.data
  if (eventData.type === 'init') {
    const data = eventData.data
    if (!data || !data.title) {
      return
    }
    titleValue.value = data.title.substring(0, data.title.lastIndexOf('.'))
    handleUpdateTitle()
    if (data.content) {
      if (data.content instanceof ArrayBuffer) {
        const fileName = 'pptx_' + md5(arrayBufferToBase64(data.content))
        const loc = localStorage.getItem(fileName)
        // console.log(loc)
        if (loc) {
          const slides = JSON.parse(loc)
          slidesStore.setSlides(slides)
        }
        else {
          importData(data.content)
        }
      }
      else {
        if (isBase64(data.content)) {
          data.content = base64ToBuffer(data.content)
          importData(data.content)
        }       
      }
      
    }
  }
  if (eventData.type === 'start') {
    if (eventData.title) {
      titleValue.value = eventData.title.substring(
        0,
        eventData.title.lastIndexOf('.')
      )
      handleUpdateTitle()
    }
  }
}
const debouncedHandleKeyDown = (event: KeyboardEvent) => {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "s") {
    event.stopPropagation() // 先阻止事件冒泡
    event.preventDefault() // 再阻止默认行为
    saveData()
  }
}
const saveData = async () => {
  const baseTitle = titleValue.value || title.value
  if (!baseTitle || baseTitle === '') {
    message.error('名称不能为空')
    titleInputRef.value?.focus()
    return
  }
  const buffer = await exportBuffer()
  const save = {
    data: JSON.stringify({ content: buffer, title: baseTitle }),
    type: 'exportPPTX',
  }
  window.parent.postMessage(save, '*')
}
document.addEventListener('keydown', debouncedHandleKeyDown);
window.addEventListener('load', () => {
  window.parent.postMessage({ type: 'initSuccess' }, "*");
  window.addEventListener('message', eventHandler);
})
window.addEventListener('unload', () => {
  window.removeEventListener('message', eventHandler);
  document.removeEventListener('keydown', debouncedHandleKeyDown);
})

</script>

<style lang="scss" scoped>
.editor-header {
  background-color: #fff;
  user-select: none;
  border-bottom: 1px solid $borderColor;
  display: flex;
  justify-content: space-between;
  padding: 0 5px;
}
.left,
.right {
  display: flex;
  justify-content: center;
  align-items: center;
}
.menu-item {
  height: 30px;
  display: flex;
  justify-content: center;
  align-items: center;
  font-size: 14px;
  padding: 0 10px;
  border-radius: $borderRadius;
  cursor: pointer;

  .icon {
    font-size: 18px;
    color: #666;
  }

  &:hover {
    background-color: #f1f1f1;
  }
}
.group-menu-item {
  height: 30px;
  display: flex;
  margin: 0 8px;
  padding: 0 2px;
  border-radius: $borderRadius;

  &:hover {
    background-color: #f1f1f1;
  }

  .menu-item {
    padding: 0 3px;
  }
  .arrow-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    cursor: pointer;
  }
}
.title {
  height: 32px;
  margin-right: 20px;
  font-size: 13px;

  .title-input {
    width: 200px;
    height: 100%;
    text-align: center;
    padding-left: 0;
    padding-right: 0;
  }
  .title-text {
    min-width: 20px;
    max-width: 400px;
    line-height: 32px;
    padding: 0 6px;
    border-radius: $borderRadius;
    cursor: pointer;

    @include ellipsis-oneline();

    &:hover {
      background-color: #f1f1f1;
    }
  }
}
.github-link {
  display: inline-block;
  height: 30px;
}
</style>