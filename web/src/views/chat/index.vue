<template>
  <div class="chat-view">
    <div v-show="isConversationListVisible" class="conversation-list-wrapper">
      <div class="conversation-head">
        <a-button
          type="link"
          @click="toggleConversationList"
          class="toggle-btn collapse-btn"
          size="large"
        >
          <MenuFoldOutlined />
        </a-button>
      </div>
      <a-button type="primary" style="width: 100%" :icon="h(PlusOutlined)">新建会话</a-button>
    </div>
    <div class="chat-wrapper">
      <div class="chat-header-wrapper">
        <a-button
          type="link"
          v-show="!isConversationListVisible"
          @click="toggleConversationList"
          size="large"
          class="toggle-btn expand-btn"
        >
          <MenuUnfoldOutlined />
        </a-button>
        <a-button size="large" type="link" v-show="!isConversationListVisible">
          <PlusOutlined />
        </a-button>
      </div>
      <div class="chat-messages-wrapper">
        <ChatMessage
          v-for="message in mockMessages"
          :key="message.id"
          :content="message.content"
          :role="message.role"
          :timestamp="message.timestamp"
        />
      </div>
      <div class="chat-input-wrapper">
        <ChatInput />
      </div>
    </div>
  </div>
</template>

<script setup>
  import { h, ref } from 'vue';
  import { MenuUnfoldOutlined, MenuFoldOutlined, PlusOutlined } from '@ant-design/icons-vue';
  import ChatInput from '@/components/chat/ChatInput.vue';
  import ChatMessage from '@/components/chat/ChatMessage.vue';

  const isConversationListVisible = ref(true);

  const toggleConversationList = () => {
    isConversationListVisible.value = !isConversationListVisible.value;
  };

  const mockMessages = ref([
    {
      id: 1,
      content: 'Hello, I am a human.',
      role: 'user',
      timestamp: new Date(),
    },
    {
      id: 2,
      content:
        'Hello, I am a AI assistant. How can I help you?',
      role: 'assistant',
      timestamp: new Date(),
    },
  ]);
</script>

<style scoped>
  .chat-view {
    display: flex;
    height: 100%;
    flex-direction: row;
    justify-content: flex-start;
    gap: 0;
  }

  .conversation-list-wrapper {
    width: 200px;
    height: 100%;
    background-color: #eaeaea;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    padding: 5px;
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
</style>
