import { ElMessage, ElMessageBox, ElNotification } from 'element-plus'
export function errMsg(message: string) {
  msg('error', message)
}
export function successMsg(message: string) {
  msg('success', message)
}
export function infoMsg(message: string) {
  msg('info', message)
}
export function warningMsg(message: string) {
  msg('warning', message)
}
export function noticeMsg(message: string, title: string = '提示', type: any = 'success') {
  if (document.getElementsByClassName('el-notification__group').length > 0) return
  ElNotification({
    title: title,
    type: type,
    message: message,
    //position: 'top-right',
    offset: 550,
  })
}
type MessageType = 'success' | 'warning' | 'info' | 'error'
export function msg(type: MessageType, message: string, single: boolean = true) {
  if (single && document.getElementsByClassName(`el-message--${type}`).length > 0) {
    return
  }
  ElMessage({
    type,
    message,
  })
}
export function confirmMsg(message: string, title = '提示', callback?: (action: boolean) => void): Promise<any> {
  return ElMessageBox.confirm(message, title, {
    type: 'warning',
    showCancelButton: true,
    roundButton: true,
    confirmButtonText: '确定',
    cancelButtonText: '取消',
  }).then(() => {
    if (callback) {
      callback(true);
    }
    return true;
  }).catch(() => {
    if (callback) {
      callback(false);
    }
    return false;
  })
}
export function promptMsg(message: string, title = '提示', initval: any, callback?: (value: any) => void): Promise<any> {
  return ElMessageBox.prompt(message, title, {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    roundButton: true,
    inputValue: initval
  }).then(({ value }) => {
    if (callback) {
      callback(value);
    }
    return value; // 返回输入的值
  })
}
export function promptPwd(callback?: (value: any) => void): Promise<any> {
  return ElMessageBox.prompt('', '密码', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    roundButton: true,
    inputPlaceholder: '请输入文件密码',
    inputValue: '',
    inputType: 'password'
  }).then(({ value }) => {
    if (callback) {
      callback(value);
    }
    return value; // 返回输入的值
  })
}