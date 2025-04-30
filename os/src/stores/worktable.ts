import { defineStore } from 'pinia'
import { ref } from 'vue'
export const useWorkTableStore = defineStore('worktable', () => {
  const currentTab = ref('0')
  const setCurrentTab = (tab: string) => {
    currentTab.value = tab
  }
  return {
    currentTab,
    setCurrentTab,
  }
})
