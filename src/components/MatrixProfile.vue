<template>
  <div class="hello">
    <b-container fluid>
      <b-row>
        <b-col cols="7">
          <TimeSeries />
          <b-input-group prepend="m">
            <b-form-input type="number" placeholder="subsequence length">
            </b-form-input>
            <b-input-group-append><b-btn @click="calculateMP">Calculate</b-btn></b-input-group-append>
          </b-input-group>
        </b-col>
        <b-col cols="5">
          <b-nav tabs>
              <b-nav-item @click="enableMotifs">Motifs</b-nav-item>
              <b-nav-item @click="enableDiscords">Discords</b-nav-item>
              <b-nav-item @click="enableSegments">Segments</b-nav-item>
          </b-nav>
          <Motifs v-if="motifsActive" :store="store"/>
          <Discords v-if="discordsActive" :store="store"/>
          <Segments v-if="segmentsActive" :store="store"/>
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import TimeSeries from "./TimeSeries.vue";
import Motifs from "./Motifs.vue";
import Discords from "./Discords.vue";
import Segments from "./Segments.vue";

export default {
  name: "MatrixProfile",
  data() {
    return {
      motifsActive: true,
      discordsActive: false,
      segmentsActive: false,
      store: {
        message: "Hello!",
        motifOptions: [],
        discordOptions: [],
        segmentOptions: []
      }
    }
  },
  components: {
    TimeSeries,
    Motifs,
    Discords,
    Segments
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
      this.fillMotifOptions();
      this.fillDiscordOptions();
      this.fillSegmentOptions();
    },
    chartOption: function(title, data) {
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
        plotOptions: {
          series: {
            animation: false
          }
        },
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
      this.store.motifOptions = [
        {
          chartOptions: this.chartOption("motif 1", [
            [1, 2, 3],
            [1.1, 2.1, 3.1]
          ])
        },
        {
          chartOptions: this.chartOption("motif 2", [[1, 2, 1]])
        },
        {
          chartOptions: this.chartOption("motif 3", [[1, 2, 0]])
        }
      ];
    },
    fillDiscordOptions: function() {
      // likely makes an api call to find motifs
      this.store.discordOptions = [
        {
          chartOptions: this.chartOption("discord 1", [
            [3, 2, 1]
          ])
        }
      ];
    },
    fillSegmentOptions: function() {
      // likely makes an api call to find motifs
      this.store.segmentOptions = [
        {
          chartOptions: this.chartOption("segment 1", [
            [3, 2, 1]
          ])
        }
      ];
    }
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped></style>
