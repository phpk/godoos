import Dexie from 'dexie'

export type ChatTable =  'chatuser' | 'chatmsg' | 'chatfile'

export const dbInit:any = new Dexie('GodoOSDatabase');
dbInit.version(1).stores({
  chatuser:'++id,ip,hostname,userName,avatar,mobile,nickName,isOnline,updatedAt,createdAt',
  chatmsg:'++id,targetId,targetIp,senderInfo,reciperInfo,content,type,status,isRead,isMe,readAt,createdAt',
  chatfile:'++id,msgId,fileName,fileSize,fileType,fileUrl,filePath,fileExt'
});
export const db = {

  async getMaxId(tableName: ChatTable) {
    const data = await dbInit[tableName].orderBy('id').reverse().first()
    if (!data) {
      return 0
    } else {
      return data.id
    }
  },
  async getInsertId(tableName: ChatTable) {
    const id: any = await this.getMaxId(tableName)
    return id + 1
  },
  async getPage(tableName: ChatTable, page?: number, size?: number) {
    page = (!page || page < 1) ? 1 : page
    size = size ? size : 10
    const offset = (page - 1) * size
    return dbInit[tableName]
      .orderBy("id")
      .reverse()
      .offset(offset)
      .limit(size)
      .toArray();
  },
  async getAll(tableName: ChatTable) {
    return dbInit[tableName].toArray()
  },
  async count(tableName: ChatTable) {
    return dbInit[tableName].count()
  },
  async countSearch(tableName: ChatTable, whereObj?: any) {
    return dbInit[tableName].where(whereObj).count()
  },
  async pageSearch(tableName: ChatTable, page?: number, size?: number, whereObj?: any) {
    page = (!page || page < 1) ? 1 : page
    size = size ? size : 10
    const offset = (page - 1) * size
    //console.log(whereObj)
    return dbInit[tableName]
      .where(whereObj)
      .reverse()
      .offset(offset)
      .limit(size)
      .toArray();
  },
  async filter(tableName: ChatTable, filterFunc : any) {
    return dbInit[tableName].filter(filterFunc).toArray()
  },
  table(tableName: ChatTable) {
    return dbInit[tableName]
  },
  async getOne(tableName: ChatTable, Id: number) {
    return dbInit[tableName].get(Id)
  },
  async getRow(tableName: ChatTable, fieldName: string, val: any){
    return dbInit[tableName].where(fieldName).equals(val).first()
  },
  async get(tableName: ChatTable, whereObj : any) {
    //console.log(whereObj)
    const data = await dbInit[tableName].where(whereObj).first()
    //console.log(data)
    return data? data : false
  },
  async rows(tableName: ChatTable, whereObj: any) {
    return dbInit[tableName].where(whereObj).toArray()
  },
  async field(tableName: ChatTable, whereObj: any, field: string) {
    const data = await this.get(tableName, whereObj)
    return data ? data[field] : false
  },
  async getValue(tableName: ChatTable, fieldName: string, val: any, fName : string) {
    const row = await this.getRow(tableName, fieldName, val);
    return row[fName]
  },
  async getByField(tableName: ChatTable, fieldName: string, val: any) {
    return dbInit[tableName].where(fieldName).equals(val).toArray()
  },
  async addOne(tableName: ChatTable, data: any) {
    return dbInit[tableName].add(data)
  },
  async addAll(tableName: ChatTable, data: any) {
    return dbInit[tableName].bulkAdd(data)
  },
  async update(tableName: ChatTable, Id?: number, updates?: any) {
    return dbInit[tableName].update(Id, updates)
  },
  async modify(tableName: ChatTable, fieldName: string, val: any, updates: any) {
    return dbInit[tableName].where(fieldName).equals(val).modify(updates)
  },
  async delete(tableName: ChatTable, Id?: number) {
    return dbInit[tableName].delete(Id)
  },
  async deleteByField(tableName: ChatTable, fieldName: string, val: any) {
    return dbInit[tableName].where(fieldName).equals(val).delete()
  },
  async clear(tableName: ChatTable) {
    return dbInit[tableName].clear()
  },
}
