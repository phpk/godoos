import { getSystemConfig } from "@/system/config";
export async function OpenDirDialog() {
    if ((window as any).go) {
        return (window as any)['go']['app']['App']['OpenDirDialog']();
    } else {
        return ""
    }
}
export async function ChooseFileDialog() {
    if ((window as any).go) {
        return (window as any)['go']['app']['App']['ChooseFileDialog']();
    } else {
        return ""
    }
}
export async function checkUrl(url: string) {
    try {
        await fetch(url, {
            method: "GET",
            mode: "no-cors",
        });
        return true;
    } catch (error) {
        return false;
    }
}
export function RestartApp() {
    if (!(window as any).go) {
        const config = getSystemConfig();
        if (config.userType == 'person') {
            fetch(config.apiUrl + "/system/restart").then(() => {
                setTimeout(() => {
                    window.location.reload();
                }, 1000);
            })
        } else {
            window.location.reload();
        }

        //window.location.reload();
    } else {
        return (window as any)['go']['app']['App']['RestartApp']();
    }

}
