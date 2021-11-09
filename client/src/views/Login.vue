<template>
  <div class="login">
    <el-form
      class="form"
      ref="form"
      :model="form"
      label-width="80px"
      label-position="left"
    >
      <el-form-item label="Username">
        <el-input v-model="form.account"></el-input>
      </el-form-item>
      <el-form-item label="Password">
        <el-input v-model="form.password"></el-input>
      </el-form-item>
      <el-form-item label="Store">
        <el-select v-model="form.storeId" placeholder="Choose a store">
          <el-option
            v-for="(store, i) in storeslist"
            :key="i"
            :label="store"
            :value="i"
          >
          </el-option>
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="login">Login</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script>
export default {
  name: "Login",
  data() {
    return {
      form: {
        account: "",
        password: "",
        storeId: "",
      },
      storeslist: [],
    };
  },
  mounted() {
    this.$http.get("/stores/list").then((result) => {
      this.storeslist = result.data.data;
    });
  },
  methods: {
    login() {
      this.$http
        .post("/auth/login", this.form)
        .then((result) => {
          let data = result.data;
          this.saveUserInfo({
            token: "Bearer " + data.token,
            store: data.data.store,
            user: data.data.user,
          });
          this.$router.push("/");
        })
        .catch((error) => {
          let data = error.response.data;
          this.$message.error(data.error);
        });
    },
  },
};
</script>

<style lang="scss" scoped>
.login {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
}
.form {
  box-sizing: border-box;
  min-width: 320px;
  padding: 10px;
  .el-select {
    width: 100%;
  }
}
</style>