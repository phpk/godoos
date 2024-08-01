import { defineStore } from 'pinia'
import { ref, inject } from "vue";
import { t } from "@/i18n";
// import storeInitList from "@/assets/store.json";
import { System } from "@/system";
import { getSystemKey } from "@/system/config";
export const useStoreStore = defineStore('storeStore', () => {
    const sys: any = inject<System>("system");
    const currentCateId = ref(0)
    const currentTitle = ref(t("store.hots"))
    const currentCate = ref('hots')
    const categoryList = ['hots', 'work', 'development', 'games', 'education', 'news', 'shopping', 'social', 'utilities', 'others', 'add']
    const categoryIcon = ['HomeFilled', 'Odometer', 'Postcard', 'TrendCharts', 'School', 'HelpFilled', 'ShoppingCart', 'ChatLineRound', 'MessageBox', 'Ticket', 'CirclePlusFilled']
    const isready = ref(false);
    const apiUrl = getSystemKey("apiUrl");
    const installedList: any = ref([]);
    const storeList: any = ref([])
    const outList: any = ref([])
    async function getList() {
        if (currentCate.value == 'add') return;
        //storeList.value = storeInitList
        const storeUrl = apiUrl + '/store/storelist?cate=' + currentCate.value
        const res = await fetch(storeUrl)
        if (!res.ok) {
            return []
        }
        const json: any = await res.json()
        let list = json.data || []
        //console.log(data)
        if (outList.value.length > 0 && currentCate.value == 'hots' && list && list.length > 0) {
            const names = list.map((item: any) => item.name)
            const adds: any = []
            outList.value.forEach((item: any) => {
                if (!names.includes(item.name)) {
                    adds.push(item)
                }
            })
            list = adds.concat(list)
        }
        storeList.value = list
        await checkProgress()
        isready.value = true;
    }
    async function addOutList(item: any) {
        const has = outList.value.find((i: any) => i.name === item.name)
        if (!has) {
            item.isOut = true
            outList.value.push(item)
            await getList()
            return true
        } else {
            return false
        }
    }
    async function changeCate(index: number, item: string) {
        currentCateId.value = index
        currentCate.value = item
        currentTitle.value = t("store." + item)
        await getList()
    }
    async function checkProgress() {
        const completion: any = await fetch(apiUrl + '/store/listporgress')
        if (!completion.ok) {
            return
        }
        let res: any = await completion.json()
        if (!res || res.length < 1) {
            res = []
        }
        storeList.value.forEach((item: any, index: number) => {
            const pitem: any = res.find((i: any) => i.name == item.name)
            //console.log(pitem)
            if (pitem) {
                storeList.value[index].isRuning = pitem.running
            } else {
                storeList.value[index].isRuning = false
            }
        })
    }
    async function addDesktop(item: any) {
        if (item.webUrl) {
            await sys.fs.writeFile(
                `${sys._options.userLocation}Desktop/${item.name}.url`,
                `link::url::${item.webUrl}::${item.icon}`
            );
            setTimeout(() => {
                sys.refershAppList();
            }, 1000);
        }
        installedList.value.push(item.name);
    }
    async function removeDesktop(item: any) {
        if (item.webUrl) {
            await sys.fs.unlink(`${sys._options.userLocation}Desktop/${item.name}.url`);
            setTimeout(() => {
                sys.refershAppList();
            }, 1000);
        }
        delete installedList.value[installedList.value.indexOf(item.name)];
    }
    return {
        currentCateId,
        categoryIcon,
        currentTitle,
        currentCate,
        categoryList,
        isready,
        installedList,
        storeList,
        outList,
        apiUrl,
        changeCate,
        getList,
        addOutList,
        checkProgress,
        addDesktop,
        removeDesktop
    }
}, {
    persist: {
        enabled: true,
        strategies: [
            {
                storage: localStorage,
                paths: [
                    "outList",
                    "installedList"
                ]
            }, // name 字段用localstorage存储
        ],
    }
})