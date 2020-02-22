import Vue from 'vue'
import Router from 'vue-router'
import Stroke from '@/components/stroke'
import CharLib from '@/components/charLib'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/stroke',
      name: 'stroke',
      component: Stroke
    },
    {
      path: '/charlib',
      name: 'charlib',
      component: CharLib
    }
  ]
})
