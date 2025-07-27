import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import ProtectedRoute from './ProtectedRoute';
import { useAuthStore } from '../../store/authStore';

// Mock the auth store
vi.mock('../../store/authStore');
const mockUseAuthStore = vi.mocked(useAuthStore);

// Mock react-router-dom Navigate component
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    Navigate: ({ to, state, replace }: { to: string; state?: any; replace?: boolean }) => {
      mockNavigate(to, state, replace);
      return <div data-testid="navigate">Redirecting to {to}</div>;
    },
  };
});

const TestComponent = () => <div>Protected Content</div>;

const renderProtectedRoute = (initialEntries = ['/protected']) => {
  return render(
    <MemoryRouter initialEntries={initialEntries}>
      <ProtectedRoute>
        <TestComponent />
      </ProtectedRoute>
    </MemoryRouter>
  );
};

describe('ProtectedRoute', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockNavigate.mockClear();
  });

  it('should show loading when not initialized', () => {
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: false,
      isInitialized: false,
      isLoading: false,
      user: null,
      token: null,
      error: null,
      setUser: vi.fn(),
      setToken: vi.fn(),
      setLoading: vi.fn(),
      setError: vi.fn(),
      login: vi.fn(),
      logout: vi.fn(),
      clearError: vi.fn(),
      initialize: vi.fn(),
      checkTokenExpiry: vi.fn(),
      refreshToken: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByText('加载中...')).toBeInTheDocument();
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });

  it('should show loading when loading is true', () => {
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: false,
      isInitialized: true,
      isLoading: true,
      user: null,
      token: null,
      error: null,
      setUser: vi.fn(),
      setToken: vi.fn(),
      setLoading: vi.fn(),
      setError: vi.fn(),
      login: vi.fn(),
      logout: vi.fn(),
      clearError: vi.fn(),
      initialize: vi.fn(),
      checkTokenExpiry: vi.fn(),
      refreshToken: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByText('加载中...')).toBeInTheDocument();
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
  });

  it('should redirect to login when not authenticated', () => {
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: false,
      isInitialized: true,
      isLoading: false,
      user: null,
      token: null,
      error: null,
      setUser: vi.fn(),
      setToken: vi.fn(),
      setLoading: vi.fn(),
      setError: vi.fn(),
      login: vi.fn(),
      logout: vi.fn(),
      clearError: vi.fn(),
      initialize: vi.fn(),
      checkTokenExpiry: vi.fn(),
      refreshToken: vi.fn(),
    });

    renderProtectedRoute();

    // The component should not render the protected content
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    
    // Should show redirect message
    expect(screen.getByTestId('navigate')).toBeInTheDocument();
    expect(screen.getByText('Redirecting to /login')).toBeInTheDocument();
  });

  it('should render children when authenticated', () => {
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: true,
      isInitialized: true,
      isLoading: false,
      user: { id: '1', username: 'testuser', online: true },
      token: { token: 'mock-token', expires_at: '2024-01-01', user_id: '1' },
      error: null,
      setUser: vi.fn(),
      setToken: vi.fn(),
      setLoading: vi.fn(),
      setError: vi.fn(),
      login: vi.fn(),
      logout: vi.fn(),
      clearError: vi.fn(),
      initialize: vi.fn(),
      checkTokenExpiry: vi.fn(),
      refreshToken: vi.fn(),
    });

    renderProtectedRoute();

    expect(screen.getByText('Protected Content')).toBeInTheDocument();
  });

  it('should use custom redirect path', () => {
    mockUseAuthStore.mockReturnValue({
      isAuthenticated: false,
      isInitialized: true,
      isLoading: false,
      user: null,
      token: null,
      error: null,
      setUser: vi.fn(),
      setToken: vi.fn(),
      setLoading: vi.fn(),
      setError: vi.fn(),
      login: vi.fn(),
      logout: vi.fn(),
      clearError: vi.fn(),
      initialize: vi.fn(),
      checkTokenExpiry: vi.fn(),
      refreshToken: vi.fn(),
    });

    render(
      <MemoryRouter initialEntries={['/protected']}>
        <ProtectedRoute redirectTo="/custom-login">
          <TestComponent />
        </ProtectedRoute>
      </MemoryRouter>
    );

    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument();
    expect(screen.getByText('Redirecting to /custom-login')).toBeInTheDocument();
  });
});