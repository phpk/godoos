import { get, post } from '@/utils/request'

export function joinknowledge(path: string) {
    return get(`user/files/knowledge`, { path })
}
export function askknowledge(data: any) {
    return post(`user/files/ask`, data).then((res) => res.data)
}
