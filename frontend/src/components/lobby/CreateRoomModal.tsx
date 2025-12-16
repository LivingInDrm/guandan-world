import React, { useState } from 'react';
import { Loader2, Users, Crown, ArrowRightLeft, DoorClosed } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, Button, Card, CardContent } from '../ui';
import { cn } from '@/lib/utils';

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

  const rules = [
    { icon: Users, text: '房间最多容纳4名玩家' },
    { icon: Crown, text: '房主可以在人数满足时开始游戏' },
    { icon: ArrowRightLeft, text: '房主离开时会自动转移给其他玩家' },
    { icon: DoorClosed, text: '所有玩家离开后房间自动关闭' },
  ];

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && handleClose()}>
      <DialogContent className="max-w-sm mobile-landscape:max-w-xs">
        <DialogHeader>
          <DialogTitle className={cn(
            "font-display",
            "bg-gradient-to-r from-[hsl(158,55%,30%)] to-[hsl(158,55%,42%)]",
            "bg-clip-text text-transparent",
          )}>
            创建新房间
          </DialogTitle>
        </DialogHeader>

        <p className="text-fg-primary text-sm leading-relaxed">
          确认创建新房间？您将成为房主，负责管理房间和开始游戏。
        </p>

        <Card
          variant="base"
          interactive={false}
          className={cn(
            "mt-4 overflow-hidden",
            "bg-surface-elevated/50",
            "border border-stroke/30",
          )}
        >
          <CardContent className="p-4">
            <h4 className={cn(
              "font-display font-medium text-fg-primary mb-3",
              "flex items-center gap-2",
            )}>
              <span className="text-state-active">📋</span>
              房间规则
            </h4>
            <ul className="space-y-2.5">
              {rules.map((rule, index) => (
                <li
                  key={index}
                  className="flex items-start gap-2.5 text-sm text-fg-secondary"
                >
                  <rule.icon className="w-4 h-4 mt-0.5 text-[hsl(158,55%,42%)] shrink-0" />
                  <span>{rule.text}</span>
                </li>
              ))}
            </ul>
          </CardContent>
        </Card>

        <div className="flex justify-end gap-3 mt-6">
          <Button intent="neutral" onClick={handleClose} disabled={isCreating}>
            取消
          </Button>
          <Button
            intent="primary"
            onClick={handleConfirm}
            disabled={isCreating}
            className={cn(
              "bg-gradient-to-r from-[hsl(158,55%,35%)] to-[hsl(158,55%,42%)]",
              "hover:from-[hsl(158,55%,38%)] hover:to-[hsl(158,55%,45%)]",
              "shadow-[var(--elevation-2),var(--shadow-relief)]",
              "hover:shadow-[var(--elevation-3),var(--shadow-relief),var(--shadow-jade)]",
            )}
          >
            {isCreating && <Loader2 className="animate-spin" />}
            {isCreating ? '创建中...' : '确认创建'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default CreateRoomModal;
