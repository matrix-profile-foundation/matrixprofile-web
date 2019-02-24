<template>
  <div id="app">
    <b-container fluid>
      <b-row>
        <b-col cols="7">
          <TimeSeries ref="timeseries" :store="store" />
          <MatrixProfile ref="matrixprofile" :store="store" />
          <b-input-group prepend="m">
            <b-form-input v-model="m" type="number" placeholder="subsequence length">
            </b-form-input>
            <b-input-group-append>
              <b-btn @click="calculateMP">Calculate</b-btn>
            </b-input-group-append>
          </b-input-group>
        </b-col>
        <b-col cols="5">
          <b-nav tabs>
            <b-nav-item @click="enableMotifs">Motifs</b-nav-item>
            <b-nav-item @click="enableDiscords">Discords</b-nav-item>
            <b-nav-item @click="enableSegments">Segments</b-nav-item>
          </b-nav>
          <div v-if="motifsActive">
            <b-form inline>
              <b-input-group size="sm" class="mb-2 mr-sm-2 mb-sm-0" prepend="top-k">
                <b-form-input v-model="kmotifs" type="number" placeholder="max number of motifs">
                </b-form-input>
              </b-input-group>
              <b-input-group size="sm" class="mb-2 mr-sm-2 mb-sm-0" prepend="radius">
                <b-form-input v-model="r" type="number" placeholder="radius">
                </b-form-input>
              </b-input-group>
              <b-btn size="sm" @click="getMotifs">Find</b-btn>
            </b-form>

            <Motifs ref="motifs" :store="store" />
          </div>
          <div v-if="discordsActive">
            <b-form inline>
              <b-input-group size="sm" class="mb-2 mr-sm-2 mb-sm-0" prepend="top-k">
                <b-form-input v-model="kdiscords" type="number" placeholder="max number of discords">
                </b-form-input>
              </b-input-group>
              <b-btn size="sm" @click="getDiscords">Find</b-btn>
            </b-form>

            <Discords :store="store" />
          </div>
          <div v-if="segmentsActive">
            <Segments :store="store" />
          </div>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import TimeSeries from "./components/TimeSeries.vue";
import Motifs from "./components/Motifs.vue";
import Discords from "./components/Discords.vue";
import Segments from "./components/Segments.vue";
import MatrixProfile from "./components/MatrixProfile.vue";
import axios from "axios";

