import { defineStore } from "pinia";
import { ref } from "vue";
import { getSystemConfig,fetchGet } from "@/system/config";
export const useNotifyStore = defineStore('notifyStore', () => {
    const config = getSystemConfig();
    const notifyList:any = ref([]);
    const page = ref({
        current: 1,
        size: 10,
        total: 0,
        pages: 0,
    })
    async function getList(){
       if(config.userType == 'person'){
            return
       }
       const complition = await fetchGet(config.userInfo.url + '/news/list?page=' + page.value.current + '&limit='+ page.value.size)
       if(!complition.ok){
        return;
       }
       const data = await complition.json();
       if(data.success){
        page.value.total = data.data.total;
        notifyList.value = data.data.list;
        page.value.pages = Math.ceil(page.value.total / page.value.size)
       }
       
    }
    const pageClick = async (pageNum: any) => {
        //console.log(pageNum)
        page.value.current = pageNum
        await getList()
    }
    return {
        notifyList,
        getList,
        page,
        pageClick
    }
})