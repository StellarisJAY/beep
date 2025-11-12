import httpReq from './http.js';

export const login = async (data) => await httpReq.post('/login', data);

export const register = async (data) => await httpReq.post('/register', data);

export const listWorkspace = async (query) => await httpReq.get('/workspace/list', query);

export const getWorkspaceMembers = async (query) => await httpReq.get('/workspace/members', query);

export const switchWorkspace = async (id) => await httpReq.post(`/workspace/switch/${id}`);

export const setWorkspaceRole = async (data) => await httpReq.post("/workspace/role", data);

export const getLoginInfo = async () => await httpReq.get('/user/loginInfo');
