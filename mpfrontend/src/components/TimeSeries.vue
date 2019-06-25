<template>
  <div>
    <div>
      <highcharts ref="highchartA" :options="store.tsAOption"></highcharts>
    </div>
    <div v-if="!store.selfJoin">
      <highcharts ref="highchartB" :options="store.tsBOption"></highcharts>
    </div>
  </div>
</template>

<script>
export default {
  name: "TimeSeries",
  props: {
    store: {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
