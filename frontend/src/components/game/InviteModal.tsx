import React, { useState, useCallback } from 'react';
import { Copy, Check, Link } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, Button, Input, Label } from '../ui';

interface InviteModalProps {
  open: boolean;
  onClose: () => void;
  roomCode: string;
}

const InviteModal: React.FC<InviteModalProps> = ({ open, onClose, roomCode }) => {
  const [copied, setCopied] = useState(false);
  const [copyError, setCopyError] = useState(false);

  const inviteLink = `${window.location.origin}/join?code=${roomCode}`;

  const handleCopy = useCallback(async () => {
    try {
      await navigator.clipboard.writeText(inviteLink);
      setCopied(true);
      setCopyError(false);
      setTimeout(() => setCopied(false), 2000);
    } catch (err) {
      console.error('Failed to copy:', err);
      setCopyError(true);
      setTimeout(() => setCopyError(false), 3000);
    }
  }, [inviteLink]);

  const handleClose = () => {
    setCopied(false);
    setCopyError(false);
    onClose();
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && handleClose()}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Link className="w-5 h-5" />
            邀请伙伴
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label>房间码</Label>
            <div className="flex items-center justify-center p-4 bg-surface-elevated rounded-lg">
              <span className="text-3xl font-mono font-bold tracking-widest text-fg-primary">
                {roomCode}
              </span>
            </div>
          </div>

          <div className="space-y-2">
            <Label>邀请链接</Label>
            <div className="flex gap-2">
              <Input
                value={inviteLink}
                readOnly
                className="flex-1 font-mono text-sm"
              />
              <Button
                intent={copied ? 'primary' : 'secondary'}
                onClick={handleCopy}
                className="shrink-0"
              >
                {copied ? (
                  <>
                    <Check className="w-4 h-4" />
                    已复制
                  </>
                ) : (
                  <>
                    <Copy className="w-4 h-4" />
                    复制
                  </>
                )}
              </Button>
            </div>
            {copyError && (
              <p className="text-xs text-error">复制失败，请手动复制链接</p>
            )}
          </div>

          <p className="text-sm text-fg-secondary">
            将链接发送给朋友，点击即可加入房间
          </p>
        </div>

        <div className="flex justify-end mt-4">
          <Button intent="neutral" onClick={handleClose}>
            关闭
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default InviteModal;
