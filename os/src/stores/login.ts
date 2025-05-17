import * as auth from '@/api/auth'
import router from '@/router'
import { useDesktopStore } from '@/stores/desktop'
import { loadScript } from '@/utils/load'
import { errMsg, successMsg } from '@/utils/msg'
import { getClientId } from '@/utils/uuid'
import { clearToken } from '@/utils/request'
import { defineStore } from 'pinia'
import { ref, Ref } from 'vue'

export const useLoginStore = defineStore('login', () => {
  // 第三方登录方式
  const thirdPartyPlatform: Ref<string | null> = ref(null)
  const isLoginState: Ref<boolean> = ref(false)
  const isLogining = ref(false)
  // 第三方登录方法
  const thirdPartyLoginMethod: Ref<string> = ref('password')

  // 第三方列表
  const thirdpartyList: Ref<any[]> = ref([])
  const desktopStore = useDesktopStore()
  const loginForm = ref({
    username: "",
    password: "",
    rememberMe: false
  });
  // 用户信息
  const userInfo: Ref<any> = ref({})

  const checkLogin = async () => {
    isLoginState.value = await auth.isLogin()
  }
  const loginOut = async () => {
    isLoginState.value = false
    userInfo.value = {}
    clearToken()
    
    await auth.logout()
  }

  const backToLogin = () => {
    thirdPartyLoginMethod.value = 'password'
  }
  const onRegister = async (params: any) => {
    //delete params.param.confirmPassword
    params.client_id = getClientId()
    const res = await auth.register(params)
    if (res.success) {
      successMsg('注册成功')
      thirdPartyLoginMethod.value = 'password'
    } else {
      errMsg(res.message)
    }
  }
  const onLogin = async (params: { loginType: string; params: any }) => {
    const postData = {
      login_type: params.loginType,
      client_id: getClientId(),
      action: 'login',
      param: params.params,
    }
    
    auth.loginIn(postData).then((res: any) => {
      //console.log(res)
      if (!params.params.rememberMe) {
        // loginForm.value.rememberMe = false
        // loginForm.value.username = ''
        // loginForm.value.password = ''
      }
      if (res.success) {
        successMsg('登录成功')
        userInfo.value = res.data.user
        isLoginState.value = true
        desktopStore.initDesktop().then(() => {
          router.push('/')
        })
      } else {
        errMsg(res.message)
      }
    })
  }
  const onThirdPartyLogin = async (platform: string) => {
    //console.log(platform)
    thirdPartyLoginMethod.value = platform
    switch (platform) {
      case 'github':
        return await authWithGithub()
      case 'gitee':
        return await authWithGitee()
      case 'dingding':
        return await authWithDingDingScan()
      case 'qyweixin':
        return await authWithWechat()
      default:
        return
    }
  }
  const initThirdPartyLogin = async () => {
    
    const list = await auth.getThirdpartyList()
    // console.log(list)
    if (list && list.length > 0) {
      let res: any = []
      const namekey: any = {
        password: '用户名密码登录',
        github: 'GitHub',
        gitee: 'Gitee',
        dingding: '钉钉扫码登录',
        qyweixin: '企业微信扫码登录',
        phone: '手机号登录',
        email: '邮箱登录',
        ldap: 'LDAP登录',
      }
      list.forEach((item: string) => {
        res.push({
          name: item,
          icon: `/os/images/login/${item}.png`,
          content: namekey[item],
        })
      })
      //console.log(res)
      thirdpartyList.value = res
    }
  }
  const authWithGithub = async (): Promise<boolean> => {
    // 传递state用于防止CSRF攻击,使用时间戳加随机字符串
    const state = Date.now() + Math.random().toString(36).substring(2, 15)
    // 获取当前页面url当做回调参数
    const currentUrl = window.location.href
    const url = '/user/github/authorize?state=' + state
    const res: any = await fetch(url, {
      method: 'POST',
      body: JSON.stringify({
        state: state,
        redirect_url: currentUrl,
      }),
    })
    if (!res.ok) {
      return false
    }
    const data = await res.json()
    if (data && data.data && data.data.url) {
      // 使用正则表达式检查URL格式
      const urlPattern = /client_id=[^&]+/
      if (urlPattern.test(data.data.url)) {
        window.location.href = data.data.url
        return true
      } else {
        errMsg('请先在系统配置中设置github登陆配置')
        return false
      }
    } else {
      errMsg('获取授权URL失败')
      return false
    }
  }

  const authWithWechat = async (): Promise<boolean> => {
    await loadScript(
      'http://rescdn.qqmail.com/node/ww/wwopenmng/js/sso/wwLogin-1.0.0.js'
    )
    const res = await fetch('/user/qyweixin/conf')

    if (res.ok) {
      const data = await res.json()
      if (data.success) {
          (window as any).WwLogin({
            id: 'qywechat-qr-code',
            appid: data.data.corp_id,
            agentid: data.data.agent_id,
            redirect_uri: data.data.redirect,
            state: 'WWLogin',
          })
        return true
      }
      return false
    } else {
      errMsg('网络错误，无法获取二维码')
      return false
    }
  }

  const authWithGitee = async (): Promise<boolean> => {
    // 传递state用于防止CSRF攻击,使用时间戳加随机字符串
    const state = Date.now() + Math.random().toString(36).substring(2, 15)
    // 获取当前页面url当做回调参数
    const currentUrl = window.location.href
    const url = '/user/gitee/authorize?state=' + state
    const res: any = await fetch(url, {
      method: 'POST',
      body: JSON.stringify({
        state: state,
        redirect_url: currentUrl,
      }),
    })
    if (!res.ok) {
      return false
    }
    const data = await res.json()
    if (data && data.data && data.data.url) {
      // 使用正则表达式检查URL格式
      const urlPattern = /client_id=[^&]+/
      if (urlPattern.test(data.data.url)) {
        window.location.href = data.data.url
        return true
      } else {
        errMsg('请先在系统配置中设置gitee登陆配置')
        return false
      }
    } else {
      errMsg('获取授权URL失败')
      return false
    }
    // store.page = "phone";
    // return true;
  }
  const authWithDingDingScan = async () => {
    try {
      // 加载钉钉登录脚本
      await loadScript(
        'https://g.alicdn.com/dingding/h5-dingtalk-login/0.21.0/ddlogin.js'
      )

      const res = await fetch('/user/ding/conf')
      const data = await res.json()

        // 在这里可以调用DTFrameLogin或其他依赖于该脚本的方法
        ; (window as any).DTFrameLogin(
          {
            id: 'dd-qr-code',
            width: 230,
            height: 240,
          },
          {
            redirect_uri: encodeURIComponent(data.data.host),
            client_id: data.data.client_id,
            scope: 'openid',
            response_type: 'code',
            state: 'xxxxxxxxx',
            prompt: 'consent',
          },
          (loginResult: any) => {
            const { authCode } = loginResult
            onLogin(authCode)
          },
          (errMsg: any) => {
            console.log('二维码获取错误', errMsg)
          }
        )
    } catch (error) {
      console.error(error)
    }
  }

  // 获取基础登录方式
  const getBasicLoginMethods = () => {
    return thirdpartyList.value.filter((method) =>
      ['password', 'phone', 'email'].includes(method.name)
    )
  }

  // 获取第三方登录方式
  const getThirdPartyLoginMethods = () => {
    return thirdpartyList.value.filter(
      (method) => !['password', 'phone', 'email'].includes(method.name)
    )
  }

  return {
    userInfo,
    loginForm,
    isLoginState,
    isLogining,
    thirdPartyPlatform,
    thirdPartyLoginMethod,
    thirdpartyList,
    initThirdPartyLogin,
    onThirdPartyLogin,
    checkLogin,
    onRegister,
    onLogin,
    loginOut,
    backToLogin,
    getBasicLoginMethods,
    getThirdPartyLoginMethods,
  }
}, {
  persist: {
    key: 'loginStore',
    pick: ['userInfo', 'isLoginState', 'loginForm'],
  },
})
