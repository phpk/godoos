import { noticeMsg} from "@/utils/msg";
import { get, post } from '@/utils/request'
import { onMounted, reactive, watch } from 'vue'

export const useHandler = (domRef: any) => {
  const state = reactive({ popoverVisible: false, type: "default" })
  const cData = reactive({ image: "", thumb: "", captKey: "", thumbX: 0, thumbY: 0, thumbWidth: 0, thumbHeight: 0 })

  const clickEvent = () => {
    state.popoverVisible = true
  }

  const visibleChangeEvent = (visible: boolean) => {
    state.popoverVisible = visible
  }

  const closeEvent = () => {
    state.popoverVisible = false
  }

  const requestCaptchaData = () => {
    domRef.value.clear && domRef.value.clear()

    get('user/getcaptcha').then((response) => {
      //console.log(response)
      const { data = {} } = response;
      if (response.code === 0) {
        cData.image = data['image_base64'] || ''
        cData.thumb = data['tile_base64'] || ''
        cData.captKey = data['captcha_key'] || ''
        cData.thumbX = data['tile_x'] || 0
        cData.thumbY = data['tile_y'] || 0
        cData.thumbWidth = data['tile_width'] || 0
        cData.thumbHeight = data['tile_height'] || 0
      } else {
        noticeMsg(`获取验证码失败`, '提示', 'error')
      }
    }).catch((e) => {
      console.warn(e)
    })
  }

  const refreshEvent = () => {
    requestCaptchaData()
  }

  const confirmEvent = (point: any, clear: any) => {
    //console.log(cData)
    post('user/checkcaptcha', {
      point: [point.x, point.y].join(','),
      key: cData.captKey || ''
    }).then((response) => {
      //console.log(response)
      if (response.success) {
        //successMsg(`check data success`)
        state.popoverVisible = false
        state.type = "success"
        return;
      } else {
        noticeMsg(`校验失败`, '提示', 'error')
        //warningMsg(`check data failed`)
        state.type = "error"
      }

      setTimeout(() => {
        requestCaptchaData()
      }, 1000)
    }).catch((e) => {
      console.warn(e)
    })
  }

  watch(() => state.popoverVisible, () => {
    if (state.popoverVisible) {
      requestCaptchaData()
    }
  })

  onMounted(() => {
    requestCaptchaData()
  })

  return {
    state,
    data: cData,
    visibleChangeEvent,
    clickEvent,
    closeEvent,
    refreshEvent,
    confirmEvent,
  }
}