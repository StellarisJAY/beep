<template>
  <div class="model-setting-view">
    <h3>已配置供应商</h3>
    <a-space direction="horizontal">
      <a-input v-model:value="modelFactoryQuery.name" placeholder="搜索供应商" style="width: 200px"/>
      <a-select
        v-model:value="modelFactoryQuery.type"
        :options="modelFactoryTypeOptions"
        placeholder="搜索供应商类型"
        allow-clear
        style="width: 200px"
      />
      <a-button type="primary" @click="listModelFactory">搜索</a-button>
      <a-button type="primary" @click="openModelFactoryConfigModal">添加供应商</a-button>
    </a-space>
    <div class="model-factory-card" v-for="item in modelFactoryList" :key="item.id">
      <div class="model-factory-card-header">
        <a-space direction="horizontal">
          <img :src="modelFactoryIconMap[item.type]" alt="供应商图标" width="40" height="40" />
          <h4>{{ item.name }}</h4>
          <a-button type="link" @click="openModelFactoryConfigModal(item)">编辑</a-button>
        </a-space>
      </div>
      <a-collapse>
        <a-collapse-panel header="模型列表">
          <a-table
            :columns="modelTableColumns"
            :data-source="item.models || []"
            :pagination="false"
          >
            <template #bodyCell="{ column, record }">
              <a-switch
                v-if="column.dataIndex === 'status'"
                v-model:checked="record.status"
                checked-children="启用"
                unchecked-children="停用"
              />
              <a-space direction="horizontal" v-if="column.dataIndex === 'tags'">
                <a-tag v-for="tag in record.tags.split(',')" :key="tag">{{ tag }}</a-tag>
              </a-space>
              <a-tag v-if="column.dataIndex === 'type'">{{ record.type }}</a-tag>
            </template>
          </a-table>
        </a-collapse-panel>
      </a-collapse>
    </div>

    <ModelFactoryConfigModal ref="modelFactoryConfigModalRef" @close="listModelFactory" />
  </div>
</template>

<script setup>
  import { onMounted, ref } from 'vue';
  import { getModelFactoryList } from '@/services/model.js';
  import openai_icon from '@/assets/svg/openai_icon.svg';
  import dashscope_icon from '@/assets/svg/dashscope_icon.svg';
  import ollama_icon from '@/assets/svg/ollama_icon.svg';
  import ModelFactoryConfigModal from '@/components/model/ModelFactoryConfigModal.vue';

  const modelFactoryConfigModalRef = ref();

  const modelFactoryTypeOptions = ref([
    { label: 'OpenAI', value: 'openai' },
    { label: '通义千问', value: 'dashscope' },
    { label: 'Ollama', value: 'ollama' },
  ]);

  const modelFactoryIconMap = {
    openai: openai_icon,
    dashscope: dashscope_icon,
    ollama: ollama_icon,
  };

  const modelFactoryQuery = ref({
    name: '',
    type: null,
  });

  const modelFactoryList = ref([]);
  const modelTableColumns = [
    {
      title: '模型名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '模型类型',
      dataIndex: 'type',
      key: 'type',
    },
    {
      title: '标签',
      dataIndex: 'tags',
      key: 'tags',
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
    },
  ];

  const listModelFactory = async () => {
    try {
      const res = await getModelFactoryList(modelFactoryQuery.value);
      modelFactoryList.value = res.data || [];
    } catch (error) {
      console.error('listMcpServers error:', error);
    }
  };

  const openModelFactoryConfigModal = (record) => {
    modelFactoryConfigModalRef.value.open(record);
  };

  onMounted(() => listModelFactory());
</script>

<style scoped>
  .model-setting-view {
    padding: 10px;
    height: 100%;
    width: 100%;
    display: flow;
    gap: 10px;
    overflow-y: auto;
  }
  .model-factory-card {
    border-radius: 5px;
    padding: 10px;
    background-color: #f5f5f5;
    width: 100%;
    height: fit-content;
    border: #eaeaea 1px solid;
    margin-top: 10px;
  }
</style>
