import apiClient from './client';

export const getUserProfile = async () => {
  const response = await apiClient.get('/users/me');
  return response.data;
};