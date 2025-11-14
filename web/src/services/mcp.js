import httpReq from '@/services/http.js';

export const getMcpServerList = async (query) => await httpReq.get('/mcp/list', query);

export const createMcpServer = async (data) => await httpReq.post('/mcp/create', data);

export const updateMcpServer = async (data) => await httpReq.post('/mcp/update', data);

export const deleteMcpServer = async (id) => await httpReq.delete(`/mcp/delete/${id}`);

export const getMcpServerListWithoutTool = async (query) => await httpReq.get('/mcp/listWithoutTools', query);
