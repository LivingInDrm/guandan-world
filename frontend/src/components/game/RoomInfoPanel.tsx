import React, { useState, useCallback } from 'react';
import { Copy, Check } from 'lucide-react';
import { Card, Button } from '@/components/ui';

export interface RoomInfoPanelProps {
  roomCode: string;
}

const RoomInfoPanel: React.FC<RoomInfoPanelProps> = ({ roomCode }) => {
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

  return (
    <div className="absolute top-4 left-4 z-10 flex gap-2">
      <Card
        variant="base"
        interactive={false}
        className="w-28 flex flex-col items-center rounded-md overflow-hidden shadow-elevation-2 bg-gradient-to-b from-[hsl(var(--primitive-primary-500))] to-[hsl(var(--primitive-primary-700))] text-fg-inverse"
      >
        <div className="w-full text-center py-1 text-sm font-semibold backdrop-brightness-95 bg-black/10">
          房间码
        </div>
        <div className="py-1.5 text-2xl font-mono font-extrabold tracking-widest drop-shadow-sm">
          {roomCode}
        </div>
        <div className="w-full pb-2 px-2">
          <Button
            intent="secondary"
            size="sm"
            onClick={handleCopy}
            className="w-full min-w-0 px-3 py-1.5 text-sm"
          >
            {copied ? (
              <>
                <Check className="w-3.5 h-3.5" />
                已复制
              </>
            ) : (
              <>
                <Copy className="w-3.5 h-3.5" />
                邀请
              </>
            )}
          </Button>
          {copyError && (
            <p className="text-xs text-center mt-1 text-red-200">复制失败</p>
          )}
        </div>
      </Card>
    </div>
  );
};

export default RoomInfoPanel;
