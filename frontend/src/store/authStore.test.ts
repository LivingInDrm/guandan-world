import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { useAuthStore } from './authStore';
import { apiClient } from '../services/api';
import type { User, AuthToken } from '../types';

// Mock the API client
vi.mock('../services/api');
const mockApiClient = vi.mocked(apiClient);

// Mock localStorage for zustand persist
const localStorageMock = {
  getItem: vi.fn(),
  setItem: vi.fn(),
  removeItem: vi.fn(),
  clear: vi.fn(),
};
Object.defineProperty(window, 'localStorage', {
  value: localStorageMock
});

describe('AuthStore', () => {
  const mockUser: User = {
    id: '1',
    username: 'testuser',
    online: true
  };

  const mockToken: AuthToken = {
    token: 'mock-jwt-token',
    expires_at: new Date(Date.now() + 3600000).toISOString(), // 1 hour from now
    user_id: '1'
  };

  const mockExpiredToken: AuthToken = {
    token: 'expired-jwt-token',
    expires_at: new Date(Date.now() - 3600000).toISOString(), // 1 hour ago
    user_id: '1'
  };

  beforeEach(() => {
    vi.clearAllMocks();
    localStorageMock.getItem.mockReturnValue(null);
    
    // Reset store state
    useAuthStore.setState({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      isInitialized: false
    });
  });

  afterEach(() => {
    vi.clearAllTimers();
  });

  describe('login', () => {
    it('should set user and token on login', () => {
      const store = useAuthStore.getState();
      
      store.login(mockUser, mockToken);
      
      const state = useAuthStore.getState();
      expect(state.user).toEqual(mockUser);
      expect(state.token).toEqual(mockToken);
      expect(state.isAuthenticated).toBe(true);
      expect(state.error).toBeNull();
      expect(state.isInitialized).toBe(true);
      expect(mockApiClient.setToken).toHaveBeenCalledWith(mockToken.token);
    });
  });

  describe('logout', () => {
    it('should clear user and token on logout', async () => {
      const store = useAuthStore.getState();
      
      // First login
      store.login(mockUser, mockToken);
      expect(useAuthStore.getState().isAuthenticated).toBe(true);
      
      // Mock logout API call
      mockApiClient.logout.mockResolvedValue({ success: true });
      
      // Then logout
      store.logout();
      
      const state = useAuthStore.getState();
      expect(state.user).toBeNull();
      expect(state.token).toBeNull();
      expect(state.isAuthenticated).toBe(false);
      expect(state.error).toBeNull();
      expect(mockApiClient.setToken).toHaveBeenCalledWith(null);
    });

    it('should handle logout API errors gracefully', () => {
      const store = useAuthStore.getState();
      
      // First login
      store.login(mockUser, mockToken);
      
      // Mock logout API to fail
      mockApiClient.logout.mockRejectedValue(new Error('Network error'));
      
      // Logout should still work
      store.logout();
      
      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(false);
    });
  });

  describe('initialize', () => {
    it('should initialize with valid token from storage', () => {
      const store = useAuthStore.getState();
      
      // Set initial state as if loaded from storage
      useAuthStore.setState({
        user: mockUser,
        token: mockToken,
        isAuthenticated: false,
        isInitialized: false
      });
      
      store.initialize();
      
      const state = useAuthStore.getState();
      expect(state.isAuthenticated).toBe(true);
      expect(state.isInitialized).toBe(true);
      expect(mockApiClient.setToken).toHaveBeenCalledWith(mockToken.token);
    });

    it('should clear expired token on initialize', () => {
      const store = useAuthStore.getState();
      
      // Set initial state with expired token
      useAuthStore.setState({
        user: mockUser,
        token: mockExpiredToken,
        isAuthenticated: false,
        isInitialized: false
      });
      
      store.initialize();
      
      const state = useAuthStore.getState();
      expect(state.user).toBeNull();
      expect(state.token).toBeNull();
      expect(state.isAuthenticated).toBe(false);
      expect(state.isInitialized).toBe(true);
    });

    it('should initialize without token', () => {
      const store = useAuthStore.getState();
      
      store.initialize();
      
      const state = useAuthStore.getState();
      expect(state.isInitialized).toBe(true);
      expect(state.isAuthenticated).toBe(false);
    });
  });

  describe('checkTokenExpiry', () => {
    it('should return false for valid token', () => {
      const store = useAuthStore.getState();
      
      useAuthStore.setState({
        token: mockToken
      });
      
      const isExpired = store.checkTokenExpiry();
      expect(isExpired).toBe(false);
    });

    it('should return true for expired token', () => {
      const store = useAuthStore.getState();
      
      useAuthStore.setState({
        token: mockExpiredToken
      });
      
      const isExpired = store.checkTokenExpiry();
      expect(isExpired).toBe(true);
    });

    it('should return true for null token', () => {
      const store = useAuthStore.getState();
      
      const isExpired = store.checkTokenExpiry();
      expect(isExpired).toBe(true);
    });

    it('should return true for token expiring within buffer time', () => {
      const store = useAuthStore.getState();
      
      // Token expires in 2 minutes (less than 5 minute buffer)
      const soonToExpireToken: AuthToken = {
        token: 'soon-to-expire-token',
        expires_at: new Date(Date.now() + 2 * 60 * 1000).toISOString(),
        user_id: '1'
      };
      
      useAuthStore.setState({
        token: soonToExpireToken
      });
      
      const isExpired = store.checkTokenExpiry();
      expect(isExpired).toBe(true);
    });
  });

  describe('refreshToken', () => {
    it('should return true for valid token', async () => {
      const store = useAuthStore.getState();
      
      useAuthStore.setState({
        user: mockUser,
        token: mockToken
      });
      
      const result = await store.refreshToken();
      expect(result).toBe(true);
    });

    it('should logout and return false for expired token', async () => {
      const store = useAuthStore.getState();
      const logoutSpy = vi.spyOn(store, 'logout');
      
      useAuthStore.setState({
        user: mockUser,
        token: mockExpiredToken
      });
      
      const result = await store.refreshToken();
      expect(result).toBe(false);
      expect(logoutSpy).toHaveBeenCalled();
    });

    it('should return false when no token exists', async () => {
      const store = useAuthStore.getState();
      
      const result = await store.refreshToken();
      expect(result).toBe(false);
    });
  });

  describe('state management', () => {
    it('should set loading state', () => {
      const store = useAuthStore.getState();
      
      store.setLoading(true);
      expect(useAuthStore.getState().isLoading).toBe(true);
      
      store.setLoading(false);
      expect(useAuthStore.getState().isLoading).toBe(false);
    });

    it('should set error state', () => {
      const store = useAuthStore.getState();
      
      store.setError('Test error');
      expect(useAuthStore.getState().error).toBe('Test error');
      
      store.clearError();
      expect(useAuthStore.getState().error).toBeNull();
    });

    it('should set user', () => {
      const store = useAuthStore.getState();
      
      store.setUser(mockUser);
      expect(useAuthStore.getState().user).toEqual(mockUser);
    });

    it('should set token', () => {
      const store = useAuthStore.getState();
      
      store.setToken(mockToken);
      expect(useAuthStore.getState().token).toEqual(mockToken);
    });
  });
});