<template>
  <div class="page">
    <div class="left-side">
      <div class="logo-wrapper"></div>
      <div class="menu-wrapper">
        <div
          v-for="item in menuItems"
          :key="item.path"
          class="menu-item"
          @click="router.push(item.path)"
        >
          <component :is="item.icon" class="menu-icon" />
          {{ item.name }}
        </div>
      </div>
      <a-dropdown>
        <a-avatar :src="loginUserInfo.avatar" />
        <template #overlay>
          <a-menu>
            <a-menu-item key="0" @click="showSwitchWorkspaceModal">
              <a> 切换工作空间 </a>
            </a-menu-item>
            <a-menu-item key="1" @click="doLogout">
              <a> 退出登录 </a>
            </a-menu-item>
          </a-menu>
        </template>
      </a-dropdown>
    </div>
    <div class="right-side">
      <div class="content-card">
        <RouterView />
      </div>
    </div>
    <SwitchWorkspaceModal ref="switchWorkspaceModalRef" />
  </div>
</template>

<script setup lang="js">
  import { RouterView, useRouter } from 'vue-router';
  import { useUserStore } from '@/stores/UserStore.js';
  import { computed, ref } from 'vue';
  import SwitchWorkspaceModal from '@/components/workspace/SwitchWorkspaceModal.vue';

  const switchWorkspaceModalRef = ref();
  const userStore = useUserStore();
  const router = useRouter();
  const menuItems = router
    .getRoutes()
    .filter((route) => route.meta.showInMenu)
    .map((route) => {
      return {
        name: route.name,
        path: route.path,
        icon: route.meta.icon,
      };
    });

  const loginUserInfo = computed(() => userStore.getLoginUserInfo());

  const doLogout = () => {
    userStore.clearLoginInfo();
    router.push('/login');
  };

  const showSwitchWorkspaceModal = () => {
    switchWorkspaceModalRef.value.show();
  };
</script>

<style scoped>
  .page {
    display: flex;
    justify-content: flex-start;
    height: 100vh;
    background-color: #2f70ec;
    align-items: center;
  }

  .left-side {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 60px;
    height: 100%;
    padding-left: 10px;
    padding-top: 10px;
    padding-bottom: 10px;
  }

  .right-side {
    width: 100%;
    height: 100%;
    padding: 10px;
  }

  .content-card {
    background-color: white;
    border-radius: 10px;
    width: 100%;
    height: 100%;
  }

  .logo-wrapper {
    width: 100%;
    height: 80px;
  }

  .menu-wrapper {
    width: 100%;
    height: calc(100% - 160px);
    display: flex;
    flex-direction: column;
    gap: 10px;

    .menu-item {
      width: 60px;
      height: 60px;
      text-align: center;
      cursor: pointer;
      color: #eaeaea;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      gap: 5px;
    }

    .menu-item:hover {
      color: #ffffff;
    }

    .menu-icon {
      width: 40px;
      height: 40px;
      border-radius: 20px;
    }
  }
  .user-wrapper {
    width: 100%;
    height: 80px;
    display: flow;
    align-items: center;
    justify-content: center;
    color: #eaeaea;
    gap: 10px;
    text-align: center;
    .user-name {
      font-size: 14px;
      text-align: center;
    }
  }

  .user-wrapper:hover {
    color: #ffffff;
  }
</style>
