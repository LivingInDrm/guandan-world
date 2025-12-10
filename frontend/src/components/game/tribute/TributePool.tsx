import React from 'react';
import type { Card } from '../../../types';
import { TributeType } from '../../../types/generated/event';
import CardDisplay from '../CardDisplay';

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
        <div className="text-sm font-semibold text-foreground">
          {getTributeTypeName(tributeType)}
        </div>
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
              className={`w-16 h-24 border border-dashed rounded-lg flex items-center justify-center shadow-sm
                ${shouldShowCard ? 'border-table-400 bg-card/50' : 'border-table-400/50 bg-table-300/10'}
                ${canSelect && shouldShowCard ? 'cursor-pointer hover:border-accent hover:shadow-lg transition-all' : ''}
              `}
              onClick={() => canSelect && shouldShowCard && onSelectCard(card)}
            >
              {shouldShowCard ? (
                <CardDisplay card={card} size="small" />
              ) : (
                <span className="text-muted-foreground text-xs">空</span>
              )}
            </div>
          );
        })}
      </div>

      <div className="border-t border-table-400/30 pt-2 mt-1 max-h-24 overflow-y-auto">
        <div className="text-xs text-foreground space-y-1">
          {messages.slice(-3).map((msg, i) => (
            <div key={i} className="truncate">{msg}</div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TributePool;
