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
export function isBase64(str:string) {
  if (str === '' || str.trim() === '') {
    return false;
  }
  try {
    return btoa(atob(str)) == str;
  } catch (err) {
    return false;
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
export function base64ToBuffer(base64String : string) {
  
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
export function decodeBase64(base64String :string) {
  // 将Base64字符串分成每64个字符一组
  const bytes = base64ToBuffer(base64String)
  // 将TypedArray转换为字符串
  return new TextDecoder('utf-8').decode(bytes);
}
export function stringToBinary(str:string) {
  const binaryArr = Array.from(str, c => c.charCodeAt(0).toString(2).padStart(8, '0'));
  return binaryArr.join('');
}

