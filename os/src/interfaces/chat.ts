export interface OnlineUserInfoType {
  id: string
  login_ip: string
  avatar: string
  online: boolean
  type: string
  chatId: string
  username: string
  nickname: string
}

// 文件消息类型
export interface ChatMessageType {
  type: string
  content_type: string
  time: number
  userId: number
  toUserId: number
  message: string
  to_groupid: string
  userInfo: {}
}

// 文件发送模型
export const Message: ChatMessageType = {
  type: '',
  content_type: '',
  time: 0,
  userId: 0,
  toUserId: 0,
  message: '',
  to_groupid: '',
  userInfo: {},
}