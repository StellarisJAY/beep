import httpReq from './http.js';

export const getAgentList = async (query) => await httpReq.get('/agent/list', query);

export const getAgentDetail = async (id) => await httpReq.get(`/agent/detail/${id}`);
