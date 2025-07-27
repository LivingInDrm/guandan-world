import React, { useEffect } from 'react';
import { useAuthStore } from '../../store/authStore';

interface AuthProviderProps {
  children: React.ReactNode;
}

const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const { 
    initialize, 
    checkTokenExpiry, 
    refreshToken, 
    isAuthenticated, 
    isInitialized 
  } = useAuthStore();

  useEffect(() => {
    // Initialize auth state on app start
    if (!isInitialized) {
      initialize();
    }
  }, [initialize, isInitialized]);

  useEffect(() => {
    if (!isAuthenticated) {
      return;
    }

    // Set up token expiry check interval
    const checkInterval = setInterval(() => {
      if (checkTokenExpiry()) {
        refreshToken();
      }
    }, 60000); // Check every minute

    // Initial check
    if (checkTokenExpiry()) {
      refreshToken();
    }

    return () => {
      clearInterval(checkInterval);
    };
  }, [isAuthenticated, checkTokenExpiry, refreshToken]);

  // Handle page visibility change to refresh token when user returns
  useEffect(() => {
    if (!isAuthenticated) {
      return;
    }

    const handleVisibilityChange = () => {
      if (document.visibilityState === 'visible') {
        if (checkTokenExpiry()) {
          refreshToken();
        }
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);

    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [isAuthenticated, checkTokenExpiry, refreshToken]);

  return <>{children}</>;
};

export default AuthProvider;