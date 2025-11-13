<template>
  <div :class="`message-container message-container__${role}`">
    <div v-if="role === 'user'">
      {{ content }}
    </div>
    <div v-else-if="isToolCallRequest">
      <div v-html="DOMPurify.sanitize(marked.parse(content))" />
      <a-collapse>
        <a-collapse-panel>
          <template #header> 工具调用请求 </template>
          <template #extra>
            <a-button v-if="toolCallStatus === 1 && !isLoading" @click="signalToolCall(true)"
              >接受</a-button
            >
            <a-button v-if="toolCallStatus === 1 && !isLoading" @click="signalToolCall(false)"
              >拒绝</a-button
            >
            <a-tag v-if="toolCallStatus === 2" type="success">已接受</a-tag>
            <a-tag v-if="toolCallStatus === 3" type="error">已拒绝</a-tag>
          </template>
          <a-descriptions :column="1">
            <a-descriptions-item label="工具名称">{{ toolCall }}</a-descriptions-item>
            <a-descriptions-item label="工具参数">{{ toolParams }}</a-descriptions-item>
          </a-descriptions>
        </a-collapse-panel>
      </a-collapse>
    </div>
    <!-- 工具调用结果 -->
    <div v-else-if="role === 'tool'">
      <a-collapse>
        <a-collapse-panel header="工具调用结果">{{ content }}</a-collapse-panel>
      </a-collapse>
    </div>
    <div v-else v-html="DOMPurify.sanitize(marked.parse(content))" />
  </div>
</template>

<script setup>
  import DOMPurify from 'dompurify';
  import { marked } from 'marked';
  import { computed } from 'vue';
  import { useChatStore } from '@/stores/ChatStore.js';

  const props = defineProps({
    content: {
      type: String,
      default: '',
    },
    role: {
      type: String,
      default: 'user',
    },
    toolCall: {
      type: String,
      default: null,
    },
    toolParams: {
      type: String,
      default: null,
    },
    toolCallStatus: {
      type: Number,
      default: 0,
    },
  });

  const chatStore = useChatStore();
  const isLoading = computed(() => chatStore.isLoading);
  const isToolCallRequest = computed(() => props.toolCall && props.role === 'assistant');

  const signalToolCall = (accept) => {
    chatStore.signalToolCall(accept);
  };
</script>

<style scoped>
  .message-container {
    width: 100%;
    display: flex;
    flex-direction: column;
    gap: 1px;

    padding: 10px;
    border-radius: 10px;
    margin-bottom: 5px;
  }

  .message-container__user {
    align-items: end;
    background-color: #f5f7ff;
  }

  .message-container__assistant {
    align-items: start;
  }
</style>
