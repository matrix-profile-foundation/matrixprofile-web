<template>
  <div class="motifs">
    <div v-for="(data, index) in motifOptions" :key="index">
      <highcharts :options="data.chartOptions"></highcharts>
    </div>
  </div>
</template>

<script>
export default {
  name: "Motifs",
  data() {
    return {
      motifOptions: [
        {
          chartOptions: {
            chart: { height: "300px" },
            title: { text: "motif 1" },
            xAxis: [{ labels: { enabled: false } }],
            yAxis: [
              {
                title: "",
                labels: { enabled: false }
              }
            ],
            series: [
              {
                showInLegend: false,
                data: [1, 2, 3]
              },
              {
                showInLegend: false,
                data: [1.1, 2.1, 3.1]
              }
            ]
          }
        },
        {
          chartOptions: {
            chart: { height: "300px" },
            title: { text: "motif 2" },
            xAxis: [{ labels: { enabled: false } }],
            yAxis: [
              {
                title: "",
                labels: { enabled: false }
              }
            ],
            series: [
              {
                showInLegend: false,
                data: [1, 2, 1]
              }
            ]
          }
        },
        {
          chartOptions: {
            chart: { height: "300px" },
            title: { text: "motif 3" },
            xAxis: [{ labels: { enabled: false } }],
            yAxis: [
              {
                title: "",
                labels: { enabled: false }
              }
            ],
            series: [
              {
                showInLegend: false,
                data: [1, 2, 0]
              }
            ]
          }
        }
      ]
    };
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
