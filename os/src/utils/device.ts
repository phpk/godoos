export const isMobileDevice=() =>{
    const userAgent = navigator.userAgent.toLowerCase();
    const mobileKeywords = ['iphone', 'android', 'mobile', 'blackberry', 'iemobile', 'opera mini'];
    for (const keyword of mobileKeywords) {
        if (userAgent.includes(keyword)) {
            return true;
        }
    }
    return false;
}