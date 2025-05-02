// src/api/local/db.ts
import Database from '@tauri-apps/plugin-sql';
import { migrationDatabase } from './migration';
let dbInstance: Database | null = null;

export async function connectDatabase(dbPath: string = 'sqlite:godoos.db'): Promise<void> {
  if (!dbInstance) {
    dbInstance = await Database.load(dbPath);
    console.log('✅ 数据库已连接');
  }
}
export async function setupDatabase() {
  await connectDatabase();
  await migrationDatabase();
}
export function getDatabase(): Database {
  if (!dbInstance) {
    throw new Error('数据库尚未连接，请先调用 connectDatabase');
  }
  return dbInstance;
}