import { defineStore } from "pinia";
import { ref } from "vue";
export const useTabsStore = defineStore("tabs", () => {
  const activeId = ref<any>("");
  const tabsList = ref<any>([]);
  const addTab = (win: any) => {
    const exists = tabsList.value.find((tab:any) => tab.id === win.id)
    if(!exists){
        tabsList.value.push(win)
    }
    activeId.value = win.id
  }
  const removeTab = (targetId: string) => {
    const tabs = tabsList.value
    let activeid = activeId.value
    if (activeid === targetId) {
        tabs.forEach((tab:any, index:number) => {
        if (tab.id === targetId) {
          const nextTab = tabs[index + 1] || tabs[index - 1]
          if (nextTab) {
            activeid= nextTab.id
          }
        }
      })
    }
  
    activeId.value = activeid
    tabsList.value = tabs.filter((tab:any) => tab.id !== targetId)
  }
  return {
    activeId,
    tabsList,
    removeTab,
    addTab
  };
});
