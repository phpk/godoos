import emojiList from "@/assets/emoji.json";
import { fetchGet, fetchPost, getSystemConfig } from '@/system/config';
import { defineStore } from 'pinia';
import { db } from "./db";




export const useChatStore = defineStore('chatStore', () => {

  interface ChatMessage {
    id?: any;
    type: any;   // 消息类型,0表示文字消息，1表示图片消息，2表示文件消息
    time: any; // 消息发送时间
    message: any; // 消息内容
    userId: any; // 发送者id
    toUserId: any; // 接收者id
    // receiver: any; // 消息接收者
    // to_groupid: any; // 群组id
    userInfo: { // 发送者信息
    }
  };

  // 发起群聊对话框显示
  const groupChatInvitedDialogVisible = ref(false);

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

  // 定义用户类型
  type User = {
    id: number;
    ip: string;
    isOnline: boolean;
    avatar: string;
    nickname: string;
    username: string;
    updatedAt?: number;
  };

  // 将 userList 的类型设置为 User[]
  const userList = ref<User[]>([]);

  // 聊天列表
  const chatList: any = ref([]);

  // 聊天消息记录列表
  const chatHistory: any = ref([]);

  const targetNickname: any = ref('');


  // 群组数据
  const groupList: any = ref([

  ]);

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
  const targetUserId = ref();
  const search = ref('');

  // 目标群组id
  const targetGroupId = ref(0);
  
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

  


  const initChat = () => {
    if (config.userInfo.avatar == '') {
      config.userInfo.avatar = '/logo.png';
    }
    userInfo.value = config.userInfo;
    getUserList()
    getDepartmentList()
    initUserList()
    initChatList()
    // initOnlineUserList()
    console.log(userList.value);
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
    const userchatList = await db.getAll('conversationList');
    // 获取群数据
    const groupChatListawait = await db.getAll("groupChatList")

    // 合并两个数组
    chatList.value = [...userchatList, ...groupChatListawait];
  };

  const setCurrentNavId = (id: number) => {
    currentNavId.value = id;
  };

  const sendMessage = async () => {
    const chatSendUrl = config.userInfo.url + '/chat/send';
    // 封装成消息历史记录
    console.log(chatSendUrl);
    const messageHistory: ChatMessage = {
      type: 'user',
      time: null,
      message: message.value,
      userId: userInfo.value.id,
      toUserId: targetUserId.value,
      userInfo: {
      },
    };

    console.log(messageHistory);

    // 创建没有 `id` 属性的副本
    const { id, ...messageHistoryWithoutId } = messageHistory;
    console.log(messageHistoryWithoutId);
    // 发送消息
    const res = await fetchPost(chatSendUrl, JSON.stringify(messageHistoryWithoutId));
    if (res.ok) {
      // 本地存储一份聊天记录
      await db.addOne('chatRecord', messageHistory);

      // 更新聊天历史
      chatHistory.value.push(messageHistory);

      // 更新 chatList 和 conversationList
      // await changeChatListAndGetChatHistory(userInfo.value.userId);

      // 清空输入框
      clearMessage();

      // 聊天框滚动到底部
      await setScrollToBottom();
      return;
    }
  };

  // 更新聊天和聊天记录
  const changeChatListAndChatHistory = async (data: any) => {
    try {
      // 从 conversationList 数据库中查找是否存在对应的会话
      const conversation = await db.getByField('conversationList', 'userId', data.userId);

      // 准备会话更新数据
      const updatedConversation = {
        userId: data.userId,
        avatar: data.userInfo.avatar || "logo.png", // 如果没有头像使用默认图片
        toUserId: data.toUserId,
        messages: data.message,
        nickname: data.userInfo.nickname,
        time: data.time || Date.now(),
        previewMessage: data.message,
        previewTimeFormat: formatTime(Date.now()), // 时间格式化函数
        createdAt: Date.now()
      };

      if (conversation.length === 0) {
        // 如果会话不存在，则添加到数据库和 chatList
        await db.addOne('conversationList', updatedConversation);

        chatList.value.push(updatedConversation);
      } else {
        // 如果会话存在，则更新数据库和 chatList
        // 只更新变化的字段，而非全部覆盖,以减少写入数据的量
        await db.update('conversationList', conversation[0].id, {
          avatar: data.userInfo.avatar || "logo.png",
          nickname: data.userInfo.nickname,
          previewMessage: data.message,
          time: data.time || Date.now(),
          previewTimeFormat: formatTime(Date.now())
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
    } catch (error) {
      console.error("更新聊天和聊天记录时出错:", error);
    }
  };

  const changeChatListAndGetChatHistory = async (userId: number) => {

    // todo 优化,需要加唯一标识
    console.log(chatList.value);
    let userExists = false;
    // 检查 chatList 中是否存在 userId
    for (let i = 0; i < chatList.value.length; i++) {
      console.log(chatList);
      if (chatList.value[i].id === userId) {
        userExists = true;
        break;
      }
    }
    // 检查 userId 是否已经存
    if (userExists) {
      // 如果存在，获取聊天记录添加到 historyList
      console.log("存在");
      chatHistory.value = await getHistory(userId, userInfo.value.id);
      return;
    } else {
      // 如果不存在，则从用户表中获取该用户的基本信息
      const user = await db.getOne("workbenchusers", userId);
      if (user) {
        // 将新用户信息添加到 chatList
        chatList.value.push({
          id: user.id,
          type: "user",
          nickname: user.nickname,
          avatar: user.avatar,
          previewTimeFormat: formatTime(Date.now()),
          previewMessage: "",
        });

        // 持久化
        await db.addOne("conversationList", {
          userId: user.id,
          type: "user",
          username: user.username,
          nickname: user.nickname,
          avatar: user.avatar,
          toUserId: userInfo.value.id,
          time: Date.now(),
          previewMessage: "",
          createdAt: Date.now(),
        });
      } else {
        console.warn("User not found in workbenchusers with userId:", userId);
      }
    }
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
    allUserList.value = await db.getAll("workbenchusers");
    console.log(allUserList.value)
  };

  // 创建群聊
  const createGroupChat = async (userIds: number[]) => {
    try {

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
      const groupData = await res.json();
      console.log(groupData)
      // 构建数据入库
      // 群数据
      const group_id = groupData.data.group_id
      const gourpData = {
        name: departmentName.value,
        avatar: "./logo.png",
        groupId: group_id,
        creator: currUserId,
        createdAt: new Date()
      }

      // 群成员数据
      const groupMembers = {
        userId: currUserId,
        groupId: group_id,
        createdAt: new Date()
      }
      // 添加数据库
      db.addOne("group", gourpData)
      db.addOne("groupMembers", groupMembers)

      // 添加到会话列表中
      const groupConversation = {
        group_id: group_id,
        name: departmentName.value,
        avatar: "./logo.png",
        messages: "",
        type: "group",
        previewMessage: "",
        previewTimeFormat: formatTime(Date.now()),
        createdAt: new Date()
      }
      // 
      db.addOne("groupChatList", groupConversation)

      chatList.value.push(groupConversation)

      // 关闭对话弹窗
      setGroupChatInvitedDialogVisible(false)

    } catch (error) {
      console.log(error);
    }



  };

  // 处理用户消息
  const userChatMessage = async (data: any) => {
    // 先判断数据库是否有该用户
    // 更新聊天记录表
    // 更新会话列表数据库
    // 更新chatlist
    console.log(data)
    const isPresence = await db.getByField('workbenchusers', 'id', data.userId)
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
      createdAt: Date.now(),
      // 用户信息
      userInfo: {
        id: data.userId,
        nickname: data.userInfo.nickname || "未知用户",
        avatar: data.userInfo.avatar || "logo.png", // 使用默认头像。
        email: data.userInfo.email,
        phone: data.userInfo.phone,
        remark: data.userInfo.remark,
        role_id: data.userInfo.role_id,
      }
    }

    await db.addOne('chatRecord', addMessageHistory)


    // 更新 chatList 和 conversationList表
    changeChatListAndChatHistory(data)
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

  const handleUserData = async (data: any[]) => {
    // 创建一个用户数组，将所有在线的用户提取出来
    const users: any[] = [];
    // 遍历每个数据项
    data.forEach((item: any) => {
      if (item.id && item.login_ip) {
        users.push({
          id: item.id,
          ip: item.login_ip,
          avatar: item.avatar,
          username: item.username,
          nickname: item.nickname
        });
      }
    });


    // 将提取到的用户数据传递给 setUserList
    if (users.length > 0) {
      await setUserList(users);
    }
  };



  const setUserList = async (data: any[]) => {
    if (data.length < 1) {
      return;
    }

    // 从当前用户列表中获取已有用户的 IP 和完整用户映射
    const existingIps = new Set(userList.value.map((d: any) => d.ip));
    const userMap = new Map<string, User>(
      userList.value.map((user: User) => [user.ip, user])
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
      await db.table('workbenchusers').bulkUpdate(updates);
    }
    if (newEntries.length > 0) {
      await db.table('workbenchusers').bulkPut(newEntries);
    }

    // 刷新用户列表
    await getUserList();

  };

  const getUserList = async () => {
    try {
      // 从数据库中获取所有用户信息
      const list = await db.getAll("workbenchusers");
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
      userList.value = uniqueIdList;
    } catch (error) {
      console.error("获取用户列表失败:", error);
    }
  };

  // 初始化统一用户列表状态
  const initUserList = async () => {
    if (!userList.value.length) return;
    // 获取需要更新的用户数据（只选取在线的用户并设为离线状态）
    const updates = userList.value.reduce((acc: any[], user: any) => {
      console.log(user);
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
      await db.table('workbenchusers').bulkUpdate(updates);
    }
  };


  const initOnlineUserList = async () => {
    const msgAll = await db.getAll('workbenchusers');

    const list = msgAll.reduce((acc: any, msg: any) => {
      if (!msg.isMe) {
        const targetId = msg.targetId;
        if (!acc[targetId]) {
          acc[targetId] = { chatArr: [], readNum: 0 };
        }
        acc[targetId].chatArr.push(msg);
        // 计算未读消息数量
        if (!msg.isRead) {
          acc[targetId].readNum++;
        }
      }
      return acc;
    }, {});

    const res = Object.keys(list).map(targetId => {
      const { chatArr, readNum } = list[targetId];
      const lastMessage = chatArr[chatArr.length - 1];
      if (lastMessage) {
        lastMessage.readNum = readNum;
        return lastMessage;
      }
      return null; // 防止返回空值
    }).filter(Boolean); // 过滤掉空值

    userList.value = res.sort((a, b) => b.createdAt - a.createdAt);
  };


  const changeChatList = async (id) => {
    // 设置 targetUserId
    console.log(id)
    targetUserId.value = id
    

    // 获取当前用户和目标用户的聊天记录

    const history = await getHistory(id, userInfo.value.id)
    chatHistory.value = history;
    // 设置目标用户的信息
    await setTargetUserInfo(id);      
  };

  const getHistory = async (userId: any, toUserId: any) => {
    const messagesList = await db.filter('chatRecord', (record: any) => {
      return (
        (record.userId === userId && record.toUserId === toUserId) ||
        (record.toUserId === userId && record.userId === toUserId)
      );
    });

    return messagesList
  }


  const setTargetUserInfo = async (id: number) => {
    targetUserInfo.value = await db.getOne('workbenchusers', id);
  };

  const handleContextMenu = async () => {
    contextMenu.value.visible = false;
  };


  const showContextMenu = (event: any, id: number) => {
    contextMenu.value.visible = true;
    contextMenu.value.chatMessageId = id;
    contextMenu.value.x = event.x;
    contextMenu.value.y = event.y;
  };

  return {
    emojiList,
    userList,
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
    targetUserId,
    message,
    contextMenu,
    activeNames,
    groupChatInvitedDialogVisible,
    groupInfoSettingDrawerVisible,
    targetNickname,
    departmentList,
    allUserList,
    departmentName,
    initChat,
    showContextMenu,
    setCurrentNavId,
    sendMessage,
    changeChatList,
    handleContextMenu,
    changeChatListAndGetChatHistory,
    handleUserData,
    initUserList,
    setGroupChatInvitedDialogVisible,
    setGroupInfoDrawerVisible,
    createGroupChat,
    userChatMessage,
    initOnlineUserList,
    getDepartmentList,
    getAllUser
  };
});