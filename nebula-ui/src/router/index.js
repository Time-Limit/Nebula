import Vue from 'vue'
import Router from 'vue-router'
import Stroke from '@/components/stroke'
import CharLib from '@/components/charLib'
import InvalidBillImg from '@/components/invalidBillImg'

Vue.use(Router)

export default new Router({
  routes: [
    {
      path: '/stroke',
      name: 'stroke',
      component: Stroke
    },
    {
      path: '/invalidbillimg',
      name: 'invalidbillimg',
      component: InvalidBillImg
    },
    {
      path: '/charlib',
      name: 'charlib',
      component: CharLib
    }
  ]
})
