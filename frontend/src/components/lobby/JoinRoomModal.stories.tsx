import type { Story } from "@ladle/react";
import { useState } from "react";
import JoinRoomModal from "./JoinRoomModal";

export const Default: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开加入房间模态框
      </button>
      <JoinRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onJoin={async (roomCode) => {
          await new Promise((resolve) => setTimeout(resolve, 1000));
          console.log("加入房间:", roomCode);
          setOpen(false);
        }}
      />
    </div>
  );
};

export const Closed: Story = () => {
  const [open, setOpen] = useState(false);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开加入房间模态框
      </button>
      <JoinRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onJoin={async (roomCode) => {
          await new Promise((resolve) => setTimeout(resolve, 1000));
          console.log("加入房间:", roomCode);
          setOpen(false);
        }}
      />
    </div>
  );
};

export const JoiningState: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开加入房间模态框
      </button>
      <p className="mt-4 text-gray-400 text-sm">
        输入4位数字房间码，点击"加入房间"按钮查看加载状态
      </p>
      <JoinRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onJoin={async (roomCode) => {
          await new Promise((resolve) => setTimeout(resolve, 3000));
          console.log("加入房间:", roomCode);
          setOpen(false);
        }}
      />
    </div>
  );
};

export const JoinFailure: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开加入房间模态框
      </button>
      <p className="mt-4 text-gray-400 text-sm">
        输入任意4位数字，点击"加入房间"会显示错误
      </p>
      <JoinRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onJoin={async () => {
          await new Promise((resolve) => setTimeout(resolve, 1000));
          throw new Error("房间不存在或已满员");
        }}
      />
    </div>
  );
};

export const Interactive: Story = () => {
  const [open, setOpen] = useState(false);
  const [status, setStatus] = useState<string>("等待操作");

  return (
    <div className="p-8">
      <button
        onClick={() => {
          setOpen(true);
          setStatus("模态框已打开");
        }}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开加入房间模态框
      </button>
      <p className="mt-4 text-gray-400">状态: {status}</p>
      <JoinRoomModal
        open={open}
        onClose={() => {
          setOpen(false);
          setStatus("用户取消加入");
        }}
        onJoin={async (roomCode) => {
          setStatus(`正在加入房间 ${roomCode}...`);
          await new Promise((resolve) => setTimeout(resolve, 1500));
          setOpen(false);
          setStatus(`成功加入房间 ${roomCode}!`);
        }}
      />
    </div>
  );
};
