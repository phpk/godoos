import { defineStore } from 'pinia'
import { db } from './db.ts'
import { getLang } from "@/i18n/index.ts"
import { promptAction,promptsZh,promptsEn } from "./prompt/index.ts"
import { ref } from "vue"

export const useAssistantStore = defineStore('assistant', () => {
  const currentLang = getLang()
  const showAdd = ref(false)
  const showLeft = ref(false)
  const page = ref({
    current: 1,
    size: 10,
    total: 0,
    pages: 0,
    visible: 5
  })
  const promptList = ref([])
  const currentCate = ref('all')
  const editId = ref(0)

  const addPrompt = async (prompt: any) => {
    prompt.lang = currentLang
    await db.addOne('prompts', prompt)
  }
  const handlerLeft = () => {
    showLeft.value = !showLeft.value
  }
  const savePromptData = async (saveData: any) => {
    saveData.createdAt = new Date()
    //console.log(saveData)
    if(saveData.isdef > 0) {
      saveData.isdef = 1
      //console.log(saveData)
      await db.modify('prompts', "action",saveData.action, {isdef: 0})
    }else{
      saveData.isdef = 0
    }
    //console.log(saveData)
    showAdd.value = false;
    if (editId.value > 0) {
      if(saveData.isdef < 1) {
        const has = await getPromptById(editId.value)
        if(has.isdef > 0){
          return false
        }
      }
      await updatePrompt(saveData);
      editId.value = 0
    } else {
      await addPrompt(saveData);
    }

    await getPromptList();
    
    return true
  }
  const getPromptById = async (id: number) => {
    return await db.getOne('prompts', id)
  }
  const getPrompt = async (action : string) => {
    const data:any = await db.get("prompts",{
      action,
      isdef:1,
      lang:currentLang
    })
    if(data){
      return data.prompt
    }else{
      return ''
    }
  }
  const getPrompts = async (action: string) => {
    const list = await db.rows("prompts", {
      action,
      lang: currentLang
    })
    const promptData = list.find((item:any) => item.isdef == 1)
    return { list, current : promptData }
  }
  const getWhere = () => {
    const where:any = {}
    if(currentCate.value != 'all'){
      where.action = currentCate.value
    }
    where.lang = currentLang
    return where
  }
  const getPromptList = async () => {
    //promptList.value = await db.getAll('prompt')
    const wsql = getWhere()
    promptList.value = await db.pageSearch('prompts',
    page.value.current,
    page.value.size,
      wsql)
    if (promptList.value.length == 0) {
      page.value.current = page.value.current > 1 ? page.value.current - 1 : 1
      promptList.value = await db.pageSearch('prompts', page.value.current, page.value.size, wsql)
    }
    await getPageCount()

  }
  const getPageCount = async () => {
    page.value.total = await db.countSearch('prompts', getWhere())
    page.value.pages = Math.floor(page.value.total / 10)
    // 检查是否有余数
    if (page.value.total % 10 !== 0) {
      // 如果有余数，则加1
      page.value.pages++;
    }
    //console.log(pageCount.value)
    return page.value
  }
  const pageClick = async (pageNum: any) => {
    //console.log(pageNum)
    page.value.current = pageNum
    await getPromptList()
  }
  const updatePrompt = async (prompt: any) => {
    //console.log(prompt)
    await db.update('prompts', editId.value, prompt)
  }
  const changeCate = async (catename: string) => {
    currentCate.value = catename
    showLeft.value = false
    await getPromptList()
  }
  const deletePrompt = async (id: number) => {
    const data = await db.getOne('prompts', id)
    if(data.isdef > 0) {
      return false
    }
    await db.delete('prompts', id)
    await getPromptList()
    return true
  }
  async function initPrompt() {
    await db.clear("prompts")
    promptsZh.forEach((d: any) => {
      d.lang = "zh-cn"
      if (!d.action) {
        d.action = "chat"
      }

      if (!d.isdef) {
        d.isdef = 0
      }
    })
    promptsEn.forEach((d: any) => {
      d.lang = "en"
      if (!d.action) {
        d.action = "chat"
      }
      if (!d.isdef) {
        d.isdef = 0
      }
    })

    const save = [...promptsZh, ...promptsEn]
    await db.addAll("prompts", save)
  }
  return {
    showAdd,
    showLeft,
    page,
    currentCate,
    promptList,
    promptAction,
    editId,
    handlerLeft,
    addPrompt,
    savePromptData,
    getPromptById,
    getPrompt,
    getPrompts,
    getPromptList,
    pageClick,
    changeCate,
    deletePrompt,
    initPrompt
  }
})
//export const assistantStore = useAssistantStore()
