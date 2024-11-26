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
  const hasDefaultPwd = ref(false)
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
    await checkDefault()
  }
  const delPwd = async (id: number) => {
    await db.delete('filePwdBox', id)
    //如果之前存在默认密码，则清除
    // if (params.isDefault == 1) {
    //   const list = await db.getAll('filePwdBox')
    //   if (list && list?.length > 0) {
    //     const item = list[list.length - 1]
    //     item.isDefault = 1
    //     await db.update('filePwdBox', item.id, item)
    //   }
    // }
  }
  const pageChange = async (current: number) => {
    page.value.current = current
    getPage()
  }
  const getPageCount = async () => {
    page.value.total = await db.countSearch('filePwdBox')
    page.value.pages = Math.floor(page.value.total / page.value.size)
    
    // 检查是否有余数
    if (page.value.total % page.value.size !== 0) {
      // 如果有余数，则加1
      page.value.pages++;
    }
    return page.value
  }
  // 判断是否有默认密码
  const checkDefault = async () => {
    const res = await db.getByField('filePwdBox', 'isDefault', 1)
    if (res && res.length > 0) {
      hasDefaultPwd.value = true
      return
    } 
    hasDefaultPwd.value = false
  }
  //将其余默认密码设置为空
  const setDefaultPwd = async () => {
    const res = await db.getByField('filePwdBox', 'isDefault', 1)
    if (!res && res.length == 0) {
      return
    };
    for (const item of res) {
      item.isDefault = 0
      db.update('filePwdBox', item.id, item)
    }
  }
  return {
    page,
    pwdList,
    hasDefaultPwd,
    addPwd,
    getPage,
    pageChange,
    delPwd,
    setDefaultPwd
  }
})