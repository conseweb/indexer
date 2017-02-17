<template>
  <router-view></router-view>
</template>

<script>
import api from './api/api'
import Alert from './utils/alert'

export default {
  name: 'app',
  created: function () {
    let vm = this
    api.status().end(function (err, resp) {
      if (err) {
        Alert.error('服务器错误', '请求错误')
        return
      }

      if (resp.body.user === 0) {
        vm.$router.push('/signup')
      }
    })
  }
}
</script>

<style>
#app {
  font-family: 'Avenir', Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
}
</style>
