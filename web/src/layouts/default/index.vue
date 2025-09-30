<template>
  <div class="page">
    <div class="left-side">
      <div class="logo-wrapper"></div>
      <div class="menu-wrapper">
        <div v-for="item in menuItems" :key="item.path" class="menu-item" @click="router.push(item.path)">
          {{ item.name }}
        </div>
      </div>
      <div class="user-wrapper"></div>
    </div>
    <div class="right-side">
      <div class="content-card">
        <RouterView />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { RouterView, useRouter } from 'vue-router';

const router = useRouter();
const menuItems = router.getRoutes().filter(route=>route.meta.showInMenu).map(route=>{
  return {
    name: route.name,
    path: route.path,
  };
});
</script>

<style scoped>
  .page {
    display: flex;
    justify-content: flex-start;
    height: 100vh;
    background-color: cornflowerblue;
  }

  .left-side {
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 80px;
    height: calc(100% - 20px);
    padding-left: 10px;
    padding-top: 10px;
    padding-bottom: 10px;
  }

  .right-side {
    width: calc(100% - 80px);
    height: calc(100% - 20px);
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
  }

  .user-wrapper {
    width: 100%;
    height: 80px;
  }

  .menu-item {
    width: 100%;
    height: 40px;
    line-height: 40px;
    text-align: center;
    cursor: pointer;
  }
</style>
