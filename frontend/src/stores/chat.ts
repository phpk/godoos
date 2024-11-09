import emojiList from "@/assets/emoji.json";
import { fetchGet, fetchPost, getSystemConfig } from '@/system/config';
import { notifyError, notifySuccess } from "@/util/msg";
import { defineStore } from 'pinia';
import { db } from "./db";

export const useChatStore = defineStore('chatStore', () => {

  interface OnlineUserInfoType {
    id: string;
    login_ip: string;
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
    time: number;
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
    time: 0,
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
    // groupTitle.value = '创建群聊'
  };
  // 设置群信息抽屉状态
  const setGroupInfoDrawerVisible = (visible: boolean) => {
    groupInfoSettingDrawerVisible.value = visible
  }

  const groupMembers = ref([])

  // 群名
  const departmentName = ref('');

  const sendInfo: any = ref()
  // 在线用户列表
  const onlineUserList = ref<OnlineUserInfoType[]>([]);

  // 添加成员
  const addMemberDialogVisible = ref(false)
  // 聊天列表
  const chatList: any = ref([]);

  // 聊天消息记录列表
  const chatHistory: any = ref([]);

  // 群组l列表
  const groupList: any = ref([]);
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
  const searchList = ref([]);
  const groups: any = ref([])
  const searchInput = ref('');
  // 定义消息发送、接受的状态用于控制滚动条的滚动
  const messageSendStatus = ref(false)
  const messageReceiveStatus = ref(false)
  // 群成员列表抽屉是否显示
  const groupMemberDrawerVisible = ref(false)

  // 群成员列表
  const groupMemberList = ref<any[]>([])

  // 所有用户列表
  const allUserList = ref([])

  // 部门列表
  const departmentList = ref([])

  // groupTitle
  // const groupTitle = ref('')

  // 群组信息系统消息
  const groupSystemMessage = ref('')

  // 邀请好友
  const inviteFriendDialogVisible = ref(false)
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
    // 获取所有列表
    getAllList()
  };

  // 获取部门列表
  const getAllList = async () => {
    const res = await fetchGet(userInfo.value.url + "/chat/user/list")
    console.log(res);
    if (res.ok) {
      const data = await res.json();
      console.log(data.data.users)
      groups.value = data.data.groups;
      departmentList.value = data.data.users;

      // 新增代码：提取部门成员并去重，按指定字段保存
      const uniqueUsers = new Set();
      console.log("部门成员", data.data.users)
      data.data.users.forEach((department: { users: any[]; }) => {
        department.users?.forEach(async (user) => {
          if (!uniqueUsers.has(user.user_id)) {
            uniqueUsers.add(user.user_id);
            const userToProcess = {
              id: user.user_id,
              login_ip: user.login_ip,
              type: "user",
              chatId: user.user_id,
              online: user.is_online,
              avatar: user.avatar,
              email: user.email,
              phone: user.phone,
              desc: user.desc,
              jobNumber: user.job_number,
              workPlace: user.work_place,
              hiredDate: user.hired_date,
              username: user.user_name,
              nickname: user.user_name
            };

            const existingUser = await db.getOne("workbenchChatUser", user.user_id);
            if (existingUser) {
              await db.update("workbenchChatUser", user.user_id, userToProcess);
            } else {
              await db.addOne("workbenchChatUser", userToProcess);
            }
          }
        });
      });
    }
  }

  // 初始化用户列表
  const initChatList = async () => {
    console.log("收到消息被刷新了！！！！")
    const userSessionList = await db.getAll("workbenchSessionList");
    // 给userSessionList去一次重
    const uniqueUserSessionList = userSessionList.filter((item: any, index: number, self: any[]) =>
      index === self.findIndex((t: any) => t.chatId === item.chatId)
    );

    // 确保groupList已加载
    if (groupList.value.length > 0) {
      // 合并两个数组
      chatList.value = [...uniqueUserSessionList, ...groupList.value];
    } else {
      // 如果groupList为空，只使用userSessionList
      chatList.value = [...uniqueUserSessionList];
    }
  };

  const setCurrentNavId = (id: number) => {
    currentNavId.value = id;
  };

  const sendMessage = async (messageType: string) => {

    if (messageType == "text" && message.value.trim() == '') {
      return
    }
    messageSendStatus.value = false

    if (messageType == 'applyfile') {
      // 根据文件扩展名调整消息类型
      const fileExtension = sendInfo.value[0]?.split('.').pop().toLowerCase();  // 确保扩展名比较时不区分大小写
      if (['png', 'jpg', 'jpeg', 'gif', 'bmp'].includes(fileExtension)) {
        messageType = 'image';
      } else if (['txt', 'doc', 'pdf', 'xls', 'xlsx', 'ppt', 'pptx'].includes(fileExtension)) {
        messageType = 'applyfile';
      }
    }

    if (messageType == 'text') {
      await sendTextMessage()
    } else if (messageType == 'image') {
      await sendImageMessage()
    } else if (messageType == 'applyfile') {
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

    // 发送文件消息
    const res = await fetchPost(config.userInfo.url + '/chat/send', JSON.stringify(Message));
    if (!res.ok) {
      fileSendActive.value = false;
      return;
    }

    // 封装成消息历史记录
    const messageHistory = {
      ...Message,
      isMe: true,
      message: await getImageSrc(Message.message),
      previewMessage: "[图片消息]",
      previewTimeFormat: formatTime(Date.now()),
      displayName: userInfo.value.nickname,
      chatId: Message.type === 'group' ? targetGroupInfo.value.chatId : targetUserInfo.value.chatId,
      avatar: Message.type === 'group' ? targetUserInfo.value.avatar : '',
      createdAt: Date.now()
    };


    // 添加到聊天记录
    if (Message.type === 'group') {
      await db.addOne("workbenchGroupChatRecord", messageHistory);
      updateGroupSessionList(messageHistory.previewMessage as string)
    } else if (Message.type === 'user') {
      await db.addOne("workbenchChatRecord", messageHistory);
      updateUserSessionList(messageHistory.previewMessage as string)
    }
    chatHistory.value.push(messageHistory);
    fileSendActive.value = true;
    notifySuccess('文件发送成功');

    messageSendStatus.value = true

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
    // 发送文件消息
    const res = await fetchPost(config.userInfo.url + '/chat/send', JSON.stringify(Message));
    console.log(res)
    if (!res.ok) {
      fileSendActive.value = false;
      return;
    }
    const fileInfo = await res.json()

    // 封装成消息历史记录
    const messageHistory = {
      ...Message,
      isMe: true,
      file_path: Message.message,
      file_name: Message.message.split('/').pop(),
      file_size: fileInfo.data.file_info.size,
      previewMessage: "[文件消息]",
      previewTimeFormat: formatTime(Date.now()),
      displayName: userInfo.value.nickname,
      chatId: Message.type === 'group' ? targetGroupInfo.value.chatId : targetUserInfo.value.chatId,
      avatar: userInfo.value.avatar,
      createdAt: Date.now()
    };


    // 添加到聊天记录
    if (Message.type === 'group') {
      await db.addOne("workbenchGroupChatRecord", messageHistory);
      updateGroupSessionList(messageHistory.previewMessage as string)
    } else if (Message.type === 'user') {
      await db.addOne("workbenchChatRecord", messageHistory);
      updateUserSessionList(messageHistory.previewMessage as string)
    }
    chatHistory.value.push(messageHistory);
    fileSendActive.value = true;
    notifySuccess('文件发送成功');

    messageSendStatus.value = true
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
      console.log(await res.json())
      // 封装成消息历史记录
      var messageHistory
      // 本地存储一份聊天记录
      if (Message.type === 'user') {
        messageHistory = {
          ...Message,
          isMe: true,
          previewMessage: message.value,
          previewTimeFormat: formatTime(Date.now()),
          displayName: userInfo.value.nickname,
          chatId: targetUserInfo.value.chatId,
          avatar: userInfo.value.avatar,
          createdAt: Date.now()
        }
        await db.addOne("workbenchChatRecord", messageHistory);
        // 更新sessionlist表的字段
        updateUserSessionList(messageHistory.message)

      } else if (Message.type === 'group') {
        messageHistory = {
          ...Message,
          previewMessage: message.value,
          previewTimeFormat: formatTime(Date.now()),
          displayName: targetGroupInfo.value.name,
          isMe: true,
          chatId: targetGroupInfo.value.chatId,
          createdAt: Date.now()
        }

        await db.addOne("workbenchGroupChatRecord", messageHistory);
        // 更新groupSessionList
        updateGroupSessionList(message.value)

      }

      // 更新聊天历史
      chatHistory.value.push(messageHistory);
      // 清空输入框
      clearMessage();
      messageSendStatus.value = true
    }

  }

  const updateGroupSessionList = async (message: string) => {
    const groupSessionList = await db.getByField("groupSessionList", "chatId", targetChatId.value)

    await db.update("groupSessionList", groupSessionList[0].id, {
      previewMessage: message,
      time: Date.now(),
      previewTimeFormat: Date.now()
    })

    // 更新grouplist的预览字段
    groupList.value.forEach((item: any) => {
      if (item.chatId === targetChatId.value) {
        item.previewMessage = message
        item.previewTimeFormat = formatTime(Date.now())
      }
    })
    initChatList()
  }


  const updateUserSessionList = async (message: string) => {
    // 更新sessionlist表的字段
    const sessionList = await db.getByField("workbenchSessionList", "chatId", targetUserInfo.value.chatId)
    await db.update("workbenchSessionList", sessionList[0].id, {
      previewMessage: message,
      previewTimeFormat: formatTime(Date.now())
    })

    initChatList()
  }



  // 更新聊天和聊天记录
  const changeChatListAndChatHistory = async (data: any, message: string) => {

    // 从 conversationList 数据库中查找是否存在对应的会话
    const conversation = await db.getByField("workbenchSessionList", 'chatId', data.userId);
    // 准备会话更新数据
    const updatedConversation = {
      userId: data.userId,
      avatar: data.userInfo.avatar, // 如果没有头像使用默认图片
      toUserId: data.toUserId,
      chatId: data.userId,
      type: data.type,
      online: data.online,
      messages: data.message,
      displayName: data.userInfo.nickname,
      nickname: data.userInfo.nickname,
      time: Date.now(),
      previewMessage: message,
      previewTimeFormat: formatTime(Date.now()), // 时间格式化函数
      createdAt: Date.now(),
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
        previewMessage: message,
        time: Date.now(),
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
      initChatList()
    }
  };

  // 点击目标用户时，将其添加到 chatList，并添加到数据库,以及获取到聊天记录
  const addChatListAndGetChatHistory = async (chatId: string, type: string) => {
    const chatIdSet = new Set(chatList.value.map((chat: { chatId: any }) => chat.chatId));
    // getSessionInfo(chatId, type);
    // 如果会话存在于 chatList，则获取聊天记录并更新 chatHistory
    messageSendStatus.value = false;

    if (chatIdSet.has(chatId)) {
      console.log("存在");

      if (type == "group") {
        getInviteUserList()
        chatHistory.value = await db.getByField("workbenchGroupChatRecord", "chatId", chatId);
        messageSendStatus.value = true;
        return
      }

      chatHistory.value = await getHistory(chatId, userInfo.value.id, type);
      messageSendStatus.value = true;
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
      online: user.online,
      toUserId: userInfo.value.id,
      previewMessage: message.value || "快开始打招呼吧！",
      previewTimeFormat: formatTime(Date.now()),
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
      online: user.online,
      toUserId: userInfo.value.id,
      previewMessage: message.value || "快开始打招呼吧！",
      previewTimeFormat: formatTime(Date.now()),
      createdAt: Date.now(),
    });

    // 获取聊天记录
    chatHistory.value = await getHistory(chatId, userInfo.value.id, type);

    console.log(chatHistory.value)
    messageSendStatus.value = true;
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
    const users = await db.getAll("workbenchChatUser");
    // 过滤掉当前用户
    allUserList.value = users.filter((user: { id: any; }) => user.id !== userInfo.value.id);
  };

  // 创建群聊
  const createGroupChat = async (userIds: string[]) => {

    if (userIds.length === 0) {
      notifyError('请选择用户')
      return false
    }

    if (departmentName.value === '') {
      notifyError('请输入群名')
      return false
    }

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
    // 构建数据入库
    // 群数据
    const group_id = groupData.data.group_id

    // 添加到会话列表中
    const groupConversation = {
      group_id: group_id,
      avatar: "",
      messages: message.value,
      chatId: group_id,
      type: "group",
      time: Date.now(),
      displayName: departmentName.value,
      previewMessage: "快开始打招呼吧！",
      previewTimeFormat: Date.now(),
      createdAt: Date.now()
    }
    groupList.value.push(groupConversation)

    // 这里的message是群主发送的邀请信息，需要根据userIds去拿到被邀请用户的nickname
    const userList = await db.getByIds("workbenchChatUser", userIds)

    const inviteUserName = userList.map((user: { nickname: string; }) => user.nickname).join(',')

    // 封装邀请信息
    const groupInviteMessage = {
      content_type: "invite_group_message",
      isMe: true,
      chatId: group_id,
      time: Date.now(),
      message: `你邀请 ${inviteUserName} 加入群聊`,
      createdAt: Date.now()
    }
    // 添加到数据库
    await db.addOne("workbenchGroupChatRecord", groupInviteMessage)


    // groupConversation添加到groupSessionList
    await db.addOne("groupSessionList", groupConversation)


    // 更新会话列表
    initChatList()
    // 关闭对话弹窗
    setGroupChatInvitedDialogVisible(false)
    // 更新群组列表
    // await getAllList()
    getGroupList()
    notifySuccess('创建群聊成功')
  };

  // 处理用户消息
  const userChatMessage = async (data: any) => {
    // 先判断数据库是否有该用户
    // 更新聊天记录表
    // 更新会话列表数据库
    // 更新chatlist
    // 更新聊天记录
    console.log(data)

    messageReceiveStatus.value = false

    // 判断是否是自己发的消息
    if (data.userId === userInfo.value.id) {
      return
    }

    const isPresence = await db.getByField('workbenchChatUser', 'chatId', data.userId)
    if (isPresence[0].id !== data.userId) {
      return
    }
    // 添加消息记录
    const addMessageHistory = {
      type: data.type,
      time: data.time,
      userId: data.userId,
      message: await getImageSrc(data.message),
      toUserId: data.toUserId,
      chatId: data.toUserId,
      isMe: false,
      file_path: data.message,
      file_name: data.file_info.origin_name,
      file_size: data.file_info.size,
      content_type: data.content_type,
      to_groupid: data.to_groupid,
      previewMessage: "",
      previewTimeFormat: formatTime(Date.now()),
      displayName: data.userInfo.nickname,
      avatar: data.userInfo.avatar,
      createdAt: Date.now(),
    }

    if (data.content_type === 'image') {
      addMessageHistory.message = await getImageSrc(data.message)
      addMessageHistory.previewMessage = "[图片消息]"
    } else if (data.content_type === 'text') {
      addMessageHistory.message = data.message
      addMessageHistory.previewMessage = data.message
    } else if (data.content_type === 'file') {
      addMessageHistory.previewMessage = "[文件消息]"
      addMessageHistory.message = data.message
    }

    await db.addOne("workbenchChatRecord", addMessageHistory)

    // 更新 chatList 和 conversationList表
    changeChatListAndChatHistory(data, addMessageHistory.previewMessage)




    // 更改聊天记录
    chatHistory.value.push(addMessageHistory)
    messageReceiveStatus.value = true
  };


  // 接收 innerRef 和 scrollbarRef 作为参数
  const setScrollToBottom = async (innerRef: any, scrollbarRef: any) => {
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
  // 获取群列表信息
  const getGroupList = async () => {

    const res = await fetchGet(userInfo.value.url + '/chat/group/list');
    if (!res.ok) {
      console.warn("Error fetching group list:", res);
      return false;
    }
    const list = await res.json()

    if (list.data.groups == null) {
      list.data.groups = []
    }

    // 从groupSessionList中获取群信息
    const groupSessionList = await db.getAll("groupSessionList")

    // 合并查找和封装逻辑到一个循环中
    const formattedGroups = list.data.groups.map((group: any) => {
      const groupSession = groupSessionList.find((item: { chatId: string; }) => item.chatId === group.id);
      return {
        group_id: group.id,
        avatar: group.avatar || '', // 使用默认头像
        displayName: group.name,
        chatId: group.id,
        type: 'group',
        time: groupSession ? groupSession.time : "",
        previewMessage: groupSession ? groupSession.previewMessage : "", // 如果找到对应的会话则使用其预览消息，否则为空
      };
    });
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

    if (data.length === 0) {
      onlineUserList.value = []
      // 更新会话列表用户在线状态
      chatList.value = chatList.value.map((chat: { type: string; }) => {
        if (chat.type === 'user') {
          return { ...chat, online: false };
        }
        return chat;
      });

      return
    }
    // 创建一个新的用户数组，用于更新在线用户列表
    const updatedOnlineUsers = data.map(item => ({
      id: item.id,
      login_ip: item.login_ip,
      type: "user",
      chatId: item.id,
      online: true,
      avatar: item.avatar || '/default_avatar.png', // 使用默认头像如果没有提供
      username: item.username,
      nickname: item.nickname
    })).filter(item => item.id && item.login_ip); // 确保所有项都有有效的id和ip

    // 更新在线用户列表，只添加不存在的用户
    updatedOnlineUsers.forEach(newUser => {
      if (!onlineUserList.value.some(existingUser => existingUser.id === newUser.id)) {
        onlineUserList.value.push(newUser);
      }
    });

    // 新增代码：从onlineUserList中移除不在data中的用户
    onlineUserList.value = onlineUserList.value.filter(user =>
      data.some(newUser => newUser.id === user.id)
    );

    // 更新数据库中的用户信息
    await updateOrAddUsers(updatedOnlineUsers);

    // 更新会话列表用户在线状态
    const updateOnlineStatus = async (data: any[]) => {
      const updatedChats = chatList.value.map((chat: { type: string; chatId: any; }) => {
        if (chat.type === 'user') {
          const isOnline = data.some(user => user.id === chat.chatId);
          return { ...chat, online: isOnline };
        }
        return chat;
      });
      return updatedChats;
    };
    chatList.value = await updateOnlineStatus(updatedOnlineUsers);
  };

  const updateOrAddUsers = async (users: OnlineUserInfoType[]) => {
    // 从数据库中获取所有用户信息
    const allUsers = await db.getAll("workbenchChatUser");
    // 添加或更新在线用户
    for (const user of users) {
      const existingUser = allUsers.find((u: { id: string; }) => u.id === user.id);
      if (existingUser) {
        // 更新现有用户
        await db.table("workbenchChatUser").update(user.id, {
          ...existingUser,
          isOnline: true,
          avatar: user.avatar,
          nickname: user.nickname,
          username: user.username,
          updatedAt: Date.now()
        });
      } else {
        // 添加新用户
        await db.table("workbenchChatUser").add({
          id: user.id,
          login_ip: user.login_ip,
          type: "user",
          chatId: user.id,
          isOnline: true,
          avatar: user.avatar,
          nickname: user.nickname,
          username: user.username,
          createdAt: Date.now(),
          updatedAt: Date.now()
        });
      }
    }
  };

  // 邀请好友
  const inviteFriend = async (groupId: string, userIds: string[]) => {
    if (userIds.length === 0) {
      notifyError('请选择用户')
      return false
    }
    // 邀请加入群聊
    const url = config.userInfo.url + '/chat/group/join';
    const res = await fetchPost(url, JSON.stringify({ group_id: groupId, user_ids: userIds }));
    if (!res.ok) {
      return false;
    }
    console.log(await res.json())
    // 关闭对话框
    inviteFriendDialogVisible.value = false
    notifySuccess('邀请成功')
  }

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
    messageSendStatus.value = false
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
      chatHistory.value = [...history];
      // 设置目标用户的信息
      await setTargetUserInfo(chatId);
      messageSendStatus.value = true
    } else if (type === 'group') {
      console.log('group')
      console.log(userInfo.value.id, chatId, type)
      // 获取当前用户和目标用户的聊天记录
      getInviteUserList()
      const history = await getHistory(userInfo.value.id, chatId, type)
      // await getGroupInviteMessage(chatId)
      chatHistory.value = history;
      // 设置目标用户的信息
      await setTargetGrouprInfo(chatId);
      messageSendStatus.value = true
    }
  };

  const getHistory = async (sendUserId: string, toUserId: string, type: string) => {
    var messagesHistory
    console.log(sendUserId, toUserId, type)
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
    var userInfoArray = await db.getByField("workbenchChatUser", "chatId", id);

    // 封装用户信息
    const userInfo = {
      type: "user",
      avatar: userInfoArray[0].avatar || "",
      displayName: userInfoArray[0].nickname || "",
      jobNumber: userInfoArray[0].jobNumber || "",
      desc: userInfoArray[0].desc || "",
      email: userInfoArray[0].email || "",
      phone: userInfoArray[0].phone || "",
      hiredDate: userInfoArray[0].hiredDate || "",
      toUserId: config.userInfo.id,
      chatId: userInfoArray[0].chatId
    }
    targetUserInfo.value = userInfoArray.length > 0 ? userInfo : {};
    targetGroupInfo.value = {}
  };

  // 设置目标群信息
  const setTargetGrouprInfo = async (id: string) => {
    for (var group of groupList.value) {
      if (group.group_id === id) {
        // 模拟群信息
        // group = {
        //   chatId: id,
        //   displayName: "湖南果度科技有限公司",
        //   avatar: "./src/assets/icons/group.png",
        //   memberCount: 150, // 假设成员数
        //   createdAt: "2023-01-01" // 假设创建日期
        // }
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

    messageReceiveStatus.value = false


    // 构建消息记录
    const messageRecord: any = {
      userId: data.userId,
      groupId: data.to_groupid,
      content_type: data.content_type,
      time: data.time,
      type: data.type,
      chatId: data.to_groupid,
      isMe: false,
      file_path: data.message,
      file_name: data.file_info.origin_name,
      file_size: data.file_info.size,
      previewTimeFormat: formatTime(Date.now()),
      displayName: data.userInfo.nickname, // 发送者昵称
      avatar: data.userInfo.avatar,
      role_id: data.userInfo.role_id,
      createdAt: Date.now(),
    };

    if (data.content_type === 'image') {
      messageRecord.message = await getImageSrc(data.message)
      messageRecord.previewMessage = "[图片消息]"
      updateRecipientGroupSessionList(messageRecord.previewMessage, data.to_groupid)
    } else if (data.content_type === 'text') {
      messageRecord.message = data.message
      messageRecord.previewMessage = data.message
      updateRecipientGroupSessionList(messageRecord.previewMessage, data.to_groupid)
    } else if (data.content_type === 'file') {
      messageRecord.message = data.message
      messageRecord.previewMessage = "[文件消息]"
      updateRecipientGroupSessionList(messageRecord.previewMessage, data.to_groupid)
    }

    // 判断接受的消息是否是自己发送的
    if (messageRecord.userId === userInfo.value.id) {
      return;
    }


    // 将消息记录添加到数据库
    await db.addOne("workbenchGroupChatRecord", messageRecord);
    // // push进groupList
    // groupList.value.push(messageRecord)

    // 更改聊天记录
    chatHistory.value.push(messageRecord)
    messageReceiveStatus.value = true
  };

  // 接受方更新groupSessionList
  const updateRecipientGroupSessionList = async (message: string, chatId: string) => {
    const group = await db.getByField("groupSessionList", "chatId", chatId)
    // 更新groupSessionList
    await db.update("groupSessionList", group[0].id, {
      previewMessage: message,
      time: Date.now(),
      previewTimeFormat: Date.now(),
    })
    await getGroupList()
    initChatList()
  }

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
    await db.deleteByField("workbenchGroupUserList", "group_id", group_id)
    // 获取群列表
    await getGroupList()
    // 初始化聊天列表
    initChatList()
    targetGroupInfo.value = {}
    targetChatId.value = ''
    notifySuccess("退出群聊成功")
  }

  // 获取群成员
  const getGroupMember = async (group_id: string) => {
    const member = await db.getByField("workbenchGroupUserList", "group_id", group_id)
    groupMemberList.value = member[0].userIdArray
  }

  const inviteUserList = ref([])// 添加群成员


  // 获取图片消息预览
  // 获取图片消息预览
  const getImageSrc = async (imageMessage: string) => {
    // 确保路径以 '/' 开始
    if (!imageMessage.startsWith('/')) {
      imageMessage = '/' + imageMessage;
    }

    const path = config.userInfo.url + "/chat/image/view?path=" + imageMessage;
    const response = await fetchGet(path);
    if (!response.ok) {
      return '';
    }

    // 检查返回的内容类型
    const contentType = response.headers.get("Content-Type");
    console.log("Content-Type:", contentType);


    const blob = await response.blob(); // 获取 Blob 对象
    return new Promise((resolve) => {
      const reader = new FileReader();
      reader.onloadend = () => {
        const base64data = reader.result; // 转换为 base64
        resolve(base64data);
      };
      reader.readAsDataURL(blob);
    });
  }

  // 设置群组信息(todo:如果群的系统消息更新，则更新群组信息)
  const groupInviteMessage = async (data: any) => {

    const groupInviteMessage = {
      group_id: data.group_id,
      message: data.message,
      chatId: data.group_id,
      isMe: false,
      previewMessage: "快开始打招呼吧！",
      content_type: "invite_group_message",
      createdAt: Date.now()
    };
    await db.addOne("workbenchGroupChatRecord", groupInviteMessage)
    // 添加到groupSessionList
    await db.addOne("groupSessionList", groupInviteMessage)

    const groupExists = groupList.value.some((group: { group_id: string; }) => group.group_id === data.group_id);

    // 如果不存在，则获取群聊列表
    if (!groupExists) {
      await getGroupList();
      // 合并groupList到chatList，同时去重
      const existingGroupIds = new Set(chatList.value.map((item: { chatId: string }) => item.chatId));
      const newGroups = groupList.value.filter((group: { chatId: string }) => !existingGroupIds.has(group.chatId));
      chatList.value = [...chatList.value, ...newGroups];
    }
  }

  // 获取群成员
  const getGroupMemberList = async (group_id: string) => {
    console.log(group_id)

    const res = await fetchGet(userInfo.value.url + '/chat/group/info?gid=' + group_id);
    if (!res.ok) {
      return false;
    }

    const memberData = await res.json();
    console.log(memberData.data.group.members)
    // 封装 avatar 和 nickname 到 groupMembers
    groupMembers.value = memberData.data.group.members.map((member: any) => ({
      id: member.id,
      avatar: member.avatar,
      nickname: member.nickname,
    }));



  }

  // 邀请群聊获取用户
  const getInviteUserList = async () => {

    const userList = await db.getAll("workbenchChatUser");

    await getGroupMemberList(targetChatId.value);

    // 筛选出不在 groupMembers 中的用户
    inviteUserList.value = userList.filter(
      (user: { id: string; }) => !groupMembers.value.some((member: { id: string; }) => member.id === user.id)
    );
  };


  // 清空发送的聊天记录
  const clearSentMessages = async () => {
    const whereObjSent = {
      toUserId: targetChatId.value,
      userId: userInfo.value.id
    };
    const resSent = await db.deleteByWhere("workbenchChatRecord", whereObjSent);
    if (resSent >= 0) {
      // 更新chatList中的预览消息
      chatList.value.forEach((item: any) => {
        if (item.chatId === targetChatId.value) {
          item.previewMessage = "快开始打招呼吧！";
        }
      });
      return true
    } else {
      return false
    }
  };

  // 清空接收的聊天记录
  const clearReceivedMessages = async () => {
    const whereObjReceived = {
      userId: targetChatId.value,
      toUserId: userInfo.value.id
    };
    const resReceived = await db.deleteByWhere("workbenchChatRecord", whereObjReceived);
    if (resReceived >= 0) {

      // 更新chatList中的预览消息
      chatList.value.forEach((item: any) => {
        if (item.chatId === targetChatId.value) {
          item.previewMessage = "快开始打招呼吧！";
        }
      });
      return true
    } else {
      return false
    }
  };

  // 清空群消息
  const clearGroupMessages = async () => {
    const res = await db.deleteByField("workbenchGroupChatRecord", "chatId", targetChatId.value);
    if (res >= 0) {
      // 还需要找到
      getSessionInfo(targetChatId.value, "group")
      // 更新chatList中的预览消息
      chatList.value.forEach((item: any) => {
        if (item.chatId === targetChatId.value) {
          item.previewMessage = "快开始打招呼吧！";
        }
      });
      notifySuccess("删除成功");
    } else {
      notifyError("删除失败");
    }
  }


  return {
    emojiList,
    groupSystemMessage,
    onlineUserList,
    chatList,
    groupList,
    chatHistory,
    userInfo,
    searchList,
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
    searchInput,
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
    groups,
    inviteFriendDialogVisible,
    addMemberDialogVisible,
    inviteUserList,
    messageSendStatus,
    messageReceiveStatus,
    groupMembers,
    groupMemberDrawerVisible,
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
    getAllList,
    getAllUser,
    quitGroup,
    inviteFriend,
    getGroupMember,
    getImageSrc,
    groupInviteMessage,
    setScrollToBottom,
    getGroupMemberList,
    getInviteUserList,
    clearReceivedMessages,
    clearSentMessages,
    clearGroupMessages
  };
});