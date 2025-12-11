import type { Story } from "@ladle/react";
import { Card } from "./Card";
import { Button } from "./Button";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Card className="p-6">
      <p className="text-ds-text-primary">基础面板</p>
    </Card>
  </div>
);

export const AllVariants: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-6">
    <div className="text-center">
      <Card variant="base" className="p-6 w-40">
        <p className="text-ds-text-primary">Base</p>
      </Card>
      <p className="mt-2 text-white text-sm">base</p>
    </div>
    <div className="text-center">
      <Card variant="elevated" className="p-6 w-40">
        <p className="text-ds-text-primary">Elevated</p>
      </Card>
      <p className="mt-2 text-white text-sm">elevated</p>
    </div>
    <div className="text-center">
      <Card variant="emphasis" className="p-6 w-40">
        <p className="text-ds-text-primary">Emphasis</p>
      </Card>
      <p className="mt-2 text-white text-sm">emphasis</p>
    </div>
  </div>
);

export const Nested: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Card variant="emphasis" className="p-6">
      <p className="text-ds-text-primary mb-4">强调层</p>
      <Card variant="elevated" className="p-4">
        <p className="text-ds-text-primary mb-4">抬升层</p>
        <Card variant="base" className="p-4">
          <p className="text-ds-text-primary">基础层</p>
        </Card>
      </Card>
    </Card>
  </div>
);

export const WithContent: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Card variant="elevated" className="p-6 w-80">
      <h3 className="text-ds-text-primary text-lg font-bold mb-2">游戏设置</h3>
      <p className="text-ds-text-secondary text-sm mb-4">
        调整游戏规则和玩法
      </p>
      <div className="flex gap-2">
        <Button intent="neutral" size="sm">取消</Button>
        <Button intent="primary" size="sm">确认</Button>
      </div>
    </Card>
  </div>
);

export const GlassEffect: Story = () => (
  <div 
    className="p-8 flex items-center justify-center min-h-[300px]"
    style={{
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
    }}
  >
    <Card variant="emphasis" className="p-8 w-72">
      <h3 className="text-ds-text-primary text-lg font-bold mb-2">毛玻璃效果</h3>
      <p className="text-ds-text-secondary text-sm">
        emphasis 变体带有 backdrop-blur 效果，在彩色背景上可以看到毛玻璃效果。
      </p>
    </Card>
  </div>
);

export const Interactive: Story = () => (
  <div className="p-8 bg-gray-800 flex gap-6">
    <div className="text-center">
      <Card variant="elevated" className="p-6 w-48">
        <p className="text-ds-text-primary">Interactive</p>
        <p className="text-ds-text-secondary text-xs mt-1">hover 会放大</p>
      </Card>
      <p className="mt-2 text-white text-sm">interactive=true (默认)</p>
    </div>
    <div className="text-center">
      <Card variant="elevated" className="p-6 w-48" interactive={false}>
        <p className="text-ds-text-primary">Non-Interactive</p>
        <p className="text-ds-text-secondary text-xs mt-1">hover 无效果</p>
      </Card>
      <p className="mt-2 text-white text-sm">interactive=false</p>
    </div>
  </div>
);
