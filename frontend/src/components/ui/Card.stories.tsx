import type { Story } from "@ladle/react";
import { Card } from "./Card";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800">
    <Card>这是一个默认的 Card 组件</Card>
  </div>
);

export const AllVariants: Story = () => (
  <div className="p-8 bg-gray-800 flex gap-6">
    <div className="text-center">
      <Card variant="default">Default 样式</Card>
      <p className="mt-2 text-white text-sm">default</p>
    </div>
    <div className="text-center">
      <Card variant="glass">Glass 样式</Card>
      <p className="mt-2 text-white text-sm">glass</p>
    </div>
    <div className="text-center">
      <Card variant="gradient">Gradient 样式</Card>
      <p className="mt-2 text-white text-sm">gradient</p>
    </div>
  </div>
);

export const AllPaddings: Story = () => (
  <div className="p-8 bg-gray-800 flex items-start gap-6">
    <div className="text-center">
      <Card padding="none">
        <span className="bg-blue-200 px-1">none</span>
      </Card>
      <p className="mt-2 text-white text-sm">none</p>
    </div>
    <div className="text-center">
      <Card padding="sm">sm padding</Card>
      <p className="mt-2 text-white text-sm">sm</p>
    </div>
    <div className="text-center">
      <Card padding="md">md padding</Card>
      <p className="mt-2 text-white text-sm">md</p>
    </div>
    <div className="text-center">
      <Card padding="lg">lg padding</Card>
      <p className="mt-2 text-white text-sm">lg</p>
    </div>
  </div>
);

export const Combined: Story = () => (
  <div className="p-8 bg-gray-800 space-y-6">
    <div className="flex gap-4">
      <Card variant="default" padding="sm">default + sm</Card>
      <Card variant="glass" padding="md">glass + md</Card>
      <Card variant="gradient" padding="lg">gradient + lg</Card>
    </div>
    <div className="flex gap-4">
      <Card variant="glass" padding="none">
        <div className="bg-blue-100 p-2">glass + none</div>
      </Card>
      <Card variant="gradient" padding="sm">gradient + sm</Card>
    </div>
  </div>
);

export const WithCustomClassName: Story = () => (
  <div className="p-8 bg-gray-800">
    <Card className="border-2 border-blue-500 w-64">
      自定义 className 会与默认样式合并
    </Card>
  </div>
);
