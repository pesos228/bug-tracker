import apiClient from './client';

export const getUserProfile = async () => {
  const response = await apiClient.get('/users/me');
  return response.data;
};

export const searchUsers = async (page = 1, pageSize = 10, fullName = '') => {
  const params = new URLSearchParams({
    page,
    pageSize,
    fullName,
  });
  const response = await apiClient.get(`/users?${params.toString()}`);
  return response.data;
};