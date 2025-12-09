import type { Story } from "@ladle/react";
import { useState } from "react";
import { Input, Label } from "./Input";

export const Default: Story = () => (
  <div className="p-8 max-w-md">
    <Input placeholder="请输入内容" />
  </div>
);

export const WithLabel: Story = () => (
  <div className="p-8 max-w-md space-y-2">
    <Label htmlFor="username">用户名</Label>
    <Input id="username" placeholder="请输入用户名" />
  </div>
);

export const WithError: Story = () => (
  <div className="p-8 max-w-md space-y-2">
    <Label htmlFor="email">邮箱</Label>
    <Input id="email" placeholder="请输入邮箱" className="border-destructive" />
    <p className="text-sm text-destructive">邮箱格式不正确</p>
  </div>
);

export const WithHelperText: Story = () => (
  <div className="p-8 max-w-md space-y-6">
    <div className="space-y-2">
      <Label htmlFor="username2">用户名</Label>
      <Input id="username2" placeholder="请输入用户名" />
      <p className="text-xs text-table-300">3-20个字符，支持字母、数字、下划线和中文</p>
    </div>
    <div className="space-y-2">
      <Label htmlFor="password">密码</Label>
      <Input id="password" type="password" placeholder="请输入密码" />
      <p className="text-xs text-table-300">至少6个字符</p>
    </div>
  </div>
);

export const HelperTextWithError: Story = () => (
  <div className="p-8 max-w-md space-y-2">
    <Label htmlFor="username3">用户名</Label>
    <Input id="username3" placeholder="请输入用户名" className="border-destructive" />
    <p className="text-sm text-destructive">用户名已存在</p>
    <p className="mt-4 text-xs text-gray-500">
      当存在错误时，显示错误提示而非辅助文本
    </p>
  </div>
);

export const Disabled: Story = () => (
  <div className="p-8 max-w-md space-y-6">
    <div className="space-y-2">
      <Label htmlFor="disabled1">禁用状态</Label>
      <Input id="disabled1" placeholder="不可编辑" disabled />
    </div>
    <div className="space-y-2">
      <Label htmlFor="disabled2">禁用带值</Label>
      <Input id="disabled2" value="已有内容" disabled />
    </div>
  </div>
);

export const AllStates: Story = () => (
  <div className="p-8 max-w-md space-y-6">
    <div className="space-y-2">
      <Label htmlFor="normal">正常状态</Label>
      <Input id="normal" placeholder="请输入内容" />
    </div>
    <div className="space-y-2">
      <Label htmlFor="error">错误状态</Label>
      <Input id="error" placeholder="请输入内容" className="border-destructive" />
      <p className="text-sm text-destructive">这是错误提示</p>
    </div>
    <div className="space-y-2">
      <Label htmlFor="disabled3">禁用状态</Label>
      <Input id="disabled3" placeholder="不可编辑" disabled />
    </div>
  </div>
);

export const Controlled: Story = () => {
  const [value, setValue] = useState("");

  return (
    <div className="p-8 max-w-md space-y-4">
      <div className="space-y-2">
        <Label htmlFor="controlled">受控输入</Label>
        <Input
          id="controlled"
          value={value}
          onChange={(e) => setValue(e.target.value)}
          placeholder="输入内容会在下方显示"
        />
      </div>
      <p className="text-sm text-gray-600">当前值: {value || "(空)"}</p>
    </div>
  );
};

export const InputTypes: Story = () => (
  <div className="p-8 max-w-md space-y-6">
    <div className="space-y-2">
      <Label htmlFor="text">文本</Label>
      <Input id="text" type="text" placeholder="文本输入" />
    </div>
    <div className="space-y-2">
      <Label htmlFor="pwd">密码</Label>
      <Input id="pwd" type="password" placeholder="密码输入" />
    </div>
    <div className="space-y-2">
      <Label htmlFor="num">数字</Label>
      <Input id="num" type="number" placeholder="数字输入" />
    </div>
    <div className="space-y-2">
      <Label htmlFor="mail">邮箱</Label>
      <Input id="mail" type="email" placeholder="邮箱输入" />
    </div>
  </div>
);
