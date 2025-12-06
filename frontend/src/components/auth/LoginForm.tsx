import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import type { LoginRequest } from '../../types';
import { Input, Button } from '../ui';

const LoginForm: React.FC = () => {
  const [formData, setFormData] = useState<LoginRequest>({
    username: '',
    password: ''
  });
  const [errors, setErrors] = useState<Partial<LoginRequest>>({});
  const navigate = useNavigate();
  const { login, setLoading, setError, isLoading, error } = useAuthStore();

  const validateForm = (): boolean => {
    const newErrors: Partial<LoginRequest> = {};

    if (!formData.username.trim()) {
      newErrors.username = '用户名不能为空';
    } else if (formData.username.length < 3) {
      newErrors.username = '用户名至少3个字符';
    }

    if (!formData.password) {
      newErrors.password = '密码不能为空';
    } else if (formData.password.length < 6) {
      newErrors.password = '密码至少6个字符';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));
    
    // Clear field error when user starts typing
    if (errors[name as keyof LoginRequest]) {
      setErrors(prev => ({
        ...prev,
        [name]: undefined
      }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await apiClient.login(formData);
      
      if (response.success && response.data) {
        // Set token in API client
        apiClient.setToken(response.data.token.token);
        
        // Update auth store
        login(response.data.user, response.data.token);
        
        // Navigate to lobby
        navigate('/lobby');
      } else {
        setError(response.error || '登录失败');
      }
    } catch (err: any) {
      console.error('Login error:', err);
      setError(err.message || '登录失败，请检查网络连接');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <Input
        id="username"
        name="username"
        type="text"
        label="用户名"
        value={formData.username}
        onChange={handleInputChange}
        error={errors.username}
        placeholder="请输入用户名"
        disabled={isLoading}
      />

      <Input
        id="password"
        name="password"
        type="password"
        label="密码"
        value={formData.password}
        onChange={handleInputChange}
        error={errors.password}
        placeholder="请输入密码"
        disabled={isLoading}
      />

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-card p-3">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      <Button
        type="submit"
        variant="primary"
        fullWidth
        loading={isLoading}
      >
        {isLoading ? '登录中...' : '登录'}
      </Button>
    </form>
  );
};

export default LoginForm;