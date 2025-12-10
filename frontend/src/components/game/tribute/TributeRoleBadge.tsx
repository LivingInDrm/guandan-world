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
      <div className="text-xs px-2 py-1 rounded bg-destructive/20 text-destructive whitespace-nowrap">
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
          ? 'bg-accent/20 text-accent-foreground animate-pulse' 
          : 'bg-primary/20 text-primary'
      )}>
        {getText()}
      </div>
    );
  }

  return null;
};

export default React.memo(TributeRoleBadge);
