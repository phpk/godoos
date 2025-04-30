// src/stores/selectionBox.ts
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useClickingStore } from '@/stores/clicking';
import { useDragFilesStore } from '@/stores/dragfiles';
export const useSelectionBoxStore = defineStore('selectionBox', () => {
  const selectionBoxRef = ref<HTMLDivElement | null>(null);
  const startX = ref(0);
  const startY = ref(0);
  const endX = ref(0);
  const endY = ref(0);
  const isDragging = ref(false);
  const selectedIcons = ref<any[]>([]);
  let mousedownTimeout: any;
  let desktopAreaRef: HTMLDivElement | null = null;
  const fileList: any = ref([]);
  const clickingStore = useClickingStore();
  const dragFilesStore = useDragFilesStore();

  function handleMouseDown(event: MouseEvent, desktopArea: any, files: any[]) {
    //console.log(dragFilesStore.isDragFile)
    if (dragFilesStore.isDragFile) return; // 如果正在拖动文件，则不处理选择框
    if (event.button !== 0) return; // 只处理左键
    desktopAreaRef = desktopArea;
    mousedownTimeout = setTimeout(() => {
      if (!desktopAreaRef) return;

      fileList.value = files;
      isDragging.value = true;
      const offset = getOffset(desktopAreaRef!);
      startX.value = event.clientX - offset.left;
      startY.value = event.clientY - offset.top;
      endX.value = startX.value;
      endY.value = startY.value;
      if (!(event.ctrlKey || event.metaKey)) {
        clickingStore.clickedIcons = []; // 如果没有按住 Control 或 Command 键，清除现有的选中状态
      }
      createSelectionBox();
      //mousedownTimeout = null; // 清除定时器，防止在拖动过程中再次触发
      desktopAreaRef!.addEventListener('mousemove', handleMouseMove)
      desktopAreaRef!.addEventListener('mouseup', handleMouseUp)
    }, 100);

  }

  function handleMouseMove(event: MouseEvent) {
    //if (dragFilesStore.isDragFile) return;
    if (!isDragging.value) return;
    const offset = getOffset(desktopAreaRef!);
    endX.value = event.clientX - offset.left;
    endY.value = event.clientY - offset.top;
    updateSelectionBox();
  }

  function handleMouseUp() {
    //if (dragFilesStore.isDragFile) return;
    if (!isDragging.value) return;
    isDragging.value = false;
    updateSelectedIcons();
    removeSelectionBox();
    clearTimeout(mousedownTimeout);
    //removeEvents(desktopAreaRef!)
  }

  function createSelectionBox() {
    const selectionBox = document.createElement('div');
    selectionBox.className = 'rect selection-box';
    selectionBox.style.position = 'absolute';
    selectionBox.style.border = '1px dashed #007bff';
    selectionBox.style.background = 'rgba(0, 123, 255, 0.1)';
    selectionBox.style.pointerEvents = 'none';
    selectionBox.style.zIndex = '9999';
    desktopAreaRef?.appendChild(selectionBox);
    selectionBoxRef.value = selectionBox;
  }

  function updateSelectionBox() {
    if (!selectionBoxRef.value || !desktopAreaRef) return;
    const left = Math.min(startX.value, endX.value);
    const top = Math.min(startY.value, endY.value);
    const width = Math.abs(endX.value - startX.value);
    const height = Math.abs(endY.value - startY.value);

    selectionBoxRef.value.style.left = `${left}px`;
    selectionBoxRef.value.style.top = `${top}px`;
    selectionBoxRef.value.style.width = `${width}px`;
    selectionBoxRef.value.style.height = `${height}px`;
  }

  function removeSelectionBox() {
    const selectionBoxes = document.querySelectorAll('.selection-box');
    selectionBoxes.forEach(box => {
      if (box.parentElement) {
        box.parentElement.removeChild(box);
      }
    });
    selectionBoxRef.value = null;
  }

  function getOffset(element: HTMLElement) {
    const rect = element.getBoundingClientRect();
    return {
      top: rect.top + window.scrollY,
      left: rect.left + window.scrollX
    };
  }

  function updateSelectedIcons() {
    if (!desktopAreaRef) return;
    const selectionRect = {
      left: Math.min(startX.value, endX.value),
      top: Math.min(startY.value, endY.value),
      right: Math.max(startX.value, endX.value),
      bottom: Math.max(startY.value, endY.value),
    };
    const iconElements = desktopAreaRef.querySelectorAll(".desktop-icon") as NodeListOf<HTMLElement>;
    if (!iconElements) return;


    const desktopOffset = getOffset(desktopAreaRef!);
    const checkedIcons: any = [];
    iconElements.forEach((iconElement, index) => {
      const rect = iconElement.getBoundingClientRect();
      const iconRect = {
        left: rect.left - desktopOffset.left,
        top: rect.top - desktopOffset.top,
        right: rect.right - desktopOffset.left,
        bottom: rect.bottom - desktopOffset.top,
      };

      if (
        iconRect.left >= selectionRect.left &&
        iconRect.right <= selectionRect.right &&
        iconRect.top >= selectionRect.top &&
        iconRect.bottom <= selectionRect.bottom &&
        fileList.value[index]?.path
      ) {
        checkedIcons.push(fileList.value[index]?.path);
      }
    });

    if (checkedIcons.length > 0) {
      //selectedIcons.value = []; // 清空选中的图标
      selectedIcons.value = checkedIcons;
      clickingStore.setClickIcons(checkedIcons);
      //console.log('Assigned to clickingStore.clickedIcons:', clickingStore.clickedIcons); // 调试信息
    }
  }

  return {
    selectionBoxRef,
    startX,
    startY,
    endX,
    endY,
    isDragging,
    selectedIcons,
    fileList,
    //addEvents,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    removeSelectionBox,
    //removeEvents
  };
});