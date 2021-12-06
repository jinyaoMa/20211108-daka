<template>
  <el-container class="store">
    <el-header>
      <div class="title">Twtrade Timesheet Tool</div>
      <div class="right">
        <div class="username">{{ username }}</div>
        <el-button type="danger" @click="logout">Logout</el-button>
      </div>
    </el-header>
    <el-main>
      <el-form class="form" ref="form">
        <el-form-item label="Store">
          <el-select v-model="storeId" placeholder="Choose a store">
            <el-option
              v-for="(store, i) in storeslist"
              :key="i"
              :label="store"
              :value="i"
            >
            </el-option>
          </el-select>
          <el-button class="margin-l-20" type="success" @click="viewAll">
            View All
          </el-button>
        </el-form-item>
      </el-form>
      <div class="datePicker">
        <div class="custom-label">Pick a date to view</div>
        <el-calendar v-model="datePicked" :first-day-of-week="7">
          <template slot="dateCell" slot-scope="{ date, data }">
            <div
              :class="data.isSelected ? 'dateCell is-selected' : 'dateCell'"
              @click="showDate(date)"
            >
              {{ data.day.split("-")[2] }}
            </div>
          </template>
        </el-calendar>
      </div>
      <div class="export">
        <el-form label-width="100px">
          <el-form-item label="Export from">
            <el-date-picker
              v-model="exportStart"
              type="date"
              placeholder="Start Date"
            >
            </el-date-picker>
          </el-form-item>
          <el-form-item label="To">
            <el-date-picker
              v-model="exportEnd"
              type="date"
              placeholder="End Date"
            >
            </el-date-picker>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" @click="downloadExcel">
              Export .xlsx
            </el-button>
          </el-form-item>
          <a
            ref="download"
            style="display: none"
            :download="downloadFilename"
          ></a>
        </el-form>
        <div v-if="canAddUser">
          <el-button class="margin-l-20" type="success" @click="handleAddUser">
            Add User
          </el-button>
        </div>
      </div>
    </el-main>
  </el-container>
</template>

<script>
export default {
  name: "Store",
  data() {
    return {
      storeId: "",
      datePicked: new Date(),
      exportStart: new Date(),
      exportEnd: new Date(),
      storeslist: [],
      username: "",
      usertype: "",
    };
  },
  mounted() {
    const userInfo = this.getUserInfo();
    if (typeof userInfo.token !== "string" || userInfo.token == "") {
      this.$router.replace("/login");
      return;
    }

    this.username = userInfo.user.name;
    this.usertype = userInfo.user.type;
    if (userInfo.store) {
      this.storeId = userInfo.store.id;
    }

    this.$http
      .get("/stores/list", {
        headers: {
          Authorization: userInfo.token,
        },
      })
      .then((result) => {
        const list = result.data.data;
        if (userInfo.user.type == "Office") {
          this.storeslist = list
            .slice(0, this.OFFICE_EXCEPT_INDEX)
            .concat(list.slice(this.OFFICE_EXCEPT_INDEX + 1));
        } else {
          this.storeslist = list;
        }
        userInfo.store.name = this.storeslist[parseInt(this.storeId)];
        this.saveUserInfo(userInfo);
      });
  },
  watch: {
    datePicked: {
      handler(newValue) {
        newValue = new Date(newValue.setDate(newValue.getDate() + 1));
      },
    },
    storeId: {
      handler(newValue) {
        const userInfo = this.getUserInfo();
        userInfo.store = {
          id: newValue,
          name: this.storeslist[parseInt(newValue)],
        };
        this.saveUserInfo(userInfo);
      },
    },
  },
  computed: {
    canAddUser() {
      return (
        this.storeId == 0 &&
        (this.usertype == "Admin" || this.usertype == "Office")
      );
    },
    downloadFilename() {
      return (
        this.storeslist[parseInt(this.storeId)] +
        "_" +
        this.date2YMD(this.exportStart) +
        "_" +
        this.date2YMD(this.exportEnd) +
        ".xlsx"
      );
    },
  },
  methods: {
    logout() {
      this.saveUserInfo({});
      this.$router.push("/login");
    },
    viewAll() {
      if (typeof this.storeId === "string" && this.storeId === "") {
        this.$message.error("Please choose a store.");
        return;
      }

      this.$router.push({
        path: "/",
        query: {
          showDate: "All",
        },
      });
    },
    showDate(date) {
      if (typeof this.storeId === "string" && this.storeId === "") {
        this.$message.error("Please choose a store.");
        return;
      }

      date = new Date(date.setDate(date.getDate() + 1));
      this.$router.push({
        path: "/",
        query: {
          showDate: this.date2YMD(date),
        },
      });
    },
    downloadExcel() {
      if (typeof this.storeId === "string" && this.storeId === "") {
        this.$message.error("Please choose a store.");
        return;
      }
      const userInfo = this.getUserInfo();
      this.$http
        .get("/timesheet/download", {
          params: {
            storeId: this.storeId,
            startDate: this.date2YMD(this.exportStart),
            endDate: this.date2YMD(this.exportEnd),
          },
          responseType: "blob",
          headers: {
            Authorization: userInfo.token,
          },
        })
        .then((result) => {
          if (result.data) {
            let a = this.$refs.download;
            let url = window.URL.createObjectURL(new Blob([result.data]));
            a.href = url;
            a.click();
          }
        });
    },
    handleAddUser() {
      this.$router.push("/adduser");
    },
  },
};
</script>

<style lang="scss" scoped>
.store {
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
    max-width: 1024px;
    margin: 0 auto;
    .export {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      @media (max-width: 512px) {
        flex-direction: column;
        gap: 20px;
      }
    }
  }
}
.margin-l-20 {
  margin-left: 20px;
}
.custom-label {
  vertical-align: middle;
  font-size: 14px;
  color: #606266;
  line-height: 40px;
  box-sizing: border-box;
}
.is-selected {
  color: #1989fa;
}
.dateCell {
  padding: 8px;
  height: 100%;
  box-sizing: border-box;
}
</style>