export default {
  name: "app",
  data() {
    return {
      motifsActive: true,
      discordsActive: false,
      segmentsActive: false,
      ts: [],
      n: 0,
      m: 800,
      kmotifs: 3,
      r: 2,
      motifs: [],
      kdiscords: 3,
      discords: [],
      err: "",
      store: {
        tsOption: createChartOption("Time Series", []),
        matrixProfileOption: createChartOption("Matrix Profile", []),
        motifOptions: [],
        discordOptions: [],
        segmentOptions: []
      }
    };
  },
  components: {
    TimeSeries,
    Motifs,
    Discords,
    Segments,
    MatrixProfile
  },
  created: function() {
    this.getTimeSeries();
  },
  methods: {
    enableMotifs: function() {
      this.motifsActive = true;
      this.discordsActive = false;
      this.segmentsActive = false;
    },
    enableDiscords: function() {
      this.motifsActive = false;
      this.discordsActive = true;
      this.segmentsActive = false;
    },
    enableSegments: function() {
      this.motifsActive = false;
      this.discordsActive = false;
      this.segmentsActive = true;
    },
    calculateMP: function() {
      axios.get("http://localhost:8081/calculate", {
        params: {
          m: this.m
        }
      }).then(
        result => {
          var option = createChartOption(
            "Matrix Profile",
            result.data
          );

          option.xAxis[0].max = this.n;
          this.store.matrixProfileOption = option;

          this.getMotifs();
          this.getDiscords();
        },
        error => {
          this.err = JSON.stringify(error);
        }
      );
      this.getSegments();
    },
    getTimeSeries: function() {
      axios.get("http://localhost:8081/data").then(
        result => {
          this.ts = result.data;
          this.n = result.data.length;
          this.store.tsOption = createChartOption(
            "Time Series",
            result.data
          );
          //console.log(this.$refs.timeseries.$refs.highchart.chart);
        },
        error => {
          this.err = JSON.stringify(error);
        }
      );
    },
    getMotifs: function() {
      axios.get("http://localhost:8081/topkmotifs", {
        params: {
          m: this.m,
          k: this.kmotifs,
          r: this.r
        }
      }).then(
        result => {
          this.motifs = result.data;

          var options = [];

          for (var i in this.motifs.series) {
            var motifGroup = this.motifs.series[i];
            if (motifGroup.length != 0) {
              options.push({
                chartOptions: this.createMotifChartOption(
                  "motif "+i+": "+this.motifs.groups[i].MinDist.toFixed(2),
                  motifGroup.slice(0, Math.min(10, motifGroup.length)),
                  this.motifs.groups[i].Idx.slice(0, Math.min(10, motifGroup.length))
                )
              })
            } else {
              break;
            }
          }

          this.store.motifOptions = options;
        },
        error => {
          this.err = JSON.stringify(error);
        }
      );
    },
    getDiscords: function() {
      axios.get("http://localhost:8081/topkdiscords", {
        params: {
          m: this.m,
          k: this.kdiscords
        }
      }).then(
         result => {
          this.discords = result.data;

          var options = [];

          for (var i in this.discords.series) {
            options.push({
              chartOptions: this.createMotifChartOption(
                "discord "+i,
                [this.discords.series[i]],
                [this.discords.groups[i]]
              )
            })
          }

          this.store.discordOptions = options;
        },
        error => {
          this.err = JSON.stringify(error);
        }
      );
    },
    getSegments: function() {
      // likely makes an api call to find motifs
      this.store.segmentOptions = [
        { chartOptions: this.createMotifChartOption("segment 1", [[3, 2, 1]], [1]) }
      ];
    },
    getM: function() {
      return this.m;
    },
    createMotifChartOption: function(title, data, startIndices) {
      var self = this;
      var option = {
        chart: {
          height: "300px",
          events: {
            click: function(e) {
              var motifNum = e.path[3].id;
              var motifs = self.$refs.motifs.$refs;
              var series = motifs[motifNum][0].chart.series;
              self.store.tsOption.xAxis[0].plotBands.length = 0;

              var points = {
                name: startIdx,
                showInLegend: false,
                type: "scatter",
                lineWidth: 0,
                marker: {
                  symbol: "circle",
                  fillColor: "#FF0000",
                  lineWidth: 1
                },
                data: []
              }

              for (var i in series) {
                var startIdx = series[i].name;
                self.store.tsOption.xAxis[0].plotBands.push({
                  from: startIdx,
                  to: startIdx+parseInt(self.getM(), 10),
                  color: "#FCFFC5"
                })

                points.data.push([startIdx, self.$refs.matrixprofile.$refs.highcharts.chart.series[0].data[startIdx].y])

             }

              if (self.store.matrixProfileOption.series.length > 1) {
                self.store.matrixProfileOption.series.pop()
              }
              self.store.matrixProfileOption.series.push(points);

            }
          }
        },
        tooltip: { enabled: false },
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

      option.plotOptions = {
        stickyTracking: false,
        series: {
          lineWidth: 1,
          animation: false,
          events: {
            click: function(e) {
              self.store.tsOption.xAxis[0].plotBands.length = 0;
              var startIdx = parseInt(e.point.series.userOptions.id, 10)
              self.store.tsOption.xAxis[0].plotBands.push({
                from: startIdx,
                to: startIdx+parseInt(self.getM(), 10),
                color: LightenDarkenColor(e.point.series.color, 20)
              })

              var point = {
                name: startIdx,
                showInLegend: false,
                type: "scatter",
                marker: {
                  symbol: "circle",
                  fillColor: "#FF0000",
                  lineWidth: 1
                },
                data: [[startIdx, self.$refs.matrixprofile.$refs.highcharts.chart.series[0].data[startIdx].y]]
              }
              if (self.store.matrixProfileOption.series.length > 1) {
                self.store.matrixProfileOption.series.pop()
              }
              self.store.matrixProfileOption.series.push(point);
            }
          }
        }
      };

      for (var i in data) {
        option.series.push({
          name: startIndices[i],
          id: startIndices[i],
          showInLegend: false,
          data: data[i]
        });
      }

      return option;
    }
  }
};

function createChartOption(title, data) {
  var option = {
    chart: { height: "300", zoomType: "x" },
    title: { text: title },
    plotOptions: {
      stickyTracking: false,
      series: {
        lineWidth: 1,
        animation: false
      }
    },
    series: [
      {
        showInLegend: false,
        data: data
      }
    ],
    xAxis: [
      {
        plotBands: []
      }
    ]
  };

  return option;
}

function LightenDarkenColor(col, amt) {
    var usePound = false;
    if (col[0] == "#") {
        col = col.slice(1);
        usePound = true;
    }
    var num = parseInt(col,16);
    var r = (num >> 16) + amt;
    if (r > 255) r = 255;
    else if  (r < 0) r = 0;
    var b = ((num >> 8) & 0x00FF) + amt;
    if (b > 255) b = 255;
    else if  (b < 0) b = 0;
    var g = (num & 0x0000FF) + amt;
    if (g > 255) g = 255;
    else if (g < 0) g = 0;
    return (usePound?"#":"") + (g | (b << 8) | (r << 16)).toString(16);
}
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  height: 100%;
}
</style>
