import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { BrowserRouter } from 'react-router-dom';
import RegisterForm from './RegisterForm';
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

const renderRegisterForm = () => {
  return render(
    <BrowserRouter>
      <RegisterForm />
    </BrowserRouter>
  );
};

describe('RegisterForm', () => {
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

  it('renders register form with all fields', () => {
    renderRegisterForm();
    
    expect(screen.getByLabelText('用户名')).toBeInTheDocument();
    expect(screen.getByLabelText('密码')).toBeInTheDocument();
    expect(screen.getByLabelText('确认密码')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: '注册' })).toBeInTheDocument();
  });

  it('shows validation errors for empty fields', async () => {
    const user = userEvent.setup();
    renderRegisterForm();
    
    const submitButton = screen.getByRole('button', { name: '注册' });
    await user.click(submitButton);
    
    expect(screen.getByText('用户名不能为空')).toBeInTheDocument();
    expect(screen.getByText('密码不能为空')).toBeInTheDocument();
    expect(screen.getByText('请确认密码')).toBeInTheDocument();
  });

  it('shows validation errors for short inputs', async () => {
    const user = userEvent.setup();
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'ab');
    await user.type(passwordInput, '12345');
    await user.type(confirmPasswordInput, '12345');
    await user.click(submitButton);
    
    expect(screen.getByText('用户名至少3个字符')).toBeInTheDocument();
    expect(screen.getByText('密码至少6个字符')).toBeInTheDocument();
  });

  it('shows error for mismatched passwords', async () => {
    const user = userEvent.setup();
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123');
    await user.type(confirmPasswordInput, 'password456');
    await user.click(submitButton);
    
    expect(screen.getByText('两次输入的密码不一致')).toBeInTheDocument();
  });

  it('shows error for invalid username characters', async () => {
    const user = userEvent.setup();
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'test@user');
    await user.type(passwordInput, 'password123');
    await user.type(confirmPasswordInput, 'password123');
    await user.click(submitButton);
    
    expect(screen.getByText('用户名只能包含字母、数字、下划线和中文')).toBeInTheDocument();
  });

  it('shows error for username too long', async () => {
    const user = userEvent.setup();
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'a'.repeat(21)); // 21 characters
    await user.type(passwordInput, 'password123');
    await user.type(confirmPasswordInput, 'password123');
    await user.click(submitButton);
    
    expect(screen.getByText('用户名不能超过20个字符')).toBeInTheDocument();
  });

  it('clears field errors when user starts typing', async () => {
    const user = userEvent.setup();
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    // Trigger validation error
    await user.click(submitButton);
    expect(screen.getByText('用户名不能为空')).toBeInTheDocument();
    
    // Start typing to clear error
    await user.type(usernameInput, 'test');
    expect(screen.queryByText('用户名不能为空')).not.toBeInTheDocument();
  });

  it('submits form with valid data and auto-logins', async () => {
    const user = userEvent.setup();
    const mockResponse = {
      success: true,
      data: {
        user: { id: '1', username: 'testuser', online: true },
        token: { token: 'mock-token', expires_at: '2024-01-01', user_id: '1' }
      }
    };
    
    mockApiClient.register.mockResolvedValue(mockResponse);
    mockApiClient.setToken = vi.fn();
    
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123');
    await user.type(confirmPasswordInput, 'password123');
    await user.click(submitButton);
    
    await waitFor(() => {
      expect(mockApiClient.register).toHaveBeenCalledWith({
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

  it('handles username already exists error', async () => {
    const user = userEvent.setup();
    const mockError = { status: 409, message: 'Username already exists' };
    
    mockApiClient.register.mockRejectedValue(mockError);
    
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'existinguser');
    await user.type(passwordInput, 'password123');
    await user.type(confirmPasswordInput, 'password123');
    await user.click(submitButton);
    
    await waitFor(() => {
      expect(mockSetError).toHaveBeenCalledWith('用户名已存在，请选择其他用户名');
    });
  });

  it('handles general registration error', async () => {
    const user = userEvent.setup();
    const mockError = new Error('Network error');
    
    mockApiClient.register.mockRejectedValue(mockError);
    
    renderRegisterForm();
    
    const usernameInput = screen.getByLabelText('用户名');
    const passwordInput = screen.getByLabelText('密码');
    const confirmPasswordInput = screen.getByLabelText('确认密码');
    const submitButton = screen.getByRole('button', { name: '注册' });
    
    await user.type(usernameInput, 'testuser');
    await user.type(passwordInput, 'password123');
    await user.type(confirmPasswordInput, 'password123');
    await user.click(submitButton);
    
    await waitFor(() => {
      expect(mockSetError).toHaveBeenCalledWith('Network error');
    });
  });

  it('shows loading state during submission', async () => {
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
    
    renderRegisterForm();
    
    expect(screen.getByText('注册中...')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /注册中/ })).toBeDisabled();
  });

  it('displays API error message', () => {
    mockUseAuthStore.mockReturnValue({
      user: null,
      token: null,
      isAuthenticated: false,
      isLoading: false,
      error: '注册失败，请重试',
      login: mockLogin,
      setLoading: mockSetLoading,
      setError: mockSetError,
      logout: vi.fn(),
      setUser: vi.fn(),
      setToken: vi.fn(),
      clearError: vi.fn(),
    });
    
    renderRegisterForm();
    
    expect(screen.getByText('注册失败，请重试')).toBeInTheDocument();
  });
});