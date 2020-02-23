<template>
  <div class="charLib">
    <Button @click="getCharLib">展示charLib</Button>
    <div v-for="value in charLibData">
      <h1>{{value.Id}}</h1>
      <img :src="value.ImgUrl"/>
      <hr>
    </div>
  </div>
</template>

<script>
import Ajax from '../ajax'
export default {
  name: 'charLib',

  data () {
    return {
      charLibData: []
    }
  },

  methods: {
    getCharLib: function () {
      let _this = this
      Ajax({
        method: 'GET',
        url: '/api/charlib'
      }).then((res) => {
        _this.charLibData = res.data.DataList
        console.log(_this.charLibData)
      }, (error) => {
        console.log('error', '获取图片库失败: ' + error)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
