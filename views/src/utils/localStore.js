export default {
  setItem: (key, data) => {
    let str = ''
    if (typeof data === 'string') {
      str = data
    } else {
      str = JSON.stringify(data)
    }
    window.localStorage.setItem(key, str)
    console.log('setItem:', key, str)
  },
  getItem: (key) => {
    let data = window.localStorage.getItem(key)
    try {
      return JSON.parse(data)
    } catch (e) {
      return data
    }
  },
  rmItem: (key) => {
    window.localStorage.removeItem(key)
  }
}
