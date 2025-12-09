import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import type { LoginRequest } from '../../types';
import { Input, Label, Button } from '../ui';

const PASSWORD_MIN_LENGTH = 8;
const USERNAME_MIN_LENGTH = 4;

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
    } else if (formData.username.length < USERNAME_MIN_LENGTH) {
      newErrors.username = `用户名至少${USERNAME_MIN_LENGTH}个字符`;
    }

    if (!formData.password) {
      newErrors.password = '密码不能为空';
    } else if (formData.password.length < PASSWORD_MIN_LENGTH) {
      newErrors.password = `密码至少${PASSWORD_MIN_LENGTH}个字符`;
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
      
      if (response.user && response.token) {
        apiClient.setToken(response.token.access_token);
        login(response.user, response.token);
        navigate('/lobby');
      } else {
        setError('登录失败');
      }
    } catch (err: unknown) {
      console.error('Login error:', err);
      const message = err instanceof Error ? err.message : '登录失败，请检查网络连接';
      setError(message);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      <div className="space-y-2">
        <Label htmlFor="username">用户名</Label>
        <Input
          id="username"
          name="username"
          type="text"
          value={formData.username}
          onChange={handleInputChange}
          placeholder="请输入用户名"
          disabled={isLoading}
          className={errors.username ? 'border-destructive' : ''}
        />
        {errors.username && <p className="text-sm text-destructive">{errors.username}</p>}
      </div>

      <div className="space-y-2">
        <Label htmlFor="password">密码</Label>
        <Input
          id="password"
          name="password"
          type="password"
          value={formData.password}
          onChange={handleInputChange}
          placeholder="请输入密码"
          disabled={isLoading}
          className={errors.password ? 'border-destructive' : ''}
        />
        {errors.password && <p className="text-sm text-destructive">{errors.password}</p>}
      </div>

      {error && (
        <div className="bg-red-50 border border-red-200 rounded-sm p-3">
          <p className="text-sm text-red-600">{error}</p>
        </div>
      )}

      <Button type="submit" variant="primary" disabled={isLoading} className="w-full">
        {isLoading && <Loader2 className="animate-spin" />}
        {isLoading ? '登录中...' : '登录'}
      </Button>
    </form>
  );
};

export default LoginForm;
