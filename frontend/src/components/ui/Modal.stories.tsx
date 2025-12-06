import type { Story } from "@ladle/react";
import { useState } from "react";
import { Modal } from "./Modal";

export const Default: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开模态框
      </button>
      <Modal open={open} onClose={() => setOpen(false)} title="默认模态框">
        <p className="text-table-400">这是一个默认尺寸 (md) 的模态框。</p>
      </Modal>
    </div>
  );
};

export const AllSizes: Story = () => {
  const [size, setSize] = useState<"sm" | "md" | "lg">("md");
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <div className="flex gap-2 mb-4">
        <button
          onClick={() => {
            setSize("sm");
            setOpen(true);
          }}
          className="px-3 py-1 bg-green-500 text-white rounded text-sm"
        >
          Small
        </button>
        <button
          onClick={() => {
            setSize("md");
            setOpen(true);
          }}
          className="px-3 py-1 bg-blue-500 text-white rounded text-sm"
        >
          Medium
        </button>
        <button
          onClick={() => {
            setSize("lg");
            setOpen(true);
          }}
          className="px-3 py-1 bg-purple-500 text-white rounded text-sm"
        >
          Large
        </button>
      </div>
      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title={`${size.toUpperCase()} 尺寸`}
        size={size}
      >
        <p className="text-table-400">
          当前尺寸: <strong>{size}</strong>
        </p>
        <p className="text-table-300 mt-2">
          sm = max-w-sm, md = max-w-lg, lg = max-w-2xl
        </p>
      </Modal>
    </div>
  );
};

export const WithTitle: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开带标题的模态框
      </button>
      <Modal open={open} onClose={() => setOpen(false)} title="带标题的模态框">
        <p className="text-table-400">标题栏包含关闭按钮。</p>
        <p className="text-table-300 mt-2">点击右上角 X 或按 ESC 可关闭。</p>
      </Modal>
    </div>
  );
};

export const WithoutTitle: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开无标题模态框
      </button>
      <Modal open={open} onClose={() => setOpen(false)}>
        <p className="text-table-400 text-center">这是一个无标题的模态框。</p>
        <p className="text-table-300 mt-2 text-center">
          点击遮罩层或按 ESC 可关闭。
        </p>
      </Modal>
    </div>
  );
};

export const Interactive: Story = () => {
  const [open, setOpen] = useState(false);
  const [count, setCount] = useState(0);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开模态框
      </button>
      <p className="mt-4 text-gray-600">模态框已打开次数: {count}</p>
      <Modal
        open={open}
        onClose={() => {
          setOpen(false);
          setCount((c) => c + 1);
        }}
        title="交互演示"
      >
        <p className="text-table-400">尝试以下操作关闭模态框:</p>
        <ul className="list-disc list-inside text-table-300 mt-2 space-y-1">
          <li>点击右上角关闭按钮</li>
          <li>点击遮罩层</li>
          <li>按 ESC 键</li>
        </ul>
      </Modal>
    </div>
  );
};

export const DisableOverlayClick: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开模态框
      </button>
      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="禁用遮罩点击"
        closeOnOverlayClick={false}
      >
        <p className="text-table-400">点击遮罩层不会关闭此模态框。</p>
        <p className="text-table-300 mt-2">只能通过关闭按钮或 ESC 键关闭。</p>
      </Modal>
    </div>
  );
};

export const DisableEscapeClose: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <button
        onClick={() => setOpen(true)}
        className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        打开模态框
      </button>
      <Modal
        open={open}
        onClose={() => setOpen(false)}
        title="禁用 ESC 关闭"
        closeOnEscape={false}
      >
        <p className="text-table-400">按 ESC 键不会关闭此模态框。</p>
        <p className="text-table-300 mt-2">只能通过关闭按钮或点击遮罩关闭。</p>
      </Modal>
    </div>
  );
};
