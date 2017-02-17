import Vue from 'vue'
import Vuex from 'vuex'

import middlewares from './middlewares'
import user from './modules/user'
import userAuth from './modules/userAuth'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    user,
    userAuth
  },
  strict: false,
  debug: process.env.NODE_ENV !== 'production',
  middlewares
})
