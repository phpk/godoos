import { getSystemConfig,setSystemKey,parseJson } from '@/system/config'
import { RestartApp } from './goutil';
import { Dialog } from '@/system';
import { ElMessage } from 'element-plus'
export async function checkUpdate() {
  const config = getSystemConfig();
  const updateGiteeUrl = `${config.apiUrl}/system/updateInfo`
  const releaseRes = await fetch(updateGiteeUrl)
  if (!releaseRes.ok) return;
  const releaseData = await releaseRes.json()
  const versionTag = releaseData.version;
  if (!versionTag) return;
  if (versionTag <= config.version) return;
  const updateUrl = releaseData.url
  if (!updateUrl || updateUrl == '') return;
  const dialogRes: any = await Dialog.showMessageBox({
    title: '更新提示',
    message: `发现新版本：${versionTag}，是否更新？`
  })
  //console.log(dialogRes)
  if (dialogRes.response !== -1) {
    return;
  }
  const { setProgress,dialogwin } = Dialog.showProcessDialog({
    message: '正在更新',
  });
  const upUrl = `${config.apiUrl}/system/update?url=${updateUrl}`
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
    console.log(json)
    if(json){
      if(json.progress){
        setProgress(json.progress)
      }
      if(json.updateCompleted){
        dialogwin.close()
        ElMessage({
          type: 'success',
          message: "update completed!"
        })
        setSystemKey('version',versionTag)
        RestartApp()
        break;
      }
    }
  }
}
