import request from '@/utils/request'

const testPrefix = '/v1/test'

const api = {
  profile: `${testPrefix}/profile`,
  tasks: `${testPrefix}/tasks`,
  envs: `${testPrefix}/envs`,

  user: `${testPrefix}/user`,
  role: `${testPrefix}/role`,
  service: `${testPrefix}/service`,
  permission: `${testPrefix}/permission`,
  permissionNoPager: `${testPrefix}/permission/no-pager`,
  orgTree: `${testPrefix}/org/tree`
}

const WebSocketPath = 'api/v1/ws'
export const WebBaseDev = 'http://127.0.0.1:8085/'
export const WebSocketBaseDev = ''
export function GetWebSocketApi () {
  const isProd = process.env.NODE_ENV === 'production'

  let wsUri = ''
  if (!isProd) {
    wsUri = WebBaseDev.replace('http', 'ws')
  } else {
    const loc = window.location

    if (loc.protocol === 'https:') {
      wsUri = 'wss:'
    } else {
      wsUri = 'ws:'
    }
    wsUri += '//' + loc.host
    wsUri += loc.pathname
  }

  wsUri += WebSocketPath

  return wsUri
}

export function requestSuccess (code) {
  return code === 200
}

export function getProfile (parameter) {
  return request({
    url: api.profile,
    method: 'get',
    data: parameter
  })
}

// Task
export function listTask (params) {
  return request({
    url: api.tasks,
    method: 'get',
    params: params
  })
}
export function getTask (id, withIntents) {
  return request({
    url: api.tasks + '/' + id,
    method: 'get',
    params: { withIntents: withIntents }
  })
}
export function saveTask (model) {
  return request({
    url: !model.id ? api.tasks : api.tasks + '/' + model.id,
    method: !model.id ? 'post' : 'put',
    data: model
  })
}
export function disableTask (model) {
  return request({
    url: api.tasks + '/' + model.id + '/disable',
    method: 'post',
    params: {}
  })
}
export function removeTask (model) {
  return request({
    url: api.tasks + '/' + model.id,
    method: 'delete',
    params: {}
  })
}

// Environment
export function getTestEnvs (env) {
  return request({
    url: api.envs,
    method: 'post',
    data: env
  })
}

export function getUserList (parameter) {
  return request({
    url: api.user,
    method: 'get',
    params: parameter
  })
}

export function getRoleList (parameter) {
  return request({
    url: api.role,
    method: 'get',
    params: parameter
  })
}

export function getServiceList (parameter) {
  return request({
    url: api.service,
    method: 'get',
    params: parameter
  })
}

export function getPermissions (parameter) {
  return request({
    url: api.permissionNoPager,
    method: 'get',
    params: parameter
  })
}

// id == 0 add     post
// id != 0 update  put
export function saveService (parameter) {
  return request({
    url: api.service,
    method: parameter.id === 0 ? 'post' : 'put',
    data: parameter
  })
}
