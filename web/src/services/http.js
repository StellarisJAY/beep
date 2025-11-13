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

const postEventStream = async (path, data, onDataReceived) => {
  const baseURL = import.meta.env.VITE_API_URL;
  const url = baseURL + path;
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Access-Token': localStorage.getItem('beep_token'),
      'Accept': 'text/event-stream',
    },
    body: JSON.stringify(data),
  });
  const reader = response.body.getReader();
  while (true) {
    const { done, value } = await reader.read();
    if (done) {
      break;
    }
    const chunk = new TextDecoder().decode(value);
    // 截取data: 后的内容
    const dataChunk = chunk.split('data:', 2)[1];
    if (dataChunk) {
      onDataReceived(dataChunk);
    }
  }
};

export default  {
  get: (url, params) => http.get(url, { params }),
  post: (url, data) => http.post(url, data),
  put: (url, data) => http.put(url, data),
  delete: (url) => http.delete(url),
  postEventStream: postEventStream,
};
