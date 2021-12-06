<template>
  <el-container class="home">
    <el-header>
      <div class="title">Twtrade Timesheet Tool</div>
      <div class="right">
        <div class="username">{{ username }}</div>
        <div class="storename">(in {{ storename }}, {{ showDate }})</div>
        <el-button-group>
          <el-button type="warning" @click="changeStore">
            Change Store
          </el-button>
          <el-button type="danger" @click="logout">Logout</el-button>
        </el-button-group>
      </div>
    </el-header>
    <el-main>
      <el-table :data="tableData" @sort-change="handleSortChange">
        <el-table-column
          prop="Username"
          label="Username"
          min-width="120"
          sortable="custom"
        >
        </el-table-column>
        <el-table-column
          prop="SigninTime"
          label="Sign In Time"
          min-width="200"
          sortable="custom"
        >
          <template slot-scope="scope">
            <span v-if="!scope.row.isEdit">
              {{ data2String(scope.row.SigninTime) }}
            </span>
            <el-date-picker
              v-if="scope.row.isEdit"
              v-model="tableDataChange[scope.row.key].SigninTime"
              type="datetime"
              size="mini"
            >
            </el-date-picker>
          </template>
        </el-table-column>
        <el-table-column
          prop="SignoutTime"
          label="Sign Out Time"
          min-width="200"
          sortable="custom"
        >
          <template slot-scope="scope">
            <span v-if="!scope.row.isEdit">
              {{ data2String(scope.row.SignoutTime) }}
            </span>
            <el-date-picker
              v-if="scope.row.isEdit"
              v-model="tableDataChange[scope.row.key].SignoutTime"
              type="datetime"
              size="mini"
            >
            </el-date-picker>
          </template>
        </el-table-column>
        <el-table-column
          prop="Total"
          label="Daily Total"
          :min-width="usertype == 'Account' ? 120 : 160"
          :fixed="usertype == 'Account' ? 'right' : false"
          sortable="custom"
        >
          <template slot-scope="scope">
            <span v-if="!scope.row.isEdit">
              {{ scope.row.Total }}
            </span>
            <el-input-number
              v-if="scope.row.isEdit"
              v-model="tableDataChange[scope.row.key].Total"
              :precision="2"
              :step="0.1"
              size="mini"
            ></el-input-number>
          </template>
        </el-table-column>
        <el-table-column
          v-if="usertype != 'Account'"
          min-width="180"
          label="Operation"
          fixed="right"
        >
          <template slot-scope="scope">
            <el-button
              v-if="!scope.row.isEdit"
              @click="handleEdit(scope.row)"
              type="primary"
              size="small"
            >
              Edit
            </el-button>
            <el-button
              v-if="scope.row.isEdit"
              @click="handleUpdate(tableDataChange[scope.row.key])"
              type="success"
              size="small"
            >
              Update
            </el-button>
            <el-button
              v-if="scope.row.isEdit"
              @click="handleCancel(scope.row)"
              size="small"
            >
              Cancel
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        @current-change="handleCurrentChange"
        :page-size="limit"
        layout="prev, pager, next, jumper"
        :total="count"
      >
      </el-pagination>
      <div v-if="count > 0" class="export">
        <el-button type="success" @click="downloadExcel">
          Export .xlsx
        </el-button>
        <a
          ref="download"
          style="display: none"
          :download="downloadFilename"
        ></a>
      </div>
    </el-main>
  </el-container>
</template>

<script>
// @ is an alias to /src

export default {
  name: "Home",
  props: {
    showDate: {
      type: String,
    },
  },
  data() {
    return {
      username: "",
      storename: "",
      usertype: "",
      offset: 0,
      limit: 10,
      currentPage: 1,
      count: 0,
      tableData: [],
      tableDataChange: [],
      downloadFilename: "",
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
    this.usertype = userInfo.user.type;

    if (this.showDate == "All") {
      this.downloadFilename = userInfo.store.name + "_all.xlsx";
    } else {
      this.downloadFilename =
        userInfo.store.name + "_" + this.showDate + ".xlsx";
    }

    this.loadList();
  },
  methods: {
    loadList(field = "SigninTime", order = "asc") {
      const userInfo = this.getUserInfo();
      this.$http
        .get("/timesheet/list", {
          params: {
            storeId: userInfo.store.id,
            showDate: this.showDate,
            offset: this.offset,
            limit: this.limit,
            order: order.startsWith("asc") ? "asc" : "desc",
            orderby: field,
          },
          headers: {
            Authorization: userInfo.token,
          },
        })
        .then((result) => {
          result.data.data.forEach((item, key) => {
            item.key = key;
            item.isEdit = false;
          });
          this.tableData = result.data.data;
          this.tableDataChange = result.data.data;
          this.count = result.data.count;
        });
    },
    logout() {
      this.saveUserInfo({});
      this.$router.push("/login");
    },
    changeStore() {
      this.$router.push("/store");
    },
    handleEdit(row) {
      row.isEdit = true;
    },
    handleCancel(row) {
      row.isEdit = false;
    },
    handleSortChange(sort) {
      this.loadList(sort.prop, sort.order || "asc");
    },
    handleCurrentChange(page) {
      this.offset = (page - 1) * this.limit;
      this.loadList();
    },
    handleUpdate(row) {
      const userInfo = this.getUserInfo();
      this.$http
        .post(
          "/timesheet/update",
          { ...row, StoreID: userInfo.store.id },
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
              message: "Updated!",
            });
            this.tableData = this.tableDataChange;
            this.tableData[row.key].isEdit = false;
          }
        })
        .catch((error) => {
          let data = error.response.data;
          this.$message.error(data.error);
        });
    },
    downloadExcel() {
      const userInfo = this.getUserInfo();
      this.$http
        .get("/timesheet/download", {
          params: {
            storeId: userInfo.store.id,
            startDate: this.showDate,
            endDate: this.showDate,
            all: this.showDate == "All",
            r: Math.random(),
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
  },
};
</script>

<style lang="scss" scoped>
.home {
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
    @media (max-width: 800px) {
      flex-direction: column;
      height: auto !important;
    }
    @media (max-width: 520px) {
      .right {
        flex-direction: column;
      }
    }
  }
  main {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    .el-date-editor.el-input,
    .el-date-editor.el-input__inner {
      width: 180px;
    }
    .el-pagination {
      margin: 10px 0;
    }
    .export {
      border-top: 2px dashed;
      padding: 10px 0;
      > * {
        margin-left: 10px;
      }
      .label {
        display: inline-block;
        line-height: 40px;
      }
      .el-input-number {
        width: 120px;
      }
    }
  }
}
</style>