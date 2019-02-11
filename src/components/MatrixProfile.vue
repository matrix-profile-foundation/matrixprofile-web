<template>
  <div class="hello">
    <b-container fluid>
      <b-row>
        <b-col cols="7"><TimeSeries /></b-col>
        <b-col cols="5">
          <b-navbar toggleable="md">
            <b-navbar-nav>
              <b-nav-item>Motifs</b-nav-item>
              <b-nav-item>Discords</b-nav-item>
              <b-nav-item>Segments</b-nav-item>
            </b-navbar-nav>
          </b-navbar>
          <Motifs />
        </b-col>
      </b-row>
    </b-container>
  </div>
</template>

<script>
import TimeSeries from "./TimeSeries.vue";
import Motifs from "./Motifs.vue";

export default {
  name: "MatrixProfile",
  components: {
    TimeSeries,
    Motifs
  }
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>

</style>
