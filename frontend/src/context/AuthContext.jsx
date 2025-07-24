import React, { createContext, useState, useContext, useEffect, useCallback } from 'react';
import { getUserProfile } from '../api/userService.js';
import apiClient from '../api/client.js';
import { useSnackbar } from 'notistack';

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [isInitializing, setIsInitializing] = useState(true);
  const [authError, setAuthError] = useState(null); 

  const { enqueueSnackbar } = useSnackbar();
  
  const fetchUser = useCallback(async () => {
    setIsInitializing(true);
    setAuthError(null);
    try {
      const userData = await getUserProfile();
      setUser(userData);
    } catch (error) {
      setUser(null);

      if (error.code === 'ERR_NETWORK') {
        setAuthError('network');
        enqueueSnackbar('Ошибка соединения с сервером. Попробуйте позже.', { 
          variant: 'error',
        });
      } else if (error.response && error.response.status >= 500) {
        setAuthError('auth');
        enqueueSnackbar('На сервере произошла непредвиденная ошибка.', { 
          variant: 'error' 
        });
      } else{
        setAuthError('auth');
      }
    } finally {
      setIsInitializing(false);
    }
  }, [enqueueSnackbar]); 

  useEffect(() => {
    fetchUser();
  }, [fetchUser]);

  const login = useCallback(async () => {
    try {
      const response = await apiClient.get('/auth/login-url'); 
      window.location.href = response.data.login_url;
    } catch (error) {
      console.error('Failed to get login URL', error);
      enqueueSnackbar('Не удалось начать процесс входа. Проверьте соединение.', { variant: 'error' });
    }
  }, [enqueueSnackbar]); 

  const logout = useCallback(async () => {
    try {
      const response = await apiClient.get('/auth/logout-url'); 
      setUser(null);
      window.location.href = response.data.logout_url;
    } catch (error) {
      console.error('Failed to logout', error);
      enqueueSnackbar('Произошла ошибка при выходе из системы.', { variant: 'error' });
    }
  }, [enqueueSnackbar]);

  useEffect(() => {
    if (!isInitializing) {
      const preloader = document.getElementById('preloader');
      if (preloader) {
        preloader.classList.add('app-preloader--hidden');
      }
    }
  }, [isInitializing])

  const value = {
    user,
    isInitializing,
    isAuthenticated: !!user,
    isAdmin: !!user?.IsAdmin,
    login,
    logout,
    fetchUser,
    authError
  };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === null) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};