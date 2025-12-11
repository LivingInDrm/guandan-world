import type { Story } from "@ladle/react";
import { Avatar } from "./Avatar";

const sampleAvatarUrl = "https://api.dicebear.com/7.x/avataaars/svg?seed=Felix";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Avatar src={sampleAvatarUrl} alt="用户头像" fallback="用户" />
  </div>
);

export const AllSizes: Story = () => (
  <div className="p-8 bg-gray-800 flex items-end gap-4">
    <div className="text-center">
      <Avatar size="sm" src={sampleAvatarUrl} fallback="SM" />
      <p className="mt-2 text-white text-xs">sm</p>
    </div>
    <div className="text-center">
      <Avatar size="md" src={sampleAvatarUrl} fallback="MD" />
      <p className="mt-2 text-white text-xs">md</p>
    </div>
    <div className="text-center">
      <Avatar size="lg" src={sampleAvatarUrl} fallback="LG" />
      <p className="mt-2 text-white text-xs">lg</p>
    </div>
    <div className="text-center">
      <Avatar size="xl" src={sampleAvatarUrl} fallback="XL" />
      <p className="mt-2 text-white text-xs">xl</p>
    </div>
    <div className="text-center">
      <Avatar size="2xl" src={sampleAvatarUrl} fallback="2X" />
      <p className="mt-2 text-white text-xs">2xl</p>
    </div>
  </div>
);

export const AllRingStates: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-6">
    <div className="text-center">
      <Avatar ringState="none" src={sampleAvatarUrl} fallback="无" />
      <p className="mt-2 text-white text-xs">none</p>
    </div>
    <div className="text-center">
      <Avatar ringState="normal" src={sampleAvatarUrl} fallback="普" />
      <p className="mt-2 text-white text-xs">normal</p>
    </div>
    <div className="text-center">
      <Avatar ringState="active" src={sampleAvatarUrl} fallback="活" />
      <p className="mt-2 text-white text-xs">active (脉冲)</p>
    </div>
    <div className="text-center">
      <Avatar ringState="teamUs" src={sampleAvatarUrl} fallback="我" />
      <p className="mt-2 text-white text-xs">teamUs</p>
    </div>
    <div className="text-center">
      <Avatar ringState="teamThem" src={sampleAvatarUrl} fallback="敌" />
      <p className="mt-2 text-white text-xs">teamThem</p>
    </div>
  </div>
);

export const SizeRingMatrix: Story = () => (
  <div className="p-8 bg-gray-800">
    <table className="border-separate border-spacing-4">
      <thead>
        <tr className="text-white text-xs">
          <th></th>
          <th>sm</th>
          <th>md</th>
          <th>lg</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td className="text-white text-xs pr-4">none</td>
          <td><Avatar size="sm" ringState="none" src={sampleAvatarUrl} /></td>
          <td><Avatar size="md" ringState="none" src={sampleAvatarUrl} /></td>
          <td><Avatar size="lg" ringState="none" src={sampleAvatarUrl} /></td>
        </tr>
        <tr>
          <td className="text-white text-xs pr-4">active</td>
          <td><Avatar size="sm" ringState="active" src={sampleAvatarUrl} /></td>
          <td><Avatar size="md" ringState="active" src={sampleAvatarUrl} /></td>
          <td><Avatar size="lg" ringState="active" src={sampleAvatarUrl} /></td>
        </tr>
        <tr>
          <td className="text-white text-xs pr-4">teamUs</td>
          <td><Avatar size="sm" ringState="teamUs" src={sampleAvatarUrl} /></td>
          <td><Avatar size="md" ringState="teamUs" src={sampleAvatarUrl} /></td>
          <td><Avatar size="lg" ringState="teamUs" src={sampleAvatarUrl} /></td>
        </tr>
      </tbody>
    </table>
  </div>
);

export const WithFallback: Story = () => (
  <div className="p-8 bg-gray-800 flex flex-wrap gap-4">
    <div className="text-center">
      <Avatar fallback="张三" />
      <p className="mt-2 text-white text-xs">张三 → 张三</p>
    </div>
    <div className="text-center">
      <Avatar fallback="李四五" />
      <p className="mt-2 text-white text-xs">李四五 → 李四</p>
    </div>
    <div className="text-center">
      <Avatar fallback="user" />
      <p className="mt-2 text-white text-xs">user → US</p>
    </div>
    <div className="text-center">
      <Avatar fallback="A" />
      <p className="mt-2 text-white text-xs">A → A</p>
    </div>
  </div>
);

export const ActiveAnimation: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <div className="text-center">
      <Avatar 
        size="xl" 
        ringState="active" 
        src={sampleAvatarUrl} 
        fallback="活跃" 
      />
      <p className="mt-4 text-white text-sm">当前回合玩家 - 脉冲动画</p>
    </div>
  </div>
);
