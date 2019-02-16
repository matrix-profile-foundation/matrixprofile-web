<template>
  <div>
    <highcharts :options="chartOptions"></highcharts>
  </div>
</template>

<script>
import ts from "../data/testdata.json";
export default {
  name: "TimeSeries",
  data() {
    return {
      chartOptions: {
        chart: { zoomType: "x" },
        tooltip: false,
        title: { text: "Time Series" },
        plotOptions: {
          series: {
            lineWidth: 1,
            states: {
              hover: { enabled: false }
            }
          }
        },
        series: [
          {
            showInLegend: false,
            data: ts.data
          }
        ]
      }
    };
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
