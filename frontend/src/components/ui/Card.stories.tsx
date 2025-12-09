import type { Story } from "@ladle/react";
import { Card, CardHeader, CardTitle, CardDescription, CardContent, CardFooter } from "./Card";
import { Button } from "./Button";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800">
    <Card className="p-5">这是一个默认的 Card 组件</Card>
  </div>
);

export const AllVariants: Story = () => (
  <div className="p-8 bg-gray-800 flex gap-6">
    <div className="text-center">
      <Card variant="default" className="p-5">Default 样式</Card>
      <p className="mt-2 text-white text-sm">default</p>
    </div>
    <div className="text-center">
      <Card variant="glass" className="p-5">Glass 样式</Card>
      <p className="mt-2 text-white text-sm">glass</p>
    </div>
    <div className="text-center">
      <Card variant="gradient" className="p-5">Gradient 样式</Card>
      <p className="mt-2 text-white text-sm">gradient</p>
    </div>
  </div>
);

export const WithPadding: Story = () => (
  <div className="p-8 bg-gray-800 flex items-start gap-6">
    <div className="text-center">
      <Card>
        <span className="bg-blue-200 px-1">无 padding</span>
      </Card>
      <p className="mt-2 text-white text-sm">none</p>
    </div>
    <div className="text-center">
      <Card className="p-3">p-3 padding</Card>
      <p className="mt-2 text-white text-sm">p-3</p>
    </div>
    <div className="text-center">
      <Card className="p-5">p-5 padding</Card>
      <p className="mt-2 text-white text-sm">p-5</p>
    </div>
    <div className="text-center">
      <Card className="p-8">p-8 padding</Card>
      <p className="mt-2 text-white text-sm">p-8</p>
    </div>
  </div>
);

export const WithSubComponents: Story = () => (
  <div className="p-8 bg-gray-800">
    <Card className="w-80">
      <CardHeader>
        <CardTitle>卡片标题</CardTitle>
        <CardDescription>这是卡片的描述文字</CardDescription>
      </CardHeader>
      <CardContent>
        <p>卡片内容区域，可以放置任何内容。</p>
      </CardContent>
      <CardFooter className="gap-2">
        <Button variant="ghost" size="sm">取消</Button>
        <Button variant="primary" size="sm">确认</Button>
      </CardFooter>
    </Card>
  </div>
);

export const Combined: Story = () => (
  <div className="p-8 bg-gray-800 space-y-6">
    <div className="flex gap-4">
      <Card variant="default" className="p-3">default + p-3</Card>
      <Card variant="glass" className="p-5">glass + p-5</Card>
      <Card variant="gradient" className="p-8">gradient + p-8</Card>
    </div>
    <div className="flex gap-4">
      <Card variant="glass">
        <div className="bg-blue-100 p-2">glass + 无 padding</div>
      </Card>
      <Card variant="gradient" className="p-3">gradient + p-3</Card>
    </div>
  </div>
);

export const WithCustomClassName: Story = () => (
  <div className="p-8 bg-gray-800">
    <Card className="border-2 border-blue-500 w-64 p-5">
      自定义 className 会与默认样式合并
    </Card>
  </div>
);
