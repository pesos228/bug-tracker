import apiClient from './client';

export const getTasksByFolderId = async (folderId, { page = 1, pageSize = 10, checkStatus = '', requestID = '' }) => {
  const params = new URLSearchParams({
    page,
    pageSize,
    checkStatus,
    requestID,
  });

  const response = await apiClient.get(`/folders/${folderId}/tasks?${params.toString()}`);
  return response.data;
};

export const getTaskDetails = async (taskId, view = 'short') => {
  const response = await apiClient.get(`/tasks/${taskId}?view=${view}`);
  return response.data;
};

export const updateTaskByAdmin = async (taskId, data) => {
  const response = await apiClient.patch(`/tasks/${taskId}`, data);
  return response.data;
};

export const updateTaskByUser = async (taskId, data) => {
  const response = await apiClient.patch(`/tasks/${taskId}/review`, data);
  return response.data;
};

export const deleteTask = async (taskId) => {
  const response = await apiClient.delete(`/tasks/${taskId}`);
  return response.data;
};

export const createTask = async (folderId, taskData) => {
  const response = await apiClient.post(`/folders/${folderId}/tasks`, taskData);
  return response.data;
};

export const getMyTasks = async ({ page = 1, pageSize = 10, checkStatus = '', requestID = '' }) => {
  const params = new URLSearchParams({
    page,
    pageSize,
    checkStatus,
    requestID,
  });

  const response = await apiClient.get(`/tasks/my?${params.toString()}`);
  return response.data;
};