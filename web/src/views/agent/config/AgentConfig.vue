<template>
  <div class="setting-view">
    <div class="side-bar">

    </div>
    <div class="content-wrapper">
      <RouterView/>
    </div>
  </div>
</template>

<script setup>

  import { useRouter } from 'vue-router';
  import { getAgentDetail } from '@/services/agent.js';
  import { onMounted, ref, watch } from 'vue';
  import { useAgentConfigStore } from '@/stores/AgentStore.js';

  const router = useRouter();
  const id = router.currentRoute.value.query.id;
  const agentDetail = ref({});
  const agentConfigStore = useAgentConfigStore();

  const getAgent = async () => {
    if (!id) return;
    try {
      const res = await getAgentDetail(id);
      agentDetail.value = res.data;
      agentConfigStore.setAgentDetail(res.data);
    } catch (error) {
      console.log(error);
    }
  };

  onMounted(()=>getAgent());
  watch(id, getAgent);
</script>

<style scoped>
  .setting-view {
    display: flex;
    height: 100%;
    flex-direction: row;
    justify-content: flex-start;
    gap: 0;
  }

  .side-bar {
    width: 200px;
    height: 100%;
    border-top-left-radius: 10px;
    border-bottom-left-radius: 10px;
    padding: 5px;
    display: flow;
    flex-direction: column;
    gap: 5px;
    justify-content: center;
    align-items: center;
    border-right: #eaeaea 1px solid;
  }

  .content-wrapper {
    height: 100%;
    width: 80%;
    padding: 10px;
    display: flex;
    flex-direction: column;
    justify-content: flex-start;
    align-items: center;
    gap: 10px;
  }
</style>
