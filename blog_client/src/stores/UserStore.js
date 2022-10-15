import { defineStore } from "pinia"; // 引入pinia

export const UserStore = defineStore("admin", {
  state: () => {
    return {
      token: "",
    };
  },
  actions: {},
  getters: {},
});