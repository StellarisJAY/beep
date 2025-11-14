<template>
  <div class="model-setting-view">
    <h3>MCP服务配置</h3>
    <a-space direction="horizontal">
      <a-input
        v-model:value="mcpServerQuery.name"
        placeholder="搜索MCP服务"
        style="width: 200px"
        :disabled="loading"
      />
      <a-select
        v-model:value="mcpServerQuery.available"
        :options="statusOptions"
        placeholder="请选择状态"
        :disabled="loading"
      />
      <a-button type="primary" @click="listMcpServers" :disabled="loading">搜索</a-button>
      <a-button type="primary" @click="openModelFactoryConfigModal" :disabled="loading"
        >新增MCP</a-button
      >
    </a-space>
    <a-spin :spinning="loading">
      <div class="model-factory-card" v-for="item in mcpServerList" :key="item.id">
        <div class="model-factory-card-header">
          <a-space direction="horizontal">
            <h4>{{ item.name }}</h4>
            <a-tag v-if="item.available">可用</a-tag>
            <a-tag v-else>连接失败</a-tag>
            <a-tooltip title="编辑">
              <a-button type="link" @click="openModelFactoryConfigModal(item)"
                ><EditOutlined
              /></a-button>
            </a-tooltip>
            <a-tooltip title="删除">
              <a-button type="link" @click="handleDelete(item.id)"><DeleteOutlined /></a-button>
            </a-tooltip>
          </a-space>
        </div>
        <a-collapse v-if="item.available">
          <a-collapse-panel header="工具列表">
            <a-table
              :columns="mcpToolTableColumns"
              :data-source="item.tools || []"
              :pagination="false"
            >
            </a-table>
          </a-collapse-panel>
        </a-collapse>
      </div>
    </a-spin>

    <McpConfigModal ref="mcpConfigModalRef" @close="listMcpServers" />
  </div>
</template>

<script setup>
  import { onMounted, ref } from 'vue';
  import { deleteMcpServer, getMcpServerList } from '@/services/mcp.js';
  import { EditOutlined, DeleteOutlined } from '@ant-design/icons-vue';
  import McpConfigModal from '@/components/mcp/McpConfigModal.vue';
  import { message } from 'ant-design-vue';

  const loading = ref(false);
  const mcpConfigModalRef = ref();

  const statusOptions = ref([
    { label: '可用', value: true },
    { label: '连接失败', value: false },
  ]);

  const mcpServerQuery = ref({
    name: '',
    available: null,
  });

  const mcpServerList = ref([]);
  const mcpToolTableColumns = [
    {
      title: '工具名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '工具描述',
      dataIndex: 'description',
      key: 'description',
    },
    {
      title: '参数',
      dataIndex: 'inputSchema',
      key: 'inputSchema',
    },
  ];

  const listMcpServers = async () => {
    try {
      loading.value = true;
      const res = await getMcpServerList(mcpServerQuery.value);
      mcpServerList.value = res.data || [];
    } catch (error) {
      console.error('listMcpServers error:', error);
    } finally {
      loading.value = false;
    }
  };

  const openModelFactoryConfigModal = (record) => {
    mcpConfigModalRef.value.open(record);
  };

  const handleDelete = async (id) => {
    try {
      loading.value = true;
      await deleteMcpServer(id);
      message.success('删除成功');
      await listMcpServers();
    } catch (error) {
      console.error('deleteMcpServer error:', error);
    } finally {
      loading.value = false;
    }
  };

  onMounted(() => listMcpServers());
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
