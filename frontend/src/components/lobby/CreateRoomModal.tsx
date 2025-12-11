import React, { useState } from 'react';
import { Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, Button, Card } from '../ui';

interface CreateRoomModalProps {
  open: boolean;
  onClose: () => void;
  onConfirm: () => void;
}

const CreateRoomModal: React.FC<CreateRoomModalProps> = ({ open, onClose, onConfirm }) => {
  const [isCreating, setIsCreating] = useState(false);

  const handleConfirm = async () => {
    setIsCreating(true);
    try {
      await onConfirm();
    } finally {
      setIsCreating(false);
    }
  };

  const handleClose = () => {
    if (!isCreating) {
      onClose();
    }
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && handleClose()}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle>创建新房间</DialogTitle>
        </DialogHeader>

        <p className="text-fg-primary mb-4">
          确认创建新房间？您将成为房主，负责管理房间和开始游戏。
        </p>

        <Card variant="base" className="p-3 mb-6" interactive={false}>
          <h4 className="font-medium text-fg-primary mb-2">房间规则</h4>
          <ul className="text-sm text-fg-secondary space-y-1">
            <li>• 房间最多容纳4名玩家</li>
            <li>• 房主可以在人数满足时开始游戏</li>
            <li>• 房主离开时会自动转移给其他玩家</li>
            <li>• 所有玩家离开后房间自动关闭</li>
          </ul>
        </Card>

        <div className="flex justify-end gap-3">
          <Button intent="neutral" onClick={handleClose} disabled={isCreating}>
            取消
          </Button>
          <Button intent="primary" onClick={handleConfirm} disabled={isCreating}>
            {isCreating && <Loader2 className="animate-spin" />}
            {isCreating ? '创建中...' : '确认创建'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default CreateRoomModal;
