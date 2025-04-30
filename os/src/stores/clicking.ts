// clicking.ts
import { defineStore } from 'pinia';
import { ref } from 'vue';
import { useFileSystemStore } from './filesystem';
import { useSettingsStore } from "@/stores/settings";
export const useClickingStore = defineStore('clicking', () => {
    const clickedIcons = ref<string[]>([]);
    const copiedIcons = ref<string[]>([]); // 用于存储复制的图标路径
    const cutedIcons = ref<string[]>([]); // 用于存储剪切的图标路径
    const fsStore = useFileSystemStore();
    const settingsStore = useSettingsStore();
    const resetClickedIcons = (event: MouseEvent) => {
        const target = event.target as HTMLElement;
        if(!target)return;
        const iconElement = target.closest(".desktop-icon"); // 使用 closest 方法查找最近的匹配类名的父元素
        if (!iconElement) {
            clickedIcons.value = [];
        }
    };
   
    const setClickIcons = (list: any) => {
        clickedIcons.value = list;
    };
    
    const handleClick = (event: MouseEvent, icon: any) => {
        if (event.button !== 0) return; // 只处理左键
        const target = event.target as HTMLElement;
        const iconElement = target.closest(".desktop-icon"); // 使用 closest 方法查找最近的匹配类名的父元素
        if (!iconElement) {
            clickedIcons.value = [];
            return; // 如果没有找到匹配的元素，直接返回
        }
        event.preventDefault(); // 阻止默认的点击事件，例如文本选择等

        const iconPath = icon?.path;
        if (iconPath) {
            const indexInSelected = clickedIcons.value.indexOf(iconPath);
            if (indexInSelected === -1) {
                clickedIcons.value.push(iconPath);
            } else {
                clickedIcons.value.splice(indexInSelected, 1);
            }
            // if (event.ctrlKey || event.metaKey) { // 检查 Control 或 Command 键
            //     if (indexInSelected === -1) {
            //         clickedIcons.value.push(iconPath);
            //     } else {
            //         clickedIcons.value.splice(indexInSelected, 1);
            //     }
            // } else {
            //     // 如果没有按住 Control 或 Command 键，清除现有的选中状态并选择当前点击的图标
            //     clickedIcons.value = [iconPath];
            // }
        }
        //console.log(clickedIcons.value);
    };

    const handleCopy = (event: KeyboardEvent) => {
        if ((event.ctrlKey || event.metaKey) && event.key === 'c') { // 检查 Control 或 Command 键和 C 键
            copiedIcons.value = [...clickedIcons.value]; // 将选中的图标路径存储到 copiedIcons 中
            //console.log('Copied icons:', copiedIcons.value);
            
        }
    };
    const pasteFiles = (toPath: string) => {
        if(copiedIcons.value.length > 0) {
            fsStore.copyFiles(copiedIcons.value, toPath)
            copiedIcons.value = []
        }
        if(cutedIcons.value.length > 0) {
            fsStore.moveFiles(cutedIcons.value, toPath)
            cutedIcons.value = []
        }
    };
    const checkPaste = () => {
        return copiedIcons.value.length > 0 || cutedIcons.value.length > 0;
    }
    const handlePaste = (event: KeyboardEvent, toPath: string) => {
        if ((event.ctrlKey || event.metaKey) && event.key === 'v') { // 检查 Control 或 Command 键和 V 键
            // 这里可以添加粘贴逻辑，例如将 copiedIcons 中的图标路径移动到 toPath
            //console.log('Pasted icons to:', toPath);
            pasteFiles(toPath)
        }
    };

    const handleCut = (event: KeyboardEvent) => {
        if ((event.ctrlKey || event.metaKey) && event.key === 'x') { // 检查 Control 或 Command 键和 X 键
            cutedIcons.value = [...clickedIcons.value]; // 将选中的图标路径存储到 cutedIcons 中
            //console.log('Cut icons:', cutedIcons.value);
        }
    };
    const handleUpload = (event: KeyboardEvent) => {
        if ((event.ctrlKey || event.metaKey) && event.key === 'q') { // 检查 Control 或 Command 键和 O 键
            // 这里可以添加上传逻辑，例如打开文件选择对话框并处理选中的文件
            fsStore.uploadFile(); // 需要提供具体的 toPath
        }
    };

    const handleAll = (event: KeyboardEvent, fileList: any[]) => {
        if ((event.ctrlKey || event.metaKey) && event.key === 'a') { // 检查 Control 或 Command 键和 A 键
            clickedIcons.value = fileList.map(item => item.path); // 将选中的图标路径存储到 clickedIcons 中
            //console.log('Selected all icons:', clickedIcons.value);
        }
    };
    const handleDelete = (event: KeyboardEvent) => {
        if (event.key === 'Delete') { // 检查 Control 或 Command 键和 Delete 键
            // 这里可以添加删除逻辑，例如从文件系统中删除 clickedIcons 中的图标路径
            fsStore.handleDeleteFiles(clickedIcons.value);
            clickedIcons.value = []; // 清空选中的图标路径

        }
    };
    let handleKeyDown: (event: KeyboardEvent) => void;

    const addEvents = (toPath: string, fileList: any[]) => {
        //if (!desktopAreaRef) return;
        removeEvents()
        handleKeyDown = (event: KeyboardEvent) => {
            //console.log(event.key)
            settingsStore.setLockTime();
            handleCopy(event);
            handleCut(event);
            handleAll(event, fileList);
            handleDelete(event);
            handleUpload(event);
            if (event.key === 'v') {
                handlePaste(event, toPath); // 需要提供具体的 toPath
            }
        };

        document.addEventListener('keydown', handleKeyDown);
        document.addEventListener('click', settingsStore.setLockTime);
    };

    const removeEvents = () => {
        if (handleKeyDown) {
            document.removeEventListener('keydown', handleKeyDown);
            resetClickedIcons(new MouseEvent('click'));
        }
        document.removeEventListener('click', settingsStore.setLockTime);
    };

    return {
        clickedIcons,
        copiedIcons,
        cutedIcons,
        setClickIcons,
        handleClick,
        handleCopy,
        handlePaste,
        handleCut,
        handleAll,
        resetClickedIcons,
        checkPaste,
        pasteFiles,
        addEvents,
        removeEvents,
    };
});