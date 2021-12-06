<template>
  <div id="app">
    <router-view></router-view>
  </div>
</template>

<script>
export default {
  name: "app",
};
</script>

<style lang="scss">
body {
  margin: 0;
}

#app {
  font-family: "Avenir", Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

.el-message {
  min-width: 320px !important;
  max-width: 100vw;
}

.el-calendar-table .el-calendar-day {
  padding: 0 !important;
}

@media (max-width: 512px) {
  .el-calendar__header,
  .el-calendar__body {
    padding-left: 0 !important;
    padding-right: 0 !important;
  }
}
</style>
