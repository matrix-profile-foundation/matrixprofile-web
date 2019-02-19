<template>
  <div>
    <highcharts ref="highchart" :options="store.tsOption"></highcharts>
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
