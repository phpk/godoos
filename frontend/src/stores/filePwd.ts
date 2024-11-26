import { defineStore } from "pinia"
import { db } from './db.ts'
import { ref } from "vue"

export const useFilePwdStore = defineStore('filePwd', () => {
  const page = ref({
    current: 1,
    size: 6,
    total: 0,
    pages: 0,
    visible: 5
  })
  const pwdList = ref([])

  const addPwd = async (params: any) => {
    await db.addOne('filePwdBox', params)
  }
  const getPage = async () => {
    pwdList.value = await db.getPage('filePwdBox', page.value.current, page.value.size)
    if (pwdList.value.length == 0) {
      page.value.current = page.value.current > 1 ? page.value.current -1 : 1
      pwdList.value = await db.getPage('filePwdBox', page.value.current, page.value.size)
    }
    await getPageCount()
  }
  const getPageCount = async () => {
    page.value.total = await db.countSearch('filePwdBox')
    page.value.pages = Math.floor(page.value.total / page.value.size)
    
    // 检查是否有余数
    if (page.value.total % page.value.size !== 0) {
      // 如果有余数，则加1
      page.value.pages++;
    }
    //console.log(pageCount.value)
    return page.value
  }
  return {
    page,
    pwdList,
    addPwd,
    getPage
  }
})