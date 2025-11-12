import axios  from 'axios';
import { message } from 'ant-design-vue';
import {router} from '@/router/router.js';

const http =axios.create({
  baseURL: import.meta.env.VITE_API_URL,
  timeout: 10000,
});

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('beep_token');
  if (token) {
    config.headers['Access-Token'] = token;
  }
  return config;
});

http.interceptors.response.use((response)=>{
  const body = response.data;

  if (body.code === 401) {
    localStorage.removeItem('beep_token');
    message.error('登录过期，请重新登录').then();
    router.push('/login').then();
  }
  if (body.code === 200) {
    if (response.headers['access-token']) {
      localStorage.setItem('beep_token', response.headers['access-token']);
    }
    return body;
  } else {
    message.error(body.message || '请求失败').then();
    return Promise.reject(body);
  }
}, (error)=>{
  return Promise.reject(error);
});

export default  {
  get: (url, params) => http.get(url, { params }),
  post: (url, data) => http.post(url, data),
  put: (url, data) => http.put(url, data),
  delete: (url) => http.delete(url),
};
