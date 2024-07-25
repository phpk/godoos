import { defineStore } from 'pinia'
import emojiList from "@/assets/emoji.json"
import { ref, toRaw, inject } from "vue";
import { db } from './db'
import { System,dirname } from "@/system";
import { getSystemConfig } from "@/system/config";
import { isBase64, base64ToBuffer } from "@/util/file";
import { notifyError, notifySuccess } from "@/util/msg";
export const useLocalChatStore = defineStore('localChatStore', () => {
  const config = getSystemConfig();
  const sys = inject<System>("system");
  const userList:any = ref([])
  const msgList:any = ref([])
  const contentList:any = ref([])
  const hostInfo:any = ref({})
  const showChooseFile = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(50)
  const navList = ref([
    { index: 1, lable: "消息列表", icon: "ChatDotRound", type:"success" },
    { index: 2, lable: "用户列表", icon: "UserFilled", type: "info" },
  ])
  const navId = ref(1)
  const sendInfo = ref("")
  const chatTargetId = ref(0)
  const chatTargetIp = ref("")
  const handleSelect = (key: number) => {
    navId.value = key;
  };
  const setChatId = async(ip : string) => {
    //console.log(ip)
    //chatTargetId.value = id
    chatTargetIp.value = ip
    const data = await db.get("chatuser", {ip : ip})
    if(!data)return;
    chatTargetId.value = data.id
    clearContentList(data.id)
    currentPage.value = 1
    await getMsgList()
  }
  const initContentList = async () => {
    const list:any = {}
    const msgAll = await db.getAll('chatmsg')
    msgAll.forEach((d : any) => {
      if(!d.isMe){
        if (!list[d.targetIp]){
          list[d.targetIp] = []
        }
        list[d.targetIp].push(d)
      }
      
    })
    const res = []
    for(const p in list){
      const chatArr = list[p]
      let readNum = 0
      chatArr.forEach((d:any) => {
        if(!d.isRead){
          readNum++
        }
      })
      const last = chatArr.pop()
      last.readNum = readNum
      res.push(last)
    }
    contentList.value = res.sort((a, b) => b.createdAt - a.createdAt);
  }
  const clearContentList = (targetId:number) => {
    contentList.value.forEach((d : any) => {
      if(d.targetId === targetId) {
        d.readNum = 0
      }
    })
  }
  const clearMsg = async () => {
    if (chatTargetIp.value === '') return
    await db.deleteByField('chatmsg','targetIp', chatTargetIp.value)
    msgList.value = []
  }

  const updateContentList = async (msg:any) => {
    if(msg.isMe)return;
    const has = contentList.value.find((d:any) => d.targetIp === msg.targetIp)
    if(has){   
      contentList.value.forEach((d: any, index: number) => {
        if (d.targetIp === msg.targetIp) {
          if (!msg.isRead) {
            msg.readNum = d.readNum + 1;
          } else {
            msg.readNum = 0;
          }
          // 直接替换数组中的元素以触发更新
          contentList.value.splice(index, 1, msg);
        }
      });
      //console.log(contentList.value)
    }else{
      if(msg.isRead){
        msg.readNum = 0
      }else{
        msg.readNum = 1
      }
      contentList.value.unshift(msg)
      //console.log(contentList.value)
    }
    contentList.value = contentList.value.sort((a: any, b: any) => b.createdAt - a.createdAt)
  }
  const init = async() => {
    await getUserList()
    await initUserList()
    await initContentList()
  }
  const initUserList = async() => {
    if(userList.value.length > 0) {
      const updates : any  = []
      userList.value.forEach((d: any) => {
        if (d.isOnline) {
          updates.push({
            key: d.id,
            changes: {
              isOnline: false
            }
          })
        }
      });
      await db.table('chatuser').bulkUpdate(updates)
    }
  }
  const refreshUserList = async () => {
    await db.clear('chatuser')
    userList.value = []
  }
  const getMsgList = async () => {
    if (chatTargetId.value < 1)return;
    //msgList.value = await db.getByField('chatmsg', 'targetId', chatTargetId.value)
    //msgList.value = await db.pageSearch('chatmsg', currentPage.value, pageSize.value, { targetId: chatTargetId.value })
    const offset = (currentPage.value - 1) * pageSize.value
    const list = await db.table('chatmsg')
      .where({ targetIp: chatTargetIp.value })
      .desc()
      .offset(offset)
      .limit(pageSize.value)
      .toArray();
    list.sort((a:any,b:any) => a.id > b.id)
    msgList.value = list;
  }
  const moreMsgList = async() => {
    if (chatTargetId.value < 1) return;
    //const list = await db.pageSearch('chatmsg', currentPage.value + 1, pageSize.value, { targetId: chatTargetId.value })
    const offset = currentPage.value * pageSize.value
    const list = await db.table('chatmsg')
      .where({ targetIp: chatTargetIp.value })
      .desc()
      .offset(offset)
      .limit(pageSize.value)
      .toArray();
    if(list && list.length > 0) {
      list.sort((a: any, b: any) => a.id > b.id)
      currentPage.value = currentPage.value + 1
      msgList.value = [...list, ...msgList.value];
    }
  }
  const getUserList = async () => {
    const list = await db.getAll('chatuser')
    list.sort((a: any, b: any) => a.updatedAt > b.updatedAt)
    userList.value = list
  }

  const setUserList = async (data:any) => {
    //console.log(data)
    if(data.length < 1){
      return
    }
    hostInfo.value = data[0]
    data.shift()
    const ips:any = []
    userList.value.forEach((d : any) => {
      ips.push(d.ip)
    });
    const has:any = []
    const nothas:any = []
    data.forEach((d:any) => {
      if(ips.includes(d.ip)){
        has.push(d.ip)
      }else{
        nothas.push(d)
      }
    })
    if(has.length > 0) {
      const updates:any = []
      userList.value.forEach((d: any) => {
        if(has.includes(d.ip)){
          updates.push({
            key : d.id,
            changes : {
              isOnline : true,
              updatedAt:Date.now()
            }
          })
        }
      });
      await db.table('chatuser').bulkUpdate(updates)
    }
    if(nothas.length > 0) {
      nothas.forEach((d:any) => {
          d.isOnline = true
          d.username = data.hostname
          d.createdAt = Date.now()
          d.updatedAt = Date.now()
      })
      await db.table('chatuser').bulkAdd(nothas)
    }
    await getUserList()
  }
  const addFile = async (data:any) => {
    const targetUser:any = await getTargetUser(data)
    const files:any = []
    data.fileList.forEach((d:any) => {
      d.save_path = d.save_path.replace(/\\/g, "/");
      files.push({
        name : d.name,
        path : d.save_path,
        ext : d.save_path.split('.').pop(),
        content: d.content
      })
    })
    const saveMsg: any = {
      type: 'file',
      targetId: targetUser.id,
      targetIp: targetUser.ip,
      content: files,
      reciperInfo: data.senderInfo,
      createdAt: Date.now(),
      isMe: false,
      isRead: false,
      status: 'reciped'
    }
    if (targetUser.id === chatTargetId.value) {
      saveMsg.readAt = Date.now()
      saveMsg.isRead = true
      msgList.value.push(saveMsg)
    }
    console.log(saveMsg)
    await db.addOne('chatmsg', saveMsg)
    //await getMsgList()
    
    await updateContentList(saveMsg)
    if (config.storeType === 'browser') {
      await storeFile(files)
    }
    handleSelect(1)
  }
  const storeFile = async(fileList : any) => {
    if (fileList.length < 1) return;
    console.log(fileList)
    for (let i = 0; i < fileList.length; i++) {
      let content = fileList[i].content
      if(typeof content === 'string') {
        if(isBase64(content)){
          content = base64ToBuffer(content);
        }
        const path = dirname(fileList[i].path)
        //console.log(path)
        await sys?.fs.mkdir(path);
        await sys?.fs.writeFile(fileList[i].path, content);
      }
    }
  }
  const getTargetUser = async (data:any) => {
    let targetUser:any = userList.value.find((d: any) => d.ip === data.senderInfo.ip)
    if (!targetUser){
      targetUser = {
        isOnline : true,
        ip:data.senderInfo.ip,
        hostname: data.senderInfo.hostname,
        username : data.senderInfo.hostname,
        createdAt : Date.now(),
        updatedAt : Date.now()
      }
      targetUser.id = await db.addOne("chatuser",targetUser)
      userList.value.unshift(targetUser)
    }
    return targetUser
  }
  const addText = async (data:any) => {
    const targetUser:any = await getTargetUser(data)

    const saveMsg:any = {
      type: 'text',
      targetId: targetUser.id,
      targetIp: targetUser.ip,
      content: data.content,
      reciperInfo: data.senderInfo,
      createdAt: Date.now(),
      isMe: false,
      isRead: false,
      status: 'reciped'
    }
    if (targetUser.id === chatTargetId.value){
      saveMsg.readAt = Date.now()
      saveMsg.isRead = true
      msgList.value.push(saveMsg)
    }
    //console.log(saveMsg)
    await db.addOne('chatmsg', saveMsg)
    //await getMsgList()
    
    await updateContentList(saveMsg)
    handleSelect(1)
  }
  const sendMsg = async () => {
    if(chatTargetId.value < 1) {
      return
    }
    const saveMsg:any = {
      type : 'text',
      targetId : chatTargetId.value,
      targetIp: chatTargetIp.value,
      content: sendInfo.value.trim(),
      senderInfo: toRaw(hostInfo.value),
      createdAt:Date.now(),
      isMe:true,
      isRead:false,
      status:'sending'
    }
    //console.log(saveMsg)
    const msgId = await db.addOne('chatmsg', saveMsg)
    //await getMsgList()
    msgList.value.push(saveMsg)
    const targetUser = userList.value.find((d: any) => d.id === chatTargetId.value)
    //console.log(targetUser)
    if(targetUser.isOnline) {
      const postUrl = `http://${targetUser.ip}:56780/localchat/message`
      const completion = await fetch(postUrl, {
        method: "POST",
        body: JSON.stringify(saveMsg),
      })
      if (!completion.ok) {
        console.log(completion)
      }else{
        saveMsg.isRead = true
        saveMsg.status = 'sended'
        saveMsg.readAt = Date.now()
        await db.update('chatmsg', msgId, saveMsg)

      }
    }
    sendInfo.value = ""
    await updateContentList(saveMsg)
    
  }
  //上传文件资源
  async function uploadFile(paths:any) {
    if (chatTargetId.value < 1) {
      return
    }
    const targetUser = userList.value.find((d: any) => d.id === chatTargetId.value)
    if (!targetUser.isOnline) {
      notifyError("The user is not online!");
      return;
    }
    if (!hostInfo.value || !hostInfo.value.ip) {
      notifyError("Please wait for a moment");
      return;
    }
    //console.log(paths)
    const formData = new FormData();
    const errstr:any = []
    const files:any = []
    for (let i = 0; i < paths.length; i++) {
      const content = await sys?.fs.readFile(paths[i]);
      let blobContent;
      if(!content || content == ''){
        errstr.push(paths[i] + " is empty")
        continue
      }
      if (content instanceof ArrayBuffer) {
        blobContent = new Blob([content]);
      }
      else if (typeof content === 'string') {
        if (isBase64(content)) {
          const base64 = base64ToBuffer(content);
          blobContent = new Blob([base64]);
        } else {
          blobContent = new Blob([content], { type: "text/plain;charset=utf-8" });
        }
      }
      else {
        errstr.push(paths[i] + " type is error")
        continue
      }
      const fileName = paths[i].split("/").pop()
      files.push({
        name: fileName,
        path: paths[i],
        ext : fileName.split(".").pop(),
      })
      //files.push(blobContent);
      formData.append(`files`, blobContent, fileName);
    }
    if(errstr.length > 0) {
      errstr.forEach((d:any) => {
        notifyError(d);
      })
      return
    }
    //formData.append("files", files);
    formData.append("ip", hostInfo.value.ip);
    formData.append("hostname", hostInfo.value.hostname);
    //console.log(formData)
    const postUrl = `http://${targetUser.ip}:56780/localchat/upload`
    const res = await fetch(postUrl, {
      method: "POST",
      body: formData,
    });
    if (!res.ok) {
      console.log(res);
      notifyError("Upload error!");
      return;
    }
    const saveMsg: any = {
      type: 'file',
      targetId: targetUser.id,
      targetIp: targetUser.ip,
      content: files,
      reciperInfo: toRaw(targetUser),
      createdAt: Date.now(),
      isMe: false,
      isRead: true,
      status: 'reciped'
    }
    console.log(saveMsg)
    await db.addOne('chatmsg', saveMsg)
    msgList.value.push(saveMsg)

    notifySuccess("upload success!");
  }
  return { 
    userList, 
    navList, 
    sendInfo,
    navId, 
    chatTargetId,
    chatTargetIp,
    msgList,
    hostInfo,
    contentList,
    emojiList,
    showChooseFile,
    pageSize,
    init,
    setUserList,
    getUserList,
    handleSelect,
    setChatId,
    sendMsg,
    addText,
    addFile,
    uploadFile,
    moreMsgList,
    refreshUserList,
    clearMsg
   }
})
