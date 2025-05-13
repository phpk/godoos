// import { md5 } from "js-md5";
// // const token = sign({ id: 1 }, "secret", { expiresIn: "1h" });
// // console.log(token);

// // try {
// //     const decoded = verify(token, "secret");
// //     console.log(decoded);
// // } catch (err) {
// //     console.error(err.message); // token 已过期 或 无效的 token
// // }
// // 支持 expiresIn 参数
// export function sign(payload: any, secret: string, options: { expiresIn?: string; refreshExpiresIn?: string } = {}): string {
//     const { expiresIn = '24h', refreshExpiresIn = '7d' } = options;

//     const now = Math.floor(Date.now() / 1000);
//     payload.exp = now + parseExpiresIn(expiresIn);
//     payload.refreshExp = now + parseExpiresIn(refreshExpiresIn);

//     const header = { alg: "HS256", typ: "JWT" };
//     const base64Header = btoa(JSON.stringify(header));
//     const base64Payload = btoa(JSON.stringify(payload));
//     const signature = md5(`${base64Header}.${base64Payload}${secret}`);
//     return `${base64Header}.${base64Payload}.${signature}`;
// }

// export function verify(token: string, secret: string): any {
//     const [base64Header, base64Payload, signature] = token.split(".");
//     const expectedSignature = md5(`${base64Header}.${base64Payload}${secret}`);

//     if (signature !== expectedSignature) {
//         throw new Error("无效的 token");
//     }

//     const decodedPayload = JSON.parse(atob(base64Payload));

//     const now = Math.floor(Date.now() / 1000);
//     if (decodedPayload.exp && decodedPayload.exp < now) {
//         throw new Error("token 已过期");
//     }

//     return decodedPayload;
// }
// export function verifyWithRefresh(token: string, secret: string, options: { expiresIn?: string } = {}): { success: boolean, token: string | null } {
//     try {
//         verify(token, secret);
//         return { success: true, token: token }; // 未过期，无需刷新
//     } catch (err: any) {
//         if (err.message === "token 已过期") {
//             try {
//                 // 尝试刷新 token
//                 const newToken = refreshToken(token, secret, options);
//                 return { success: true, token: newToken }; // 刷新成功
//             } catch (refreshErr: any) {
//                 if (refreshErr.message === "刷新 token 已过期") {
//                     return { success: false, token: null }; // refresh token 也过期了
//                 }
//                 return { success: false, token: null }; // 其他错误
//             }
//         } else {
//             return { success: false, token: null }; // 校验失败（签名错误等）
//         }
//     }
// }
// export function refreshToken(oldToken: string, secret: string, options: { expiresIn?: string } = {}): string {
//     const decoded = verifyRefreshToken(oldToken, secret);
//     const newPayload = { id: decoded.id }; // 可根据需要提取更多字段
//     return sign(newPayload, secret, options);
// }

// function verifyRefreshToken(token: string, secret: string): any {
//     const decoded = verify(token, secret);

//     const now = Math.floor(Date.now() / 1000);
//     if (!decoded.refreshExp || decoded.refreshExp < now) {
//         throw new Error("刷新 token 已过期");
//     }

//     return decoded;
// }
// function parseExpiresIn(expiresIn: string): number {
//     const regex = /^(\d+)([smhdw])$/;
//     const match = expiresIn.match(regex);
//     if (!match) throw new Error("Invalid expiresIn format");

//     const value = parseInt(match[1], 10);
//     const unit = match[2];

//     switch (unit) {
//         case 's': return value;
//         case 'm': return value * 60;
//         case 'h': return value * 3600;
//         case 'd': return value * 86400;
//         case 'w': return value * 604800;
//         default: throw new Error("Unknown time unit");
//     }
// }

// export function generate16BitString(): string {
//     const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
//     let result = '';
//     for (let i = 0; i < 16; i++) {
//         const randomIndex = Math.floor(Math.random() * chars.length);
//         result += chars.charAt(randomIndex);
//     }
//     return result;
// }