import type { Story } from "@ladle/react";
import { useState } from "react";
import { Input } from "./Input";

export const Default: Story = () => (
  <div className="p-8 max-w-md">
    <Input placeholder="请输入内容" />
  </div>
);

export const WithLabel: Story = () => (
  <div className="p-8 max-w-md">
    <Input label="用户名" placeholder="请输入用户名" />
  </div>
);

export const WithError: Story = () => (
  <div className="p-8 max-w-md">
    <Input label="邮箱" placeholder="请输入邮箱" error="邮箱格式不正确" />
  </div>
);

export const WithHelperText: Story = () => (
  <div className="p-8 max-w-md space-y-4">
    <Input
      label="用户名"
      placeholder="请输入用户名"
      helperText="3-20个字符，支持字母、数字、下划线和中文"
    />
    <Input
      label="密码"
      type="password"
      placeholder="请输入密码"
      helperText="至少6个字符"
    />
  </div>
);

export const HelperTextWithError: Story = () => (
  <div className="p-8 max-w-md">
    <Input
      label="用户名"
      placeholder="请输入用户名"
      helperText="3-20个字符"
      error="用户名已存在"
    />
    <p className="mt-4 text-xs text-gray-500">
      当存在错误时，helperText 不显示
    </p>
  </div>
);

export const Disabled: Story = () => (
  <div className="p-8 max-w-md space-y-4">
    <Input label="禁用状态" placeholder="不可编辑" disabled />
    <Input label="禁用带值" value="已有内容" disabled />
  </div>
);

export const NotFullWidth: Story = () => (
  <div className="p-8">
    <Input label="固定宽度" placeholder="不占满宽度" fullWidth={false} />
  </div>
);

export const AllStates: Story = () => (
  <div className="p-8 max-w-md space-y-6">
    <Input label="正常状态" placeholder="请输入内容" />
    <Input label="错误状态" placeholder="请输入内容" error="这是错误提示" />
    <Input label="禁用状态" placeholder="不可编辑" disabled />
  </div>
);

export const Controlled: Story = () => {
  const [value, setValue] = useState("");

  return (
    <div className="p-8 max-w-md space-y-4">
      <Input
        label="受控输入"
        value={value}
        onChange={(e) => setValue(e.target.value)}
        placeholder="输入内容会在下方显示"
      />
      <p className="text-sm text-gray-600">当前值: {value || "(空)"}</p>
    </div>
  );
};

export const InputTypes: Story = () => (
  <div className="p-8 max-w-md space-y-4">
    <Input label="文本" type="text" placeholder="文本输入" />
    <Input label="密码" type="password" placeholder="密码输入" />
    <Input label="数字" type="number" placeholder="数字输入" />
    <Input label="邮箱" type="email" placeholder="邮箱输入" />
  </div>
);
