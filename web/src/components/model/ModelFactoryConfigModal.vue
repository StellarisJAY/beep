<template>
  <a-modal title="供应商配置" :open="visible" @cancel="close" @ok="submit">
    <a-form :model="factoryForm" :rules="rules" ref="formRef">
      <a-form-item label="供应商类型" name="type" v-if="!factoryForm.id">
        <a-select
          v-model:value="factoryForm.type"
          :options="modelFactoryTypeOptions"
          placeholder="请选择供应商类型"
          @change="handleTypeChange"
        />
      </a-form-item>
      <a-form-item label="供应商名称" name="name">
        <a-input v-model:value="factoryForm.name" placeholder="请输入供应商名称" />
      </a-form-item>
      <a-form-item label="Base URL" name="base_url" v-if="needBaseUrl(factoryForm.type)">
        <a-input v-model:value="factoryForm.base_url" placeholder="请输入API地址" />
      </a-form-item>
      <a-form-item label="API Key" name="api_key">
        <a-input-password
          v-model:value="factoryForm.api_key"
          placeholder="请输入 API Key"
          :visibility-toggle="false"
        />
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<script setup>
  import { ref } from 'vue';
  import { createModelFactory } from '@/services/model.js';
  import { message } from 'ant-design-vue';

  const visible = ref(false);

  const emit = defineEmits(['close']);

  const modelFactoryTypeOptions = ref([
    { label: 'OpenAI', value: 'openai' },
    { label: '通义千问', value: 'dashscope' },
    { label: 'Ollama', value: 'ollama' },
  ]);

  const formRef = ref();
  const factoryForm = ref({
    id: null,
    type: '',
    name: '',
    api_key: '',
  });
  const rules = ref({
    type: [{ required: true, message: '请选择供应商类型' }],
    name: [{ required: true, message: '请输入供应商名称' }],
    api_key: [{ required: true, message: '请输入 API Key' }],
    base_url: [{ required: true, message: '请输入API地址' }],
  });

  const open = (record) => {
    if (record && record.id) {
      factoryForm.value = { ...record };
    } else {
      factoryForm.value = {
        id: null,
        type: '',
        name: '',
        api_key: '',
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
        // 更新供应商配置
        message.success('更新成功');
        close();
      } else {
        // 创建供应商配置
        await createModelFactory(factoryForm.value);
        message.success('创建成功');
        close();
      }
    } catch (error) {
      console.log(error);
    }
  };

  const needBaseUrl = (type) => {
    return type === 'openai' || type === 'ollama';
  };

  const handleTypeChange = (type) => {
    if (needBaseUrl(type)) {
      factoryForm.value.base_url = defaultBaseUrlMap[type];
    }
  };

  const defaultBaseUrlMap = {
    openai: 'https://api.openai.com/v1',
    ollama: 'http://localhost:11434/v1',
  };

  defineExpose({
    open,
  });
</script>
