<template>
  <div class="login">
    <div class="title">Twtrade Timesheet Tool</div>
    <el-form
      class="form"
      ref="form"
      :model="form"
      label-width="80px"
      label-position="left"
    >
      <el-form-item label="Username">
        <el-input
          v-model="form.account"
          @keyup.enter.native="login()"
        ></el-input>
      </el-form-item>
      <el-form-item label="Password">
        <el-input
          type="password"
          v-model="form.password"
          @keyup.enter.native="login()"
        ></el-input>
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
      },
    };
  },
  methods: {
    login() {
      this.$http
        .post("/auth/login", this.form)
        .then((result) => {
          let data = result.data;
          this.saveUserInfo({
            token: "Bearer " + data.token,
            user: data.data.user,
          });
          this.$router.push("/store");
        })
        .catch((error) => {
          let data = error.response.data;
          this.$message.error(data.error);
        });
      return false;
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
  .title {
    line-height: 3;
    font-size: 1.8em;
    font-weight: bold;
  }
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