import React from 'react';
import type { Card } from '../../../types';
import { TributeType } from '../../../types/generated/event';
import CardDisplay from '../CardDisplay';
import { Badge } from '@/components/ui';

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
    <div className="relative p-2 mobile-landscape:p-1">
      <div className="text-center mb-2 mobile-landscape:mb-1">
        <Badge variant="neutral" size="md">
          {getTributeTypeName(tributeType)}
        </Badge>
      </div>

      <div className="flex justify-center gap-2 mobile-landscape:gap-1 mb-2 mobile-landscape:mb-1">
        {Array.from({ length: maxSlots }).map((_, index) => {
          const card = poolCards[index] ?? null;
          const isFlying = card && flyingToPoolDeckIndexes.includes(card.deckIndex);
          const shouldShowCard = card && !isFlying;
          return (
            <div
              key={index}
              ref={(el) => onSlotRefReady?.(index, el)}
              style={{
                width: 'var(--card-width-small)',
                height: 'var(--card-height-small)',
              }}
              className={`border border-dashed rounded-sm flex items-center justify-center shadow-elevation-1
                ${shouldShowCard ? 'border-stroke-emphasis bg-surface-elevated/50' : 'border-stroke/50 bg-surface-base/10'}
                ${canSelect && shouldShowCard ? 'cursor-pointer hover:border-state-active hover:shadow-elevation-3 transition-all' : ''}
              `}
              onClick={() => canSelect && shouldShowCard && onSelectCard(card)}
            >
              {shouldShowCard ? (
                <CardDisplay card={card} size="small" />
              ) : (
                <span className="text-fg-secondary text-xs mobile-landscape:text-[10px]">空</span>
              )}
            </div>
          );
        })}
      </div>

      <div className="border-t border-stroke/30 pt-1 mt-1 max-h-16 mobile-landscape:max-h-12 overflow-y-auto">
        <div className="text-xs mobile-landscape:text-[10px] text-fg-primary space-y-0.5">
          {messages.slice(-3).map((msg, i) => (
            <div key={i} className="truncate">{msg}</div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TributePool;
