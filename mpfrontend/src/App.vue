<template>
  <div id="app">
    <b-container fluid>
      {{ store.tsOption.xAxis[0].min }}
      {{ store.tsOption.xAxis[0].max }}
      <b-row>
        <b-col cols="7">
          <TimeSeries :store="store" />
          <MatrixProfile :store="store" />
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
              <b-input-group class="mb-2 mr-sm-2 mb-sm-0" prepend="top-k">
                <b-form-input v-model="k" type="number" placeholder="max number of motifs">
                </b-form-input>
              </b-input-group>
              <b-input-group class="mb-2 mr-sm-2 mb-sm-0" prepend="radius">
                <b-form-input v-model="r" type="number" placeholder="radius">
                </b-form-input>
              </b-input-group>
              <b-btn @click="getMotifs">Find</b-btn>
            </b-form>

            <Motifs :store="store" />
          </div>
          <div v-if="discordsActive">
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
      m: 128,
      k: 3,
      r: 1,
      motifs: [],
      err: "",
      store: {
        message: "Hello!",
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
        },
        error => {
          this.err = JSON.stringify(error);
        }
      );
      this.getDiscords();
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
        },
        error => {
          this.err = JSON.stringify(error);
        }
      );
    },
    getMotifs: function() {
      axios.get("http://localhost:8081/topkmotifs", {
        params: {
          k: this.k,
          r: this.r,
          m: this.m
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
                  "motif "+i+": "+this.motifs.Groups[i].MinDist.toFixed(2),
                  motifGroup.slice(0, Math.min(10, motifGroup.length)),
                  this.motifs.Groups[i].Idx.slice(0, Math.min(10, motifGroup.length))
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
      // likely makes an api call to find motifs
      this.store.discordOptions = [
        { chartOptions: this.createMotifChartOption("discord 1", [[3, 2, 1]], [1]) }
      ];
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
          height: "200px",
          events: {
            click: function(e) {
              console.log(e);
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
        plotOptions: {
          series: {
            lineWidth: 1,
            animation: false,
            events: {
              click: function(e) {
                var startIdx = parseInt(e.point.series.userOptions.id, 10)
                self.store.tsOption.xAxis[0].plotBands[0].from = startIdx;
                self.store.tsOption.xAxis[0].plotBands[0].to = startIdx+parseInt(self.getM(), 10);
                self.store.matrixProfileOption.xAxis[0].plotBands[0].from = startIdx;
                self.store.matrixProfileOption.xAxis[0].plotBands[0].to = startIdx+parseInt(self.getM(), 10);
              }
            }
          }
        },
        series: []
      };

      for (var i in data) {
        option.series.push({
          name: "idx_"+startIndices[i],
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
        plotBands: [
          {
            from: 0,
            to: 0,
            color: "#FCFFC5",
            id: "plot-band-1"
          }
        ]
      }
    ]
  };

  return option;
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
