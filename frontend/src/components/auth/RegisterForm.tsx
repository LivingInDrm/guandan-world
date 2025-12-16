import React, { useState } from 'react';
import { Loader2, User, Lock, ShieldCheck } from 'lucide-react';
import { useAuthStore } from '../../store/authStore';
import { apiClient } from '../../services/api';
import type { RegisterRequest } from '../../types';
import { Input, Label, Button, Card, CardHeader, CardTitle, CardContent } from '../ui';
import { cn } from '@/lib/utils';

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
    <form onSubmit={handleSubmit} className="space-y-5">
      {/* 用户名输入 */}
      <div className="space-y-2">
        <Label htmlFor="username" className="text-fg-primary font-medium">
          用户名
        </Label>
        <div className="relative">
          <User className={cn(
            "absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4",
            "text-fg-secondary/60 transition-colors duration-fast",
            errors.username && "text-action-danger",
          )} />
          <Input
            id="username"
            name="username"
            type="text"
            value={formData.username}
            onChange={handleInputChange}
            placeholder="请输入用户名"
            disabled={isLoading}
            className={cn(
              "pl-10",
              errors.username && "border-action-danger focus:ring-action-danger/30",
            )}
          />
        </div>
        {errors.username ? (
          <p className="text-sm text-action-danger flex items-center gap-1">
            <span className="inline-block w-1 h-1 rounded-full bg-action-danger" />
            {errors.username}
          </p>
        ) : (
          <p className="text-xs text-fg-secondary/70">
            {USERNAME_MIN_LENGTH}-{USERNAME_MAX_LENGTH}个字符，支持字母、数字和下划线
          </p>
        )}
      </div>

      {/* 密码输入 */}
      <div className="space-y-2">
        <Label htmlFor="password" className="text-fg-primary font-medium">
          密码
        </Label>
        <div className="relative">
          <Lock className={cn(
            "absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4",
            "text-fg-secondary/60 transition-colors duration-fast",
            errors.password && "text-action-danger",
          )} />
          <Input
            id="password"
            name="password"
            type="password"
            value={formData.password}
            onChange={handleInputChange}
            placeholder="请输入密码"
            disabled={isLoading}
            className={cn(
              "pl-10",
              errors.password && "border-action-danger focus:ring-action-danger/30",
            )}
          />
        </div>
        {errors.password ? (
          <p className="text-sm text-action-danger flex items-center gap-1">
            <span className="inline-block w-1 h-1 rounded-full bg-action-danger" />
            {errors.password}
          </p>
        ) : (
          <p className="text-xs text-fg-secondary/70">
            至少{PASSWORD_MIN_LENGTH}个字符，包含大小写字母、数字和特殊字符
          </p>
        )}
      </div>

      {/* 确认密码输入 */}
      <div className="space-y-2">
        <Label htmlFor="confirmPassword" className="text-fg-primary font-medium">
          确认密码
        </Label>
        <div className="relative">
          <ShieldCheck className={cn(
            "absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4",
            "text-fg-secondary/60 transition-colors duration-fast",
            errors.confirmPassword && "text-action-danger",
          )} />
          <Input
            id="confirmPassword"
            name="confirmPassword"
            type="password"
            value={formData.confirmPassword}
            onChange={handleInputChange}
            placeholder="请再次输入密码"
            disabled={isLoading}
            className={cn(
              "pl-10",
              errors.confirmPassword && "border-action-danger focus:ring-action-danger/30",
            )}
          />
        </div>
        {errors.confirmPassword && (
          <p className="text-sm text-action-danger flex items-center gap-1">
            <span className="inline-block w-1 h-1 rounded-full bg-action-danger" />
            {errors.confirmPassword}
          </p>
        )}
      </div>

      {/* 全局错误提示 */}
      {error && (
        <div className={cn(
          "p-4 rounded-lg",
          "bg-action-danger/10 border border-action-danger/30",
          "animate-in fade-in-0 slide-in-from-top-2 duration-normal",
        )}>
          <p className="text-sm text-action-danger text-center">{error}</p>
        </div>
      )}

      {/* 提交按钮 */}
      <Button
        type="submit"
        intent="primary"
        disabled={isLoading}
        className={cn(
          "w-full mt-2",
          "bg-gradient-to-r from-[hsl(158,55%,35%)] to-[hsl(158,55%,42%)]",
          "hover:from-[hsl(158,55%,38%)] hover:to-[hsl(158,55%,45%)]",
          "shadow-[var(--elevation-2),var(--shadow-relief)]",
          "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-jade)]",
        )}
      >
        {isLoading && <Loader2 className="animate-spin" />}
        {isLoading ? '注册中...' : '注册'}
      </Button>
    </form>
  );

  if (!showCard) {
    return formContent;
  }

  return (
    <Card
      variant="elevated"
      interactive={false}
      className={cn(
        "overflow-hidden",
        "bg-surface-base/95 backdrop-blur-sm",
        "shadow-[var(--elevation-3),var(--shadow-relief)]",
        "border border-stroke/50",
      )}
    >
      {/* 顶部装饰线 */}
      <div className="h-1 bg-gradient-to-r from-transparent via-state-active to-transparent" />

      <CardHeader className="pb-2 pt-6">
        <CardTitle className={cn(
          "text-2xl font-display font-bold text-center",
          "bg-gradient-to-r from-[hsl(158,55%,30%)] to-[hsl(158,55%,42%)]",
          "bg-clip-text text-transparent",
        )}>
          {title}
        </CardTitle>
      </CardHeader>

      <CardContent className="px-8 pb-8">
        {formContent}
      </CardContent>
    </Card>
  );
};

export default RegisterForm;
