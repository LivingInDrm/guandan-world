import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import LoginForm from './LoginForm';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';

// Mock the auth store
vi.mock('../../store/authStore');
const mockUseAuthStore = vi.mocked(useAuthStore);

// Mock the API client
vi.mock('../../services/api');
const mockApiClient = vi.mocked(apiClient);

// Mock navigate
const mockNavigate = vi.fn();
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual('react-router-dom');
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

const renderLoginForm = () => {
  return render(
    <BrowserRouter>
      <LoginForm />
    </BrowserRouter>
  );
};

describe('LoginForm', () => {
  const mockLogin = vi.fn();
  const mockSetLoading = vi.fn();
  const mockSetError = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    
    mockUseAuthStore.mockReturnValue({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      login: mockLogin,
      setLoading: mockSetLoading,
      setError: mockSetError,
      logout: vi.fn(),
      setUser: vi.fn(),
      setToken: vi.fn(),
      clearError: vi.fn(),
    });
  });

  it('renders login form with all fields', () => {
    renderLoginForm();
    
    expect(screen.getByLabelText('用户名')).toBeInTheDocument();
    expect(screen.getByLabelText('密码')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '登录' })).toBeInTheDocument();
  });

  it('shows validation errors for empty fields', async () => {
    const user = userEvent.setup();
    renderLoginForm();
    
    const submitButton = screen.getByRole('button', { name: '登录' });
    await user.click(submitButton);
    
    expect(screen.getByText('用户名不能为空')).toBeInTheDocument();
    expect(screen.getByText('密码不能为空')).toBeInTheDocument();
  });

  it('shows validation errors for short inputs', async () => {
    const user = userEvent.setup();
    renderLoginForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const submitButton = screen.getByRole('button', { name: '登录' });
    
    await user.type(usernameInput, 'ab');
    await user.type(passwordInput, '12345');
    await user.click(submitButton);
    
    expect(screen.getByText('用户名至少3个字符')).toBeInTheDocument();
    expect(screen.getByText('密码至少6个字符')).toBeInTheDocument();
  });

  it('clears field errors when user starts typing', async () => {
    const user = userEvent.setup();
    renderLoginForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const submitButton = screen.getByRole('button', { name: '登录' });
    
    // Trigger validation error
    await user.click(submitButton);
    expect(screen.getByText('用户名不能为空')).toBeInTheDocument();
    
    // Start typing to clear error
    await user.type(usernameInput, 'test');
    expect(screen.queryByText('用户名不能为空')).not.toBeInTheDocument();
  });

  it('submits form with valid data', async () => {
    const user = userEvent.setup();
    const mockResponse = {
      success: true,
      data: {
        user: { id: '1', username: 'testuser', online: true },
        token: { token: 'mock-token', expires_at: '2024-01-01', user_id: '1' }
      }
    };
    
    mockApiClient.login.mockResolvedValue(mockResponse);
    mockApiClient.setToken = vi.fn();
    
    renderLoginForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const submitButton = screen.getByRole('button', { name: '登录' });
    
    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123');
    await user.click(submitButton);
    
    await waitFor(() => {
      expect(mockApiClient.login).toHaveBeenCalledWith({
        username: 'testuser',
        password: 'password123'
      });
    });
    
    expect(mockApiClient.setToken).toHaveBeenCalledWith('mock-token');
    expect(mockLogin).toHaveBeenCalledWith(
      mockResponse.data.user,
      mockResponse.data.token
    );
    expect(mockNavigate).toHaveBeenCalledWith('/lobby');
  });

  it('handles login API error', async () => {
    const user = userEvent.setup();
    const mockError = new Error('Invalid credentials');
    
    mockApiClient.login.mockRejectedValue(mockError);
    
    renderLoginForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const submitButton = screen.getByRole('button', { name: '登录' });
    
    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'wrongpassword');
    await user.click(submitButton);
    
    await waitFor(() => {
      expect(mockSetError).toHaveBeenCalledWith('Invalid credentials');
    });
  });

  it('shows loading state during submission', async () => {
    const user = userEvent.setup();
    
    // Mock loading state
    mockUseAuthStore.mockReturnValue({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: true,
      error: null,
      login: mockLogin,
      setLoading: mockSetLoading,
      setError: mockSetError,
      logout: vi.fn(),
      setUser: vi.fn(),
      setToken: vi.fn(),
      clearError: vi.fn(),
    });
    
    renderLoginForm();
    
    expect(screen.getByText('登录中...')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /登录中/ })).toBeDisabled();
  });

  it('displays API error message', () => {
    mockUseAuthStore.mockReturnValue({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: '用户名或密码错误',
      login: mockLogin,
      setLoading: mockSetLoading,
      setError: mockSetError,
      logout: vi.fn(),
      setUser: vi.fn(),
      setToken: vi.fn(),
      clearError: vi.fn(),
    });
    
    renderLoginForm();
    
    expect(screen.getByText('用户名或密码错误')).toBeInTheDocument();
  });
});