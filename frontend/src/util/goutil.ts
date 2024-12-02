import { getSystemKey } from "@/system/config";
export async function OpenDirDialog() {
    if ((window as any).go) {
        return (window as any)['go']['app']['App']['OpenDirDialog']();
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
        const apiUrl = getSystemKey("apiUrl");
        fetch(apiUrl + "/system/restart").then(() => {
            setTimeout(() => {
                window.location.reload();
            }, 1000);
        })
        //window.location.reload();
    } else {
        return (window as any)['go']['app']['App']['RestartApp']();
    }

}
