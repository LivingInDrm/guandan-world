import type { Story } from "@ladle/react";
import { useState } from "react";
import { Button } from "./Button";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Button>出牌</Button>
  </div>
);

export const AllIntents: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <Button intent="primary">出牌</Button>
    <Button intent="secondary">提示</Button>
    <Button intent="neutral">不出</Button>
    <Button intent="danger">退出</Button>
  </div>
);

export const AllSizes: Story = () => (
  <div className="p-8 bg-gray-800 flex items-end gap-4">
    <div className="text-center text-white">
      <Button size="sm">小</Button>
      <p className="mt-2 text-sm">sm</p>
    </div>
    <div className="text-center text-white">
      <Button size="md">中</Button>
      <p className="mt-2 text-sm">md</p>
    </div>
    <div className="text-center text-white">
      <Button size="lg">大</Button>
      <p className="mt-2 text-sm">lg</p>
    </div>
  </div>
);

export const IntentSizeMatrix: Story = () => (
  <div className="p-8 bg-gray-800">
    <table className="border-separate border-spacing-4">
      <thead>
        <tr className="text-white text-sm">
          <th></th>
          <th>sm</th>
          <th>md</th>
          <th>lg</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td className="text-white text-sm pr-4">primary</td>
          <td><Button intent="primary" size="sm">出牌</Button></td>
          <td><Button intent="primary" size="md">出牌</Button></td>
          <td><Button intent="primary" size="lg">出牌</Button></td>
        </tr>
        <tr>
          <td className="text-white text-sm pr-4">secondary</td>
          <td><Button intent="secondary" size="sm">提示</Button></td>
          <td><Button intent="secondary" size="md">提示</Button></td>
          <td><Button intent="secondary" size="lg">提示</Button></td>
        </tr>
        <tr>
          <td className="text-white text-sm pr-4">neutral</td>
          <td><Button intent="neutral" size="sm">不出</Button></td>
          <td><Button intent="neutral" size="md">不出</Button></td>
          <td><Button intent="neutral" size="lg">不出</Button></td>
        </tr>
        <tr>
          <td className="text-white text-sm pr-4">danger</td>
          <td><Button intent="danger" size="sm">退出</Button></td>
          <td><Button intent="danger" size="md">退出</Button></td>
          <td><Button intent="danger" size="lg">退出</Button></td>
        </tr>
      </tbody>
    </table>
  </div>
);

export const DisabledState: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <Button intent="primary" disabled>出牌</Button>
    <Button intent="secondary" disabled>提示</Button>
    <Button intent="neutral" disabled>不出</Button>
    <Button intent="danger" disabled>退出</Button>
  </div>
);

export const DangerVariant: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <Button intent="danger" size="sm">删除</Button>
    <Button intent="danger" size="md">退出游戏</Button>
    <Button intent="danger" size="lg">确认退出</Button>
  </div>
);

export const Interactive: Story = () => {
  const [count, setCount] = useState(0);

  return (
    <div className="p-8 bg-gray-800">
      <div className="flex items-center gap-4">
        <Button intent="primary" onClick={() => setCount((c) => c + 1)}>
          出牌次数: {count}
        </Button>
        <Button intent="neutral" onClick={() => setCount(0)}>
          重置
        </Button>
      </div>
    </div>
  );
};
