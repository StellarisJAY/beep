<template>
  <div class="chat-input-container">
    <div class="input-wrapper">
      <textarea
        ref="textareaRef"
        v-model="inputValue"
        class="chat-textarea"
        placeholder="请输入消息..."
        @input="handleInput"
        @keydown="handleKeydown"
      ></textarea>
    </div>
    <div class="input-toolbar">
      <div class="toolbar-left">
        <!-- 工具选项可以在这里添加 -->
        <button class="tool-button">📎</button>
        <button class="tool-button">😊</button>
      </div>
      <div class="toolbar-right">
        <button
          class="send-button"
          :disabled="sendDisabled"
          :loading="isLoading"
          @click="sendMessage"
        >
          发送
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, computed } from 'vue';
  import {useChatStore} from '@/stores/ChatStore.js';

  const chatStore = useChatStore();
  const isLoading = computed(() => chatStore.isLoading);
  const agentId = computed(() => chatStore.agentId);
  const inputValue = ref('');
  const textareaRef = ref(null);

  const sendDisabled = computed(()=>!agentId.value || isLoading.value || !inputValue.value.trim());

  // 处理输入事件，自动调整文本框高度
  const handleInput = () => {
    if (textareaRef.value) {
      textareaRef.value.style.height = 'auto'
      textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px'
    }
  };

  // 处理键盘事件，支持Enter发送消息，Shift+Enter换行
  const handleKeydown = (event) => {
    if (event.key === 'Enter' && !event.shiftKey) {
      event.preventDefault()
      sendMessage()
    }
  };

  // 发送消息
  const sendMessage = () => {
    if (inputValue.value.trim()) {
      // 这里可以添加发送消息的逻辑
      chatStore.sendMessage(inputValue.value);
      inputValue.value = ''
      // 重置文本框高度
      if (textareaRef.value) {
        textareaRef.value.style.height = 'auto'
      }
    }
  };
</script>

<style scoped>
  .chat-input-container {
    border: 1px solid #e0e0e0;
    border-radius: 8px;
    padding: 12px;
    background-color: white;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    width: 100%;
  }

  .input-wrapper {
    margin-bottom: 12px;
  }

  .chat-textarea {
    width: 100%;
    min-height: 80px;
    max-height: 200px;
    padding: 8px 12px;
    border: 1px solid #e0e0e0;
    border-radius: 4px;
    resize: none; /* 禁止用户手动调整大小 */
    overflow-y: auto;
    font-family: inherit;
    font-size: 14px;
    line-height: 1.5;
    box-sizing: border-box;
  }

  .chat-textarea:focus {
    outline: none;
    border-color: #409eff;
    box-shadow: 0 0 0 2px rgba(64, 158, 255, 0.2);
  }

  .input-toolbar {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .toolbar-left {
    display: flex;
    gap: 8px;
  }

  .tool-button {
    background: none;
    border: none;
    font-size: 16px;
    cursor: pointer;
    padding: 4px;
    border-radius: 4px;
  }

  .tool-button:hover {
    background-color: #f5f5f5;
  }

  .toolbar-right {
    display: flex;
  }

  .send-button {
    background-color: #409eff;
    color: white;
    border: none;
    border-radius: 4px;
    padding: 8px 16px;
    font-size: 14px;
    cursor: pointer;
  }

  .send-button:hover:not(:disabled) {
    background-color: #66b1ff;
  }

  .send-button:disabled {
    background-color: #a0cfff;
    cursor: not-allowed;
  }
</style>
