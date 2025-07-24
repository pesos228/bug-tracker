import axios from 'axios';

const apiClient = axios.create({
  baseURL: '/api',
});

apiClient.interceptors.response.use(
  (response) => {
    return response;
  },
  (error) => {
    if (!error.response) {
      console.error('Network Error:', error);
      return Promise.reject(error);
    }

    const { status } = error.response;

    if (status === 401) {
      const publicPaths = ['/', '/login'];
      if (!publicPaths.includes(window.location.pathname) && !window.location.pathname.startsWith('/login')) {
        sessionStorage.setItem('redirectAfterLogin', window.location.pathname);
        window.location.href = '/login';
      }
    }

    return Promise.reject(error);
  }
);

export default apiClient;