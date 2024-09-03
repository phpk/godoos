import { defineStore } from "pinia";
import { ref } from "vue";
import { getSystemConfig, fetchGet } from "@/system/config";
import { useUpgradeStore } from "./upgrade";
export const useNotifyStore = defineStore('notifyStore', () => {
    const config = getSystemConfig();
    const upgradeStore = useUpgradeStore()
    const notifyList: any = ref([]);
    const readList: any = ref([])
    const showNews = ref(false)
    const newsDetail = ref({})
    const page = ref({
        current: 1,
        size: 10,
        total: 0,
        pages: 0,
    })
    async function getList() {
        if (config.userType == 'person') {
            return
        }
        const complition = await fetchGet(config.userInfo.url + '/news/list?page=' + page.value.current + '&limit=' + page.value.size)
        if (!complition.ok) {
            return;
        }
        const data = await complition.json();
        if (data.success) {
            page.value.total = data.data.total;
            notifyList.value = data.data.list;
            page.value.pages = Math.ceil(page.value.total / page.value.size)
            checkNotify()
        }

    }
    const checkNotify = ()=> {
        const list = notifyList.value.filter((item: any) => {
            return readList.value.indexOf(item.id) < 0
        })
        if(list.length < 1)return;
        const centerList = list.filter((item: any) => {
            return item.position == 'center'
        })
        const bottomList = list.filter((item: any) => {
            return item.position == 'bottom'
        })
        if(centerList.length > 0){
            upgradeStore.hasAd = true
            upgradeStore.adList = [...centerList, ...upgradeStore.adList]
        }
        if(bottomList.length > 0){
            upgradeStore.hasNotice = true
            upgradeStore.noticeList = [...bottomList, ...upgradeStore.noticeList]
        }

    }
    const addRead = (id: number) => {
        if (readList.value.indexOf(id) < 0) {
            readList.value.push(id)
        }
    }
    const pageClick = async (pageNum: any) => {
        //console.log(pageNum)
        page.value.current = pageNum
        await getList()
    }
    const viewContent = (item :any) =>{
        const timestamp = item.add_time * 1000; // 将 10 位时间戳转换为 13 位
        const date = new Date(timestamp);
        item.showtime = date.toLocaleString();
        newsDetail.value = item
        addRead(item.id)
        showNews.value = true
    }
    return {
        notifyList,
        showNews,
        newsDetail,
        getList,
        addRead,
        page,
        pageClick,
        readList,
        viewContent
    }
}, {
    persist: {
        enabled: true,
        strategies: [
            {
                storage: localStorage,
                paths: [
                    "readList",
                ]
            }, // name 字段用localstorage存储
        ],
    }
})