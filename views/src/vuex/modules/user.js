import * as types from '../types'
import localStore from '../../utils/localStore'

const state = {
  id: '',
  role: ''
}

const getters = {
  isAdmin: state => state.role && state.role === 'admin',
  getUser: (state) => {
    if (state.id === '') {
      let localUser = localStore.getItem('sloth.user')
      if (localUser) {
        for (var k in localUser) {
          state[k] = localUser[k]
        }
        console.log('store getUser from localStore', state)
        return state
      }
    }
    console.log('store getUser', state)
    return state
  }
}

const actions = {
  setAccount ({ commit }, user) {
    commit(types.SET_USER, { user })
  },
  removeAccount ({ commit }) {
    commit(types.REMOVE_USER)
  },
  getAccount ({ commit }) {
    commit(types.REMOVE_USER)
  }
}

const mutations = {
  [types.SET_USER] (state, {user}) {
    console.log('store set user', user)
    for (var k in user) {
      state[k] = user[k]
    }
    localStore.setItem('sloth.user', user)
  },
  [types.REMOVE_USER] (state) {
    localStore.rmItem('sloth.user')
    state = null
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
