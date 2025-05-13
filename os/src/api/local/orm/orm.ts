// // src/api/local/orm/orm.ts
// import { getDatabase } from './db';
// import Database from '@tauri-apps/plugin-sql';
// // const userORM = createORM<{ id: number; name: string; role: string }>('users');
// // // 查询并分页
// // await userORM.where({ role: 'admin' }).page(1, 10).select();
// // // 更新数据
// // await userORM.save({ role: 'admin' }).where({ id: 1 }).select();
// // // 统计数量
// // const total = await userORM.where({ role: 'admin' }).count();
// type Model = Record<string, any>;

// interface ORM<T extends Model> {
//   create: (data: T) => Promise<void>;
//   update: (id: any, data: Partial<T>) => Promise<void>;
//   save: (data: Partial<T>) => ORMQueryBuilder<T>;
//   select: () => Promise<T[]>;
//   find: () => Promise<T | null>;
//   findById: (id: any) => Promise<T | null>;
//   delete: (id: any) => Promise<void>;
//   where: (conditions: Record<string, any>) => ORMQueryBuilder<T>;
//   count: () => Promise<number>;
//   page: (page: number, pageSize: number) => ORMQueryBuilder<T>;
// }

// interface ORMQueryBuilder<T extends Model> {
//   where: (conditions: Record<string, any>) => ORMQueryBuilder<T>;
//   count: () => Promise<number>;
//   page: (page: number, pageSize: number) => ORMQueryBuilder<T>;
//   select: () => Promise<T[]>;
//   find: () => Promise<T | null>; 
// }

// function createORMQueryBuilder<T extends Model>(
//   db: Database,
//   table: string
// ): ORMQueryBuilder<T> {
//   let conditions: Record<string, any> = {};
//   let limit: number | null = null;
//   let offset: number | null = null;

//   function buildWhereClause(): { clause: string; values: any[] } {
//     const clauses = [];
//     const values = [];

//     for (const key in conditions) {
//       if (conditions.hasOwnProperty(key)) {
//         clauses.push(`${key} = $${values.length + 1}`);
//         values.push(conditions[key]);
//       }
//     }

//     const whereClause = clauses.length > 0 ? `WHERE ${clauses.join(' AND ')}` : '';
//     return { clause: whereClause, values };
//   }

//   const queryBuilder: ORMQueryBuilder<T> = {
//     where: (newConditions: Record<string, any>): ORMQueryBuilder<T> => {
//       conditions = { ...conditions, ...newConditions };
//       return queryBuilder;
//     },
//     count: async (): Promise<number> => {
//       const { clause, values } = buildWhereClause();
//       const query = `SELECT COUNT(*) as count FROM ${table} ${clause}`;
//       const result: any = await db.select(query, values);
//       return result[0]?.count || 0;
//     },
//     page: (page: number, pageSize: number): ORMQueryBuilder<T> => {
//       const newLimit = pageSize;
//       const newOffset = (page - 1) * pageSize;
//       limit = newLimit;
//       offset = newOffset;
//       return queryBuilder;
//     },
//     select: async (): Promise<T[]> => {
//       const { clause, values } = buildWhereClause();

//       let limitClause = '';
//       if (limit !== null && offset !== null) {
//         limitClause = `LIMIT ${limit} OFFSET ${offset}`;
//       }

//       const query = `SELECT * FROM ${table} ${clause} ${limitClause}`;
//       const result = await db.select(query, values);
//       return result as T[];
//     },
//     find: async (): Promise<T | null> => {
//       const { clause, values } = buildWhereClause();
    
//       // 加 LIMIT 1 只取一条记录
//       const query = `SELECT * FROM ${table} ${clause} LIMIT 1`;
//       const result:any = await db.select(query, values);
    
//       return result.length > 0 ? (result[0] as T) : null;
//     }
//   };

//   return queryBuilder;
// }

// function buildWhereClauseFromConditions(
//   conditions: Record<string, any>
// ): { clause: string; values: any[] } {
//   const clauses = [];
//   const values = [];

