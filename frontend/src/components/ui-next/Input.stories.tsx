import type { Story } from "@ladle/react";
import { Input, Label } from "./Input";
import { Button } from "./Button";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Input placeholder="请输入内容..." className="w-64" />
  </div>
);

export const WithLabel: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <div className="space-y-2 w-64">
      <Label htmlFor="nickname">昵称</Label>
      <Input id="nickname" placeholder="请输入昵称" />
    </div>
  </div>
);

export const Disabled: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <div className="space-y-4 w-64">
      <div className="space-y-2">
        <Label htmlFor="normal">正常状态</Label>
        <Input id="normal" placeholder="可编辑" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="disabled">禁用状态</Label>
        <Input id="disabled" placeholder="不可编辑" disabled />
      </div>
    </div>
  </div>
);

export const AllTypes: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <div className="space-y-4 w-64">
      <div className="space-y-2">
        <Label htmlFor="text">文本</Label>
        <Input id="text" type="text" placeholder="文本输入" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="password">密码</Label>
        <Input id="password" type="password" placeholder="密码输入" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="number">数字</Label>
        <Input id="number" type="number" placeholder="数字输入" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="email">邮箱</Label>
        <Input id="email" type="email" placeholder="example@email.com" />
      </div>
    </div>
  </div>
);

export const FormExample: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <form className="space-y-4 w-72 p-6 bg-ds-surface-elevated rounded-ds-lg">
      <h3 className="text-ds-text-primary text-lg font-bold">登录</h3>
      <div className="space-y-2">
        <Label htmlFor="username">用户名</Label>
        <Input id="username" placeholder="请输入用户名" />
      </div>
      <div className="space-y-2">
        <Label htmlFor="pwd">密码</Label>
        <Input id="pwd" type="password" placeholder="请输入密码" />
      </div>
      <Button intent="primary" className="w-full">登录</Button>
    </form>
  </div>
);
