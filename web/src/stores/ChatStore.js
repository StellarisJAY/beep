import { defineStore } from 'pinia';
import httpReq from '@/services/http.js';
import { getConversationList, getConversationMessages } from '@/services/conversation.js';

export const useChatStore =defineStore('beep-chat-store', {
  state: () => ({
    messages: [],
    conversationId: "",
    agentId: "",
    isLoading: false,
    conversationList: [],
    conversationQuery: {
      title: '',
    }
  }),
  actions: {
    async sendMessage(message) {
      try {
        const data = {
          query: message,
          agent_id: this.agentId,
          conversation_id: this.conversationId,
        };
        this.isLoading = true;
        await httpReq.postEventStream('/chat/send', data, this.onDataReceived);
      } catch (error) {
        console.error(error);
      } finally {
        this.isLoading = false;
      }
    },

    async signalToolCall(accept) {
      try {
        const data = {
          accept: accept,
          conversation_id: this.conversationId,
          agent_id: this.agentId,
        };
        this.isLoading = true;
        await httpReq.postEventStream('/chat/signal', data, this.onDataReceived);
      } catch (error) {
        console.error(error);
      } finally {
        this.isLoading = false;
      }
    },

    async listConversations() {
      try {
        const res = await getConversationList(this.conversationQuery);
        this.conversationList = res.data || [];
      } catch (error) {
        console.error(error);
      }
    },

    onDataReceived(chunk) {
      const message = JSON.parse(chunk);
      const lastMessage = this.messages.find(item => item.id === message.id);
      if (message.conversation_id && message.conversation_id !== this.conversationId) {
        this.conversationId = message.conversation_id;
        this.listConversations().then();
      }
      if (lastMessage) {
        lastMessage.content += message.content;
        lastMessage.tool_call += message.tool_call;
        lastMessage.tool_call_params += message.tool_call_params;
        lastMessage.tool_call_id += message.tool_call_id;
        if (message.tool_status !== 0) {
          lastMessage.tool_status = message.tool_status;
        }
      } else {
        this.messages.push(message);
      }
    },

    clearSession() {
      this.messages = [];
      this.conversationId = "";
      this.agentId = null;
    },
    setAgentId(agentId) {
      this.agentId = agentId;
    },
    isInConversation() {
      return this.conversationId !== "";
    },
    async openHistoryConversation(conversation) {
      this.conversationId = conversation.id;
      this.agentId = conversation.agent_id;
      try {
        const res = await getConversationMessages(conversation.id);
        this.messages = res.data || [];
      } catch (error) {
        console.error(error);
      }
    }
  }
})
