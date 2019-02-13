<template>
  <div class="motifs">
    <div v-for="(data, index) in motifOptions" :key="index">
      {{ store.message }}
      <highcharts :options="data.chartOptions"></highcharts>
    </div>
  </div>
</template>

<script>
export default {
  name: "Motifs",
  props: {
    store: {}
  },
  created() {
    this.fillMotifOptions();
  },
  data() {
    return {
      motifOptions: []
    };
  },
  methods: {
    motifChartOption: function(title, data) {
      var option = {
        chart: { height: "300px" },
        title: { text: title },
        xAxis: [{ labels: { enabled: false } }],
        yAxis: [
          {
            title: "",
            labels: { enabled: false }
          }
        ],
        series: []
      };

      for (var i in data) {
        option.series.push({
          showInLegend: false,
          data: data[i]
        });
      }
      return option;
    },
    fillMotifOptions: function() {
      // likely makes an api call to find motifs
      this.motifOptions = [
        {
          chartOptions: this.motifChartOption("motif 1", [
            [1, 2, 3],
            [1.1, 2.1, 3.1]
          ])
        },
        {
          chartOptions: this.motifChartOption("motif 2", [[1, 2, 1]])
        },
        {
          chartOptions: this.motifChartOption("motif 3", [[1, 2, 0]])
        }
      ]
    }
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
