import httpReq from '@/services/http.js';

export const getModelFactoryList = async (query) => await httpReq.get('/model/factory/list', query);

export const createModelFactory = async (data) => await httpReq.post('/model/factory/create', data);
