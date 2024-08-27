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
    const versionTag = ref('')
    const upgradeDesc = ref('')
    const currentVersion = ref('')
    const progress = ref(0)
    const noticeList:any = ref([])
    const adList:any = ref([])
    function compareVersions(version1:string, version2:string) {
        // 将版本号字符串按"."分割成数组
        const parts1 = version1.split('.').map(Number);
        const parts2 = version2.split('.').map(Number);
    
        // 确保两个数组长度相同
        const maxLength = Math.max(parts1.length, parts2.length);
        while (parts1.length < maxLength) parts1.push(0);
        while (parts2.length < maxLength) parts2.push(0);
    
        // 比较每个部分
        for (let i = 0; i < maxLength; i++) {
            if (parts1[i] > parts2[i]) return 1;
            if (parts1[i] < parts2[i]) return -1;
        }
    
        // 如果所有部分都相等，则返回0
        return 0;
    }
    async function checkUpdate() {
        const config = getSystemConfig();
        currentVersion.value = config.version;
        const releaseRes = await fetch(`${config.apiUrl}/system/updateInfo`)
        if (!releaseRes.ok) return;
        const releaseData = await releaseRes.json()
        if(releaseData.code !== 200){
            return
        }
        const res = releaseData.data
        if(!res.version || res.version == ""){
            return
        }
        versionTag.value = res.version
        if (compareVersions(versionTag.value, config.version) > 0) {
            upgradeDesc.value = res.desc ?? t('upgrade.msg')
            hasUpgrade.value = true
            updateUrl.value = res.url
        }
        if (res.adList) {
            if(res.adList.center && res.adList.center.length > 0){
                hasNotice.value = true
                noticeList.value = res.adList.center
            }
            if(res.adList.bottom && res.adList.bottom.length > 0){
                hasAd.value = true
                adList.value = res.adList.bottom
            }
            
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