//   for (const key in conditions) {
//     if (conditions.hasOwnProperty(key)) {
//       clauses.push(`${key} = $${values.length + 1}`);
//       values.push(conditions[key]);
//     }
//   }

//   const whereClause = clauses.length > 0 ? `WHERE ${clauses.join(' AND ')}` : '';
//   return { clause: whereClause, values };
// }

// async function executeUpdate<T extends Model>(
//   db: Database,
//   table: string,
//   data: Partial<T>,
//   conditions: Record<string, any>
// ): Promise<void> {
//   const now = new Date().toISOString();
//   const timestampedData = { ...data, updated_at: now };
//   const setClause = Object.keys(timestampedData)
//     .map((key, i) => `${key} = $${i + 1}`)
//     .join(', ');
//   const values = Object.values(timestampedData);

//   const { clause: whereClause, values: whereValues } = buildWhereClauseFromConditions(conditions);

//   const query = `UPDATE ${table} SET ${setClause} ${whereClause}`;
//   await db.execute(query, [...values, ...whereValues]);
// }

// export function createORM<T extends Model>(table: string): ORM<T> {
//   const db = getDatabase();
//   if (!db) {
//     throw new Error('Database not initialized');
//   }

//   const ormInstance: ORM<T> = {
//     create: async (data: T): Promise<void> => {
//       const now = new Date().toISOString();
//       const timestampedData = {
//         ...data,
//         created_at: now,
//         updated_at: now,
//       };

//       const columns = Object.keys(timestampedData).join(', ');
//       const placeholders = Object.keys(timestampedData)
//         .map((_, i) => `$${i + 1}`)
//         .join(', ');
//       const values = Object.values(timestampedData);

//       const query = `INSERT INTO ${table} (${columns}) VALUES (${placeholders})`;
//       await db.execute(query, values);
//     },

//     update: async (id: any, data: Partial<T>): Promise<void> => {
//       const now = new Date().toISOString();
//       const timestampedData = {
//         ...data,
//         updated_at: now,
//       };

//       const setClause = Object.keys(timestampedData)
//         .map((key, i) => `${key} = $${i + 1}`)
//         .join(', ');

//       const values = Object.values(timestampedData);
//       const query = `UPDATE ${table} SET ${setClause} WHERE id = $${values.length + 1}`;
//       values.push(id);

//       await db.execute(query, values);
//     },

//     save: (data: Partial<T>): ORMQueryBuilder<T> => {
//       const queryBuilder = createORMQueryBuilder<T>(db, table);

//       const originalWhere = queryBuilder.where;
//       queryBuilder.where = (conditions: Record<string, any>): ORMQueryBuilder<T> => {
//         const updatedQueryBuilder = originalWhere(conditions);

//         updatedQueryBuilder.select = async (): Promise<T[]> => {
//           await executeUpdate(db, table, data, conditions);
//           return [];
//         };

//         return updatedQueryBuilder;
//       };

//       return queryBuilder;
//     },

//     select: (): Promise<T[]> => {
//       return createORMQueryBuilder<T>(db, table).select();
//     },
//     find: async (): Promise<T | null> => {
//       return createORMQueryBuilder<T>(db, table).find();
//     },
//     findById: async (id: any): Promise<T | null> => {
//       const query = `SELECT * FROM ${table} WHERE id = $1`;
//       const result: any = await db.select(query, [id]);
//       return result.length > 0 ? (result[0] as T) : null;
//     },

//     delete: async (id: any): Promise<void> => {
//       const query = `DELETE FROM ${table} WHERE id = $1`;
//       await db.execute(query, [id]);
//     },

//     where: (conditions: Record<string, any>): ORMQueryBuilder<T> => {
//       console.log('where', conditions);
//       return createORMQueryBuilder<T>(db, table).where(conditions);
//     },

//     count: (): Promise<number> => {
//       return createORMQueryBuilder<T>(db, table).count();
//     },

//     page: (page: number, pageSize: number): ORMQueryBuilder<T> => {
//       return createORMQueryBuilder<T>(db, table).page(page, pageSize);
//     },
//   };

//   return ormInstance;
// }