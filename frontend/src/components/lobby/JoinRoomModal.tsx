import React, { useState, useRef, useEffect } from 'react';
import { Loader2 } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, Button } from '../ui';

interface JoinRoomModalProps {
  open: boolean;
  onClose: () => void;
  onJoin: (roomCode: string) => Promise<void>;
}

const JoinRoomModal: React.FC<JoinRoomModalProps> = ({ open, onClose, onJoin }) => {
  const [digits, setDigits] = useState<string[]>(['', '', '', '']);
  const [isJoining, setIsJoining] = useState(false);
  const [error, setError] = useState('');
  const inputRefs = useRef<(HTMLInputElement | null)[]>([]);

  const roomCode = digits.join('');

  useEffect(() => {
    if (open) {
      setDigits(['', '', '', '']);
      setError('');
      navigator.clipboard.readText().then(text => {
        const match = text.match(/[?&]code=(\d{4})/);
        if (match) {
          setDigits(match[1].split(''));
        }
      }).catch(() => {});
      setTimeout(() => inputRefs.current[0]?.focus(), 50);
    }
  }, [open]);

  const handleDigitChange = (index: number, value: string) => {
    const digit = value.replace(/\D/g, '').slice(-1);
    const newDigits = [...digits];
    newDigits[index] = digit;
    setDigits(newDigits);
    setError('');

    if (digit && index < 3) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handleKeyDown = (index: number, e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Backspace' && !digits[index] && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
    if (e.key === 'Enter' && roomCode.length === 4 && !isJoining) {
      handleSubmit();
    }
    if (e.key === 'ArrowLeft' && index > 0) {
      inputRefs.current[index - 1]?.focus();
    }
    if (e.key === 'ArrowRight' && index < 3) {
      inputRefs.current[index + 1]?.focus();
    }
  };

  const handlePaste = (e: React.ClipboardEvent) => {
    e.preventDefault();
    const text = e.clipboardData.getData('text');
    const match = text.match(/[?&]code=(\d{4})/) || text.match(/^(\d{4})$/);
    if (match) {
      setDigits(match[1].split(''));
      setError('');
      inputRefs.current[3]?.focus();
    }
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
      setDigits(['', '', '', '']);
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : '加入房间失败');
    } finally {
      setIsJoining(false);
    }
  };

  const handleClose = () => {
    if (!isJoining) {
      setDigits(['', '', '', '']);
      setError('');
      onClose();
    }
  };

  return (
    <Dialog open={open} onOpenChange={(isOpen) => !isOpen && handleClose()}>
      <DialogContent className="max-w-sm">
        <DialogHeader>
          <DialogTitle>输入房间码</DialogTitle>
        </DialogHeader>

        <div className="space-y-4">
          <div>
            <div className="flex justify-center gap-3">
              {digits.map((digit, index) => (
                <input
                  key={index}
                  ref={(el) => { inputRefs.current[index] = el; }}
                  type="text"
                  inputMode="numeric"
                  maxLength={1}
                  value={digit}
                  onChange={(e) => handleDigitChange(index, e.target.value)}
                  onKeyDown={(e) => handleKeyDown(index, e)}
                  onPaste={handlePaste}
                  disabled={isJoining}
                  className="w-14 h-14 text-center text-2xl font-semibold
                    border border-stroke rounded-md bg-surface-elevated text-fg-primary
                    shadow-[var(--elevation-2),var(--shadow-relief)]
                    focus:outline-none focus-visible:border-state-active
                    focus-visible:shadow-[var(--elevation-2),var(--shadow-relief),0_0_12px_hsla(45,100%,51%,0.3)]
                    focus-visible:scale-[1.02]
                    disabled:opacity-50 disabled:cursor-not-allowed
                    transition-all duration-300 ease-[cubic-bezier(0.34,1.56,0.64,1)]"
                />
              ))}
            </div>
          </div>

          {error && (
            <p className="text-error text-sm text-center">{error}</p>
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
