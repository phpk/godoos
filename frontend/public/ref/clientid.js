function getOSInfo() {
  const userAgent = navigator.userAgent;
  if (/windows phone/i.test(userAgent)) {
    return 'Windows Phone';
  }
  if (/win/i.test(userAgent)) {
    return 'Windows';
  }
  if (/mac/i.test(userAgent)) {
    return 'Mac';
  }
  if (/x11/i.test(userAgent)) {
    return 'UNIX';
  }
  if (/android/i.test(userAgent)) {
    return 'Android';
  }
  if (/iphone|ipad|ipod/i.test(userAgent)) {
    return 'iOS';
  }
  return 'Unknown';
}

function bin2hex(s) {
  s = encodeURI(s); // 只会有0-127的ascii不转化
  let m = s.match(/%[\dA-F]{2}/g), a = s.split(/%[\dA-F]{2}/), i, j, n, t;
  m.push("");
  for (i in a) {
    if (a[i] === "") { a[i] = m[i]; continue; }
    n = "";
    for (j in a[i]) {
      t = a[i][j].charCodeAt().toString(16).toUpperCase();
      if (t.length === 1) t = "0" + t;
      n += "%" + t;
    }
    a[i] = n + m[i];
  }
  return a.join("").split("%").join("");
}

function getBrowserFingerprint() {
  const canvas = document.createElement('canvas');
  const ctx = canvas.getContext('2d'); // 移除类型声明
  const txt = window.location.hostname;
  ctx.textBaseline = 'top';
  ctx.font = '14px \'Arial\'';
  ctx.textBaseline = 'alphabetic';
  ctx.fillStyle = '#f60';
  ctx.fillRect(125, 1, 62, 20);
  ctx.fillStyle = '#069';
  ctx.fillText(txt, 2, 15);
  ctx.fillStyle = 'rgba(102, 204, 0, 0.7)';
  ctx.fillText(txt, 4, 17);

  const b64 = canvas.toDataURL().replace("data:image/png;base64,", "");
  const bin = window.atob(b64);
  const hash = bin2hex(bin.slice(-16, -12));

  const fingerprint = [
    navigator.platform,
    navigator.product,
    navigator.productSub,
    navigator.appName,
    navigator.appVersion,
    navigator.javaEnabled(),
    navigator.userAgent,
    screen.width,
    screen.height,
    new Date().getTimezoneOffset(),
    hash,
    screen.colorDepth,
    navigator.language,
    navigator.hardwareConcurrency,
    getOSInfo(),
    navigator.maxTouchPoints,
    navigator.doNotTrack,
    navigator.cookieEnabled
  ].join('|');

  return fingerprint;
}

function djb2Hash(str) {
  let hash = 5381;
  for (let i = 0; i < str.length; i++) {
    const char = str.charCodeAt(i);
    hash = ((hash << 5) + hash) + char; /* hash * 33 + c */
  }
  return hash.toString(36); // Convert to base 36 for shorter string
}

function GetClientId() {
  let uuid = localStorage.getItem("ClientID");
  if (!uuid) {
    const fingerprint = getBrowserFingerprint();
    uuid = djb2Hash(fingerprint).slice(0, 36);
    localStorage.setItem("ClientID", uuid);
  }

  return uuid;
}