import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Loader2 } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import type { RegisterRequest } from '../../types';
import { Input, Label, Button, Card } from '../ui-next';

const PASSWORD_MIN_LENGTH = 8;
const PASSWORD_MAX_LENGTH = 50;
const USERNAME_MIN_LENGTH = 4;
const USERNAME_MAX_LENGTH = 20;

interface RegisterFormData extends RegisterRequest {
  confirmPassword: string;
}

interface RegisterFormProps {
  showCard?: boolean;
  title?: string;
}

const RegisterForm: React.FC<RegisterFormProps> = ({ 
  showCard = true, 
  title = '注册账号' 
}) => {
  const [formData, setFormData] = useState<RegisterFormData>({
    username: '',
    password: '',
    confirmPassword: ''
  });
  const [errors, setErrors] = useState<Partial<RegisterFormData>>({});
  const navigate = useNavigate();
  const { login, setLoading, setError, isLoading, error } = useAuthStore();

  const validateForm = (): boolean => {
    const newErrors: Partial<RegisterFormData> = {};

    if (!formData.username.trim()) {
      newErrors.username = '用户名不能为空';
    } else if (formData.username.length < USERNAME_MIN_LENGTH) {
      newErrors.username = `用户名至少${USERNAME_MIN_LENGTH}个字符`;
    } else if (formData.username.length > USERNAME_MAX_LENGTH) {
      newErrors.username = `用户名不能超过${USERNAME_MAX_LENGTH}个字符`;
    } else if (!/^[a-zA-Z0-9_]+$/.test(formData.username)) {
      newErrors.username = '用户名只能包含字母、数字和下划线';
    }

    if (!formData.password) {
      newErrors.password = '密码不能为空';
    } else if (formData.password.length < PASSWORD_MIN_LENGTH) {
      newErrors.password = `密码至少${PASSWORD_MIN_LENGTH}个字符`;
    } else if (formData.password.length > PASSWORD_MAX_LENGTH) {
      newErrors.password = `密码不能超过${PASSWORD_MAX_LENGTH}个字符`;
    } else if (!/[A-Z]/.test(formData.password)) {
      newErrors.password = '密码必须包含至少一个大写字母';
    } else if (!/[a-z]/.test(formData.password)) {
      newErrors.password = '密码必须包含至少一个小写字母';
    } else if (!/[0-9]/.test(formData.password)) {
      newErrors.password = '密码必须包含至少一个数字';
    } else if (!/[!@#$%^&*()_+=[\]{}|;:,.<>?-]/.test(formData.password)) {
      newErrors.password = '密码必须包含至少一个特殊字符';
    }

    if (!formData.confirmPassword) {
      newErrors.confirmPassword = '请确认密码';
    } else if (formData.password !== formData.confirmPassword) {
      newErrors.confirmPassword = '两次输入的密码不一致';
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
    
    if (errors[name as keyof RegisterFormData]) {
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
      const registerData: RegisterRequest = {
        username: formData.username,
        password: formData.password
      };

      const response = await apiClient.register(registerData);
      
      if (response.user && response.token) {
        apiClient.setToken(response.token.access_token);
        login(response.user, response.token);
        navigate('/lobby');
      } else {
        setError('注册失败');
      }
    } catch (err: unknown) {
      console.error('Register error:', err);
      const apiErr = err as { status?: number; message?: string };
      if (apiErr.status === 409) {
        setError('用户名已存在，请选择其他用户名');
      } else {
        setError(apiErr.message || '注册失败，请检查网络连接');
      }
    } finally {
      setLoading(false);
    }
  };

  const formContent = (
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
          className={errors.username ? 'border-ds-error' : ''}
        />
        {errors.username ? (
          <p className="text-sm text-ds-error">{errors.username}</p>
        ) : (
          <p className="text-xs text-ds-text-secondary">{USERNAME_MIN_LENGTH}-{USERNAME_MAX_LENGTH}个字符，支持字母、数字和下划线</p>
        )}
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
          className={errors.password ? 'border-ds-error' : ''}
        />
        {errors.password ? (
          <p className="text-sm text-ds-error">{errors.password}</p>
        ) : (
          <p className="text-xs text-ds-text-secondary">至少{PASSWORD_MIN_LENGTH}个字符，包含大小写字母、数字和特殊字符</p>
        )}
      </div>

      <div className="space-y-2">
        <Label htmlFor="confirmPassword">确认密码</Label>
        <Input
          id="confirmPassword"
          name="confirmPassword"
          type="password"
          value={formData.confirmPassword}
          onChange={handleInputChange}
          placeholder="请再次输入密码"
          disabled={isLoading}
          className={errors.confirmPassword ? 'border-ds-error' : ''}
        />
        {errors.confirmPassword && <p className="text-sm text-ds-error">{errors.confirmPassword}</p>}
      </div>

      {error && (
        <div className="bg-ds-error/15 border border-ds-error rounded-ds-sm p-3">
          <p className="text-sm text-ds-error">{error}</p>
        </div>
      )}

      <Button type="submit" intent="primary" disabled={isLoading} className="w-full">
        {isLoading && <Loader2 className="animate-spin" />}
        {isLoading ? '注册中...' : '注册'}
      </Button>
    </form>
  );

  if (!showCard) {
    return formContent;
  }

  return (
    <Card variant="elevated" interactive={false} className="p-8">
      <h2 className="text-2xl font-bold text-ds-text-primary mb-6 text-center">
        {title}
      </h2>
      {formContent}
    </Card>
  );
};

export default RegisterForm;
