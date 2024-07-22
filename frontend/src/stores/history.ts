import { defineStore } from "pinia";
import { ref } from "vue";
export const useHistoryStore = defineStore('historyStore', () => {
    const historyList:any = ref({});
    function getList(name : string){
        return historyList.value[name];
    }
    function addList(name:string, data:any){
        //console.log(data)
        if (!historyList.value[name]){
            historyList.value[name] = []
        }
        historyList.value[name] = historyList.value[name]
                        .filter((item : any) => item.path !== data.path)
                        .slice(0, 10);
        historyList.value[name].unshift(data)
        if(historyList.value[name].length > 10){
            historyList.value[name].pop()
        }
        //console.log(historyList.value[name])
    }
    return {
        historyList,
        getList,
        addList
    }
}, {
    persist: {
        enabled: true,
        strategies: [
            {
                storage: localStorage,
                paths: [
                    "historyList"
                ]
            }, // name 字段用localstorage存储
        ],
    }
})