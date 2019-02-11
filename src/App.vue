<template>
  <div id="app">
    <MatrixProfile />
  </div>
</template>

<script>
import MatrixProfile from "./components/MatrixProfile.vue";

export default {
  name: "app",
  components: {
    MatrixProfile
  }
};
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
