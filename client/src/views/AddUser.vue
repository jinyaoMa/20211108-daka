<template>
  <el-container class="addUser">
    <el-header>
      <div class="title">Twtrade Timesheet Tool</div>
      <div class="right">
        <div class="username">{{ username }}</div>
        <div class="storename">(in {{ storename }})</div>
        <el-button-group>
          <el-button type="primary" @click="changeStore">Back</el-button>
          <el-button type="danger" @click="logout">Logout</el-button>
        </el-button-group>
      </div>
    </el-header>
    <el-main>
      <el-form class="form" ref="form" :model="form" label-position="top">
        <el-form-item label="First Name">
          <el-input v-model="form.firstname"></el-input>
        </el-form-item>
        <el-form-item label="Last Name">
          <el-input v-model="form.lastname"></el-input>
        </el-form-item>
        <el-form-item label="Password">
          <el-input v-model="form.password"></el-input>
        </el-form-item>
        <el-form-item label="Usertype">
          <el-select
            v-model="form.usertype"
            placeholder="Select a usertype"
            disabled
          >
            <el-option label="Staff" :value="3"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="add">Add</el-button>
          <el-button @click="clearForm">Clear</el-button>
        </el-form-item>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: "AddUser",
  data() {
    return {
      username: "",
      storename: "",
      form: {
        firstname: "",
        lastname: "",
        usertype: 3,
        password: "",
      },
    };
  },
  mounted() {
    const userInfo = this.getUserInfo();
    if (typeof userInfo.token !== "string" || userInfo.token == "") {
      this.$router.replace("/login");
      return;
    }

    this.username = userInfo.user.name;
    this.storename = userInfo.store.name;
  },
  methods: {
    logout() {
      this.saveUserInfo({});
      this.$router.push("/login");
    },
    changeStore() {
      this.$router.push("/store");
    },
    clearForm() {
      this.form = {
        firstname: "",
        lastname: "",
        usertype: 3,
        password: "",
      };
    },
    add() {
      const userInfo = this.getUserInfo();
      this.$http
        .post(
          "/user/add",
          { ...this.form, StoreID: userInfo.store.id },
          {
            headers: {
              Authorization: userInfo.token,
            },
          }
        )
        .then((result) => {
          if (result.data.ok) {
            this.$message({
              type: "success",
              message:
                "User (" +
                this.form.firstname +
                " " +
                this.form.lastname +
                ") created!",
            });
            this.clearForm();
          }
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
.addUser {
  header {
    border-bottom: 2px solid;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: space-between;
    .storename,
    .username,
    .title {
      margin-top: 0.2rem;
    }
    .storename {
      margin-right: 1rem;
    }
    .username {
      margin-right: 0.5rem;
    }
    .title {
      font-size: 1.5em;
      text-align: center;
    }
    .right {
      display: flex;
      flex-direction: row;
      align-items: center;
    }
    @media (max-width: 512px) {
      flex-direction: column;
      height: auto !important;
    }
  }
  main {
    margin: 20px auto;
    min-width: 320px;
    box-sizing: border-box;
  }
}
</style>