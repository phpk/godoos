//import { EventsOff, EventsOn } from '~/runtime';
//import manifest from '../../package.json';
import {isWindowsOS,getSystemConfig} from '@/system/config'
export async function checkUpdate() {
    if(!(window as any).go) return;
    const config = getSystemConfig();
    const updateGiteeUrl = `https://gitee.com/api/v5/repos/ruitao_admin/godoos/releases/`
    const releaseRes = await fetch(updateGiteeUrl)
    if(!releaseRes.ok) return;
    const releaseData = await releaseRes.json()
    const versionTag = releaseData.tag_name;
    if(!versionTag) return;
    if (versionTag.replace('v', '') <= config.version) return;
    const verifyUrl = `${updateGiteeUrl}tags/${versionTag}`;
    const verRes = await fetch(verifyUrl);
    if(!verRes.ok) return;
    const verData = await verRes.json()
    if(!verData.assets || verData.assets.length <= 0) return;
    const appName = "godoos"+ versionTag + (isWindowsOS() ? '.exe' : '');
    const updateUrl = `${updateGiteeUrl}download/${versionTag}/${appName}`;
    console.log(updateUrl)
    // fetch(`${updateGiteeUrl}latest`).then((r) => {
    //     if (r.ok) {
    //       r.json().then((data) => {
    //         if (data.tag_name) {
            //   const versionTag = data.tag_name;
            //   console.log(versionTag)
            //   if (versionTag.replace('v', '') > manifest.version) {
            //     const verifyUrl = `${updateGiteeUrl}tags/${versionTag}`;
            //   }
              /*
              if (versionTag.replace('v', '') > manifest.version) {
                const verifyUrl = `${updateGiteeUrl}tags/${versionTag}`;
  
                fetch(verifyUrl).then((r) => {
                  if (r.ok) {
                    r.json().then((data) => {
                      if (data.assets && data.assets.length > 0) {
                        const asset = data.assets.find((a: any) => a.name.toLowerCase().includes(commonStore.platform.toLowerCase().replace('darwin', 'macos')));
                        if (asset) {
                          const updateUrl = `${updateGiteeUrl}download/${versionTag}/${asset.name}`;
                          toastWithButton(t('New Version Available') + ': ' + versionTag, t('Update'), () => {
                            DeleteFile('cache.json');
                            const progressId = 'update_app';
                            const progressEvent = 'updateApp';
                            const updateProgress = (ds: DownloadStatus | null) => {
                              const content =
                                t('Downloading update, please wait. If it is not completed, please manually download the program from GitHub and replace the original program.')
                                + (ds ? ` (${ds.progress.toFixed(2)}%  ${bytesToReadable(ds.transferred)}/${bytesToReadable(ds.size)})` : '');
                              const options: ToastOptions = {
                                type: 'info',
                                position: 'bottom-left',
                                autoClose: false,
                                toastId: progressId,
                                hideProgressBar: false,
                                progress: ds ? ds.progress / 100 : 0
                              };
                              if (toast.isActive(progressId))
                                toast.update(progressId, {
                                  render: content,
                                  ...options
                                });
                              else
                                toast(content, options);
                            };
                            updateProgress(null);
                            EventsOn(progressEvent, updateProgress);
                            UpdateApp(updateUrl).then(() => {
                              toast(t('Update completed, please restart the program.'), {
                                  type: 'success',
                                  position: 'bottom-left',
                                  autoClose: false
                                }
                              );
                            }).catch((e) => {
                              toast(t('Update Error') + ' - ' + (e.message || e), {
                                type: 'error',
                                position: 'bottom-left',
                                autoClose: false
                              });
                            }).finally(() => {
                              toast.dismiss(progressId);
                              EventsOff(progressEvent);
                            });
                          }, {
                            autoClose: false,
                            position: 'bottom-left'
                          });
                        }
                      }
                    });
                  } else {
                    throw new Error('Verify response was not ok.');
                  }
                });
              } else {
                if (notifyEvenLatest) {
                  toast(t('This is the latest version'), { type: 'success', position: 'bottom-left', autoClose: 2000 });
                }
              }
                */
    //         } else {
    //           throw new Error('Invalid response.');
    //         }
    //       });
    //     } else {
    //       throw new Error('Network response was not ok.');
    //     }
    //   }
    // ).catch((e) => {
    //   //toast(t('Updates Check Error') + ' - ' + (e.message || e), { type: 'error', position: 'bottom-left' });
    // });
  }
  