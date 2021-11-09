import moment from "moment";

export default {
  methods: {
    data2String(date) {
      if (date) {
        return moment(date).format("YYYY-MM-DD HH:mm:ss");
      }
      return "N/A";
    },
    dataCompare(a, b) {
      return moment(a).isBefore(b) ? 1 : -1;
    },
    getUserInfo() {
      let userinfo = window.localStorage.getItem("userInfo");
      if (userinfo) {
        try {
          return JSON.parse(userinfo);
        } catch (error) {
          return {};
        }
      }
      return {};
    },
    saveUserInfo(userInfo) {
      window.localStorage.setItem("userInfo", JSON.stringify(userInfo));
    },
  },
};
