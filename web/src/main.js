import { createApp } from 'vue';
import { createPinia } from 'pinia';

import App from './App.vue'

const app = createApp(App)
async function bootstrap() {
  app.use(createPinia());
  app.mount('#app');
}

await bootstrap();
