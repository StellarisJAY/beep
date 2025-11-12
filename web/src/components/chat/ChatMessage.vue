<template>
  <div :class="`message-container message-container__${role}`">
    <div :class="`message-content message-content__${role}`">
      <div v-if="role === 'user'">
        {{ content }}
      </div>
      <div v-else v-html="marked.parse(content)" />
    </div>
    <div class="message-timestamp">
      {{ formatTime(timestamp) }}
    </div>
  </div>
</template>

<script setup>
  import {marked} from 'marked';

  const props = defineProps({
    content: {
      type: String,
      default: ''
    },
    role: {
      type: String,
      default: 'user'
    },
    timestamp: {
      type: Date,
      default: () => {
        return new Date();
      }
    }
  });

  const formatTime = (timestamp) => {
    const date = new Date(timestamp);
    const year = date.getFullYear();
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const day = String(date.getDate()).padStart(2, '0');
    const hours = String(date.getHours()).padStart(2, '0');
    const minutes = String(date.getMinutes()).padStart(2, '0');
    return `${year}-${month}-${day} ${hours}:${minutes}`;
  };

</script>

<style scoped>
.message-container {
  max-width: 100%;
  min-width: 100px;
  display: flex;
  flex-direction: column;
  gap: 1px;
}
.message-content {
  padding: 10px;
  border-radius: 10px;
}

.message-container__user {
  align-items: end;
}

.message-content__user {
  background-color: #f5f7ff;
}
.message-content__assistant {
  background-color: #ffffff;
}
.message-timestamp {
  font-size: 12px;
  color: #888888;
}
</style>
