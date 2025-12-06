import type { Story } from "@ladle/react";
import { useState } from "react";
import { Button } from "./Button";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Button>默认按钮</Button>
  </div>
);

export const AllVariants: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <Button variant="primary">Primary</Button>
    <Button variant="secondary">Secondary</Button>
    <Button variant="warning">Warning</Button>
    <Button variant="danger">Danger</Button>
    <Button variant="ghost">Ghost</Button>
  </div>
);

export const AllSizes: Story = () => (
  <div className="p-8 bg-gray-800 flex items-end gap-4">
    <div className="text-center text-white">
      <Button size="sm">Small</Button>
      <p className="mt-2 text-sm">sm</p>
    </div>
    <div className="text-center text-white">
      <Button size="md">Medium</Button>
      <p className="mt-2 text-sm">md</p>
    </div>
    <div className="text-center text-white">
      <Button size="lg">Large</Button>
      <p className="mt-2 text-sm">lg</p>
    </div>
  </div>
);

export const LoadingState: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <Button loading>加载中...</Button>
    <Button variant="secondary" loading>加载中...</Button>
    <Button variant="warning" loading>加载中...</Button>
    <Button variant="danger" loading>加载中...</Button>
    <Button variant="ghost" loading>加载中...</Button>
  </div>
);

export const DisabledState: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <Button disabled>Primary</Button>
    <Button variant="secondary" disabled>Secondary</Button>
    <Button variant="warning" disabled>Warning</Button>
    <Button variant="danger" disabled>Danger</Button>
    <Button variant="ghost" disabled>Ghost</Button>
  </div>
);

export const FullWidth: Story = () => (
  <div className="p-8 bg-gray-800 w-80">
    <div className="space-y-4">
      <Button fullWidth>全宽按钮</Button>
      <Button variant="secondary" fullWidth>全宽按钮</Button>
    </div>
  </div>
);

export const Interactive: Story = () => {
  const [count, setCount] = useState(0);

  return (
    <div className="p-8 bg-gray-800">
      <div className="flex items-center gap-4">
        <Button onClick={() => setCount((c) => c + 1)}>点击次数: {count}</Button>
        <Button variant="danger" onClick={() => setCount(0)}>重置</Button>
      </div>
    </div>
  );
};
