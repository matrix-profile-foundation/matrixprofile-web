<template>
  <div class="discords">
    <div v-for="(data, index) in store.discordOptions" :key="index">
      <highcharts :options="data.chartOptions"></highcharts>
    </div>
  </div>
</template>

<script>
export default {
  name: "Discords",
  props: {
    store: {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
