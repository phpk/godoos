import emojiList from "@/assets/emoji.json";
import { fetchPost, getSystemConfig } from '@/system/config';
import { notifyError } from "@/util/msg";
import { defineStore } from 'pinia';
import { db } from "./db";
import { useMessageStore } from "./message";

interface ChatMessage {
  type: string;
  createdAt: number;
  content: any;
  targetUserId: any;
  previewType: 0 | 1; // 消息类型，0表示正常消息，1表示撤回消息
  previewMessage: any;
  isMe: boolean;
  isRead: boolean;
  userInfo: {
    id: any;
    username: any;
    avatar: any;
  };
}
// 发起群聊对话框显示
const groupChatDialogVisible = ref(false);

// 设置发起群聊对话框状态
const setGroupChatDialogVisible = (visible: boolean) => {
  groupChatDialogVisible.value = visible;
};

export const useChatStore = defineStore('chatStore', () => {
  // 用户列表
  const userList: any = ref([]);

  // 定义聊天列表项的接口
  interface ChatListItem {
    id: number;
    nickname: string;
    avatar: string;
    previewTimeFormat: string;
    previewType: 0 | 1; // 消息类型，0表示正常消息，1表示撤回消息
    previewMessage: string;
  }

  // 模拟数据 - 聊天列表
  const chatList = ref<ChatListItem[]>([
    {
      id: 2,
      nickname: '朋友2',
      avatar: '/logo.png',
      previewTimeFormat: "昨天",
      previewType: 1,
      previewMessage: "测试消息",
    },
    {
      id: 3,
      nickname: '朋友2',
      avatar: '/logo.png',
      previewTimeFormat: "昨天",
      previewType: 1,
      previewMessage: "测试消息",
    },

  ]);

  // 模拟数据 - 聊天消息列表
  const chatHistory = ref<ChatMessage[]>([]);

  // 群组数据
  const groupList = ref([
    {
      id: 1,
      name: '群组1',
      avatar: '/logo.png',
      previewTimeFormat: "今天",
      previewType: 0,
      previewMessage: "这是一个示例消息。",
    },
    {
      id: 2,
      name: '群组2',
      avatar: '/logo.png',
      previewTimeFormat: "今天",
      previewType: 0,
      previewMessage: "这是一个示例消息。",
    },
    {
      id: 3,
      name: '群组3',
      avatar: '/logo.png',
      previewTimeFormat: "今天",
      previewType: 0,
      previewMessage: "这是一个示例消息。",
    }
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
  const messageStore = useMessageStore();
  const apiUrl = "http://192.168.1.10:8816";

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
    initUserList()
    initOnlineUserList()
    console.log(userList.value);
  };

  const setCurrentNavId = (id: number) => {
    currentNavId.value = id;
  };

  const sendMessage = async () => {
    const chatSendUrl = apiUrl + '/chat/send';
    const messageHistory: ChatMessage = {
      type: 'text',
      createdAt: Date.now(),
      content: message.value,
      targetUserId: targetUserId.value,
      previewType: 0, // 消息类型，0表示正常消息，1表示撤回消息
      previewMessage: message.value,
      isMe: true,
      isRead: false,
      userInfo: {
        id: config.userInfo.id,
        username: config.userInfo.username,
        avatar: config.userInfo.avatar,
      },
    };

    const res = await fetchPost(chatSendUrl, messageHistory);

    if (res.ok) {
      // 本地存储一份聊天记录
      await db.addOne('chatRecord', messageHistory);

      // 更新聊天历史
      chatHistory.value.push(messageHistory);

      // 更新 chatList 和 conversationList
      await updateConversationList(targetUserId.value);

      // 清空输入框
      clearMessage();

      // 聊天框滚动到底部
      await setScrollToBottom();
      return;
    }
    notifyError("消息发送失败");
  };

  const updateConversationList = async (id: number) => {
    // 先判断是否已经存在该会话
    const res = await db.getRow('conversationList', 'id', id);

    if (res) {
      // 更新现有会话
      const updatedConversation = {
        ...res,
        previewMessage: message.value,
        previewTimeFormat: formatTime(Date.now()),
        previewType: 0,
      };
      await db.update('conversationList', id, updatedConversation);

      // 更新 chatList
      const existingConversationIndex = chatList.value.findIndex(conversation => conversation.id === id);
      if (existingConversationIndex !== -1) {
        chatList.value[existingConversationIndex] = updatedConversation;
      } else {
        chatList.value.push(updatedConversation);
      }
    } else {
      const targetUser = await db.getOne('workbenchusers', id);
      const lastMessage: any = messageStore;

      const targetUserInfo = {
        id: targetUser.id,
        nickname: targetUser.nickname,
        avatar: targetUser.avatar,
      };

      // 计算时间差
      const now = new Date();
      const createdAt = new Date(lastMessage.createdAt);
      const diffTime = Math.abs(now.getTime() - createdAt.getTime());
      console.log(diffTime);
      // 根据时间差格式化时间
      const previewTimeFormat = formatTime(Date.now());

      const newConversation = {
        ...targetUserInfo,
        previewTimeFormat,
        previewMessage: lastMessage.content,
        previewType: lastMessage.type,
      };

      // 添加到 conversationList
      await db.addOne('conversationList', newConversation);

      // 添加到 chatList
      chatList.value.push(newConversation);
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
    ;

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

    console.log(users);

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
    const userMap = new Map(
      userList.value.map((user: any) => [user.ip, user])
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
            nickname: d.nickname,
            username: d.usernmae,
            updatedAt: Date.now()
          }
        });
      } else {
        // 若用户不存在，添加到新条目列表
        newEntries.push({
          id: d.id,
          ip: d.ip,
          isOnline: true,
          nickname: d.nickname,
          username: d.usernmae,
          createdAt: Date.now(),
          updatedAt: Date.now()
        });
      }
    });

    console.log(updates);
    console.log(newEntries);

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
        uniqueIdMap.set(item.id, item); // 使用 ID 作为键，用户对象作为值
      });

      // 将 Map 的值转换为数组（去重后的用户列表）
      const uniqueIdList = Array.from(uniqueIdMap.values());

      // 按照 updatedAt 时间进行升序排序
      uniqueIdList.sort((a: any, b: any) => a.updatedAt - b.updatedAt);

      // 更新用户列表
      userList.value = uniqueIdList;
    } catch (error) {
      console.error("获取用户列表失败:", error);
    }
  };

  // 初始化统一用户列表状态
  const initUserList = async () => {
    // 检查用户列表是否为空
    if (userList.value.length > 0) {
      // 收集需要更新的用户数据
      const updates = userList.value
        .filter((d: any) => d.isOnline) // 过滤出在线的用户
        .map((d: any) => ({
          key: d.id,
          changes: {
            isOnline: false
          }
        }));

      // 批量更新用户状态
      if (updates.length > 0) {
        await db.table('workbenchusers').bulkUpdate(updates);
      }
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




  const changeChatList = async (id: number) => {
    // 设置 targetUserId
    targetUserId.value = id;

    // 获取当前用户和目标用户的聊天记录
    const messagesList = await db.getByField('chatRecord', 'targetUserId', id);
    chatHistory.value = messagesList;

    // 设置目标用户的信息
    setTargetUserInfo(id);
  };

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
    groupChatDialogVisible,
    initChat,
    showContextMenu,
    setCurrentNavId,
    sendMessage,
    changeChatList,
    handleContextMenu,
    updateConversationList,
    handleUserData,
    initUserList,
    setGroupChatDialogVisible
  };
});