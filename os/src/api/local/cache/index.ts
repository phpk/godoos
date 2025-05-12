// import { Store } from '@tauri-apps/plugin-store';

// // 创建 Store 实例（使用异步方式获取）
// let store: Store;

// async function getStore(): Promise<Store> {
//   if (!store) {
//     store = await Store.load('store.bin');
//   }
//   return store;
// }

// // 缓存默认有效期（单位：毫秒）
// const DEFAULT_TTL = 1000 * 60 * 60 * 24; // 默认 1 天

// /**
//  * 设置带过期时间的缓存
//  * @param key 缓存键名
//  * @param value 缓存值
//  * @param ttl 缓存有效时间（毫秒），不传则使用默认值
//  */
// export async function set<T>(key: string, value: T, ttl?: number): Promise<void> {
//   const storeInstance = await getStore();
//   const expireAt = Date.now() + (ttl ?? DEFAULT_TTL);
//   await storeInstance.set(key, { value, expireAt });
//   await storeInstance.save(); // 可选：立即保存
// }

// /**
//  * 获取缓存值（如果未过期）
//  * @param key 缓存键名
//  * @returns 缓存值或 null（如果未找到或已过期）
//  */
// export async function get<T>(key: string): Promise<T | null> {
//   const storeInstance = await getStore();
//     // 使用 storeInstance.has 做初步存在性检查
//     const exists = await storeInstance.has(key);
//     if (!exists) return null;
//   const entry = await storeInstance.get<{ value: T; expireAt: number }>(key);

//   if (!entry) return null;

//   const { value, expireAt } = entry;
//   if (Date.now() > expireAt) {
//     await storeInstance.delete(key); // 删除过期缓存
//     return null;
//   }

//   return value;
// }
// /**
//  * 检查缓存中是否存在指定的键且未过期
//  * @param key 缓存键名
//  * @returns 如果存在且未过期返回 true，否则返回 false
//  */
// export async function has(key: string): Promise<boolean> {
//   const storeInstance = await getStore();
//   // 使用 storeInstance.has 做初步存在性检查
//   const exists = await storeInstance.has(key);
//   if (!exists) return false;
//   const entry = await storeInstance.get<{ value: any; expireAt: number }>(key);

//   if (!entry) return false;

//   const { expireAt } = entry;
//   if (Date.now() > expireAt) {
//     await storeInstance.delete(key); // 删除已过期的键
//     return false;
//   }

//   return true;
// }
// /**
//  * 清除所有缓存
//  */
// export async function clear(): Promise<void> {
//   const storeInstance = await getStore();
//   await storeInstance.clear();
// }