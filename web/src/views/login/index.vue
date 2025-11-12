<template>
  <div class="login-view">
    <div class="login-card">
      <div class="login-card-header">
        <div class="login-card-header-title">登录</div>
      </div>
      <a-form :model="loginForm" :rules="loginFormRule" ref="loginFormRef">
        <a-form-item label="邮箱" name="email">
          <a-input v-model:value="loginForm.email" placeholder="请输入邮箱" />
        </a-form-item>
        <a-form-item label="密码" name="password">
          <a-input v-model:value="loginForm.password" placeholder="请输入密码" type="password" />
        </a-form-item>
      </a-form>
      <a-button type="primary" @click="doLogin">登录</a-button>
    </div>
  </div>
</template>

<script setup>
  import { useRouter } from 'vue-router';
  import { ref } from 'vue';
  import {login} from '@/services/base.js';
  import { useUserStore } from '@/stores/UserStore.js';
  import { message } from 'ant-design-vue';

  const userStore = useUserStore();
  const router = useRouter();
  const currentRoute = router.currentRoute;
  const backPath = currentRoute.value.params.back || '/';

  const loginFormRef = ref();

  const loginFormRule = {
    email: [
      { required: true, message: '请输入邮箱' },
      { type: 'email', message: '请输入正确的邮箱格式' },
    ],
    password: [
      { required: true, message: '请输入密码' },
    ],
  }

  const loginForm = ref({
    email: '',
    password: '',
  });

  const doLogin = async () => {
    try {
      await loginFormRef.value.validate();
      const res = await login(loginForm.value);
      console.log(res);
      userStore.setUserInfo(res.data['user_info']);
      userStore.setWorkspaceInfo(res.data['workspace_info']);
      message.success('登录成功');
      await router.push(backPath);
    } catch (error) {
      console.log(error);
    }
  };
</script>

<style scoped>
.login-view {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100vh;
}
.login-card {
  width: 400px;
  padding: 20px;
  border-radius: 4px;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
}
</style>
