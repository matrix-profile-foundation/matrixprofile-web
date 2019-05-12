<template>
  <div id="app">
    <b-navbar toggleable="lg" type="dark" variant="info">
      <b-navbar-brand>Matrix Profiles</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>

        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-item active href="https://www.cs.ucr.edu/~eamonn/MatrixProfile.html">UCR Webpage</b-nav-item>
          <b-nav-item active href="https://github.com/aouyang1/go-matrixprofile"><Octicon :icon="markGithub" /></b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <h5>{{ err }}</h5>
    <b-container fluid>
      <b-row>
        <b-col cols="8">
          <TimeSeries ref="timeseries" :store="store" />
          <AnnotationVector ref="annotationvector" :store="store" />
          <MatrixProfile ref="matrixprofile" :store="store" />
        </b-col>
        <b-col cols="4">
          <b-input-group prepend="m" class="mb-2" >
            <b-form-input
              v-model.number="m"
              type="number"
              placeholder="subsequence length"
            >
            </b-form-input>
            <b-input-group-append>
              <b-btn @click="calculateMP">Calculate</b-btn>
            </b-input-group-append>
          </b-input-group>
          <b-nav tabs>
            <b-nav-item @click="enableMotifs">Motifs</b-nav-item>
            <b-nav-item @click="enableDiscords">Discords</b-nav-item>
            <b-nav-item @click="enableAVs">Annotation Vectors</b-nav-item>
          </b-nav>
          <div v-if="motifsActive">
            <b-form inline>
              <b-input-group
                size="sm"
                class="mb-2 mr-sm-1 mb-sm-0"
                prepend="top-k"
              >
                <b-form-input
                  v-model="kmotifs"
                  type="number"
                  placeholder="max number of motifs"
                >
                </b-form-input>
              </b-input-group>
              <b-input-group
                size="sm"
                class="mb-2 mr-sm-1 mb-sm-0"
                prepend="radius"
              >
                <b-form-input v-model="r" type="number" placeholder="radius">
                </b-form-input>
              </b-input-group>
              <b-btn size="sm" @click="getMotifs">Find</b-btn>
            </b-form>

            <Motifs ref="motifs" :store="store" />
          </div>
          <div v-if="discordsActive">
            <b-form inline>
              <b-input-group
                size="sm"
                class="mb-2 mr-sm-1 mb-sm-0"
                prepend="top-k"
              >
                <b-form-input
                  v-model="kdiscords"
                  type="number"
                  placeholder="max number of discords"
                >
                </b-form-input>
              </b-input-group>
              <b-btn size="sm" @click="getDiscords">Find</b-btn>
            </b-form>

            <Discords :store="store" />
          </div>
          <div v-if="avsActive">
            <b-form inline>
              <b-form-select v-model="selectedav" :options="avoptions" @change="avChange"></b-form-select>
            </b-form>
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
import MatrixProfile from "./components/MatrixProfile.vue";
import AnnotationVector from "./components/AnnotationVector.vue";
import Octicon, { markGithub } from "octicons-vue";
import axios from "axios";

