<template>
  <div class="agent-view">
    <div class="agent-header">
      <h3>智能体配置</h3>
      <a-row>
        <a-col :span="12">
          <a-space direction="horizontal">
            <a-input
              v-model:value="agentListQuery.name"
              placeholder="搜索智能体"
              style="width: 200px"
            />
            <a-select
              v-model:value="agentListQuery.type"
              placeholder="智能体类型"
              style="width: 200px"
              :options="agentTypeOptions"
            />
            <a-radio v-model:value="agentListQuery.create_by_me" /> 仅显示我创建的
          </a-space>
        </a-col>
        <a-col :span="12" style="text-align: right">
          <a-button type="primary" :icon="h(PlusOutlined)">创建智能体</a-button>
        </a-col>
      </a-row>
    </div>
    <div class="agent-list">
      <div class="agent-card" v-for="item in agentList" :key="item.id" @click="handleClickAgent(item)">
        <div class="agent-card-header">
          <a-space direction="horizontal">
            <img
              :src="agentIconMap[item.type]"
              alt="智能体图标"
              style="width: 30px; height: 30px; border-radius: 50%"
            />
            <h4 class="agent-card-title">{{ item.name }}</h4>
          </a-space>
        </div>
        <div class="agent-card-content">
          <span class="agent-desc-item">上次修改时间: {{ item.updated_at }}</span>
        </div>
        <div class="agent-card-footer">
          <a-dropdown trigger="click">
            <a-button type="text"><MoreOutlined /></a-button>
            <template #overlay>
              <a-menu>
                <a-menu-item key="edit">编辑</a-menu-item>
                <a-menu-item key="delete">删除</a-menu-item>
              </a-menu>
            </template>
          </a-dropdown>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
  import { ref, h, onMounted } from 'vue';
  import { PlusOutlined, MoreOutlined } from '@ant-design/icons-vue';
  import { getAgentList } from '@/services/agent.js';
  import agentIcon from '@/assets/svg/avatar_robot.svg';
  import workflowIcon from '@/assets/svg/icon_workflow.svg';
  import { useRouter } from 'vue-router';

  const router = useRouter();

  const agentListQuery = ref({
    name: '',
    type: null,
    create_by_me: false,
  });

  const agentTypeOptions = [
    {
      label: 'Reasoning&Act',
      value: 'react',
    },
    {
      label: '工作流',
      value: 'workflow',
    },
  ];

  const agentIconMap = {
    react: agentIcon,
    workflow: workflowIcon,
  };

  const agentList = ref([]);

  const listAgents = async () => {
    try {
      const res = await getAgentList(agentListQuery.value);
      agentList.value = res.data || [];
    } catch (error) {
      console.error('获取智能体列表失败', error);
    }
  };

  const handleClickAgent = (item) => {
    router.push({
      path: item.type === 'react' ? '/agent-config/react' : '/agent-config/workflow',
      query: {
        id: item.id,
      }
    });
  };

  onMounted(() => listAgents());
</script>

<style scoped>
  .agent-view {
    padding: 20px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  .agent-header {
    display: flex;
    flex-direction: column;
    gap: 10px;
    border: #eaeaea 1px solid;
    padding: 10px;
  }

  .agent-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, 300px);
    gap: 10px;
  }

  .agent-card {
    border: #cccccc 1px solid;
    padding: 10px;
    width: 300px;
    height: 150px;
    transition: all 0.3s ease;
    cursor: pointer;
    border-radius: 10px;
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .agent-card:hover {
    transform: translateY(-5px);
    box-shadow: 0 10px 20px rgba(0, 0, 0, 0.1);
    border-color: #2f70ec;
  }

  .agent-card-header {
  }

  .agent-card-content {
    width: 100%;
    height: 90%;
    display: flex;
    flex-direction: column;
    gap: 1px;
  }
  .agent-card-title {
    font-weight: bolder;
    font-size: 16px;
  }

  .agent-card-footer {
    width: 100%;
    display: flex;
    justify-content: flex-end;
  }

  .agent-desc-item {
    font-size: 12px;
    color: #666666;
  }
</style>
