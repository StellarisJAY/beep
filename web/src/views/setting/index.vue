<template>
  <div class="setting-view">
    <a-menu
      class="side-bar"
      :items="menuItems"
      :selected-keys="activeKey"
      @click="handleMenuClick"
    />
    <div class="content-wrapper">
      <RouterView />
    </div>
  </div>
</template>

<script setup>
  import { computed, h } from 'vue';
  import {
    CloudOutlined,
    DeploymentUnitOutlined,
    TeamOutlined,
    UserOutlined,
  } from '@ant-design/icons-vue';
  import { useRouter } from 'vue-router';

  const router = useRouter();
  const activeKey = computed(() => [router.currentRoute.value.path]);

  const menuItems = [
    {
      title: '用户设置',
      key: '/setting/user',
      label: '用户设置',
      icon: h(UserOutlined),
    },
    {
      title: '模型设置',
      key: '/setting/model',
      label: '模型设置',
      icon: h(DeploymentUnitOutlined),
    },
    {
      title: 'MCP设置',
      key: '/setting/mcp',
      label: 'MCP设置',
      icon: h(CloudOutlined),
    },
    {
      title: '工作空间设置',
      key: '/setting/workspace',
      label: '工作空间设置',
      icon: h(TeamOutlined),
    },
  ];

  const handleMenuClick = (item) => {
    router.replace(item.key);
  };
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

  .setting-item {
    width: 100%;
    height: 40px;
    line-height: 40px;
    text-align: center;
    border-radius: 5px;
  }
</style>
