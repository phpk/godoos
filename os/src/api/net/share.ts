import { get, post } from '@/utils/request'
export function selectUserList(page: number, nickname: string = "") {
  return get('/user/sharelist', { page, nickname }).then(res => res.data)
}
export function searchSelectUsers(ids: any) {
  return post('/user/searchuser', { ids }).then(res => res.data)
}
export function shareCreate(data: any) {
  return post('/user/files/share', data)
}

// 获取分享用户列表
export function getShareUserList(path: string) {
  return get('user/files/collaboration/editusers', { path }).then(res => res.data)
}

// 获取编辑历史记录
export function getEditHistory(path: string, page: string, size: string) {
  return get('user/files/collaboration/timeline', { path, page, size }).then(res => res.data)
}

// 还原编辑数据
export function restoreEditData(id: string) {
  return post('user/files/collaboration/recover', { id }).then(res => res)
}