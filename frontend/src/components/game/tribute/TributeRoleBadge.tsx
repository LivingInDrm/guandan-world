import React from 'react';
import { cn } from '@/lib/utils';

export interface TributeRoleBadgeProps {
  role: 'giver' | 'receiver' | null;
  isSubmitted?: boolean;
  isReceived?: boolean;
  isCurrentSelector?: boolean;
}

const TributeRoleBadge: React.FC<TributeRoleBadgeProps> = ({
  role,
  isSubmitted = false,
  isReceived = false,
  isCurrentSelector = false,
}) => {
  if (!role) {
    return null;
  }

  if (role === 'giver') {
    return (
      <div className="text-xs px-2 py-1 rounded bg-action-danger/20 text-action-danger whitespace-nowrap">
        {isSubmitted ? '已上贡' : '待上贡'}
      </div>
    );
  }

  if (role === 'receiver') {
    const getText = () => {
      if (isReceived) return '已收贡';
      if (isCurrentSelector) return '选牌中';
      return '待收贡';
    };

    return (
      <div className={cn(
        'text-xs px-2 py-1 rounded whitespace-nowrap',
        isCurrentSelector 
          ? 'bg-state-active/20 text-fg-primary animate-pulse' 
          : 'bg-action-primary/20 text-action-primary'
      )}>
        {getText()}
      </div>
    );
  }

  return null;
};

export default React.memo(TributeRoleBadge);
