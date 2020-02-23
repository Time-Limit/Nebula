<template>
  <div class="invalidBillImg">
    <Button @click="getInvalidBillImg">展示无效的账单图片</Button>
    <div v-for="value in dataList">
      <img :src="value.ImgUrl"/>
      <p>{{value.ImgUrl}}</p>
      <p>{{value.Time}}</p>
      <p>{{value.Price}}</p>
      <p>{{value.Err}}</p>
      <hr>
    </div>
  </div>
</template>

<script>
import Ajax from '../ajax'
export default {
  name: 'invalidBillImg',

  data () {
    return {
      dataList: []
    }
  },

  methods: {
    getInvalidBillImg: function () {
      let _this = this
      Ajax({
        method: 'GET',
        url: '/api/bill/image/invalid'
      }).then((res) => {
        _this.dataList = res.data.DataList
        console.log(_this.dataList)
      }, (error) => {
        console.log('error', '获取数据失败: ' + error)
      })
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
