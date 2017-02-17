import request from 'superagent'
import saPrefix from 'superagent-prefix'
import CryptoJS from 'crypto-js'

import localStore from '../utils/localStore'

const API_ROOT = process.env.API_ROOT
const prefix = saPrefix(API_ROOT)

const setSlothToken = function (req) {
  let ua = localStore.getItem('sloth.user_auth')
  if (!ua || !ua.id || ua.id === '') {
    return req
  }
  console.log('api', 'setSlothToken', ua)
  let apiKey = ua.id
  let token = ua.token
  let tsMsg = apiKey + ':' + parseInt(new Date().getTime() / 1000)
  console.log('setSlothToken', tsMsg, token)
  let sign = CryptoJS.HmacSHA256(tsMsg, token).toString(CryptoJS.enc.Hex)
  return req.set('X-Signature', tsMsg + ':' + sign)
}

export default {
  ping: () => {
    return request.get('/_ping')
      .use(prefix)
  },
  status: () => {
    return request.get('/status')
      .use(prefix)
  },
  login: (body) => {
    return request.post('/login')
      .send(body)
      .use(prefix)
      .set('Accept', 'application/json')
  },
  logout: () => {
    return request.delete('/logout')
      .use(prefix)
      .use(setSlothToken)
  },
  signup: (body) => {
    return request.post('/signup')
      .send(body)
      .use(prefix)
      .set('Accept', 'application/json')
  },
  getUser: (id) => {
    return request.get('/user/' + id)
      .use(prefix)
  },

  // settings
  addSettings: (body) => {
    return request.post('/settings')
      .send(body)
      .use(prefix)
      .use(setSlothToken)
      .set('Accept', 'application/json')
  },
  getSettings: (keys) => {
    return request.get('/settings')
      .query({key: keys.join(',')})
      .use(prefix)
      .use(setSlothToken)
  },

  // github's access url
  getGHAccessURL: () => {
    return request.get('/github/access_url')
      .use(prefix)
  },
  getGHBindURL: () => {
    return request.get('/github/bind_url')
      .use(prefix)
  },
  // github callback.
  githubAuth: (code) => {
    return request.post('/github/auth')
      .send({code: code})
      .use(prefix)
  },
  githubBind: (code) => {
    return request.post('/github/bind')
      .send({code: code})
      .use(prefix)
  },
  githubUnbind: (body) => {
    return request.delete('/github/unbind')
      .send(body)
      .use(prefix)
  }
}
