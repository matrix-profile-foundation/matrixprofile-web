<template>
  <div>
    <highcharts :options="chartOptions"></highcharts>
    <b-input-group prepend="m">
      <b-form-input type="number" placeholder="subsequence length">
      </b-form-input>
      <b-input-group-append><b-btn>Calculate</b-btn></b-input-group-append>
    </b-input-group>
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
