import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import "./plugins/element.js";
import axios from "axios";
import VueAxios from "vue-axios";
import mixin from "./mixin.js";

Vue.config.productionTip = false;

axios.defaults.baseURL = "http://localhost:8081";
Vue.use(VueAxios, axios);

Vue.mixin(mixin);

new Vue({
  router,
  render: (h) => h(App),
}).$mount("#app");
