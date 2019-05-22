<template>
  <div id="app">
    <b-modal ref="modal-mp-error" title="Matrix Profile Error" ok-only>
      <p class="my-4">{{ err }}</p>
    </b-modal>
    <b-navbar toggleable="lg" type="dark" variant="info">
      <b-navbar-brand>Matrix Profiles</b-navbar-brand>

      <b-navbar-toggle target="nav-collapse"></b-navbar-toggle>

      <b-collapse id="nav-collapse" is-nav>
        <!-- Right aligned nav items -->
        <b-navbar-nav class="ml-auto">
          <b-nav-item
            active
            target="_blank"
            href="https://github.com/aouyang1/go-matrixprofile"
          >
            <Octicon :scale="2" :icon="markGithub" />
          </b-nav-item>
          <b-nav-item
            active
            target="_blank"
            href="https://www.cs.ucr.edu/~eamonn/MatrixProfile.html"
            class="mt-1"
            >UCR Webpage
          </b-nav-item>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <b-container fluid>
      <b-row>
        <b-col cols="8">
          <b-form-select
            class="mt-2"
            v-model="selectedSource"
            :options="sourceOptions"
            @change="getTimeSeries"
          >
          </b-form-select>
          <TimeSeries class="mt-2" ref="timeseries" :store="store" />
          <b-tabs small>
            <b-tab title="Matrix Profile">
              <AnnotationVector ref="annotationvector" :store="store" />
              <MatrixProfile ref="matrixprofile" :store="store" />
            </b-tab>
            <b-tab title="Segmentation">
              <Segmentation ref="segmentation" :store="store" />
            </b-tab>
          </b-tabs>
        </b-col>
        <b-col cols="4">
          <b-input-group prepend="m" class="mb-2 mt-2">
            <b-form-input
              v-model.number="m"
              type="number"
              v-validate="'required|integer|min_value:4'"
              name="input-m"
              placeholder="subsequence length"
            >
            </b-form-input>
            <b-input-group-append>
              <b-btn
                @click="calculateMP"
                :disabled="
                  calculatingMP ||
                    (fields['input-m'] && !fields['input-m'].valid) ||
                    (fields['input-kmotifs'] &&
                      !fields['input-kmotifs'].valid) ||
                    (fields['input-rmotifs'] &&
                      !fields['input-rmotifs'].valid) ||
                    (fields['input-kdiscords'] &&
                      !fields['input-kdiscords'].valid)
                "
              >
                <template v-if="calculatingMP">
                  <b-spinner
                    small
                    variant="light"
                    class="mb-1 mr-1"
                  ></b-spinner>
                  Calculating...
                </template>
                <template v-else>
                  Calculate
                </template>
              </b-btn>
            </b-input-group-append>
          </b-input-group>
          <b-tabs>
            <b-tab title="Motifs">
              <template v-if="fields['input-m'] && fields['input-m'].valid">
                <b-form inline>
                  <b-input-group
                    size="sm"
                    class="mt-1 mb-1 mr-1"
                    prepend="top-k"
                  >
                    <b-form-input
                      v-model="kmotifs"
                      type="number"
                      v-validate="'required|integer|min_value:1'"
                      name="input-kmotifs"
                      placeholder="max number of motifs"
                    >
                    </b-form-input>
                  </b-input-group>
                  <b-input-group
                    size="sm"
                    class="mt-1 mb-1 mr-1"
                    prepend="radius"
                  >
                    <b-form-input
                      v-model="r"
                      type="number"
                      v-validate="'required|decimal|min_value:0'"
                      name="input-rmotifs"
                      placeholder="radius"
                    >
                    </b-form-input>
                  </b-input-group>
                  <b-btn
                    size="sm"
                    @click="getMotifs"
                    :disabled="
                      calculatingMotifs ||
                        (fields['input-kmotifs'] &&
                          !fields['input-kmotifs'].valid) ||
                        (fields['input-rmotifs'] &&
                          !fields['input-rmotifs'].valid) ||
                        (fields['input-kdiscords'] &&
                          !fields['input-kdiscords'].valid)
                    "
                  >
                    <template v-if="calculatingMotifs">
                      <b-spinner small variant="light" class="mr-1"></b-spinner>
                      Finding...
                    </template>
                    <template v-else>
                      Find
                    </template>
                  </b-btn>
                </b-form>

                <Motifs ref="motifs" :store="store" />
              </template>
            </b-tab>
            <b-tab title="Discords">
              <template v-if="fields['input-m'] && fields['input-m'].valid">
                <b-form inline>
                  <b-input-group
                    size="sm"
                    class="mt-1 mb-1 mr-1"
                    prepend="top-k"
                  >
                    <b-form-input
                      v-model="kdiscords"
                      type="number"
                      v-validate="'required|integer|min_value:1'"
                      name="input-kdiscords"
                      placeholder="max number of discords"
                    >
                    </b-form-input>
                  </b-input-group>
                  <b-btn
                    size="sm"
                    @click="getDiscords"
                    :disabled="
                      calculatingDiscords ||
                        (fields['input-kmotifs'] &&
                          !fields['input-kmotifs'].valid) ||
                        (fields['input-rmotifs'] &&
                          !fields['input-rmotifs'].valid) ||
                        (fields['input-kdiscords'] &&
                          !fields['input-kdiscords'].valid)
                    "
                  >
                    <template v-if="calculatingDiscords">
                      <b-spinner small variant="light" class="mr-1"></b-spinner>
                      Finding...
                    </template>
                    <template v-else>
                      Find
                    </template>
                  </b-btn>
                </b-form>

                <Discords :store="store" />
              </template>
            </b-tab>
            <b-tab title="Annotation Vectors">
              <template v-if="fields['input-m'] && fields['input-m'].valid">
                <b-form inline>
                  <b-form-select
                    size="sm"
                    class="mt-1 mb-1"
                    v-model="selectedav"
                    :options="avoptions"
                    @change="avChange"
                    :disabled="
                      (fields['input-kmotifs'] &&
                        !fields['input-kmotifs'].valid) ||
                        (fields['input-rmotifs'] &&
                          !fields['input-rmotifs'].valid) ||
                        (fields['input-kdiscords'] &&
                          !fields['input-kdiscords'].valid)
                    "
                  >
                  </b-form-select>
                  <p class="text-left">
                    {{ this.avdescription[this.selectedav] }}
                  </p>
                </b-form>
              </template>
            </b-tab>
          </b-tabs>
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
import Segmentation from "./components/Segmentation.vue";
import Octicon, { markGithub } from "octicons-vue";
import { mapFields } from "vee-validate";
import axios from "axios";

