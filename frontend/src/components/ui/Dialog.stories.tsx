import type { Story } from "@ladle/react";
import { useState } from "react";
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter } from "./Dialog";
import { Button } from "./Button";

export const Default: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <Button onClick={() => setOpen(true)} variant="primary">
        打开对话框
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>默认对话框</DialogTitle>
            <DialogDescription>这是一个默认尺寸的对话框。</DialogDescription>
          </DialogHeader>
          <p className="text-table-400">对话框内容区域。</p>
          <DialogFooter>
            <Button variant="ghost" onClick={() => setOpen(false)}>取消</Button>
            <Button variant="primary" onClick={() => setOpen(false)}>确认</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export const WithForm: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <Button onClick={() => setOpen(true)} variant="primary">
        打开表单对话框
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>编辑信息</DialogTitle>
            <DialogDescription>在此处编辑您的个人信息。</DialogDescription>
          </DialogHeader>
          <div className="space-y-4 py-4">
            <div className="space-y-2">
              <label className="text-sm font-medium text-table-400">名称</label>
              <input className="w-full px-3 py-2 border rounded-sm" placeholder="请输入名称" />
            </div>
            <div className="space-y-2">
              <label className="text-sm font-medium text-table-400">描述</label>
              <textarea className="w-full px-3 py-2 border rounded-sm" placeholder="请输入描述" rows={3} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="ghost" onClick={() => setOpen(false)}>取消</Button>
            <Button variant="primary" onClick={() => setOpen(false)}>保存</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};

export const Confirmation: Story = () => {
  const [open, setOpen] = useState(true);

  return (
    <div className="p-8">
      <Button onClick={() => setOpen(true)} variant="danger">
        删除项目
      </Button>
      <Dialog open={open} onOpenChange={setOpen}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>确认删除</DialogTitle>
            <DialogDescription>此操作无法撤销。确定要删除此项目吗？</DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="ghost" onClick={() => setOpen(false)}>取消</Button>
            <Button variant="danger" onClick={() => setOpen(false)}>删除</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
};
