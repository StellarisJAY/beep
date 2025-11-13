<template>
  <a-select
    :value="agentId"
    :options="agentOptions"
    placeholder="请选择智能体"
    @change="handleChange"
    :disabled="disabled"
  />
</template>

<script setup>
  import { getAgentList } from '@/services/agent.js';
  import { onMounted, ref, computed } from 'vue';

  defineProps({
    agentId: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  });
  const emit = defineEmits(['change']);
  const agentList = ref([]);
  const agentOptions = computed(() =>
    agentList.value.map((item) => ({
      label: item.name,
      value: item.id,
    }))
  );

  const listAgents = async () => {
    try {
      const res = await getAgentList();
      agentList.value = res.data || [];
    } catch (error) {
      console.error(error);
    }
  };

  const handleChange = (value) => {
    emit('change', value);
  };

  onMounted(() => listAgents());

  defineExpose({
    refresh: listAgents,
  });
</script>

<style scoped></style>
