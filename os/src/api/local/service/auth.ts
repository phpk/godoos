import { userDb } from "../models/users"
import { success, error } from "./msg"
import { md5 } from "js-md5"
import { toRaw } from "vue"
import * as cache from "../cache"
import { sign,generate16BitString } from './jwt'
import { setToken, setUsername } from '@/utils/request'
import { initOsSystem } from "./os"
export async function loginIn(params: any) {
    //console.log(toRaw(params.param))
    const param = toRaw(params.param)
    const user = await userDb.where({ username : param.username }).find()
    if (!user) return error("用户不存在,请注册用户！")
    if(md5(param.password) !== user.password) return error("密码错误")
    const salt = generate16BitString()
    const token = sign({id:user.id}, salt, { expiresIn: '24h' })
    //console.log(token)
    setToken(token)
    setUsername(user.username)
    cache.set("userid:"+user.id, salt)
    await initOsSystem()
    return success("登录成功", {user})
}
export async function logout() {
    
}
export async function register(params: any) {
    //console.log(params)
    const user = await userDb.where({ username : params.username }).find()
    if (user) return error("用户已存在")
    params.password = md5(params.password)
    await userDb.create(params)
    return success("注册成功", params)
}
export async function isLogin():Promise<boolean> {
    return false
}
export async function getThirdpartyList() {
    return []
}