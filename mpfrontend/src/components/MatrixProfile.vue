<template>
  <div>
    <highcharts
      ref="highcharts"
      :options="store.matrixProfileOption"
    ></highcharts>
  </div>
</template>

<script>
export default {
  name: "MatrixProfile",
  props: {
    store: {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