export default {
  name: "app",
  data() {
    return {
      markGithub: markGithub,
      motifsActive: true,
      discordsActive: false,
      avsActive: false,
      ts: [],
      n: 0,
      m: 30,
      kmotifs: 3,
      r: 2,
      motifs: [],
      kdiscords: 3,
      discords: [],
      selectedav: "default",
      avoptions: [
        { value: "default", text: "Default" },
        { value: "complexity", text: "Complexity" },
        { value: "meanstd", text: "Mean Standard Deviation" },
        { value: "clipping", text: "Clipping" }
      ],
      err: "",
      store: {
        tsOption: genTSOption([]),
        annotationVectorOption: genAVOption([]),
        matrixProfileOption: genMPOption([]),
        motifOptions: [],
        discordOptions: []
      }
    };
  },
  components: {
    TimeSeries,
    Motifs,
    Discords,
    MatrixProfile,
    AnnotationVector,
    Octicon
  },
  created: function() {
    this.getTimeSeries();
  },
  methods: {
    enableMotifs: function() {
      this.motifsActive = true;
      this.discordsActive = false;
      this.avsActive = false;
    },
    enableDiscords: function() {
      this.motifsActive = false;
      this.discordsActive = true;
      this.avsActive = false;
    },
    enableAVs: function() {
      this.motifsActive = false;
      this.discordsActive = false;
      this.avsActive = true;
    },
    getTimeSeries: function() {
      axios
        .get(process.env.VUE_APP_MPSERVER_URL + "/data", {
          withCredentials: true
        }).then(
          result => {
            this.ts = result.data;
            this.n = result.data.length;
            var option = genTSOption(result.data);
            option.xAxis[0].max = this.n;
            this.store.tsOption = option;
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
          }
        );
    },
    calculateMP: function() {
      axios
        .post(process.env.VUE_APP_MPSERVER_URL + "/calculate",
          { m: this.m },
          {
            withCredentials: true,
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
          })
        .then(
          result => {
            this.getAnnotationVector();
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
          }
        );
    },
    getMotifs: function() {
      axios
        .get(process.env.VUE_APP_MPSERVER_URL + "/topkmotifs", {
          withCredentials: true,
          params: {
            k: this.kmotifs,
            r: this.r
          }
        })
        .then(
          result => {
            this.motifs = result.data;

            var options = [];

            for (var i in this.motifs.series) {
              var motifGroup = this.motifs.series[i];
              if (motifGroup.length != 0) {
                options.push({
                  chartOptions: this.createMotifChartOption(
                    "motif " + i + ": " + this.motifs.groups[i].MinDist.toFixed(2),
                    motifGroup.slice(0, Math.min(10, motifGroup.length)),
                    this.motifs.groups[i].Idx.slice(0, Math.min(10, motifGroup.length))
                  )
                });
              } else {
                break;
              }
            }

            this.store.motifOptions = options;
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
          }
        );
    },
    getDiscords: function() {
      axios
        .get(process.env.VUE_APP_MPSERVER_URL + "/topkdiscords", {
          withCredentials: true,
          params: {
            k: this.kdiscords
          }
        })
        .then(
          result => {
            this.discords = result.data;

            var options = [];

            for (var i in this.discords.series) {
              options.push({
                chartOptions: this.createMotifChartOption(
                  "discord " + i,
                  [this.discords.series[i]],
                  [this.discords.groups[i]]
                )
              });
            }

            this.store.discordOptions = options;
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
          }
        );
    },
    getAnnotationVector: function() {
      this.avChange(this.selectedav);
    },
    avChange: function(av) {
      axios
        .post(process.env.VUE_APP_MPSERVER_URL + "/anvector",
          {name: av}, {
            withCredentials: true,
            headers: {"Content-Type": "application/x-www-form-urlencoded"}
          })
        .then(
          result => {
            var avoption = genAVOption(result.data.values);
            avoption.xAxis[0].max = this.n;
            this.store.annotationVectorOption = avoption;

            var mpoption = genMPOption(result.data.newmp);
            mpoption.xAxis[0].max = this.n;
            this.store.matrixProfileOption = mpoption;

            this.getMotifs();
            this.getDiscords();
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
          }
        );
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
              };

              for (var i in series) {
                var startIdx = series[i].name;
                self.store.tsOption.xAxis[0].plotBands.push({
                  from: startIdx,
                  to: startIdx + parseInt(self.getM(), 10),
                  color: "#FCFFC5"
                });

                points.data.push([
                  startIdx,
                  self.$refs.matrixprofile.$refs.highcharts.chart.series[0]
                    .data[startIdx].y
                ]);
              }

              if (self.store.matrixProfileOption.series.length > 1) {
                self.store.matrixProfileOption.series.pop();
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
              var startIdx = parseInt(e.point.series.userOptions.id, 10);
              self.store.tsOption.xAxis[0].plotBands.push({
                from: startIdx,
                to: startIdx + parseInt(self.getM(), 10),
                color: LightenDarkenColor(e.point.series.color, 20)
              });

              var point = {
                name: startIdx,
                showInLegend: false,
                type: "scatter",
                marker: {
                  symbol: "circle",
                  fillColor: "#FF0000",
                  lineWidth: 1
                },
                data: [
                  [
                    startIdx,
                    self.$refs.matrixprofile.$refs.highcharts.chart.series[0]
                      .data[startIdx].y
                  ]
                ]
              };
              if (self.store.matrixProfileOption.series.length > 1) {
                self.store.matrixProfileOption.series.pop();
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

function createChartOption(title, data, name, height) {
  var option = {
    chart: { height: height, zoomType: "x" },
    title: { text: title },
    plotOptions: {
      stickyTracking: false,
      series: {
        lineWidth: 1,
        animation: false
      }
    },
    tooltip: {
      shared: true
    },
    series: [{
      name: name,
      showInLegend: false,
      data: data
    }],
    xAxis: [
      {
        plotBands: []
      }
    ]
  };

  return option;
}

function genTSOption(data) {
  var option = createChartOption("Time Series", data, "value", 275);
  option.yAxis = [{
    title: {text: "value"}
  }];

  return option;
}

function genAVOption(data) {
  var option = createChartOption("Annotation Vector", data, "weight", 100);
  option.series[0].color = '#006600';
  option.yAxis = [{
    title: {text: "weight"},
    max: 1.0,
    min: 0.0
  }];
  option.title.style = {"fontSize": "12px"};

  return option;
}

function genMPOption(data) {
  var option = createChartOption("Matrix Profile", data, "distance", 275);
  option.series[0].color = '#000066';
  option.yAxis = [{
    title: {text: "distance"}
  }];

  return option;
}

function LightenDarkenColor(col, amt) {
  var usePound = false;
  if (col[0] == "#") {
    col = col.slice(1);
    usePound = true;
  }
  var num = parseInt(col, 16);
  var r = (num >> 16) + amt;
  if (r > 255) r = 255;
  else if (r < 0) r = 0;
  var b = ((num >> 8) & 0x00ff) + amt;
  if (b > 255) b = 255;
  else if (b < 0) b = 0;
  var g = (num & 0x0000ff) + amt;
  if (g > 255) g = 255;
  else if (g < 0) g = 0;
  return (usePound ? "#" : "") + (g | (b << 8) | (r << 16)).toString(16);
}
</script>

<style>
#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  height: 100%;
}
</style>
