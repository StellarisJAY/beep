import { defineStore } from 'pinia';

const useChatStore =defineStore('beep-chat-store', {
  state: () => ({
    messages: [],
    conversationId: "",
    lastMessageId: "",
    agentId: "",
  }),
  actions: {
    async sendMessage(message) {

    },
    clearSession() {
      this.messages = [];
      this.conversationId = "";
      this.lastMessageId = "";
      this.agentId = "";
    }
  }
})
