import type { Story } from "@ladle/react";
import { useState } from "react";
import Countdown from "./Countdown";

export const Default: Story = () => (
  <div className="p-8 bg-gray-800 flex items-center justify-center">
    <Countdown deadlineAtMs={Date.now() + 20000} />
  </div>
);

export const AllSizes: Story = () => (
  <div className="p-8 bg-gray-800 flex items-end gap-8">
    <div className="text-center text-white">
      <Countdown deadlineAtMs={Date.now() + 15000} size="small" />
      <p className="mt-2 text-sm">small</p>
    </div>
    <div className="text-center text-white">
      <Countdown deadlineAtMs={Date.now() + 15000} size="medium" />
      <p className="mt-2 text-sm">medium</p>
    </div>
    <div className="text-center text-white">
      <Countdown deadlineAtMs={Date.now() + 15000} size="large" />
      <p className="mt-2 text-sm">large</p>
    </div>
  </div>
);

export const WarningColors: Story = () => (
  <div className="p-8 bg-gray-800">
    <p className="text-white mb-4">观察颜色变化: 白色 → 橙色(≤10s) → 红色(≤5s)</p>
    <Countdown deadlineAtMs={Date.now() + 12000} size="large" />
  </div>
);

export const WithCallback: Story = () => {
  const [message, setMessage] = useState("");
  const [key, setKey] = useState(0);

  const handleTimeout = () => {
    setMessage("超时回调已触发!");
  };

  const handleReset = () => {
    setMessage("");
    setKey((k) => k + 1);
  };

  return (
    <div className="p-8 bg-gray-800">
      <div className="flex items-center gap-4 mb-4">
        <Countdown
          key={key}
          deadlineAtMs={Date.now() + 5000}
          size="large"
          onTimeout={handleTimeout}
        />
        <button
          onClick={handleReset}
          className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          重置 (5秒)
        </button>
      </div>
      {message && (
        <p className="text-green-400 text-lg font-bold">{message}</p>
      )}
    </div>
  );
};

export const Inactive: Story = () => (
  <div className="p-8 bg-gray-800">
    <p className="text-white mb-4">isActive=false 时显示初始值但不倒计时</p>
    <Countdown deadlineAtMs={Date.now() + 30000} isActive={false} size="large" />
  </div>
);

export const AlreadyExpired: Story = () => (
  <div className="p-8 bg-gray-800">
    <p className="text-white mb-4">deadline 已过期，显示 0</p>
    <Countdown deadlineAtMs={Date.now() - 5000} size="large" />
  </div>
);

export const DynamicDeadline: Story = () => {
  const [deadline, setDeadline] = useState(Date.now() + 10000);

  return (
    <div className="p-8 bg-gray-800">
      <div className="flex items-center gap-4 mb-4">
        <Countdown deadlineAtMs={deadline} size="large" />
        <div className="flex gap-2">
          <button
            onClick={() => setDeadline(Date.now() + 5000)}
            className="px-3 py-1 bg-red-500 text-white rounded text-sm"
          >
            5秒
          </button>
          <button
            onClick={() => setDeadline(Date.now() + 10000)}
            className="px-3 py-1 bg-orange-500 text-white rounded text-sm"
          >
            10秒
          </button>
          <button
            onClick={() => setDeadline(Date.now() + 20000)}
            className="px-3 py-1 bg-green-500 text-white rounded text-sm"
          >
            20秒
          </button>
        </div>
      </div>
      <p className="text-gray-400 text-sm">点击按钮重置倒计时</p>
    </div>
  );
};
