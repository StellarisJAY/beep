import { defineStore } from 'pinia';
import { getLoginInfo } from '@/services/base.js';

export const useUserStore = defineStore('beep-user-store', {
  state: () => ({
    userInfo: {},
    workspaceInfo: {}
  }),
  actions: {
    async queryLoginInfo() {
      console.log(this.userInfo)
      if (this.userInfo.id && this.workspaceInfo.id) {
        return true;
      }
      try {
        const res = await getLoginInfo();
        this.setUserInfo(res.data['user_info']);
        this.setWorkspaceInfo(res.data['workspace_info']);
        return true;
      } catch (error) {
        console.error(error);
        return false;
      }
    },
    setUserInfo(userInfo) {
      this.userInfo = userInfo;
    },
    setWorkspaceInfo(workspaceInfo) {
      this.workspaceInfo = workspaceInfo;
    },
    getLoginUserInfo() {
      return this.userInfo;
    },
    getWorkspaceInfo() {
      return this.workspaceInfo;
    },
    clearLoginInfo() {
      this.userInfo = {
        id: "",
        name: "",
        email: "",
      };
      this.workspaceInfo = {
        id: "",
        name: "",
        description: ""
      };
      localStorage.clear();
    }
  }
});
