<template>
  <a-modal
    title="选择MCP服务"
    :open="visible"
    @cancel="close"
    :destroy-on-close="true"
    @ok="handleOk"
  >
    <a-table
      :columns="columns"
      :data-source="mcpServerList"
      size="small"
      row-key="id"
      :row-selection="{
        selectedRowKeys: selectedRowKeys,
        onChange: onSelectChange
      }"
    />
  </a-modal>
</template>

<script setup>
  import { onMounted, ref } from 'vue';
  import { getMcpServerListWithoutTool } from '@/services/mcp.js';

  const visible = ref(false);
  const mcpServerList = ref([]);
  const selectedRowKeys = ref([]);

  // 定义选中的行数据
  const selectedRows = ref([]);

  const columns = [
    {
      title: 'MCP服务',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: 'URL',
      dataIndex: 'url',
      key: 'url',
    },
  ];

  const selected = defineModel('selected', {
    type: Array,
    default: () => [],
  });
  const selectedRecords = defineModel('selectedRecords', {
    type: Array,
    default: () => [],
  });

  // 当选择改变时更新选中项
  const onSelectChange = (selectedKeys, selectedRowsData) => {
    selectedRowKeys.value = selectedKeys;
    selectedRows.value = selectedRowsData;
  };

  const close = () => {
    visible.value = false;
  };

  // 处理确认按钮点击事件
  const handleOk = () => {
    // 更新父组件的selected值
    selected.value = [...selectedRows.value];
    selectedRecords.value = [...selectedRows.value];
    console.log(selected.value);
    console.log(selectedRecords.value);
    close();
  };

  const listMcpServers = async () => {
    try {
      const res = await getMcpServerListWithoutTool();
      mcpServerList.value = res.data || [];

      // 如果已经有选中的项，设置默认选中状态
      if (selected.value && selected.value.length > 0) {
        selectedRowKeys.value = selected.value;
      }
    } catch (error) {
      console.error('获取MCP服务列表失败', error);
    }
  };

  onMounted(() => {
    listMcpServers();
  });

  defineExpose({
    open() {
      visible.value = true;
      // 每次打开时重新加载数据并设置选中状态
      listMcpServers();
    },
  });
</script>

<style scoped></style>
