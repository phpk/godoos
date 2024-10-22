import { defineStore } from 'pinia'
import emojiList from "@/assets/emoji.json"
import { ref, toRaw } from "vue";
import { db } from './db'
import { getSystemConfig, setSystemKey } from "@/system/config";
import { isValidIP } from "@/util/common";
import { notifyError, notifySuccess } from "@/util/msg";
export const useLocalChatStore = defineStore('localChatStore', () => {
  const config = getSystemConfig();
  //const sys = inject<System>("system");
  const userList: any = ref([])
  const msgList: any = ref([])
  const contentList: any = ref([])
  //const OutUserList: any = ref([])
  const hostInfo: any = ref({})
  const showChooseFile = ref(false)
  const currentPage = ref(1)
  const pageSize = ref(50)
  const navList = ref([
    { index: 1, lable: "消息列表", icon: "ChatDotRound", type: "success" },
    { index: 2, lable: "用户列表", icon: "UserFilled", type: "info" },
  ])
  const navId = ref(1)
  const sendInfo:any = ref()
  const chatTargetId = ref(0)
  const chatTargetIp = ref("")
  const showAddUser = ref(false)
  const handlerMessage = (data : any) => {
    //console.log(data)
    if(data.onlines){
      const ips = []
      for(let ip in data.onlines){
        const info = data.onlines[ip]
        if (info.ip && isValidIP(info.ip)) {
          ips.push(info)
        }
      }
      setUserList(ips);
    }
    if(data.messages){
      const apiUrl = `${config.apiUrl}/localchat/viewimage?img=`
      for(let ip in data.messages){
        const msgList:any = data.messages[ip]
        if(!msgList || msgList.length < 1)return;
        msgList.forEach((msg: any) => {
          //console.log(msg)
          if (msg.type === "text") {
            msg.message = msg.message.replaceAll("\\n", "\n")
            //console.log(msg)
            addText(msg)
          }
          else if (msg.type === "image"){
            msg.message = msg.message.map((d: any) => `${apiUrl}${encodeURIComponent(d)}`)
            //console.log(msg)
            addText(msg)
          }
          else if (msg.type === "fileSending"){
            addText(msg)
          }
          else if( msg.type === "fileCannel") {
            changeMsg(msg)
          }
          else if( msg.type === "fileAccessed") {
            msg.message = msg.message.msgId
            changeMsg(msg)
          }
        })
      }
    }
  }
  const handleSelect = (key: number) => {
    navId.value = key;
  };
  const setChatId = async (ip: string) => {
    //console.log(ip)
    //chatTargetId.value = id
    chatTargetIp.value = ip
    const data = await db.get("chatuser", { ip: ip })
    if (!data) return;
    chatTargetId.value = data.id
    clearContentList(data.id)
    currentPage.value = 1
    await getMsgList()
  }
  const initContentList = async () => {
    const list: any = {}
    const msgAll = await db.getAll('chatmsg')
    msgAll.forEach((d: any) => {
      if (!d.isMe) {
        if (!list[d.targetIp]) {
          list[d.targetIp] = []
        }
        list[d.targetIp].push(d)
      }

    })
    const res = []
    for (const p in list) {
      const chatArr = list[p]
      let readNum = 0
      chatArr.forEach((d: any) => {
        if (!d.isRead) {
          readNum++
        }
      })
      const last = chatArr.pop()
      last.readNum = readNum
      res.push(last)
    }
    contentList.value = res.sort((a, b) => b.createdAt - a.createdAt);
  }
  const clearContentList = (targetId: number) => {
    contentList.value.forEach((d: any) => {
      if (d.targetId === targetId) {
        d.readNum = 0
      }
    })
  }
  const clearMsg = async () => {
    if (chatTargetIp.value === '') return
    await db.deleteByField('chatmsg', 'targetIp', chatTargetIp.value)
    msgList.value = []
  }

  const updateContentList = async (msg: any) => {
    if (msg.isMe) return;
    const has = contentList.value.find((d: any) => d.targetIp === msg.targetIp)
    if (has) {
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
    } else {
      if (msg.isRead) {
        msg.readNum = 0
      } else {
        msg.readNum = 1
      }
      contentList.value.unshift(msg)
      //console.log(contentList.value)
    }
    contentList.value = contentList.value.sort((a: any, b: any) => b.createdAt - a.createdAt)
  }
  const init = async () => {
    await getUserList()
    await initUserList()
    await initContentList()
  }
  const initUserList = async () => {
    if (userList.value.length > 0) {
      const updates: any = []
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
    if (chatTargetId.value < 1) return;
    const offset = (currentPage.value - 1) * pageSize.value
    const list = await db.table('chatmsg')
      .where({ targetIp: chatTargetIp.value })
      // .desc()
      .offset(offset)
      .limit(pageSize.value)
      .toArray();
    list.sort((a: any, b: any) => a.id > b.id)
    msgList.value = list;
  }
  const moreMsgList = async () => {
    if (chatTargetId.value < 1) return;
    //const list = await db.pageSearch('chatmsg', currentPage.value + 1, pageSize.value, { targetId: chatTargetId.value })
    const offset = currentPage.value * pageSize.value
    const list = await db.table('chatmsg')
      .where({ targetIp: chatTargetIp.value })
      // .desc()
      .offset(offset)
      .limit(pageSize.value)
      .toArray();
    if (list && list.length > 0) {
      list.sort((a: any, b: any) => a.id > b.id)
      currentPage.value = currentPage.value + 1
      msgList.value = [...list, ...msgList.value];
    }
  }
  const getUserList = async () => {
    const list = await db.getAll('chatuser')
    //const list = [...listAll, ...OutUserList.value]
    let uniqueIpMap = new Map<string, any>();

    // 遍历 list 并添加 IP 地址到 Map 中
    list.forEach((item: any) => {
      uniqueIpMap.set(item.ip, item);
    });

    // 将 Map 转换回数组
    const uniqueIpList: any = Array.from(uniqueIpMap.values());
    uniqueIpList.sort((a: any, b: any) => a.updatedAt > b.updatedAt)
    userList.value = uniqueIpList
  }

  const setUserList = async (data: any) => {
    //console.log(data)
    if (data.length < 1) {
      return
    }
    const existingIps = new Set(userList.value.map((d : any) => d.ip));
    const updates: any[] = [];
    const newEntries: any[] = [];
    // 创建一个映射表，将 ip 映射到 userList 中的完整对象
    const userMap = new Map(
      userList.value.map((user:any) => [user.ip, user])
    );
    //console.log(existingIps)
    data.forEach((d : any) => {
      const existingUser:any = userMap.get(d.ip);
      if (existingUser && existingIps.has(d.ip)) {
        updates.push({
          key: existingUser.id,
          changes: {
            isOnline: true,
            hostname:d.hostname,
            username: d.hostname,
            updatedAt: Date.now()
          }
        });
      } else {
        newEntries.push({
          ip: d.ip,
          isOnline: true,
          hostname: d.hostname,
          username: d.hostname,
          createdAt: Date.now(),
          updatedAt: Date.now()
        });
      }
    });
    //console.log(updates)
    //console.log(newEntries)
    if (updates.length > 0) {
      await db.table('chatuser').bulkUpdate(updates);
    }

    if (newEntries.length > 0) {
      await db.table('chatuser').bulkAdd(newEntries);
    }
    await getUserList()
  }
 
  const getTargetUser = async (data: any) => {
    let targetUser: any = userList.value.find((d: any) => d.ip === data.ip)
    if (!targetUser) {
      targetUser = {
        isOnline: true,
        ip: data.ip,
        hostname: data.hostname,
        username: data.hostname,
        createdAt: Date.now(),
        updatedAt: Date.now()
      }
      targetUser.id = await db.addOne("chatuser", targetUser)
      userList.value.unshift(targetUser)
    }
    return targetUser
  }
  const addText = async (data: any) => {
    const targetUser: any = await getTargetUser(data)

    const saveMsg: any = {
      type: data.type,
      targetId: targetUser.id,
      targetIp: targetUser.ip,
      content: data.message,
      reciperInfo:{
        hostname: data.hostname,
        username: data.hostname
      },
      createdAt: Date.now(),
      isMe: false,
      isRead: false,
      status: 'reciped'
    }
    // console.log(saveMsg)
    // console.log(chatTargetId.value)
    if (targetUser.ip === chatTargetIp.value) {
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
  const sendMsg = async (type:string) => {
    if (chatTargetId.value < 1) {
      return
    }
    const content = toRaw(sendInfo.value)
    let saves:any
    if (type === 'image') {
      const apiUrl = `${config.apiUrl}/localchat/viewimage?img=`
      saves = content.map((d: any) => `${apiUrl}${encodeURIComponent(d)}`)
    }else{
      saves = content
    }
    const saveMsg: any = {
      type: type,
      targetId: chatTargetId.value,
      targetIp: chatTargetIp.value,
      content: saves,
      createdAt: Date.now(),
      isMe: true,
      isRead: false,
      status: 'sending'
    }
    //console.log(saveMsg)
    const msgId = await db.addOne('chatmsg', saveMsg)
    //await getMsgList()
    
    const targetUser = userList.value.find((d: any) => d.ip === chatTargetIp.value)
    //console.log(targetUser)
    const messages = {
      type: type,
      message: content,
      ip: saveMsg.targetIp
    }
    if (targetUser.isOnline) {
      let postUrl = `${config.apiUrl}/localchat/message`
      if(type === 'applyfile'){
        messages.message = {
          fileList: messages.message,
          msgId: msgId,
          status: 'apply'
        }
        postUrl = `${config.apiUrl}/localchat/applyfile`
      }
      if(type === 'image'){
        postUrl = `${config.apiUrl}/localchat/sendimage`
      }
      
      const completion = await fetch(postUrl, {
        method: "POST",
        body: JSON.stringify(messages),
      })
      //console.log(completion)
      if (!completion.ok) {
        console.log(completion)
        notifyError("发送失败!")
      } else {
        saveMsg.content = messages.message
        saveMsg.isRead = true
        saveMsg.status = 'sended'
        saveMsg.readAt = Date.now()
        await db.update('chatmsg', msgId, saveMsg)
      }
     
      await updateContentList(saveMsg)
    }else{
      notifyError("对方不在线!")
    }
    msgList.value.push(saveMsg)
    sendInfo.value = ""
    
  }
  async function cannelFile(item:any){   
    const messages = {
      type: 'cannelFile',
      message: item.content.msgId,
      ip: item.targetIp
    }
    const postUrl = `${config.apiUrl}/localchat/cannelfile`
    const coms = await fetch(postUrl, {
      method: "POST",
      body: JSON.stringify(messages),
    })
    if (!coms.ok) {
      console.log(coms)
      notifyError("确认失败!")
    } else {
      item.content.status = 'cannel'
      console.log(item)
      await db.update('chatmsg', item.id, toRaw(item))
      await updateContentList(item)
      notifySuccess("确认成功!")
    }
  }
  async function changeMsg(msg:any){
    //console.log(msg)
    const msgId = msg.message
    const item = await db.getOne('chatmsg', msgId)
    //console.log(item)
    item.content.status = msg.type
    await db.update('chatmsg', item.id, item)
    await getMsgList()
  }
  async function accessFile(item:any){
    const messages = {
      type: 'accessFile',
      message: item.content,
      ip: item.targetIp
    }
    const postUrl = `${config.apiUrl}/localchat/accessfile`
    const coms = await fetch(postUrl, {
      method: "POST",
      body: JSON.stringify(messages),
    })
    if (!coms.ok) {
      //console.log(coms)
      notifyError("确认失败!")
    } else {
      const res = await coms.json()
      console.log(res)
      if(res.code > -1){
        item = toRaw(item)
          item.content.path = res.data.path
          item.content.status = 'accessed'
          console.log(item)
          await db.update('chatmsg', item.id, item)
          await getMsgList()
      }else{
        notifyError(res.message)
      }
    }
  }
  async function saveConfig(conf:any){
    conf = toRaw(conf)
    //console.log(conf)
    const postUrl = `${config.apiUrl}/localchat/setting`
    const coms = await fetch(postUrl, {
      method: "POST",
      body: JSON.stringify(conf),
    })
    if (!coms.ok) {
      //console.log(coms)
      notifyError("保存失败!")
    } else {
      setSystemKey('chatConf', conf)
      notifySuccess("保存成功!")
      showAddUser.value = false
    }

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
    showAddUser,
    init,
    setUserList,
    getUserList,
    handleSelect,
    setChatId,
    sendMsg,
    addText,
    moreMsgList,
    refreshUserList,
    clearMsg,
    handlerMessage,
    cannelFile,
    accessFile,
    saveConfig
  }
})
