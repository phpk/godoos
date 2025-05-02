import { get } from '@/utils/request'
export const getApplyList = (page: number, params?: any) => {
  if (params) {
    return get('/workflow/datalist?workType=0&knowId=0&limit=10', {
      page,
      param: 'flowId=' + params,
    }).then((res) => res.data)
  } else {
    return get('/workflow/datalist?workType=0&knowId=0&limit=10', {
      page,
    }).then((res) => res.data)
  }
}

export const getMyApply = () =>
  get('/workflow/list?workType=0&knowId=0').then((res) => res.data)

export const getMyTask = () =>
  get('/workflow/list?workType=1&knowId=0').then((res) => res.data)

export const getTaskList = (page: number, params?: any) => {
  if (params) {
    return get('/workflow/datalist?workType=1&knowId=0&limit=10', {
      page,
      param: 'flowId=' + params,
    }).then((res) => res.data)
  } else {
    return get('/workflow/datalist?workType=1&knowId=0&limit=10', {
      page,
    }).then((res) => res.data)
  }
}

export const getMyCheckList = (page: number) =>
  get('/workflow/mychecklist?knowId=0&limit=10', { page }).then(
    (res) => res.data
  )

// 我的抄送
export const getMyCopyList = (page: number, params: any[]) =>
  get('/workflow/getcopytomelist?&limit=10&knowId=0', {
    page,
    start_date: params[0],
    end_date: params[1],
  }).then((res) => res.data)

// 历史抄送
export const getCopyFromMeList = (page: number, params: any[]) =>
  get('/workflow/getcopyfrommelist?&limit=10&knowId=0', {
    page,
    start_date: params[0],
    end_date: params[1],
  }).then((res) => res.data)
