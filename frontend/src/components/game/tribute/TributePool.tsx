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
}) => {
  const maxSlots = tributeType === TributeType.TRIBUTE_TYPE_DOUBLE_DOWN ? 2 : 1;

  return (
    <div className="relative w-72 bg-green-200 border-2 border-green-400 rounded-lg p-4">
      <div className="text-center mb-3">
        <div className="text-sm font-medium text-green-800">
          贡牌池
        </div>
        <div className="text-xs text-green-700">
          {getTributeTypeName(tributeType)}
        </div>
      </div>

      <div className="flex justify-center gap-4 mb-4">
        {Array.from({ length: maxSlots }).map((_, index) => {
          const card = poolCards[index] ?? null;
          return (
            <div
              key={index}
              ref={(el) => onSlotRefReady?.(index, el)}
              className={`w-16 h-24 border-2 border-dashed rounded-lg flex items-center justify-center
                ${card ? 'border-green-500 bg-white' : 'border-green-400 bg-green-100'}
                ${canSelect && card ? 'cursor-pointer hover:border-yellow-400 hover:shadow-lg transition-all' : ''}
              `}
              onClick={() => canSelect && card && onSelectCard(card)}
            >
              {card ? (
                <CardDisplay card={card} size="small" />
              ) : (
                <span className="text-green-500 text-xs">空</span>
              )}
            </div>
          );
        })}
      </div>

      <div className="border-t border-green-400 pt-2 max-h-24 overflow-y-auto">
        <div className="text-xs text-green-800 space-y-1">
          {messages.slice(-3).map((msg, i) => (
            <div key={i} className="truncate">{msg}</div>
          ))}
        </div>
      </div>
    </div>
  );
};

export default TributePool;
