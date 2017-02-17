import * as types from '../types'
import localStore from '../../utils/localStore'

const state = {
  id: ''
}

const getters = {
  getUserAuth: (state) => {
    if (state.id === '') {
      let localUA = localStore.getItem('sloth.user_auth')
      if (localUA) {
        for (var k in localUA) {
          state[k] = localUA[k]
        }
        console.log('store getUserAuth from localStore')
        return state
      }
    }
    console.log('store getUserAuth', state)
    return state
  }
}

const actions = {
  setUserAuth ({ commit }, userauth) {
    commit(types.SET_USER_AUTH, { userauth })
  },
  removeUserAuth ({ commit }) {
    commit(types.REMOVE_USER_AUTH)
  }
}

const mutations = {
  [types.SET_USER_AUTH] (state, {userauth}) {
    console.log('store set user_auth', userauth)
    for (var k in userauth) {
      state[k] = userauth[k]
    }
    localStore.setItem('sloth.user_auth', userauth)
  },
  [types.REMOVE_USER_AUTH] (state) {
    console.log('store remove user_auth')
    localStore.rmItem('sloth.user_auth')
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
