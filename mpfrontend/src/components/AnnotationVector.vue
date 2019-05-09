<template>
  <div>
    <highcharts
      ref="highcharts"
      :options="store.annotationVectorOption"
    ></highcharts>
  </div>
</template>

<script>
export default {
  name: "AnnotationVector",
  props: {
    store: {}
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
