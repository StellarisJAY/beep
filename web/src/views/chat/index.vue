<template>
  <div class="chat-view">
    <div v-show="isConversationListVisible" class="conversation-list-wrapper">
      <div class="conversation-head">
        <a-button type="link" size="large">
          <SearchOutlined />
        </a-button>
        <a-button
          type="link"
          @click="toggleConversationList"
          class="toggle-btn collapse-btn"
          size="large"
        >
          <MenuFoldOutlined />
        </a-button>
      </div>
      <a-button type="primary" style="width: 100%" :icon="h(PlusOutlined)" @click="newConversation"
        >新建会话</a-button
      >
      <a-menu class="history-list">
        <a-menu-item
          v-for="item in conversationList"
          :key="item.id"
          @click="() => chatStore.openHistoryConversation(item)"
        >
          <span>{{ item.title }}</span>
          <a-button type="link" danger><DeleteOutlined /></a-button>
        </a-menu-item>
      </a-menu>
    </div>
    <div class="chat-wrapper">
      <div class="chat-header-wrapper">
        <a-tooltip title="展开会话列表">
          <a-button
            type="link"
            v-show="!isConversationListVisible"
            @click="toggleConversationList"
            size="large"
            class="toggle-btn expand-btn"
          >
            <MenuUnfoldOutlined />
          </a-button>
        </a-tooltip>
        <a-tooltip title="新建会话">
          <a-button
            size="large"
            type="link"
            v-show="!isConversationListVisible"
            @click="newConversation"
          >
            <PlusOutlined />
          </a-button>
        </a-tooltip>
        <AgentSelect
          class="agent-select"
          ref="agentSelectRef"
          :agent-id="agentId"
          @change="handleSelectAgent"
          :disabled="isInConversation"
        />
        <a-tooltip title="智能体信息">
          <a-button type="link" size="large">
            <InfoCircleOutlined />
          </a-button>
        </a-tooltip>
      </div>
      <div class="chat-messages-wrapper">
        <div v-if="!messages || messages.length === 0" class="welcome">
          <div class="welcome-icon">
            <img :src="agentIcon" alt="智能体" width="50" height="50" />
          </div>
          <div class="welcome-title">
            与 智能体 对话
          </div>
        </div>
        <ChatMessage
          v-for="message in messages"
          :key="message.id"
          :content="message.content"
          :role="message.role"
          :toolCall="message.tool_call"
          :toolParams="message.tool_call_params"
          :toolCallStatus="message.tool_status"
        />
      </div>
      <div class="chat-input-wrapper">
        <ChatInput />
      </div>
    </div>
  </div>
</template>

<script setup>
  import { computed, h, onMounted, ref } from 'vue';
  import {
    MenuUnfoldOutlined,
    MenuFoldOutlined,
    PlusOutlined,
    InfoCircleOutlined,
  } from '@ant-design/icons-vue';
  import ChatInput from '@/components/chat/ChatInput.vue';
  import ChatMessage from '@/components/chat/ChatMessage.vue';
  import AgentSelect from '@/components/chat/AgentSelect.vue';
  import { useChatStore } from '@/stores/ChatStore.js';
  import { DeleteOutlined, SearchOutlined } from '@ant-design/icons-vue';
  import agentIcon from '@/assets/svg/avatar_robot.svg';

  const chatStore = useChatStore();
  const agentSelectRef = ref(null);
  const isConversationListVisible = ref(true);

  const messages = computed(() => chatStore.messages);
  const agentId = computed(() => chatStore.agentId);
  const isInConversation = computed(() => chatStore.isInConversation());

  const conversationList = computed(() => chatStore.conversationList);

  const toggleConversationList = () => {
    isConversationListVisible.value = !isConversationListVisible.value;
  };
  const handleSelectAgent = (id) => {
    chatStore.setAgentId(id);
  };

  const newConversation = () => {
    chatStore.clearSession();
  };

  const listConversations = () => {
    chatStore.listConversations();
  };

  onMounted(() => {
    listConversations();
  });
</script>

<style scoped>
  .chat-view {
    display: flex;
    height: 100%;
    flex-direction: row;
    justify-content: flex-start;
    gap: 0;
    padding: 5px;
  }

  .conversation-list-wrapper {
    width: 200px;
    height: 100%;
    background-color: #eaeaea;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    padding: 5px;
    display: flex;
    gap: 5px;
    flex-direction: column;
  }
  .collapse-btn {
    float: right;
  }

  .chat-wrapper {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 0;
    justify-content: center;
    align-items: center;
  }

  .chat-header-wrapper {
    width: 100%;
    height: 50px;
    display: flex;
    justify-content: start;
    gap: 1px;
    padding: 5px;
    align-items: center;
  }

  .chat-messages-wrapper {
    width: 80%;
    height: 100%;
    display: flow;
    flex-direction: column;
    gap: 5px;
    overflow-y: auto;
  }

  .chat-input-wrapper {
    width: 80%;
    display: flex;
    justify-content: center;
    padding: 0 20px;
    box-sizing: border-box;
  }

  .agent-select {
    width: 200px;
  }

  .history-list {
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 5px;
    overflow-y: auto;
    background-color: transparent;
  }

  .toggle-btn {
    float: right;
  }
  .welcome {
    height: 100%;
    display: flex;
    flex-direction: column;
    gap: 5px;
    justify-content: center;
    align-items: center;
  }

  .welcome-icon {
  }

  .welcome-title {
    font-size: 16px;
    font-weight: bold;
  }
</style>