export default {
  name: "app",
  data() {
    return {
      markGithub: markGithub,
      selectedSource: "",
      sourceOptions: [],
      ts: [],
      n: 0,
      m: 30,
      kmotifs: 3,
      r: 2,
      motifs: [],
      kdiscords: 3,
      discords: [],
      cac: [],
      calculatingMP: false,
      calculatingMotifs: false,
      calculatingDiscords: false,
      selectedav: "default",
      avoptions: [
        { value: "default", text: "Default" },
        { value: "complexity", text: "Complexity" },
        { value: "meanstd", text: "Mean Standard Deviation" },
        { value: "clipping", text: "Clipping" }
      ],
      avdescription: {
        default:
          "This is the default annotation vector which does not modify the matrix profile whatsover.",
        complexity:
          "This biases the matrix profile towards areas of high signal complexity or variation in signal. Generally this is good for focusing motif search around areas with high signal variance. Found discords may not be useful since areas with low complexity have their matrix profile distance artificially increased.",
        meanstd:
          "This attempts to suppress regions of high fluctuation by targetting areas where the subsequence standard deviation is larger than the mean.",
        clipping:
          "This attemps to suppress regions where there are clipping effects while capture data and raises their matrix profile distances as to not find them as prominent motifs. Subsequences in the signal that have more values at the max or min of the entire signal are weighted more and their matrix profile distances are increased."
      },
      err: "",
      retries: 0,
      store: {
        tsOption: genTSOption([]),
        annotationVectorOption: genAVOption([]),
        matrixProfileOption: genMPOption([]),
        segmentationOption: genSegOption([]),
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
    Segmentation,
    Octicon
  },
  created: function() {
    this.getSources();
  },
  computed: {
    ...mapFields([
      "input-m",
      "input-kmotifs",
      "input-rmotifs",
      "input-kdiscords"
    ])
  },
  methods: {
    checkError: function() {
      if (this.err != "") {
        this.$refs["modal-mp-error"].show();
      }
    },
    getSources: function() {
      var endpoint = process.env.VUE_APP_MPSERVER_URL + "/sources";

      axios
        .get(endpoint, {
          withCredentials: true
        })
        .then(
          result => {
            var sources = result.data;
            if (sources.length > 0) {
              this.selectedSource = sources[0];
            }
            this.sourceOptions = [];
            for (var i = 0; i < sources.length; i++) {
              this.sourceOptions.push({ value: sources[i], text: sources[i] });
            }
            this.err = "";
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
            this.checkError(this.err);
          }
        );
    },
    getTimeSeries: function(source) {
      var endpoint = process.env.VUE_APP_MPSERVER_URL + "/data";

      axios
        .get(endpoint, {
          withCredentials: true,
          params: { source: source }
        })
        .then(
          result => {
            this.ts = result.data;
            this.n = result.data.length;
            var option = genTSOption(result.data);
            option.xAxis[0].max = this.n;
            this.store.tsOption = option;

            this.err = "";
          },
          error => {
            this.err = JSON.stringify(error.response.data.error);
            this.checkError(this.err);
          }
        );
    },
    calculateMP: function() {
      var endpoint = process.env.VUE_APP_MPSERVER_URL + "/calculate";
      this.calculatingMP = true;

      axios
        .post(
          endpoint,
          { m: this.m },
          {
            withCredentials: true,
            headers: { "Content-Type": "application/x-www-form-urlencoded" }
          }
        )
        .then(
          result => {
            this.cac = result.data.cac;
            var option = genSegOption(this.cac);
            option.xAxis[0].max = this.n;
            this.store.segmentationOption = option;

            this.getAnnotationVector();

            this.err = "";
          },
          error => {
            this.calculatingMP = false;
            this.err = JSON.stringify(error.response.data.error);
            this.checkError(this.err);
          }
        );
    },
    retryCalculateMP: function(errResp) {
      var maxRetries = 3;
      if (errResp.cache_expired && this.retries < maxRetries) {
        this.retries++;
        console.log(
          "Retrying to cache matrix profile with " +
            (maxRetries - this.retries).toString() +
            " retries remaining"
        );
        this.calculateMP();
      } else {
        this.err = JSON.stringify(errResp.error);
        this.checkError(this.err);
      }
    },
    getMotifs: function() {
      var endpoint = process.env.VUE_APP_MPSERVER_URL + "/topkmotifs";
      this.calculatingMotifs = true;

      axios
        .get(endpoint, {
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
                    "motif ".concat(
                      i,
                      ": ",
                      this.motifs.groups[i].MinDist.toFixed(2)
                    ),
                    motifGroup,
                    this.motifs.groups[i].Idx
                  )
                });
              } else {
                break;
              }
            }

            this.store.motifOptions = options;
            this.calculatingMotifs = false;
            this.err = "";
            this.retries = 0;
          },
          error => {
            this.calculatingMotifs = false;
            this.retryCalculateMP(error.response.data);
          }
        );
    },
    getDiscords: function() {
      var endpoint = process.env.VUE_APP_MPSERVER_URL + "/topkdiscords";
      this.calculatingDiscords = true;

      axios
        .get(endpoint, {
          withCredentials: true,
          params: { k: this.kdiscords }
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
            this.calculatingDiscords = false;
            this.err = "";
            this.retries = 0;
          },
          error => {
            this.calculatingDiscords = false;
            this.retryCalculateMP(error.response.data);
          }
        );
    },
    getAnnotationVector: function() {
      this.avChange(this.selectedav);
    },
    avChange: function(av) {
      var endpoint = process.env.VUE_APP_MPSERVER_URL + "/mp";

      axios
        .post(
          endpoint,
          { name: av },
          {
            withCredentials: true,
            headers: { "Content-Type": "application/x-www-form-urlencoded" }
          }
        )
        .then(
          result => {
            var avoption = genAVOption(result.data.annotation_vector);
            avoption.xAxis[0].max = this.n;
            this.store.annotationVectorOption = avoption;

            var mpoption = genMPOption(result.data.adjusted_mp);
            mpoption.xAxis[0].max = this.n;
            this.store.matrixProfileOption = mpoption;

            this.getMotifs();
            this.getDiscords();

            this.calculatingMP = false;
            this.err = "";
          },
          error => {
            this.retryCalculateMP(error.response.data);
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
          animation: false,
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
                  to: startIdx + parseInt(self.m, 10),
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
                to: startIdx + parseInt(self.m, 10),
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
    chart: { height: height, zoomType: "x", animation: false },
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
    series: [
      {
        name: name,
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

function genTSOption(data) {
  var option = createChartOption("Time Series", data, "value", 275);
  option.yAxis = [
    {
      title: { text: "value" }
    }
  ];

  return option;
}

function genAVOption(data) {
  var option = createChartOption("Annotation Vector", data, "weight", 100);
  option.series[0].color = "#006600";
  option.yAxis = [
    {
      title: { text: "weight" },
      max: 1.0,
      min: 0.0
    }
  ];
  option.title.style = { fontSize: "12px" };

  return option;
}

function genMPOption(data) {
  var option = createChartOption("Matrix Profile", data, "distance", 250);
  option.series[0].color = "#000066";
  option.yAxis = [
    {
      title: { text: "distance" }
    }
  ];
  option.title.style = { fontSize: "12px" };

  return option;
}

function genSegOption(data) {
  var option = createChartOption("Segmentation", data, "cac", 350);
  option.yAxis = [
    {
      title: { text: "corrected arc crossings" }
    }
  ];

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
