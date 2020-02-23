import axios from 'axios'
/*
import router from './router'
import store from './store'

axios.interceptors.request.use(
  config => {
    config.baseURL = store.state.baseURL // product env
    // config.baseURL = '/api/' // dev env
    config.timeout = 1800000
    config.headers = {'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8'}
    if (store.state.token !== '') {
      config.headers['X-USER-TOKEN'] = store.state.token
    }
    return config
  }
)
axios.interceptors.response.use(
  response => {
    return response
  },
  error => {
    if (error.response) {
      switch (error.response.status) {
        case 401:
          store.dispatch('LOGOUT')
          router.replace({
            path: '/login'
          })
      }
    }
    return Promise.reject(error.response.data)
  }
)
*/

export default axios
