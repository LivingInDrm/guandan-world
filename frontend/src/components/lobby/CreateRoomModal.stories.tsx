import type { Story } from "@ladle/react";
import { useState } from "react";
import CreateRoomModal from "./CreateRoomModal";

export const Default: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开创建房间模态框
      </button>
      <CreateRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onConfirm={async () => {
          await new Promise((resolve) => setTimeout(resolve, 1000));
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
        打开创建房间模态框
      </button>
      <CreateRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onConfirm={async () => {
          await new Promise((resolve) => setTimeout(resolve, 1000));
          setOpen(false);
        }}
      />
    </div>
  );
};

export const CreatingState: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开创建房间模态框
      </button>
      <p className="mt-4 text-gray-400 text-sm">
        点击"确认创建"按钮查看加载状态
      </p>
      <CreateRoomModal
        open={open}
        onClose={() => setOpen(false)}
        onConfirm={async () => {
          await new Promise((resolve) => setTimeout(resolve, 3000));
          setOpen(false);
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
        打开创建房间模态框
      </button>
      <p className="mt-4 text-gray-400">状态: {status}</p>
      <CreateRoomModal
        open={open}
        onClose={() => {
          setOpen(false);
          setStatus("用户取消创建");
        }}
        onConfirm={async () => {
          setStatus("创建中...");
          await new Promise((resolve) => setTimeout(resolve, 1500));
          setOpen(false);
          setStatus("房间创建成功!");
        }}
      />
    </div>
  );
};
