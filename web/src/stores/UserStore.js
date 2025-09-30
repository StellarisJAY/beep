import { defineStore } from 'pinia';

export const useUserStore = defineStore('beep-user-store', {
  state: () => ({
    user: null,
  }),
  actions: {

  }
});
