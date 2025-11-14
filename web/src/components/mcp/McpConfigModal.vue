<template>
  <a-modal title="MCP 配置" :open="visible" @cancel="close" @ok="submit">
    <a-form :model="factoryForm" :rules="rules" ref="formRef">
      <a-form-item label="MCP服务名称" name="name">
        <a-input v-model:value="factoryForm.name" placeholder="请输入MCP服务名称" />
      </a-form-item>
      <a-form-item label="MCP服务地址" name="url">
        <a-input v-model:value="factoryForm.url" placeholder="请输入MCP服务地址" />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
  import { ref } from 'vue';
  import { message } from 'ant-design-vue';
  import { createMcpServer, updateMcpServer } from '@/services/mcp.js';

  const visible = ref(false);

  const emit = defineEmits(['close']);

  const formRef = ref();
  const factoryForm = ref({
    id: null,
    name: '',
    url: '',
  });
  const rules = ref({
    name: [{ required: true, message: '请输入MCP服务名称' }],
    url: [{ required: true, message: '请输入MCP服务地址' }],
  });

  const open = (record) => {
    if (record && record.id) {
      factoryForm.value = { ...record };
    } else {
      factoryForm.value = {
        id: null,
        name: '',
        url: '',
      };
    }
    visible.value = true;
  };

  const close = () => {
    formRef.value.resetFields();
    visible.value = false;
    emit('close');
  };

  const submit = async () => {
    try {
      await formRef.value.validate();
      if (factoryForm.value.id) {
        await updateMcpServer(factoryForm.value);
        message.success('更新成功');
        close();
      } else {
        await createMcpServer(factoryForm.value);
        message.success('创建成功');
        close();
      }
    } catch (error) {
      console.log(error);
    }
  };

  defineExpose({
    open,
  });
</script>
