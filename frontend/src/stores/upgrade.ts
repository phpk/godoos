import { defineStore } from "pinia";
import { ref } from "vue";
import { getSystemKey, setSystemKey, parseJson, getSystemConfig } from '@/system/config'
import { RestartApp } from '@/util/goutil';
import { ElMessage } from 'element-plus'
import { t } from '@/i18n';
export const useUpgradeStore = defineStore('upgradeStore', () => {
    const hasUpgrade = ref(false);
    const hasNotice = ref(false);
    const hasAd = ref(false);
    const updateUrl = ref('');
    const versionTag = ref(0)
    const upgradeDesc = ref('')
    const currentVersion = ref(0)
    const progress = ref(0)
    const noticeList:any = ref([])
    const adList:any = ref([])
    async function checkUpdate() {
        const config = getSystemConfig();
        currentVersion.value = config.version;
        const releaseRes = await fetch(`${config.apiUrl}/system/updateInfo`)
        if (!releaseRes.ok) return;
        const releaseData = await releaseRes.json()
        versionTag.value = releaseData.version
        if(versionTag.value > config.version){
            upgradeDesc.value = releaseData.desc ?? t('upgrade.msg')
            hasUpgrade.value = true
            updateUrl.value = releaseData.url
        }
        if (releaseData.noticeList && releaseData.noticeList.length > 0) {
            hasNotice.value = true
            noticeList.value = releaseData.noticeList
        }
        if (!hasUpgrade.value && releaseData.adList && releaseData.adList.length > 0) {
            hasAd.value = true
            adList.value = releaseData.adList
        }
    }
    async function update() {
        const apiUrl = getSystemKey('apiUrl')
        const upUrl = `${apiUrl}/system/update?url=${updateUrl.value}`
        const upRes = await fetch(upUrl)
        if (!upRes.ok) return;
        const reader: any = upRes.body?.getReader();
        if (!reader) {
            ElMessage({
                type: 'error',
                message: "the system has not stream!"
            })
        }
        while (true) {
            const { done, value } = await reader.read();
            if (done) {
                reader.releaseLock();
                break;
            }
            const rawjson = new TextDecoder().decode(value);
            const json = parseJson(rawjson);
            //console.log(json)
            if (json) {
                if (json.progress) {
                    progress.value = json.progress
                }
                if (json.updateCompleted) {
                    hasUpgrade.value = false
                    progress.value = 0
                    ElMessage({
                        type: 'success',
                        message: "update completed!"
                    })
                    setSystemKey('version', versionTag.value)
                    currentVersion.value = versionTag.value
                    RestartApp()
                    break;
                }
            }
        }
    }
    return {
        hasUpgrade,
        hasNotice,
        hasAd,
        versionTag,
        upgradeDesc,
        updateUrl,
        noticeList,
        adList,
        progress,
        checkUpdate,
        update
    }
})