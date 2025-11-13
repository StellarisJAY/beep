import {createRouter, createWebHistory  } from 'vue-router';
import { useUserStore } from '@/stores/UserStore.js';
import { h } from 'vue';
import {
  BookOutlined,
  MessageOutlined,
  RobotOutlined,
  SettingOutlined,
} from '@ant-design/icons-vue';


const routes = [
  {
    path: '',
    component: () => import('@/layouts/default/index.vue'),
    redirect: '/chat',
    children: [
      {
        name: '聊天',
        path: '/chat',
        component: () => import('@/views/chat/index.vue'),
        meta: {
          icon: h(MessageOutlined),
          showInMenu: true,
        }
      },
      {
        name: '智能体',
        path: '/agent',
        component: () => import('@/views/agent/index.vue'),
        meta: {
          icon: h(RobotOutlined),
          showInMenu: true,
        }
      },
      {
        name: '知识库',
        path: '/knowledge',
        component: () => import('@/views/knowledge/index.vue'),
        meta: {
          icon: h(BookOutlined),
          showInMenu: true,
        },
      },
      {
        name: '设置',
        path: '/setting',
        component: () => import('@/views/setting/index.vue'),
        meta: {
          icon: h(SettingOutlined),
          showInMenu: true,
        },
        redirect: '/setting/user',
        children: [
          {
            name: '用户设置',
            path: 'user',
            component: () => import('@/views/setting/user/UserSetting.vue'),
          },
          {
            name: '模型提供商',
            path: 'model',
            component: () => import('@/views/setting/model/ModelSetting.vue'),
          },
          {
            name: 'MCP',
            path: 'mcp',
            component: () => import('@/views/setting/mcp/MCPSetting.vue'),
          },
          {
            name: '工作空间',
            path: 'workspace',
            component: () => import('@/views/setting/workspace/WorkspaceSetting.vue'),
          },
        ]
      },
    ],
  },
  {
    name: '登录',
    path: '/login',
    component: () => import('@/views/login/index.vue'),
  }
];

const beforeEach = (to, from, next) => {
  if (to.name === '登录') {
    next();
    return;
  }
  if (!localStorage.getItem('beep_token')) {
    next({ name: '登录', query: { back: to.fullPath } });
  } else {
    const userStore = useUserStore();
    userStore.queryLoginInfo().then(isLogin => {
      if (isLogin) next();
      else next({ name: '登录', query: { back: to.fullPath } });
    })
  }
};

export const router = createRouter({
  history: createWebHistory(),
  routes: routes,
});

router.beforeEach(beforeEach);
