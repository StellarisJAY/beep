import httpReq from './http.js';

export const getAgentList = async (query) => await httpReq.get('/agent/list', query);
