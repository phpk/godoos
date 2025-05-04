import { get, post } from '@/utils/request'
export function selectUserList(page: number, nickname: string = "") {
  return get('/net/sharelist', { page, nickname }).then(res => res.data)
}
export function searchSelectUsers(ids: any) {
  return post('/net/searchuser', { ids }).then(res => res.data)
}
export function shareCreate(data: any) {
  return post('/net/files/share', data)
}

// 获取分享用户列表
export function getShareUserList(path: string) {
  return get('net/files/collaboration/editusers', { path }).then(res => res.data)
}

// 获取编辑历史记录
export function getEditHistory(path: string, page: string, size: string) {
  return get('net/files/collaboration/timeline', { path, page, size }).then(res => res.data)
}

// 还原编辑数据
export function restoreEditData(id: string) {
  return post('net/files/collaboration/recover', { id }).then(res => res)
}