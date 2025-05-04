import { getDatabase } from './db';

export async function migrationDatabase(): Promise<void> {
  const db = getDatabase();
  //await db.execute('DROP TABLE IF EXISTS users');
  await db.execute(
    // `CREATE TABLE users (
    `CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      username TEXT NOT NULL UNIQUE,
      password TEXT NOT NULL,
      nickname TEXT,
      email TEXT,
      phone TEXT,
      isMaster INTEGER DEFAULT 0,
      created_at TEXT DEFAULT (DATETIME('now', 'localtime')),
      updated_at TEXT DEFAULT (DATETIME('now', 'localtime'))
    )`
  );
  console.log('✅ 用户表已初始化');
}