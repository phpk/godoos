import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

export const useDraggingStore = defineStore('dragging', () => {
  // 选中的图标
  const selectedIcons = ref<any[]>([]);
  const fileList = ref<any[]>([]);
  const desktopAreaRef = ref<HTMLElement | null>(null);
  // 框选区域相关状态
  const isDragging = ref(false);
  const isDraggingStarted = ref(false); // 新增标志来区分拖拽和单击
  const startX = ref(0);
  const startY = ref(0);
  const endX = ref(0);
  const endY = ref(0);
  

  const setFileList = (filelist: any[]) => {
    fileList.value = filelist;
  };

  // 计算框选区域的样式
  const selectionBoxStyle = computed(() => {
    const left = Math.min(startX.value, endX.value) + 5;
    const top = Math.min(startY.value, endY.value) + 5;
    const width = Math.abs(endX.value - startX.value) + 5;
    const height = Math.abs(endY.value - startY.value) + 5;
    return {
      left,
      top,
      width,
      height,
    };
  });

  // 处理鼠标按下事件
  const handleMouseDown = (event: MouseEvent) => {
    if (event.button !== 0) return; // 只处理左键
    isDragging.value = true;
    isDraggingStarted.value = false; // 初始化拖拽标志
    startX.value = event.clientX;
    startY.value = event.clientY;
    endX.value = startX.value;
    endY.value = startY.value;
    if (!(event.ctrlKey || event.metaKey)) {
      selectedIcons.value = []; // 如果没有按住 Control 或 Command 键，清除现有的选中状态
    }
  };
  
  // 处理鼠标移动事件
  const handleMouseMove = (event: MouseEvent) => {
    if (!isDragging.value) return;
    isDraggingStarted.value = true; // 设置拖拽标志
    endX.value = event.clientX;
    endY.value = event.clientY;
    updateSelectedIcons();
  };

  // 处理鼠标释放事件
  const handleMouseUp = () => {
    if (!isDragging.value) return;
    isDragging.value = false;
    if (isDraggingStarted.value) {
      // 如果是拖拽操作，不处理单击逻辑
      isDraggingStarted.value = false;
      removeEventListeners(); // 移除事件监听器
    }
   
  };

  const updateSelectedIcons = () => {
    if (!desktopAreaRef.value) return;
    const selectionRect = {
      left: Math.min(startX.value, endX.value),
      top: Math.min(startY.value, endY.value),
      right: Math.max(startX.value, endX.value),
      bottom: Math.max(startY.value, endY.value),
    };

    selectedIcons.value = []; // 清空选中的图标

    const iconElements = desktopAreaRef.value.querySelectorAll(".desktop-icon") as NodeListOf<HTMLElement>;
    if (!iconElements) return;
    //console.log(iconElements);
    iconElements.forEach((iconElement, index) => {
      const rect = iconElement.getBoundingClientRect();
      if (
        rect.left >= selectionRect.left &&
        rect.right <= selectionRect.right &&
        rect.top >= selectionRect.top &&
        rect.bottom <= selectionRect.bottom
      ) {
        selectedIcons.value.push(fileList.value[index]?.path || '');
      }
    });
    //console.log(selectedIcons.value);
  };



  const addEventListeners = (desktopArea: HTMLElement | null) => {
    if (!desktopArea) return;
    desktopAreaRef.value = desktopArea;
    desktopArea.classList.add('no-select'); // 添加 no-select 类
    desktopArea.addEventListener('mousedown', handleMouseDown);
    desktopArea.addEventListener('mousemove', handleMouseMove);
    desktopArea.addEventListener('mouseup', handleMouseUp);
    //desktopArea.addEventListener('click', handleClick); // 添加 click 事件监听器
  };

  const removeEventListeners = () => {
    if (!desktopAreaRef.value) return;
    desktopAreaRef.value.classList.remove('no-select'); // 移除 no-select 类
    desktopAreaRef.value.removeEventListener('mousedown', handleMouseDown);
    desktopAreaRef.value.removeEventListener('mousemove', handleMouseMove);
    desktopAreaRef.value.removeEventListener('mouseup', handleMouseUp);
    //desktopAreaRef.value.removeEventListener('click', handleClick); // 移除 click 事件监听器
  };
  // document.addEventListener("copy", function () {
  //   if (selectedIcons.value.length > 0) {
  //     copiedIcons.value = [...selectedIcons.value]; 
  //     console.log('Copied icons:', copiedIcons.value);
  //   }
  // });
  
  return {
    selectedIcons,
    isDragging,
    startX,
    startY,
    endX,
    endY,
    selectionBoxStyle,
    setFileList,
    handleMouseDown,
    handleMouseMove,
    handleMouseUp,
    addEventListeners,
    removeEventListeners,
  };
});