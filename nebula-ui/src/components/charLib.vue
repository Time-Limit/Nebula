<template>
  <div class="charLib">
    <Button @click="getCharLib">展示charLib</Button>
    <div v-for="value in charLibData">
      <pre>{{value.id}}</pre>
      <img :src="value.imgUrl"><img>
    </div>
  </div>
</template>

<script>
export default {
  name: 'charLib',

  data () {
    return {
      charLibData: []
    }
  },

  methods: {
    getCharLib: function () {
      this.$http.get('api/charlib').then((response) => {
         console.info(response.body)
      }, (response) => {
         console.error(response)
      });
    }
  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
</style>
