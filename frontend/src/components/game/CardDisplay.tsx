import React from 'react';
import type { Card } from '../../types';

export const getSuitSymbol = (suit: number) => {
  switch (suit) {
    case 0: return '♠';
    case 1: return '♥';
    case 2: return '♣';
    case 3: return '♦';
    default: return '';
  }
};

export const getSuitColor = (suit: number) => {
  return suit === 1 || suit === 3 ? 'text-red-600' : 'text-black';
};

export const getRankText = (rank: number) => {
  if (rank === 15) return '小王';
  if (rank === 16) return '大王';
  if (rank <= 10) return rank.toString();
  switch (rank) {
    case 11: return 'J';
    case 12: return 'Q';
    case 13: return 'K';
    case 14: return 'A';
    default: return rank.toString();
  }
};

interface CardDisplayProps {
  card: Card;
  isSelected?: boolean;
  onClick?: () => void;
  disabled?: boolean;
  stackIndex?: number;
  size?: 'small' | 'normal';
}

const CardDisplay: React.FC<CardDisplayProps> = ({ 
  card, 
  isSelected = false, 
  onClick, 
  disabled = false,
  stackIndex = 0,
  size = 'normal'
}) => {
  const getCardBackground = () => {
    if (card.isJoker) {
      return card.rank === 16 ? 'bg-red-100' : 'bg-gray-100';
    }
    return 'bg-white';
  };

  const sizeClasses = size === 'small' 
    ? 'w-8 h-11 text-[10px]' 
    : 'w-12 h-16 text-xs';

  const marginLeft = stackIndex > 0 ? (size === 'small' ? '-6px' : '-8px') : '0';

  return (
    <div
      className={`
        relative border border-gray-300 rounded transition-all duration-200
        ${sizeClasses}
        ${getCardBackground()}
        ${isSelected ? 'transform -translate-y-2 border-blue-500 shadow-lg' : onClick ? 'hover:shadow-md' : ''}
        ${disabled ? 'opacity-50 cursor-not-allowed' : onClick ? 'cursor-pointer' : ''}
      `}
      style={{ 
        marginLeft,
        zIndex: stackIndex 
      }}
      onClick={disabled ? undefined : onClick}
    >
      <div className="absolute inset-0 flex flex-col items-center justify-center">
        {card.isJoker ? (
          <div className="text-center font-bold">
            {getRankText(card.rank)}
          </div>
        ) : (
          <>
            <div className={`font-bold ${getSuitColor(card.suit)}`}>
              {getRankText(card.rank)}
            </div>
            <div className={`${size === 'small' ? 'text-sm' : 'text-lg'} ${getSuitColor(card.suit)}`}>
              {getSuitSymbol(card.suit)}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

export default CardDisplay;
