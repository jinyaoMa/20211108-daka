import Vue from "vue";
import VueRouter from "vue-router";
import Home from "../views/Home.vue";
import Store from "../views/Store.vue";
import Login from "../views/Login.vue";
import AddUser from "../views/AddUser.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
    props: (route) => ({ showDate: route.query.showDate }),
  },
  {
    path: "/store",
    name: "Store",
    component: Store,
  },
  {
    path: "/login",
    name: "Login",
    component: Login,
  },
  {
    path: "/adduser",
    name: "AddUser",
    component: AddUser,
  },
];

const router = new VueRouter({
  routes,
});

export default router;
