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

// Helper function to check if token is expired
const isTokenExpired = (token: AuthToken | null): boolean => {
  if (!token) return true;
  
  const expiryTime = new Date(token.expires_at).getTime();
  const currentTime = Date.now();
  const bufferTime = 5 * 60 * 1000; // 5 minutes buffer
  
  return currentTime >= (expiryTime - bufferTime);
};

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // State
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      isInitialized: false,

      // Actions
      setUser: (user: User) => set({ user }),
      
      setToken: (token: AuthToken) => set({ token }),
      
      setLoading: (loading: boolean) => set({ isLoading: loading }),
      
      setError: (error: string | null) => set({ error }),
      
      login: (user: User, token: AuthToken) => {
        // Set token in API client
        apiClient.setToken(token.token);
        
        set({
          user,
          token,
          isAuthenticated: true,
          error: null,
          isInitialized: true
        });
      },
      
      logout: () => {
        // Clear token from API client
        apiClient.setToken(null);
        
        // Call logout API if possible
        try {
          apiClient.logout().catch(() => {
            // Ignore logout API errors - we're logging out anyway
          });
        } catch (error) {
          // Ignore errors
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
            // Token is still valid, set it in API client
            apiClient.setToken(state.token.token);
            set({
              isAuthenticated: true,
              isInitialized: true
            });
          } else {
            // Token is expired, clear auth state
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
        
        if (!state.token || !state.user) {
          return false;
        }
        
        try {
          // For now, we don't have a refresh token endpoint
          // In a real app, you would call a refresh endpoint here
          // For this implementation, we'll just check if the token is still valid
          
          if (isTokenExpired(state.token)) {
            // Token is expired, logout
            get().logout();
            return false;
          }
          
          return true;
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
        // Initialize auth state after rehydration
        if (state) {
          state.initialize();
        }
      }
    }
  )
);