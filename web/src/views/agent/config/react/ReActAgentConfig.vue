<template>
  <div class="agent-config-view">
    <h3>思考&行动(ReAct)智能体配置</h3>
    <a-divider />
    <a-form v-if="agentConfig" :model="agentConfig" :rules="reactAgentConfigRules">
      <a-form-item label="提示词" name="prompt">
        <a-textarea v-model:value="agentConfig.prompt" rows="12" placeholder="请输入提示词" />
      </a-form-item>
      <a-form-item label="最大迭代次数" name="max_iterations">
        <a-input-number
          v-model:value="agentConfig.max_iterations"
          placeholder="请输入最大迭代次数"
        />
      </a-form-item>
      <a-form-item label="关联知识库" name="knowledge_bases">
        <a-select
          v-model:value="agentConfig.knowledge_bases"
          multiple
          placeholder="请选择关联知识库"
        />
      </a-form-item>
      <a-collapse>
        <a-collapse-panel header="MCP配置">
          <a-button type="primary" @click="handleAddMcpServer">选择MCP服务</a-button>
          <a-table
            :pagination="false"
            :columns="mcpServerColumns"
            :data-source="selectedMcpServers"
          />
        </a-collapse-panel>
        <a-collapse-panel header="知识检索配置">
          <a-form-item label="检索方式" name="retrieval_method">
            <a-radio-group v-model:value="agentConfig.retriever_option.search_type">
              <a-radio value="vector">向量检索</a-radio>
              <a-radio value="fulltext">全文检索</a-radio>
              <a-radio value="hybrid">混合检索</a-radio>
            </a-radio-group>
          </a-form-item>

          <a-form-item label="TopK" name="tok_k">
            <a-input-number
              v-model:value="agentConfig.retriever_option.top_k"
              :min="1"
              :max="999"
            />
          </a-form-item>

          <a-form-item
            label="相似度阈值"
            name="similarity_threshold"
            v-if="agentConfig.retriever_option.search_type !== 'fulltext'"
          >
            <a-slider
              v-model:value="agentConfig.retriever_option.threshold"
              :min="0"
              :max="1"
              :step="0.01"
              show-input
            />
          </a-form-item>

          <a-form-item
            v-if="agentConfig.retriever_option.search_type === 'hybrid'"
            label="混合检索类型"
            name="hybrid_type"
          >
            <a-radio-group v-model:value="agentConfig.retriever_option.hybrid_type">
              <a-radio value="weight">权重</a-radio>
              <a-radio value="rrf">RRF</a-radio>
            </a-radio-group>
          </a-form-item>

          <a-form-item
            v-if="agentConfig.retriever_option.search_type === 'hybrid'"
            label="权重"
            name="weight"
          >
            <a-slider
              v-model:value="agentConfig.retriever_option.weight"
              :min="0"
              :max="1"
              :step="0.01"
              show-input
            />
          </a-form-item>
        </a-collapse-panel>
        <a-collapse-panel header="记忆选项">
          <a-form-item label="短期记忆" name="enable_short_term_memory">
            <a-switch v-model:checked="agentConfig.memory_option.enable_short_term_memory" />
          </a-form-item>
          <a-form-item label="长期记忆" name="enable_long_term_memory">
            <a-switch v-model:checked="agentConfig.memory_option.enable_long_term_memory" />
          </a-form-item>
          <a-form-item label="记忆控制" name="memory_control">
            <a-radio-group v-model:value="agentConfig.memory_option.memory_control">
              <a-radio value="static">静态控制</a-radio>
              <a-radio value="agentic">自主控制</a-radio>
            </a-radio-group>
          </a-form-item>
          <a-form-item label="记忆窗口大小" name="memory_window_size">
            <a-input-number
              v-model:value="agentConfig.memory_option.memory_window_size"
              :min="1"
              :max="999"
            />
          </a-form-item>
        </a-collapse-panel>
      </a-collapse>
    </a-form>
    <McpSelector
      v-if="agentConfig"
      ref="mcpSelectorRef"
      v-model:selected="agentConfig.mcp_servers"
      v-model:selectedRecords="selectedMcpServers"
    />
  </div>
</template>

<script setup>
  import { defaultReActAgentConfig, useAgentConfigStore } from '@/stores/AgentStore';
  import { computed, ref, watch } from 'vue';
  import { reactAgentConfigRules } from '@/stores/AgentStore';
  import McpSelector from '@/components/mcp/McpSelector.vue';

  const agentConfigStore = useAgentConfigStore();
  const agentDetail = computed(() => agentConfigStore.agentDetail);
  const agentConfig = ref();
  const mcpSelectorRef = ref();
  const selectedMcpServers = ref([]);

  const mcpServerColumns = [
    {
      title: 'MCP名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
    },
  ];

  const handleAddMcpServer = () => {
    mcpSelectorRef.value.open();
  };

  watch(agentDetail, (newAgentDetail) => {
    if (!newAgentDetail.id) {
      newAgentDetail.config = JSON.parse(JSON.stringify(defaultReActAgentConfig));
    }
    agentConfig.value = newAgentDetail.config.re_act;
    console.log(newAgentDetail);
    selectedMcpServers.value = newAgentDetail.mcp_servers || [];
  });
</script>

<style scoped>
  .agent-config-view {
    width: 80%;
    height: 100%;
    overflow-y: auto;
    display: flow;
    flex-direction: column;
    padding: 10px;
  }
</style>
