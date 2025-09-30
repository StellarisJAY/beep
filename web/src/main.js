import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue';
import { initRouter } from '@/router/router.js';

const app = createApp(App);
async function bootstrap() {
  app.use(createPinia());
  app.use(initRouter());
  app.mount('#app');
}

await bootstrap();
