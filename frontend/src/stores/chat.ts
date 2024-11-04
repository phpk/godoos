import emojiList from "@/assets/emoji.json";
import { fetchGet, fetchPost, getSystemConfig } from '@/system/config';
import { notifySuccess } from "@/util/msg";
import { defineStore } from 'pinia';
import { db } from "./db";

export const useChatStore = defineStore('chatStore', () => {

  interface OnlineUserInfoType {
    id: string;
    ip: string;
    avatar: string;
    online: boolean;
    type: string;
    chatId: string;
    username: string;
    nickname: string;
  }
  // 文件消息类型
  interface ChatMessageType {
    type: string;
    content_type: string;
    time: string;
    userId: number;
    toUserId: number;
    message: string;
    to_groupid: string,
    userInfo: {};
  }

  // 文件发送模型
  const Message: ChatMessageType = {
    type: '',
    content_type: '',
    time: new Date().toISOString(),
    userId: 0,
    toUserId: 0,
    message: '',
    to_groupid: '',
    userInfo: {},
  }

  // 发起群聊对话框显示
  const groupChatInvitedDialogVisible = ref(false);


  const fileSendActive = ref()

  // 群信息设置抽屉
  const groupInfoSettingDrawerVisible = ref(false);
  // 设置群聊邀请对话框状态
  const setGroupChatInvitedDialogVisible = (visible: boolean) => {
    getAllUser()
    groupChatInvitedDialogVisible.value = visible;
  };
  // 设置群信息抽屉状态
  const setGroupInfoDrawerVisible = (visible: boolean) => {
    groupInfoSettingDrawerVisible.value = visible
  }

  // 群名
  const departmentName = ref('');

  const sendInfo: any = ref()
  // 在线用户列表
  const onlineUserList = ref<OnlineUserInfoType[]>([]);

  // 聊天列表
  const chatList: any = ref([]);

  // 聊天消息记录列表
  const chatHistory: any = ref([]);

  // 群组l列表
  const groupList: any = ref([

  ]);
  const targetGroupInfo: any = ref({})
  const activeNames = ref([]);
  const userInfo: any = ref({});
  const showChooseFile = ref(false);
  const currentPage = ref(1);
  const pageSize = ref(50);
  const innerRef = ref(null);
  const scrollbarRef = ref(null);
  const config = getSystemConfig();
  const currentNavId = ref(0);
  const message: any = ref('');
  const targetUserInfo: any = ref({});
  const targetChatId = ref();
  const search = ref('');

  // 群成员列表
  const groupMemberList = ref([])

  // 所有用户列表
  const allUserList = ref([])

  // 部门列表
  const departmentList = ref([

  ])

  const contextMenu = ref({
    visible: false,
    chatMessageId: 0,
    list: [
      { id: 1, label: '复制' },
      { id: 2, label: '删除' },
    ],
    x: 0,
    y: 0
  });


  const initChat = async () => {
    if (config.userInfo.avatar == '') {
      config.userInfo.avatar = '/logo.png';
    }
    userInfo.value = config.userInfo;
    // 初始化用户列表
    initUserList()
    // 获取群列表
    await getGroupList()
    // 初始化聊天列表
    initChatList()
    // 获取部门列表
    getDepartmentList()
  };

  // 获取部门列表
  const getDepartmentList = async () => {
    const res = await fetchGet(userInfo.value.url + "/chat/user/list")
    console.log(res);
    if (res.ok) {
      const data = await res.json();
      departmentList.value = data.data.users;
    }
  }

  // 初始化用户列表
  const initChatList = async () => {
    const userSessionList = await db.getAll("workbenchSessionList");
    // 确保groupList已加载
    if (groupList.value.length > 0) {
      // 合并两个数组
      chatList.value = [...userSessionList, ...groupList.value];
    } else {
      // 如果groupList为空，只使用userSessionList
      chatList.value = [...userSessionList];
    }
  };

  const setCurrentNavId = (id: number) => {
    currentNavId.value = id;
  };

  const sendMessage = async (messageType: string) => {
    console.log("sendMessage:", messageType)
    if (messageType == 'text') {
      await sendTextMessage()
    }
    if (messageType == 'image') {
      sendImageMessage()
    }

    console.log(messageType)
    if (messageType == 'applyfile') {
      await sendFileMessage()
    }
  }

  // 发送图片消息
  const sendImageMessage = async () => {

    // 判断是群聊发送还是单聊发送
    if (targetGroupInfo.value && Object.keys(targetGroupInfo.value).length > 0) {
      console.log('群聊发送文件');
      Message.type = 'group';
      Message.content_type = 'image';
      Message.userId = userInfo.value.id;
      Message.to_groupid = targetGroupInfo.value.group_id;
      Message.message = sendInfo.value[0];
      Message.userInfo = {};
      console.log("群聊发送文件", Message)
    } else if (targetUserInfo.value && Object.keys(targetUserInfo.value).length > 0) {
      console.log('单聊发送文件');
      Message.type = 'user';
      Message.content_type = 'image';
      Message.userId = userInfo.value.id;
      Message.toUserId = targetChatId.value;
      Message.message = sendInfo.value[0];
      Message.to_groupid = targetGroupInfo.value?.group_id || '';
      Message.userInfo = {};
    }
    console.log(Message)
    // 发送文件消息
    const res = await fetchPost(config.userInfo.url + '/chat/send/file', JSON.stringify(Message));
    if (!res.ok) {
      fileSendActive.value = false;
      return;
    }
    const data = await res.json();
    console.log(data);

    // 封装成消息历史记录
    const messageHistory = {
      ...Message,
      isMe: true,
      chatId: Message.type === 'group' ? targetGroupInfo.value.chatId : targetUserInfo.value.chatId,
      avatar: Message.type === 'group' ? targetUserInfo.value.avatar : '',
      createdAt: Date.now()
    };

    console.log("messageHistory------", messageHistory)

    // 添加到聊天记录
    await db.addOne("workbenchGroupChatRecord", messageHistory);
    chatHistory.value.push(messageHistory);
    fileSendActive.value = true;
    notifySuccess('文件发送成功');

    // 两秒后关闭
    setTimeout(() => {
      fileSendActive.value = false;
    }, 2000);
  }

  // 发送文件消息
  const sendFileMessage = async () => {

    // 判断是群聊发送还是单聊发送
    if (targetGroupInfo.value && Object.keys(targetGroupInfo.value).length > 0) {
      console.log('群聊发送文件');
      Message.type = 'group';
      Message.content_type = 'file';
      Message.userId = userInfo.value.id;
      Message.to_groupid = targetGroupInfo.value.group_id;
      Message.message = sendInfo.value[0];
      Message.userInfo = {};
      console.log("群聊发送文件", Message)
    } else if (targetUserInfo.value && Object.keys(targetUserInfo.value).length > 0) {
      console.log('单聊发送文件');
      Message.type = 'user';
      Message.content_type = 'file';
      Message.userId = userInfo.value.id;
      Message.toUserId = targetChatId.value;
      Message.message = sendInfo.value[0];
      Message.to_groupid = targetGroupInfo.value?.group_id || '';
      Message.userInfo = {};
    }
    console.log(Message)
    // 发送文件消息
    const res = await fetchPost(config.userInfo.url + '/chat/send/file', JSON.stringify(Message));
    if (!res.ok) {
      fileSendActive.value = false;
      return;
    }
    const data = await res.json();
    console.log(data);

    // 封装成消息历史记录
    const messageHistory = {
      ...Message,
      isMe: true,
      chatId: Message.type === 'group' ? targetGroupInfo.value.chatId : targetUserInfo.value.chatId,
      avatar: Message.type === 'group' ? targetUserInfo.value.avatar : '',
      createdAt: Date.now()
    };

    console.log("messageHistory------", messageHistory)

    // 添加到聊天记录
    await db.addOne("workbenchGroupChatRecord", messageHistory);
    chatHistory.value.push(messageHistory);
    fileSendActive.value = true;
    notifySuccess('文件发送成功');

    // 两秒后关闭
    setTimeout(() => {
      fileSendActive.value = false;
    }, 2000);
  }

  // 发送文字消息 
  const sendTextMessage = async () => {
    // 判断是群聊发送还是单聊发送
    if (targetGroupInfo.value && Object.keys(targetGroupInfo.value).length) {
      console.log('群聊发送');
      Message.type = 'group'
      Message.content_type = 'text'
      Message.to_groupid = targetGroupInfo.value?.group_id
      Message.message = message.value
      Message.userId = userInfo.value.id
      Message.userInfo = {}
      console.log(Message)

    } else if (targetUserInfo.value && Object.keys(targetUserInfo.value).length > 0) {
      console.log('单聊发送');
      Message.type = 'user'
      Message.content_type = 'text'
      Message.userId = userInfo.value.id
      Message.toUserId = targetChatId.value
      Message.message = message.value
      Message.content_type = 'text'
      Message.to_groupid = targetGroupInfo.value?.group_id || ''
      Message.userInfo = {}
    }


    // 发送消息
    const res = await fetchPost(config.userInfo.url + '/chat/send', JSON.stringify(Message));

    if (res.ok) {
      // 封装成消息历史记录
      var messageHistory
      // 本地存储一份聊天记录
      if (Message.type === 'user') {
        messageHistory = {
          ...Message,
          isMe: true,
          displayName: userInfo.value.nickname,
          chatId: targetUserInfo.value.chatId.toString(),
          avatar: userInfo.value.avatar,
          createdAt: Date.now()
        }
        await db.addOne("workbenchChatRecord", messageHistory);

      } else if (Message.type === 'group') {
        messageHistory = {
          ...Message,
          isMe: true,
          chatId: targetGroupInfo.value.chatId.toString(),
          createdAt: Date.now()
        }
        await db.addOne("workbenchGroupChatRecord", messageHistory);
      }

      // 更新聊天历史
      chatHistory.value.push(messageHistory);

      console.log(await res.json())
      // 清空输入框
      clearMessage();

      // 聊天框滚动到底部
      await setScrollToBottom();
    }

  }


  // 更新聊天和聊天记录
  const changeChatListAndChatHistory = async (data: any) => {

    // 从 conversationList 数据库中查找是否存在对应的会话
    const conversation = await db.getByField("workbenchSessionList", 'userId', data.userId);

    // 准备会话更新数据
    const updatedConversation = {
      userId: data.userId,
      avatar: data.userInfo.avatar, // 如果没有头像使用默认图片
      toUserId: data.toUserId,
      chatId: data.userId,
      type: data.type,
      messages: data.message,
      displayName: data.userInfo.nickname,
      nickname: data.userInfo.nickname,
      time: Date.now().toString,
      // previewMessage: data.message,
      // previewTimeFormat: formatTime(Date.now()), // 时间格式化函数
      createdAt: Date.now()
    };

    if (conversation.length === 0) {
      // 如果会话不存在，则添加到数据库和 chatList
      await db.addOne('workbenchSessionList', updatedConversation);

      chatList.value.push(updatedConversation);
    } else {
      // 如果会话存在，则更新数据库和 chatList
      // 只更新变化的字段，而非全部覆盖,以减少写入数据的量
      await db.update('workbenchSessionList', conversation[0].id, {
        avatar: data.userInfo.avatar,
        nickname: data.userInfo.nickname,
        displayName: data.userInfo.nickname,
        // previewMessage: data.message,
        time: data.time || Date.now().toString,
        // previewTimeFormat: formatTime(Date.now())
      });


      // 更新 chatList 中的对应项
      const existingConversationIndex = chatList.value.findIndex(
        (conv: any) => conv.userId === data.userId
      );
      if (existingConversationIndex !== -1) {
        chatList.value[existingConversationIndex] = updatedConversation;
      } else {
        chatList.value.push(updatedConversation);
      }
    }
  };

  // 点击目标用户时，将其添加到 chatList，并添加到数据库,以及获取到聊天记录
  const addChatListAndGetChatHistory = async (chatId: string) => {
    const chatIdSet = new Set(chatList.value.map((chat: { chatId: any }) => chat.chatId));

    // 如果会话存在于 chatList，则获取聊天记录并更新 chatHistory
    if (chatIdSet.has(chatId)) {
      console.log("存在");
      chatHistory.value = await getHistory(chatId, userInfo.value.id, "user");
      return;
    }
    console.log("不存在")
    // 如果会话不存在，则从用户表中获取该用户的基本信息
    const user = await db.getOne("workbenchChatUser", Number(chatId));

    // 将新用户信息添加到 chatList
    const newChat = {
      type: "user",
      chatId: user.id,
      nickname: user.nickname,
      avatar: user.avatar,
      toUserId: userInfo.value.id,
      time: Date.now(),
      displayName: user.nickname,
    };
    chatList.value.push(newChat);

    // 持久化数据到数据库
    await db.addOne("workbenchSessionList", {
      type: "user",
      chatId: user.id,
      displayName: user.nickname,
      username: user.username,
      nickname: user.nickname,
      avatar: user.avatar,
      toUserId: userInfo.value.id,
      time: Date.now(),
      createdAt: Date.now(),
    });

    // 获取聊天记录
    chatHistory.value = await getHistory(chatId, userInfo.value.id, "user");
  };



  const formatTime = (timestamp: number): string => {
    const now = new Date();
    const createdAt = new Date(timestamp);
    const diffTime = Math.abs(now.getTime() - createdAt.getTime());

    const minutes = Math.floor(diffTime / (1000 * 60));
    const hours = Math.floor(minutes / 60);
    const days = Math.floor(hours / 24);

    if (minutes < 1) {
      return '刚刚';
    } else if (minutes < 60) {
      return `${minutes}分钟前`;
    } else if (hours < 24) {
      return `${hours}小时前`;
    } else {
      return `${days}天前`;
    }
  };

  const clearMessage = () => {
    message.value = '';
  };


  // 获取用户表中所有用户
  const getAllUser = async () => {
    allUserList.value = await db.getAll("workbenchChatUser");
    console.log(allUserList.value)
  };

  // 创建群聊
  const createGroupChat = async (userIds: number[]) => {

    const currUserId = userInfo.value.id;
    // 将当前用户的 ID 添加到 user_ids 数组中
    const newUserIds = [currUserId, ...userIds]
    const data = {
      name: departmentName.value,
      user_ids: newUserIds
    }

    const url = config.userInfo.url + "/chat/group";
    const res = await fetchPost(url, JSON.stringify(data));
    console.log(res)
    if (!res.ok) {
      return false;
    }

    // getGroupList()
    // initChatList()
    const groupData = await res.json();
    console.log(groupData)
    // 构建数据入库
    // 群数据
    const group_id = groupData.data.group_id

    // 添加到会话列表中
    const groupConversation = {
      group_id: group_id,
      avatar: "./logo.png",
      messages: "",
      chatId: group_id,
      type: "group",
      displayName: departmentName.value,
      previewMessage: "",
      previewTimeFormat: formatTime(Date.now()),
      createdAt: new Date()
    }
    groupList.value.push(groupConversation)
    initChatList()
    // 关闭对话弹窗
    setGroupChatInvitedDialogVisible(false)
  };

  // 处理用户消息
  const userChatMessage = async (data: any) => {
    // 先判断数据库是否有该用户
    // 更新聊天记录表
    // 更新会话列表数据库
    // 更新chatlist
    // 更新聊天记录
    console.log("userChatMessage:", data)
    const isPresence = await db.getByField('workbenchChatUser', 'chatId', data.userId)
    if (isPresence[0].id !== data.userId) {
      return
    }

    // 添加消息记录
    const addMessageHistory = {
      type: data.type,
      time: data.time,
      userId: data.userId,
      message: data.message,
      toUserId: data.toUserId,
      chatId: data.toUserId,
      isMe: false,
      content_type: data.content_type,
      to_groupid: data.to_groupid,
      displayName: data.userInfo.nickname,
      avatar: data.userInfo.avatar,
      createdAt: Date.now(),
    }

    await db.addOne("workbenchChatRecord", addMessageHistory)

    console.log("添加成功")
    // 更新 chatList 和 conversationList表
    changeChatListAndChatHistory(data)
    // 更改聊天记录
    chatHistory.value.push(addMessageHistory)
  };


  const setScrollToBottom = async () => {
    // await nextTick(); // 确保 DOM 已经更新完毕

    // // 检查 innerRef 是否存在
    // if (!innerRef.value) {
    //   console.warn('innerRef is not defined.');
    //   return;
    // }

    // // 设置滚动条到底部
    // const max = innerRef.value.clientHeight;
    // if (scrollbarRef.value) {
    //   scrollbarRef.value.setScrollTop(max);
    // } else {
    //   console.warn('scrollbarRef is not defined.');
    // }
  };
  // 获取群列表信息
  const getGroupList = async () => {

    const res = await fetchGet(userInfo.value.url + '/chat/group/list');
    if (!res.ok) {
      console.warn("Error fetching group list:", res);
      return false;
    }
    const list = await res.json()
    console.log(list)
    if (list.data.groups == null) {
      list.data.groups = []
    }
    // 封装 list.data.groups
    const formattedGroups = list.data.groups.map((group: any) => ({
      group_id: group.id,
      avatar: group.avatar || '', // 使用默认头像
      // messages: "",
      displayName: group.name,
      chatId: group.id,
      type: 'group',
    }));
    groupList.value = formattedGroups;

    // 将群成员添加到数据库
    for (const group of list.data.groups) {
      // 查询数据库中是否已存在该群组的成员列表
      const existingGroup = await db.getByField("workbenchGroupUserList", "group_id", group.id);
      if (existingGroup.length > 0) {
        // 如果存在，更新userIdArray
        await db.update("workbenchGroupUserList", existingGroup[0].id, {
          userIdArray: group.members
        });
      } else {
        // 如果不存在，添加新记录
        await db.addOne("workbenchGroupUserList", {
          group_id: group.id,
          userIdArray: group.members
        });
      }
    }

  };

  const onlineUserData = async (data: OnlineUserInfoType[]) => {
    // 创建一个用户数组，将所有在线的用户提取出来
    const onlineUsers: OnlineUserInfoType[] = [];
    // 遍历每个数据项
    data.forEach((item: any) => {
      if (item.id && item.login_ip) {
        onlineUsers.push({
          id: item.id,
          ip: item.login_ip,
          type: "user",
          chatId: item.id,
          online: true,
          avatar: item.avatar,
          username: item.username,
          nickname: item.nickname
        });
      }
    });

    // 将提取到的用户数据传递给
    if (onlineUsers.length > 0) {
      await setOnlineUserList(onlineUsers);
    }
  };

  const setOnlineUserList = async (data: OnlineUserInfoType[]) => {
    if (data.length < 1) {
      return;
    }

    // 从当前用户列表中获取已有用户的 IP 和完整用户映射
    const existingIps = new Set(onlineUserList.value.map((d: any) => d.ip));
    const userMap = new Map<string, OnlineUserInfoType>(
      onlineUserList.value.map((user: OnlineUserInfoType) => [user.ip, user])
    );

    const updates: any[] = [];
    const newEntries: any[] = [];

    // 遍历传入的 data，每个用户根据是否存在来更新或添加
    data.forEach((d: any) => {
      const existingUser = userMap.get(d.ip);
      if (existingUser && existingIps.has(d.ip)) {
        // 若用户已存在，添加到更新列表
        updates.push({
          key: existingUser.id,
          changes: {
            isOnline: true,
            avatar: d.avatar,
            nickname: d.nickname,
            username: d.username,
            updatedAt: Date.now()
          }
        });
      } else {
        // 若用户不存在，添加到新条目列表
        newEntries.push({
          id: d.id,
          ip: d.ip,
          type: "user",
          chatId: d.id,
          isOnline: true,
          avatar: d.avatar,
          nickname: d.nickname,
          username: d.usernmae,
          createdAt: Date.now(),
          updatedAt: Date.now()
        });
      }
    });


    // 批量更新和添加用户数据
    if (updates.length > 0) {
      await db.table("workbenchChatUser").bulkUpdate(updates);
    }
    if (newEntries.length > 0) {
      await db.table("workbenchChatUser").bulkPut(newEntries);
    }
    initOnlineUserList()
  };

  const initOnlineUserList = async () => {

    // 从数据库中获取所有用户信息
    const list = await db.getAll("workbenchChatUser");
    // 创建一个 Map，用于存储每个用户的唯一 ID 地址
    let uniqueIdMap = new Map<string, any>();

    // 遍历用户列表，将每个用户添加到 Map 中（基于 ID 去重）
    list.forEach((item: any) => {
      uniqueIdMap.set(item.userId, item); // 使用 ID 作为键，用户对象作为值
    });

    // 将 Map 的值转换为数组（去重后的用户列表）
    const uniqueIdList = Array.from(uniqueIdMap.values());
    // 按照 updatedAt 时间进行升序排序
    uniqueIdList.sort((a: any, b: any) => a.time - b.time);
    // 更新用户列表
    onlineUserList.value = uniqueIdList;
  };

  // 初始化统一用户列表状态
  const initUserList = async () => {
    if (!onlineUserList.value.length) return;
    // 获取需要更新的用户数据（只选取在线的用户并设为离线状态）
    const updates = onlineUserList.value.reduce((acc: any[], user: any) => {

      if (user.isOnline) {
        console.log(user);
        acc.push({
          key: user.id,
          changes: { isOnline: false }
        });
      }
      return acc;
    }, []);

    // 如果有需要更新的用户，批量更新数据库状态
    if (updates.length) {
      await db.table('workbenchChatUser').bulkUpdate(updates);
    }
  };

  // 设置会话列表
  const getSessionInfo = async (chatId: string, type: string) => {
    // 设置 targetUserId
    // 根据type去判断
    // user会话，查发送和接收发方id
    // group会话，查群id(chatId)查到所有消息记录
    // 清空聊天记录
    chatHistory.value = []
    targetChatId.value = chatId
    if (type === 'user') {
      console.log("user")
      // 获取当前用户和目标用户的聊天记录
      const history = await getHistory(userInfo.value.id, chatId, type)
      console.log(history)
      chatHistory.value = [...history];
      // 设置目标用户的信息
      await setTargetUserInfo(chatId);
    } else if (type === 'group') {
      console.log('group')
      // 获取当前用户和目标用户的聊天记录
      console.log(userInfo.value.id, chatId, type)
      const history = await getHistory(userInfo.value.id, chatId, type)
      chatHistory.value = history;
      // 设置目标用户的信息
      await setTargetGrouprInfo(chatId);
    }
  };

  const getHistory = async (sendUserId: string, toUserId: string, type: string) => {
    var messagesHistory
    if (type === 'user') {
      messagesHistory = await db.filter("workbenchChatRecord", (record: any) => {
        return (
          (record.userId === sendUserId && record.toUserId === toUserId) ||
          (record.toUserId === sendUserId && record.userId === toUserId)
        );
      });
    } else if (type === 'group') {
      console.log('group')
      messagesHistory = await db.getByField("workbenchGroupChatRecord", "chatId", toUserId);
      console.log("messagesHistory", messagesHistory)
    }
    return messagesHistory
  }

  // 设置目标用户的信息
  const setTargetUserInfo = async (id: string) => {
    console.log(id)
    const userInfoArray = await db.getByField("workbenchChatUser", "chatId", id);
    // 封装用户信息
    const userInfo = {
      type: "user",
      avatar: userInfoArray[0].avatar,
      displayName: userInfoArray[0].nickname,
      toUserId: config.userInfo.id,
      chatId: userInfoArray[0].chatId
    }
    targetUserInfo.value = userInfoArray.length > 0 ? userInfo : {};
    targetGroupInfo.value = {}
  };

  // 设置目标群信息
  const setTargetGrouprInfo = async (id: string) => {
    for (const group of groupList.value) {
      console.log(group)
      if (group.group_id === id) {
        targetGroupInfo.value = group;
        targetUserInfo.value = {};
        break;
      }
    }
  };

  const handleContextMenu = async () => {
    contextMenu.value.visible = false;
  };

  const groupChatMessage = async (data: any) => {
    console.log(data)
    // 创建消息记录
    const messageRecord = {
      userId: data.userId,
      groupId: data.to_groupid,
      content_type: data.content_type,
      message: data.message,
      time: data.time,
      type: data.type,
      chatId: data.to_groupid,
      isMe: false,
      displayName: data.userInfo.nickname,
      avatar: data.userInfo.avatar,
      role_id: data.userInfo.role_id,
      createdAt: Date.now(),
    };

    // 判断当前消息是否是自己发送的
    if (messageRecord.userId === userInfo.value.id) {
      return;
    }

    console.log(messageRecord)

    // 将消息记录添加到数据库
    const res = await db.addOne("workbenchGroupChatRecord", messageRecord);
    console.log(res)
    // 更改聊天记录
    chatHistory.value.push(messageRecord)

  };

  const showContextMenu = (event: any, id: number) => {
    contextMenu.value.visible = true;
    contextMenu.value.chatMessageId = id;
    contextMenu.value.x = event.x;
    contextMenu.value.y = event.y;
  };

  // 退出群聊
  const quitGroup = async (group_id: string) => {
    const url = config.userInfo.url + '/chat/group/leave';
    const res = await fetchPost(url, JSON.stringify({ group_id }))
    if (!res.ok) {
      return false;
    }
    console.log(await res.json())
    // 从groupList中删除
    groupList.value = groupList.value.filter((group: any) => group.group_id !== group_id)
    const a = await db.deleteByField("workbenchGroupUserList", "group_id", group_id)
    console.log(a)
    initChatList()
    notifySuccess("退出群聊成功")

    targetGroupInfo.value = {}
    targetChatId.value = ''

    // const group = await db.getByField("workbenchGroupUserList", "group_id", group_id);
    // if (group.length > 0) {
    //   // 删除当前用户
    //   console.log(group[0])
    //   const updatedUserIdArray = group[0].userIdArray.filter((user: any) => user.id !== userInfo.value.id);
    //   await db.update("workbenchGroupUserList", group[0].id, { userIdArray: updatedUserIdArray });
    // }
  }

  return {
    emojiList,
    onlineUserList,
    chatList,
    groupList,
    chatHistory,
    userInfo,
    search,
    showChooseFile,
    currentPage,
    pageSize,
    scrollbarRef,
    innerRef,
    currentNavId,
    targetUserInfo,
    targetChatId,
    message,
    contextMenu,
    activeNames,
    groupChatInvitedDialogVisible,
    groupInfoSettingDrawerVisible,
    departmentList,
    allUserList,
    departmentName,
    targetGroupInfo,
    sendInfo,
    fileSendActive,
    groupMemberList,
    initChat,
    showContextMenu,
    setCurrentNavId,
    sendMessage,
    getSessionInfo,
    handleContextMenu,
    addChatListAndGetChatHistory,
    initUserList,
    setGroupChatInvitedDialogVisible,
    setGroupInfoDrawerVisible,
    createGroupChat,
    onlineUserData,
    groupChatMessage,
    userChatMessage,
    getDepartmentList,
    getAllUser,
    quitGroup,
  };
});