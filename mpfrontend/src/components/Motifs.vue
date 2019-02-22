<template>
  <div class="motifs">
    <div v-for="(data, index) in store.motifOptions" :key="index">
      <highcharts :id="index" :ref="index" :options="data.chartOptions"></highcharts>
    </div>
  </div>
</template>

<script>
export default {
  name: "Motifs",
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
