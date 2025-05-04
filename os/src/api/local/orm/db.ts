// src/api/local/db.ts
import Database from '@tauri-apps/plugin-sql';
import { migrationDatabase } from './migration';
import { getPlatformInfo } from "@/utils/device";
let dbInstance: Database | null = null;
const { platform } = getPlatformInfo();
export async function connectDatabase(dbPath: string = 'sqlite:godoos.db'): Promise<void> {
  if (!dbInstance) {
    dbInstance = await Database.load(dbPath);
    console.log('✅ 数据库已连接');
  }
}
export async function initDatabase() {
  if (platform === 'web')return null as any;
  await connectDatabase();
  await migrationDatabase();
}
export function getDatabase(){
  
  if (platform === 'web')return null as any;
  if (!dbInstance) {
    // initDatabase().then(() => {
    //     if (!dbInstance) throw new Error('数据库尚未连接，请先调用 connectDatabase');
    // });
    throw new Error('数据库尚未连接，请先调用 connectDatabase');
  }
  return dbInstance;
}