import { ElMessage } from 'element-plus'
export function notifyError(message : string) {
    ElMessage({
        type: 'error',
        message
    })
}
export function notifySuccess(message : string) {
    ElMessage({
        type: 'success',
        message
    })
}
export function notifyInfo(message : string) {
    ElMessage({
        type: 'info',
        message
    })
}