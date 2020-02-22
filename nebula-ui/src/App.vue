<template>
  <div id="app">
    <Menu mode="horizontal" theme="light" @on-select="changeRoute">
      <MenuItem name="stroke"> <Icon type="ios-paper" />笔画</MenuItem>
    </Menu>
    <router-view/>
  </div>
</template>

<script>
import Stroke from '@/components/stroke'
export default {
  name: 'App',
  methods: {
    changeRoute: function (routeName) {
      console.log(this.routeMap[routeName])
      this.$router.push(this.routeMap[routeName])
    }
  },
  data () {
    return {
      routeMap: {
        'stroke': '/stroke'
      }
    }
  },
  components: { Stroke }
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
