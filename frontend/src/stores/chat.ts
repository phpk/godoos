import emojiList from "@/assets/emoji.json";
import { fetchGet, fetchPost, getSystemConfig } from '@/system/config';
import { notifyError } from "@/util/msg";
import { defineStore } from 'pinia';
import { db } from "./db";
import { useMessageStore } from "./message";

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
  ]);

  // 模拟数据 - 聊天消息列表
  const chatHistory = ref([]);

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
      {
        id: 2,
        label: '撤回',
      }
    ],
    x: 0,
    y: 0
  });

  const initChat = () => {
    if (config.userInfo.avatar == '') {
      config.userInfo.avatar = '/logo.png';
    }
    userInfo.value = config.userInfo;
  };

  const setCurrentNavId = (id: number) => {
    currentNavId.value = id;
  };

  const sendMessage = async () => {
    const chatSendUrl = apiUrl + '/chat/send';
    const messageHistory = {
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
      const lastMessage = messageHistory;

      const targetUserInfo = {
        id: targetUser.id,
        nickname: targetUser.nickname,
        avatar: targetUser.avatar,
      };

      // 计算时间差
      const now = new Date();
      const createdAt = new Date(lastMessage.createdAt);
      const diffTime = Math.abs(now.getTime() - createdAt.getTime());

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

  const initSSE = async () => {
    console.log('initSSE');
    const source = new EventSource(`${apiUrl}/chat/message`);

    console.log(source);
    source.onmessage = function (event) {
      const data = JSON.parse(event.data);
      console.log(data);
      messageStore.handleMessage(data);
    };

    source.onerror = function (event) {
      console.error('EventSource error:', event);
    };
  };

  const setScrollToBottom = async () => {
    await nextTick(); // 确保 DOM 已经更新完毕

    // 检查 innerRef 是否存在
    if (!innerRef.value) {
      console.warn('innerRef is not defined.');
      return;
    }

    // 设置滚动条到底部
    const max = innerRef.value.clientHeight;
    if (scrollbarRef.value) {
      scrollbarRef.value.setScrollTop(max);
    } else {
      console.warn('scrollbarRef is not defined.');
    }
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

  const getOnlineUsers = async () => {
    const res = await fetchGet(apiUrl + '/chat/online?page=1');

    if (!res.ok) {
      notifyError("获取在线用户失败");
      return;
    }

    const data = await res.json();
    console.log(data);

    const onlineUsers = data.data.list.map((item: any) => ({
      id: item.id,
      username: item.username,
      nickname: item.nickname,
      email: item.email,
      phone: item.phone,
      desc: item.desc,
      jobNumber: item.job_number,
      workPlace: item.work_place,
      hiredDate: item.hired_date,
      avatar: item.avatar || '/logo.png',
      isOnline: true,
      ip: item.login_ip,
      updatedAt: item.updated_at,
      createdAt: item.add_time,
    }));

    // 更新或添加用户
    const existingUserIds = await db.getAll('workbenchusers').then(users => users.map(user => user.id));
    const newUserIds = onlineUsers.filter(user => !existingUserIds.includes(user.id)).map(user => user.id);

    // 更新现有用户的状态
    const updatePromises = onlineUsers.filter(user => existingUserIds.includes(user.id)).map(user => {
      return db.update('workbenchusers', user.id, {
        isOnline: true,
        updatedAt: user.updatedAt,
        ip: user.ip,
      });
    });

    // 添加新用户
    const addPromises = onlineUsers.filter((user: { id: any }) => newUserIds.includes(user.id)).map((user: any) => {
      return db.addOne('workbenchusers', user);
    });

    await Promise.all([...updatePromises, ...addPromises]);

    // 更新 userList
    userList.value = await db.getAll('workbenchusers');

    console.log(userList.value);
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
    initChat,
    initSSE,
    setCurrentNavId,
    sendMessage,
    changeChatList,
    handleContextMenu,
    showContextMenu,
    getOnlineUsers,
    updateConversationList
  };
});