// Debounce

export const debounce = (func: Function, wait: number): Function => {
  let timeout:any;
  return function (this: any,...args:any[]) {
    clearTimeout(timeout);
    timeout = setTimeout(() => func.apply(this, args), wait);
  };
};

// Throttle
export const throttle = (func: Function, wait: number): Function => {
  let timeout: any;

  return function(this: any, ...args: any[]): void {
    if (!timeout) {
      timeout = setTimeout(() => {
        timeout = null;
        func.apply(this, args);
      }, wait);
    }
  };
};

// Format file size
export const formatFileSize = (size: number): string => {
  const units = ["B", "KB", "MB", "GB", "TB"];
  let unitIndex = 0;

  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024;
    unitIndex++;
  }

  return `${size.toFixed(2)} ${units[unitIndex]}`;
};

// 判断当前设备是否为移动端
export const isMobile = (): boolean => {
  //const invoke = window.__TAURI__.invoke;
  return /Android|webOS|iPhone|iPod|BlackBerry/i.test(navigator.userAgent);
};

// scroll to top
interface ScrollOptions {
  behavior?: "auto" | "smooth";
  block?: "start" | "center" | "end" | "nearest";
  inline?: "start" | "center" | "end" | "nearest";
  top?: number;
}

export const scrollToTop = (
  element: HTMLElement | null,
  options: ScrollOptions = { top: 0, behavior: "auto" }
): void => {
  if (!element) {
    console.error("Element not found");
    return;
  }
  element.scrollTo({
    ...options,
  });
};

// scroll to bottom
export const scrollToBottom = (
  element: HTMLElement | null,
  options: ScrollOptions = { behavior: "auto" }
): void => {
  element?.scrollTo({
    ...options,
    top: element.scrollHeight,
  });
};
export function formatTime(now : any) {
  const month = now.getMonth() + 1;
  const formattedMonth = month < 10 ? '0' + month : month.toString();
  const minute = now.getMinutes()
  const formatMinute = minute < 10 ? '0' + minute : minute.toString()
  return `${now.getFullYear()}-${formattedMonth}-${now.getDate()} ${now.getHours()}:${formatMinute}`;
}
export function formatChatTime(time: any) {
  if (('' + time).length === 10) {
    time = parseInt(time) * 1000
  } else {
    time = +time
  }
  const d = +new Date(time)
  const now = Date.now()

  const diff = (now - d) / 1000

  if (diff < 30) {
    return '刚刚'
  } else if (diff < 3600) {
    // less 1 hour
    return Math.ceil(diff / 60) + '分钟前'
  } else if (diff < 3600 * 24) {
    return Math.ceil(diff / 3600) + '小时前'
  } else if (diff < 3600 * 24 * 2) {
    return '1天前'
  }
}
export function generateRandomString(length : number) {
    let result = '';
    const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
}
 
export function browsertype() {
  const userAgent = navigator.userAgent //取得浏览器的userAgent字符串
  let isOpera = false
  if (userAgent.indexOf("Edge") > -1) {
    return "Edge"
  }
  if (userAgent.indexOf(".NET") > -1) {
    return "IE"
  }
  if (userAgent.indexOf("Opera") > -1 || userAgent.indexOf("OPR") > -1) {
    isOpera = true
    return "Opera"
  } //判断是否Opera浏览器
  if (userAgent.indexOf("Firefox") > -1) {
    return "FF"
  } //判断是否Firefox浏览器
  if (userAgent.indexOf("Chrome") > -1) {
    return "Chrome"
  }
  if (userAgent.indexOf("Safari") > -1) {
    return "Safari"
  } //判断是否Safari浏览器
  if (
    userAgent.indexOf("compatible") > -1 &&
    userAgent.indexOf("MSIE") > -1 &&
    !isOpera
  ) {
    return "IE"
  } //判断是否IE浏览器
}
export function isValidIP(ip: string): boolean {
  // 正则表达式用于匹配 IPv4 地址
  const ipv4Regex = /^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$/;

  // 正则表达式用于匹配 IPv6 地址
  const ipv6Regex = /^([0-9a-fA-F]{1,4}:){7}([0-9a-fA-F]{1,4})$/;

  // 验证 IP 地址
  return ipv4Regex.test(ip) || ipv6Regex.test(ip);
}
