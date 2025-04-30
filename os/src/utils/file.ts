export function dealSize(size = 0) {
  if (size < 1024) {
    return size + 'B';
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + 'KB';
  } else if (size < 1024 * 1024 * 1024) {
    return (size / 1024 / 1024).toFixed(2) + 'MB';
  } else if (size < 1024 * 1024 * 1024 * 1024) {
    return (size / 1024 / 1024 / 1024).toFixed(2) + 'GB';
  } else {
    return (size / 1024 / 1024 / 1024 / 1024).toFixed(2) + 'TB';
  }
}
export function isBase64(str: string): boolean {
  // 去除字符串两端的空白字符
  const trimmedStr = str.trim();
  if (trimmedStr === '') {
    return false;
  }

  // 更严格的 Base64 正则表达式
  const regex = /^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=)?$/;
  const regexMatch = regex.test(trimmedStr);
  //console.log(`Regex match for "${trimmedStr}":`, regexMatch);

  // 如果正则表达式不匹配，则直接返回 false
  if (!regexMatch) {
    return false;
  }

  // 尝试解码并重新编码，检查是否与原字符串相同
  try {
    const binaryString = atob(trimmedStr);
    //console.log(`Decoded binary string for "${trimmedStr}":`, binaryString);

    // 将二进制字符串转换为 UTF-8 字符串
    const utf8String = new TextDecoder('utf-8').decode(
      Uint8Array.from(binaryString, (char) => char.charCodeAt(0))
    );
    //console.log(`Decoded UTF-8 string for "${trimmedStr}":`, utf8String);

    // 将 UTF-8 字符串重新编码为 Base64
    const reEncoded = btoa(
      new TextEncoder().encode(utf8String).reduce(
        (data, byte) => data + String.fromCharCode(byte),
        ''
      )
    );
    const reEncodedMatch = reEncoded === trimmedStr;
    //console.log(`Re-encoded match for "${trimmedStr}":`, reEncodedMatch);
    return reEncodedMatch;
  } catch (err) {
    console.error(`Error decoding "${trimmedStr}":`, err);
    return false;
  }
}


export function getContent(data: any) {
  if (typeof data === 'string') {
    if (isBase64(data)) {
      return decodeBase64(data);
    } else {
      return data;
    }
  } else if (data instanceof Blob) {
    return blobToTextOrDataURL(data);
  } else if (data instanceof ArrayBuffer) {
    return arrayBufferToTextOrBase64(data);
  } else {
    return data.toString();
  }
}

async function blobToTextOrDataURL(blob: Blob): Promise<string> {
  if (blob.type === 'application/pdf') {
    return URL.createObjectURL(blob);
  } else {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onload = () => resolve(reader.result as string);
      reader.onerror = reject;
      reader.readAsText(blob);
    });
  }
}

function arrayBufferToTextOrBase64(buffer: ArrayBuffer): string {
  const bytes = new Uint8Array(buffer);
  const text = new TextDecoder('utf-8').decode(bytes);
  if (isBase64(text)) {
    return decodeBase64(text);
  } else {
    return binaryToBase64(bytes);
  }
}
export function binaryToBase64(data: Iterable<number>) {
  let binary = "";
  const bytes = new Uint8Array(data);
  for (let i = 0; i < bytes.length; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return btoa(binary);
}
export function base64ToBuffer(base64String: string) {
  if (base64String == '') return new ArrayBuffer(0);
  // 去掉末尾的填充等号
  base64String = base64String.replace(/\=+$/, '');
  // 解码Base64字符串为二进制字符串
  const binaryString = window.atob(base64String);
  // 获取二进制字符串的长度
  const byteLength = binaryString.length;
  // 创建Uint8Array
  const uint8Array = new Uint8Array(byteLength);
  // 将二进制字符串转换为Uint8Array
  for (let i = 0; i < byteLength; i++) {
    uint8Array[i] = binaryString.charCodeAt(i);
  }

  // 返回ArrayBuffer
  return uint8Array.buffer;
}
export function base64ToBlobPdfUrl(base64: string) {
  const binStr = atob(base64);
  const len = binStr.length;
  const arr = new Uint8Array(len);
  for (let i = 0; i < len; i++) {
    arr[i] = binStr.charCodeAt(i);
  }
  const blob = new Blob([arr], { type: "application/pdf" });
  const url = URL.createObjectURL(blob);
  return url;
}
export function decodeBase64toText(base64String: string) {
  if (isBase64(base64String)) {
    // 将Base64字符串分成每64个字符一组
    const bytes = base64ToBuffer(base64String)
    // 将TypedArray转换为字符串
    return new TextDecoder('utf-8').decode(bytes);
  } else {
    return base64String;
  }

}
export function decodeBase64(base64String: string) {
  // console.log("Base64 string:", base64String);
  console.log(isBase64(base64String))
  //base64String = base64String.replace(/\=+$/, '').trim();
  // 将Base64字符串分成每64个字符一组
  if (base64String == '') {
    return base64String;
  }
  if (!isBase64(base64String)) {
    return base64String;
  }
  // 计算并添加必要的填充字符
  const padding = base64String.length % 4 === 0 ? 0 : 4 - (base64String.length % 4);
  base64String += "=".repeat(padding);

  try {
    const binaryString = atob(base64String);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    // 将TypedArray转换为字符串
    const decodedString = new TextDecoder("utf-8").decode(bytes);
    //console.log("Decoded string:", decodedString);

    return decodedString;
  } catch (error) {
    console.error("Error decoding Base64 string:", error);
    return base64String;
  }
}
export function stringToBinary(str: string) {
  const binaryArr = Array.from(str, c => c.charCodeAt(0).toString(2).padStart(8, '0'));
  return binaryArr.join('');
}

export function parsePath(path: string) {
  const sp = path.charAt(0)
  const arr = path.split(sp)
  const title = arr[arr.length - 1]
  const ext = title.split('.').pop()?.toLocaleLowerCase()
  return { title, ext }
}
export const determineFileType = (extension: string) => {
  const picExt = ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'tiff']
  if (extension === 'pdf') {
    return 'pdf';
  } else if (extension === 'xlsx' || extension === 'xls' || extension === 'csv') {
    return 'excel';
  } else if (extension === 'pptx' || extension === 'ppt') {
    return 'ppt';
  } else if (extension === 'docx') {
    return 'word';
  } else if (extension === 'md' || extension === 'txt') {
    return 'md';
  } else if (picExt.includes(extension)) {
    return 'pic';
  } else {
    return '';
  }
}