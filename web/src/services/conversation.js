import httpReq from '@/services/http.js';

// 获取历史会话列表
export const getConversationList = async (query) => await httpReq.get('/conversation/list', query);

// 获取历史会话消息
export const getConversationMessages = async (conversationId) => await httpReq.get('/conversation/messages', { id: conversationId });
