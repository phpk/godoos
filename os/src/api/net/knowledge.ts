import { get, post } from '@/utils/request'

export function joinknowledge(path: string) {
    return get(`net/files/knowledge`, { path })
}
export function askknowledge(data: any) {
    return post(`net/files/ask`, data).then((res) => res.data)
}
