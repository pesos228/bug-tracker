import apiClient from './client';

export const searchFolders = async (page = 1, pageSize = 10, query = '') => {
  const params = new URLSearchParams({
    page,
    pageSize,
    query,
  });

  const response = await apiClient.get(`/folders?${params.toString()}`);
  return response.data;
};

export const createFolder = async (name) => {
  const response = await apiClient.post('/folders', { name });
  return response.data;
};

export const getFolderDetails = async (folderId) => {
  const response = await apiClient.get(`/folders/${folderId}`);
  return response.data;
};

export const deleteFolder = async (folderId) => {
  const response = await apiClient.delete(`/folders/${folderId}`);
  return response.data;
};

export const exportFolder = async (folderId) => {
  const response = await apiClient.get(`/folders/${folderId}/reports`, {
    responseType: 'blob',
  });
  return response;
};