import moment from "moment";

export default {
  data() {
    return {
      OFFICE_EXCEPT_INDEX: 1,
    };
  },
  methods: {
    date2YMD(date) {
      if (date) {
        return moment(date).format("YYYY-MM-DD");
      }
      return "";
    },
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
