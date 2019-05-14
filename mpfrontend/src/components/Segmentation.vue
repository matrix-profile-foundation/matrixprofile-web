<template>
  <div>
    <highcharts
      ref="highcharts"
      :options="store.segmentationOption"
    ></highcharts>
  </div>
</template>

<script>
export default {
  name: "Segmentation",
  props: {
    store: {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
