import Vue from "vue";
import App from "./App.vue";
import HighchartsVue from "highcharts-vue";
import BootstrapVue from "bootstrap-vue";
import VueCookies from "vue-cookies";
import VeeValidate from "vee-validate";
import "bootstrap/dist/css/bootstrap.min.css";
import "bootstrap-vue/dist/bootstrap-vue.min.css";

Vue.config.productionTip = false;

Vue.use(HighchartsVue);
Vue.use(BootstrapVue);
Vue.use(VueCookies);
Vue.use(VeeValidate);

new Vue({
  render: h => h(App)
}).$mount("#app");
