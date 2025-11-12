<template>
  <a-modal
    title="切换工作空间"
    :open="visible"
    :footer="null"
    @cancel="handleCancel"
  >
    <a-select
      :value="currentWorkspaceId"
      :options="workspaceOptions"
      @change="handleWorkspaceChange"
      style="width: 200px"
    />
    <a-table
      :columns="workspaceMemberTableColumns"
      :data-source="workspaceMemberList"
      :pagination="false"
    />
  </a-modal>
</template>

<script setup>
  import { useUserStore } from '@/stores/UserStore.js';
  import { computed, onMounted, ref, watch } from 'vue';
  import { getWorkspaceMembers, listWorkspace, switchWorkspace } from '@/services/base.js';
  import { message } from 'ant-design-vue';

  const userStore = useUserStore();
  const workspaceOptions = ref([]);
  const workspaces = ref([]);
  const visible = ref(false);
  const currentWorkspaceId = computed(() => userStore.getWorkspaceInfo().id);

  const workspaceMemberTableColumns = [
    {
      title: '名称',
      dataIndex: 'name',
      key: 'name',
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
    },
  ];

  const workspaceMemberList = ref([]);

  // 获取工作空间列表
  const getWorkspaceList = async () => {
    try {
      const res = await listWorkspace();
      workspaces.value = res.data;
      workspaceOptions.value = res.data.map((item) => ({
        label: item.name,
        value: item.id,
      }));
    } catch (error) {
      console.error(error);
    }
  };

  // 切换工作空间
  const handleWorkspaceChange = async (id) => {
    try {
      await switchWorkspace(id);
      message.success('切换工作空间成功');
      userStore.setWorkspaceInfo(workspaces.value.find((item) => item.id === id));
      await getWorkspaceMemberList(id);
    } catch (error) {
      console.error(error);
    }
  };

  // 查询工作空间成员列表
  const getWorkspaceMemberList = async (id) => {
    try {
      const res = await getWorkspaceMembers({ id: id });
      workspaceMemberList.value = res.data;
    } catch (error) {
      console.error(error);
    }
  };

  const handleCancel = () => {
    visible.value = false;
  };

  onMounted(async () => {
    await getWorkspaceList();
    await getWorkspaceMemberList(currentWorkspaceId.value);
  });

  defineExpose({
    show: () => {
      visible.value = true;
    },
  });
</script>
