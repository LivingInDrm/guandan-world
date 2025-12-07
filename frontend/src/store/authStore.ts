import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { apiClient } from '../services/api';
import type { User, AuthToken } from '../types';

interface AuthState {
  user: User | null;
  token: AuthToken | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  isInitialized: boolean;
}

interface AuthActions {
  setUser: (user: User) => void;
  setToken: (token: AuthToken) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  login: (user: User, token: AuthToken) => void;
  logout: () => void;
  clearError: () => void;
  initialize: () => void;
  checkTokenExpiry: () => boolean;
  refreshToken: () => Promise<boolean>;
}

type AuthStore = AuthState & AuthActions;

const isTokenExpired = (token: AuthToken | null): boolean => {
  if (!token) return true;
  
  const expiryTime = new Date(token.expires_at).getTime();
  const currentTime = Date.now();
  const bufferTime = 2 * 60 * 1000;
  
  return currentTime >= (expiryTime - bufferTime);
};

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      isInitialized: false,

      setUser: (user: User) => set({ user }),
      
      setToken: (token: AuthToken) => set({ token }),
      
      setLoading: (loading: boolean) => set({ isLoading: loading }),
      
      setError: (error: string | null) => set({ error }),
      
      login: (user: User, token: AuthToken) => {
        apiClient.setToken(token.access_token);
        
        set({
          user,
          token,
          isAuthenticated: true,
          error: null,
          isInitialized: true
        });
      },
      
      logout: () => {
        const state = get();
        
        apiClient.setToken(null);
        
        if (state.token?.refresh_token) {
          apiClient.logout(state.token.refresh_token).catch(() => {});
        }
        
        set({
          user: null,
          token: null,
          isAuthenticated: false,
          error: null
        });
      },
      
      clearError: () => set({ error: null }),
      
      initialize: () => {
        const state = get();
        
        if (state.token && state.user) {
          if (!isTokenExpired(state.token)) {
            apiClient.setToken(state.token.access_token);
            set({
              isAuthenticated: true,
              isInitialized: true
            });
          } else if (state.token.refresh_token) {
            get().refreshToken().then(success => {
              if (!success) {
                set({
                  user: null,
                  token: null,
                  isAuthenticated: false,
                  isInitialized: true
                });
              }
            });
            set({ isInitialized: true });
          } else {
            set({
              user: null,
              token: null,
              isAuthenticated: false,
              isInitialized: true
            });
          }
        } else {
          set({ isInitialized: true });
        }
      },
      
      checkTokenExpiry: () => {
        const state = get();
        return isTokenExpired(state.token);
      },
      
      refreshToken: async () => {
        const state = get();
        
        if (!state.token?.refresh_token) {
          return false;
        }
        
        try {
          const response = await apiClient.refreshToken(state.token.refresh_token);
          
          if (response.user && response.token) {
            apiClient.setToken(response.token.access_token);
            set({
              user: response.user,
              token: response.token,
              isAuthenticated: true,
              error: null
            });
            return true;
          }
          
          return false;
        } catch (error) {
          console.error('Token refresh failed:', error);
          get().logout();
          return false;
        }
      }
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        isAuthenticated: state.isAuthenticated
      }),
      onRehydrateStorage: () => (state) => {
        if (state) {
          state.initialize();
        }
      }
    }
  )
);
