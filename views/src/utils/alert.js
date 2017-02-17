import toastr from 'toastr'

toastr.options = {
  closeButton: true,
  progressBar: true,
  showMethod: 'slideDown',
  positionClass: 'toast-top-full-width', // 'toast-top-center',
  timeOut: 4000
}

const debug = process.env.NODE_ENV !== 'production'

export default {
  success: (content, title) => {
    if (debug) {
      console.log('success', title, content)
    }
    toastr.success(content, title)
  },
  error: (content, title) => {
    if (debug) {
      console.log('error', title, content)
    }
    toastr.error(content, title)
  },
  warn: (content, title) => {
    if (debug) {
      console.log('warn', title, content)
    }
    toastr.warning(content, title)
  },
  // check http response.
  check: (vm, resp) => {
    if (resp) {
      if (resp.status >= 200 && resp.status < 300) {
        return true
      } else if (resp.status === 401) {
        vm.$router.push('/login')
      } else {
        toastr.error(resp.body.error, '请求错误')
      }
    } else {
      toastr.error('连接到服务器失败', '连接错误')
    }
    return false
  }
}
