import React from 'react';
import type { Card } from '../../../types';
import { TributeType } from '../../../types/generated/event';
import CardDisplay from '../CardDisplay';
import { Badge } from '@/components/ui-next';

interface TributePoolProps {
  poolCards: (Card | null)[];
  canSelect: boolean;
  onSelectCard: (card: Card) => void;
  tributeType: TributeType;
  messages: string[];
  onSlotRefReady?: (slotIndex: number, element: HTMLDivElement | null) => void;
  flyingToPoolDeckIndexes?: number[];
}

const getTributeTypeName = (type: TributeType): string => {
  switch (type) {
    case TributeType.TRIBUTE_TYPE_DOUBLE_DOWN:
      return '双下';
    case TributeType.TRIBUTE_TYPE_SINGLE_LAST:
      return '单下';
    case TributeType.TRIBUTE_TYPE_PARTNER_LAST:
      return '末游';
    case TributeType.TRIBUTE_TYPE_NONE:
      return '无需进贡';
    default:
      return '未知';
  }
};

const TributePool: React.FC<TributePoolProps> = ({
  poolCards,
  canSelect,
  onSelectCard,
  tributeType,
  messages,
  onSlotRefReady,
  flyingToPoolDeckIndexes = [],
}) => {
  const maxSlots = tributeType === TributeType.TRIBUTE_TYPE_DOUBLE_DOWN ? 2 : 1;

  return (
    <div className="relative w-72 p-4">
      <div className="text-center mb-3">
        <Badge variant="neutral" size="md">
          {getTributeTypeName(tributeType)}
        </Badge>
      </div>

      <div className="flex justify-center gap-4 mb-4">
        {Array.from({ length: maxSlots }).map((_, index) => {
          const card = poolCards[index] ?? null;
          const isFlying = card && flyingToPoolDeckIndexes.includes(card.deckIndex);
          const shouldShowCard = card && !isFlying;
          return (
            <div
              key={index}
              ref={(el) => onSlotRefReady?.(index, el)}
              className={`w-16 h-24 border border-dashed rounded-ds-sm flex items-center justify-center shadow-ds-elevation-1
                ${shouldShowCard ? 'border-ds-border-emphasis bg-ds-surface-elevated/50' : 'border-ds-border/50 bg-ds-surface-base/10'}
                ${canSelect && shouldShowCard ? 'cursor-pointer hover:border-ds-state-active hover:shadow-ds-elevation-3 transition-all' : ''}
              `}
              onClick={() => canSelect && shouldShowCard && onSelectCard(card)}
            >
              {shouldShowCard ? (
                <CardDisplay card={card} size="small" />
              ) : (
                <span className="text-ds-text-secondary text-xs">空</span>
              )}
            </div>
          );
        })}
      </div>

      <div className="border-t border-ds-border/30 pt-2 mt-1 max-h-24 overflow-y-auto">
        <div className="text-xs text-ds-text-primary space-y-1">
          {messages.slice(-3).map((msg, i) => (
            <div key={i} className="truncate">{msg}</div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TributePool;
