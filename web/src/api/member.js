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

export const findPointAccount = (params) => {
  return service({
    url: '/pointAccount/findPointAccount',
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

export const getPointTransactionList = (params) => {
  return service({
    url: '/pointTransaction/getPointTransactionList',
    method: 'get',
    params
  })
}

export const createPointProduct = (data) => {
  return service({
    url: '/pointProduct/createPointProduct',
    method: 'post',
    data
  })
}

export const deletePointProduct = (data) => {
  return service({
    url: '/pointProduct/deletePointProduct',
    method: 'delete',
    data
  })
}

export const updatePointProduct = (data) => {
  return service({
    url: '/pointProduct/updatePointProduct',
    method: 'put',
    data
  })
}

export const findPointProduct = (params) => {
  return service({
    url: '/pointProduct/findPointProduct',
    method: 'get',
    params
  })
}

export const getPointProductList = (params) => {
  return service({
    url: '/pointProduct/getPointProductList',
    method: 'get',
    params
  })
}

export const updatePointProductStatus = (data) => {
  return service({
    url: '/pointProduct/updatePointProductStatus',
    method: 'put',
    data
  })
}

export const getPointProductOptions = (params) => {
  return service({
    url: '/pointProduct/getPointProductOptions',
    method: 'get',
    params
  })
}

export const createRedemptionOrder = (data) => {
  return service({
    url: '/redemptionOrder/createRedemptionOrder',
    method: 'post',
    data
  })
}

export const findRedemptionOrder = (params) => {
  return service({
    url: '/redemptionOrder/findRedemptionOrder',
    method: 'get',
    params
  })
}

export const getRedemptionOrderList = (params) => {
  return service({
    url: '/redemptionOrder/getRedemptionOrderList',
    method: 'get',
    params
  })
}

export const completeRedemptionOrder = (data) => {
  return service({
    url: '/redemptionOrder/completeRedemptionOrder',
    method: 'post',
    data
  })
}

export const cancelRedemptionOrder = (data) => {
  return service({
    url: '/redemptionOrder/cancelRedemptionOrder',
    method: 'post',
    data
  })
}
