// import { type as osType, platform } from "@tauri-apps/plugin-os";
export const isMobileDevice = () => {
    const userAgent = navigator.userAgent.toLowerCase();
    const mobileKeywords = ['iphone', 'android', 'mobile', 'blackberry', 'iemobile', 'opera mini'];
    for (const keyword of mobileKeywords) {
        if (userAgent.includes(keyword)) {
            return true;
        }
    }
    return false;
}
// export const getPlatformInfo = () => {
//     const info = {
//         platform: "web",
//         ostype: "web",
//         isMobile: false,
//         isDesktop: false,
//         isWeb: true,
//     };
//     try {
//         info.platform = platform();
//     }
//     catch (error) {
//         info.platform = "web";
//     }
//     try {
//         info.ostype = osType();
//     }
//     catch (error) {
//         // console.warn(error);
//         info.ostype = "web";
//     }
//     info.isMobile = ["android", "ios"].includes(info.ostype);
//     info.isWeb = info.platform === "web";
//     info.isDesktop = !info.isMobile && !info.isWeb;
//     return info;
// }