import {createRouter, createWebHistory  } from 'vue-router';

const routes = [
  {
    path: '',
    component: () => import('@/layouts/default/index.vue'),
    children: [
      {
        name: '聊天',
        path: '/chat',
        component: () => import('@/views/chat/index.vue'),
        meta: {
          icon: 'assets/svg/icon_chat.svg',
          showInMenu: true,
        }
      },
      {
        name: '智能体',
        path: '/agent',
        component: () => import('@/views/agent/index.vue'),
        meta: {
          icon: 'assets/svg/icon_agent.svg',
          showInMenu: true,
        }
      },
      {
        name: '知识库',
        path: '/knowledge',
        component: () => import('@/views/knowledge/index.vue'),
        meta: {
          icon: 'assets/svg/icon_knowledge.svg',
          showInMenu: true,
        }
      },
    ],
  },
  {
    path: '/login',
    component: () => import('@/views/login/index.vue'),
  }
];

export const initRouter = () => {
  return createRouter({
    history: createWebHistory(),
    routes: routes,
  });
};
