import React, { useState } from 'react';
import { Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, Button, Input, Label } from '../ui';

interface JoinRoomModalProps {
  open: boolean;
  onClose: () => void;
  onJoin: (roomCode: string) => Promise<void>;
}

const JoinRoomModal: React.FC<JoinRoomModalProps> = ({ open, onClose, onJoin }) => {
  const [roomCode, setRoomCode] = useState('');
  const [isJoining, setIsJoining] = useState(false);
  const [error, setError] = useState('');

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.replace(/\D/g, '').slice(0, 4);
    setRoomCode(value);
    setError('');
  };

  const handleSubmit = async () => {
    if (roomCode.length !== 4) {
      setError('请输入4位数字房间码');
      return;
    }

    setIsJoining(true);
    setError('');
    try {
      await onJoin(roomCode);
      setRoomCode('');
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : '加入房间失败');
    } finally {
      setIsJoining(false);
    }
  };

  const handleClose = () => {
    if (!isJoining) {
      setRoomCode('');
      setError('');
      onClose();
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && roomCode.length === 4 && !isJoining) {
      handleSubmit();
    }
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && handleClose()}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle>加入游戏</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="room-code">房间码</Label>
            <Input
              id="room-code"
              type="text"
              inputMode="numeric"
              pattern="[0-9]*"
              maxLength={4}
              placeholder="请输入4位数字房间码"
              value={roomCode}
              onChange={handleInputChange}
              onKeyDown={handleKeyDown}
              disabled={isJoining}
              autoFocus
              className="text-center text-2xl tracking-widest"
            />
          </div>

          {error && (
            <p className="text-error text-sm">{error}</p>
          )}
        </div>

        <div className="flex justify-end gap-3 mt-4">
          <Button intent="neutral" onClick={handleClose} disabled={isJoining}>
            取消
          </Button>
          <Button
            intent="primary"
            onClick={handleSubmit}
            disabled={roomCode.length !== 4 || isJoining}
          >
            {isJoining && <Loader2 className="animate-spin" />}
            {isJoining ? '加入中...' : '加入房间'}
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
};

export default JoinRoomModal;
