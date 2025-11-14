import { defineStore } from 'pinia';

export const defaultReActAgentConfig = {
  re_act: {
    prompt: 'You are a helpful assistant, please answer the user\'s question.',
    max_iterations: 5,
    chat_model: null,
    knowledge_bases: [],
    mcp_servers: [],
    use_system_tools: true,
    memory_option: {
      enable_short_term_memory: true,
      enable_long_term_memory: true,
      memory_control: 'static',
      memory_window_size: 32,
    },
    retriever_option: {
      top_k: 10,
      threshold: 0.5,
      search_type: 'vector',
      hybrid_type: 'weight',
      weight: 0.5,
      reranker: null,
    }
  },
};

export const reactAgentConfigRules = {
  prompt: {
    required: true,
    message: '请输入提示词',
    trigger: ['submit']
  },
  max_iterations: {
    required: true,
    message: '请输入最大迭代次数',
    trigger: ['submit']
  },
  chat_model: {
    required: true,
    message: '请选择聊天模型',
    trigger: ['submit']
  },
};

export const useAgentConfigStore = defineStore('agentConfig', {
  state: () => ({
    agentDetail: {},
  }),
  actions: {
    setAgentDetail(agentDetail) {
      this.agentDetail = agentDetail;
    },
  }
})
