<template>
  <div class="stroke">
    <p v-if="showIndex == -1">no stroke data</p>
    <div v-else>
      <h3>{{ showData.pinyin }}</h3>
      <img :src="showData.imgUrl" />
    </div>
    <Button @click="updateShowIndex(-1)">上一个</Button>
    <Button @click="updateShowIndex(1)">下一个</Button>
    <Input v-model="searchData" placeholder="search" style="width: 300px" />
  </div>
</template>

<script>
export default {
  name: 'Stroke',
  data () {
    return {
      strokeData: require('../stroke/font.json'),
      showIndex: -1,
      searchData: ''
    }
  },
  computed: {
    strokeDataKeyList: function () {
      console.log(this.strokeData)
      let res = []
      for (let key in this.strokeData) {
        res.push(key)
      }
      return res
    },
    showData: function () {
      const key = this.strokeDataKeyList[this.showIndex]
      return {
        'char': key,
        'imgUrl': this.strokeData[key].imgUrl,
        'pinyin': this.strokeData[key].pinyin,
        'bihua': this.strokeData[key].bihua
      }
    }
  },
  watch: {
    searchData: function () {
      if (this.searchData === '') {
        return
      }
      for (let i = 0; i < this.strokeDataKeyList.length; i++) {
        if (this.strokeDataKeyList[i] === this.searchData) {
          this.showIndex = i
          return
        }
      }
      this.$Message['info']({
        background: true,
        content: '还没有"' + this.searchData + '"的数据。'
      })
    }
  },
  methods: {
    updateShowIndex: function (step) {
      this.showIndex += step
      if (this.showIndex >= this.strokeDataKeyList.length) {
        this.showIndex = 0
      }
      if (this.showIndex < 0) {
        this.showIndex = this.stokeDataKeyList.length
      }
      if (this.strokeDataKeyList.length === 0) {
        this.showIndex = -1
      }
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
