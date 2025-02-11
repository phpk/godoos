import Dexie from 'dexie';

export type ChatTable = 'prompts' | 'modelslabel' | 'modelslist'| 'aichatlist' | 'aichatmsg' | 'chatuser' | 'chatmsg' | 'systemChatRecord' | 'workbenchChatRecord' | 'workbenchChatUser' | 'workbenchSessionList' | 'groupSessionList' | 'workbenchGroupChatRecord' | 'workbenchGroupUserList' | 'workbenchGroupInviteMessage' | 'filePwdBox';

export const dbInit: any = new Dexie('GodoOSDatabase');
dbInit.version(1).stores({
  // ai助手
  prompts: '++id,lang,action,prompt,name,ext,isdef,createdAt,[action+lang]',
  // 模型标签
  modelslabel: '++id,name,zhdesc,endesc,family,chanel,models,action,engine',
  // 模型列表
  modelslist: '++id,model,label,engine,action,status,params,type,isdef,info,created_at',
  // ai对话列表
  aichatlist: '++id,title,model,engine,promptId,prompt,knowledgeId,createdAt',
  // ai对话消息
  aichatmsg: '++id,chatId,role,content,createdAt',
  // 用户列表
  workbenchChatUser: '++id,ip,userName,chatId,avatar,mobile,phone,nickName,isOnline,updatedAt,createdAt',
  // 会话列表
  workbenchSessionList: '++id,avatar,chatId,username,nickname,userId,toUserId,previewMessage,messages,time,createdAt',
  // 聊天记录
  workbenchChatRecord: '++id,[toUserId+userId], [userId+toUserId],toUserId,messages,messageType,time,createdAt,userInfo',
  // 群组会话列表
  groupSessionList: '++id,groupId,chatId,name,message,previewMessage,avatar,createdAt',
  // 群组聊天记录
  workbenchGroupChatRecord: '++id,chatId,userId,to_groupid,messageType,userInfo,message,time,type,createdAt',
  // 群用户列表
  workbenchGroupUserList: '++id, group_id, createdAt, userIdArray',
  // 群组邀请信息消息
  workbenchGroupInviteMessage: '++id, group_id,userId, message, createdAt',
  //系统聊天记录表
  systemChatRecord: '++id,chatId,message,time,createdAt',
  // 用户列表
  chatuser: '++id,ip,hostname,userName,avatar,mobile,nickName,isOnline,updatedAt,createdAt',
  chatmsg: '++id,toUserId,targetIp,senderInfo,reciperInfo,previewMessage,content,type,status,isRead,isMe,readAt,createdAt',
  // 文件密码箱
  filePwdBox: '++id, pwdName, pwd, isDefault'
}).upgrade((tx: {
  workbenchSessionList: any;
  workbenchChatUser: any;
  workbenchGroupChatRecord: any;
  workbenchGroupUserList: any;
  workbenchGroupInviteMessage: any;
  groupSessionList: any;
  workbenchChatRecord: any;
  systemChatRecord: any;
}) => {
  // 手动添加索引
  tx.groupSessionList.addIndex('chatId', (obj: { chatId: any; }) => obj.chatId);
  tx.workbenchSessionList.addIndex('chatId', (obj: { chatId: any; }) => obj.chatId);
  tx.workbenchChatRecord.addIndex('toUserId', (obj: { toUserId: any; }) => obj.toUserId);
  tx.systemChatRecord.addIndex('chatId', (obj: { chatId: any; }) => obj.chatId);
  // 添加复合索引
  tx.workbenchChatRecord.addIndex('[toUserId+userId]', (obj: { toUserId: any; userId: any; }) => [obj.toUserId, obj.userId]);
  tx.workbenchGroupChatRecord.addIndex('chatId', (obj: { chatId: any; }) => obj.chatId);
  tx.workbenchGroupUserList.addIndex('group_id', (obj: { group_id: any; }) => obj.group_id);
  tx.workbenchGroupInviteMessage.addIndex('group_id', (obj: { group_id: any; }) => obj.group_id);
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
    if (whereObj === undefined) {
      return dbInit[tableName].count()
    }
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
  async filter(tableName: ChatTable, filterFunc: any) {
    return dbInit[tableName].filter(filterFunc).toArray()
  },
  table(tableName: ChatTable) {
    return dbInit[tableName]
  },
  async getOne(tableName: ChatTable, Id: number) {
    return dbInit[tableName].get(Id)
  },
  async getRow(tableName: ChatTable, fieldName: string, val: any) {
    return dbInit[tableName].where(fieldName).equals(val).first()
  },
  async get(tableName: ChatTable, whereObj: any) {
    //console.log(whereObj)
    try {
      const data = await dbInit[tableName].where(whereObj).first()
      //console.log(data)
      return data ? data : false
    } catch (error) {
      return false
    }

  },
  async rows(tableName: ChatTable, whereObj: any) {
    return dbInit[tableName].where(whereObj).toArray()
  },

  async field(tableName: ChatTable, whereObj: any, field: string) {
    const data = await this.get(tableName, whereObj)
    return data ? data[field] : false
  },
  async getValue(tableName: ChatTable, fieldName: string, val: any, fName: string) {
    const row = await this.getRow(tableName, fieldName, val);
    return row[fName]
  },

  // 写一个根据id数组的in方法
  async getByIds(tableName: ChatTable, ids: string[]) {
    return dbInit[tableName].where("id").anyOf(ids).toArray()
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

  // 根据whereObj条件删除
  async deleteByWhere(tableName: ChatTable, whereObj: any) {
    return dbInit[tableName].where(whereObj).delete()
  },
  // 获取创建时间最近的记录
  async getLatest(tableName: ChatTable, fieldName: string, val: any) {
    return dbInit[tableName].where(fieldName).equals(val).reverse().first()
  },
  async clear(tableName: ChatTable) {
    return dbInit[tableName].clear()
  },
}
