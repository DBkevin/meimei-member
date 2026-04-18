import service from '@/utils/request'

export const createMember = (data) => {
  return service({
    url: '/member/createMember',
    method: 'post',
    data
  })
}

export const deleteMember = (data) => {
  return service({
    url: '/member/deleteMember',
    method: 'delete',
    data
  })
}

export const updateMember = (data) => {
  return service({
    url: '/member/updateMember',
    method: 'put',
    data
  })
}

export const findMember = (params) => {
  return service({
    url: '/member/findMember',
    method: 'get',
    params
  })
}

export const getMemberList = (params) => {
  return service({
    url: '/member/getMemberList',
    method: 'get',
    params
  })
}

export const updateMemberStatus = (data) => {
  return service({
    url: '/member/updateMemberStatus',
    method: 'put',
    data
  })
}

export const getMemberPointAccount = (params) => {
  return service({
    url: '/member/getMemberPointAccount',
    method: 'get',
    params
  })
}

export const getMemberOptions = (params) => {
  return service({
    url: '/member/getMemberOptions',
    method: 'get',
    params
  })
}

export const getPointAccountList = (params) => {
  return service({
    url: '/pointAccount/getPointAccountList',
    method: 'get',
    params
  })
}

export const manualAddPoints = (data) => {
  return service({
    url: '/pointAccount/manualAddPoints',
    method: 'post',
    data
  })
}

export const manualSubPoints = (data) => {
  return service({
    url: '/pointAccount/manualSubPoints',
    method: 'post',
    data
  })
}

export const getPointLogList = (params) => {
  return service({
    url: '/pointLog/getPointLogList',
    method: 'get',
    params
  })
}

export const createPointGoods = (data) => {
  return service({
    url: '/pointGoods/createPointGoods',
    method: 'post',
    data
  })
}

export const deletePointGoods = (data) => {
  return service({
    url: '/pointGoods/deletePointGoods',
    method: 'delete',
    data
  })
}

export const updatePointGoods = (data) => {
  return service({
    url: '/pointGoods/updatePointGoods',
    method: 'put',
    data
  })
}

export const findPointGoods = (params) => {
  return service({
    url: '/pointGoods/findPointGoods',
    method: 'get',
    params
  })
}

export const getPointGoodsList = (params) => {
  return service({
    url: '/pointGoods/getPointGoodsList',
    method: 'get',
    params
  })
}

export const updatePointGoodsStatus = (data) => {
  return service({
    url: '/pointGoods/updatePointGoodsStatus',
    method: 'put',
    data
  })
}

export const updatePointGoodsStock = (data) => {
  return service({
    url: '/pointGoods/updatePointGoodsStock',
    method: 'put',
    data
  })
}

export const getPointGoodsOptions = (params) => {
  return service({
    url: '/pointGoods/getPointGoodsOptions',
    method: 'get',
    params
  })
}

export const createExchangeOrder = (data) => {
  return service({
    url: '/exchangeOrder/createExchangeOrder',
    method: 'post',
    data
  })
}

export const findExchangeOrder = (params) => {
  return service({
    url: '/exchangeOrder/findExchangeOrder',
    method: 'get',
    params
  })
}

export const getExchangeOrderList = (params) => {
  return service({
    url: '/exchangeOrder/getExchangeOrderList',
    method: 'get',
    params
  })
}

export const verifyExchangeOrder = (data) => {
  return service({
    url: '/exchangeOrder/verifyExchangeOrder',
    method: 'post',
    data
  })
}

export const cancelExchangeOrder = (data) => {
  return service({
    url: '/exchangeOrder/cancelExchangeOrder',
    method: 'post',
    data
  })
}

export const refundExchangeOrder = (data) => {
  return service({
    url: '/exchangeOrder/refundExchangeOrder',
    method: 'post',
    data
  })
}
