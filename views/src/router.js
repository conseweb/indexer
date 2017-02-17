import Vue from 'vue'
import VueRouter from 'vue-router'

import Main from './components/Main'

Vue.use(VueRouter)

const TODO = Vue.extend({
  template: '<h2>This is developing!</h2>'
})

const router = new VueRouter({
  mode: 'history',
  // base: __dirname,
  routes: [
    {
      path: '/',
      component: Main,
      children: [
        {
          path: '',
          component: TODO
        },
        {
          path: 'dashboard',
          component: TODO
        }
      ]
    }
  ]
})

router.beforeEach((to, from, next) => {
  next()
})

export default router
