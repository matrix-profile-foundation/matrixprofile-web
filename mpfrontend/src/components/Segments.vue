<template>
  <div class="segments">
    <div v-for="(data, index) in store.segmentOptions" :key="index">
      <highcharts :options="data.chartOptions"></highcharts>
    </div>
  </div>
</template>

<script>
export default {
  name: "Segments",
  props: {
    store: {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.motifs {
  overflow-y: scroll;
  height: 500px;
}
</style>
