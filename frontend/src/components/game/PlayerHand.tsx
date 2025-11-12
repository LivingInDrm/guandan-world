import React, { useCallback } from 'react';
import type { Card } from '../../types';
import CardDisplay, { getRankText } from './CardDisplay';

interface PlayerHandProps {
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  disabled?: boolean;
}

interface CardGroupProps {
  rank: number;
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  disabled?: boolean;
}

const CardGroup: React.FC<CardGroupProps> = ({ 
  rank, 
  cards, 
  selectedCards, 
  onCardSelect, 
  disabled = false 
}) => {
  const safeCards = Array.isArray(cards) ? cards : [];

  const sortedCards = [...safeCards].sort((a, b) => {
    if (a.is_joker && b.is_joker) return b.rank - a.rank;
    if (a.is_joker) return -1;
    if (b.is_joker) return 1;
    
    const suitPriority = [3, 1, 2, 0];
    const aPriority = suitPriority.indexOf(a.suit);
    const bPriority = suitPriority.indexOf(b.suit);
    return aPriority - bPriority;
  });

  const handleCardClick = (clickedCard: Card) => {
    const isSelected = selectedCards.some(c => c.id === clickedCard.id);
    
    if (isSelected) {
      const newSelection = selectedCards.filter(c => c.id !== clickedCard.id);
      onCardSelect(newSelection);
    } else {
      const newSelection = [...selectedCards, clickedCard];
      onCardSelect(newSelection);
    }
  };

  return (
    <div className="flex flex-col items-center mb-4">
      <div className="text-xs text-gray-600 mb-1 font-medium">
        {getRankText(rank)} ({safeCards.length})
      </div>
      <div className="flex items-end">
        {sortedCards.map((card, index) => (
          <CardDisplay
            key={card.id}
            card={card}
            isSelected={selectedCards.some(c => c.id === card.id)}
            onClick={() => handleCardClick(card)}
            disabled={disabled}
            stackIndex={index}
          />
        ))}
      </div>
    </div>
  );
};

const PlayerHand: React.FC<PlayerHandProps> = ({ 
  cards, 
  selectedCards, 
  onCardSelect, 
  disabled = false 
}) => {
  // Ensure cards is an array
  const safeCards = Array.isArray(cards) ? cards : [];
  
  // Group cards by rank
  const groupedCards = safeCards.reduce((groups, card) => {
    const rank = card.rank;
    if (!groups[rank]) {
      groups[rank] = [];
    }
    groups[rank].push(card);
    return groups;
  }, {} as Record<number, Card[]>);

  // Sort ranks in descending order (big joker, small joker, A, K, Q, J, 10, 9, ..., 2)
  const sortedRanks = Object.keys(groupedCards)
    .map(Number)
    .sort((a, b) => b - a);

  const handleClearSelection = useCallback(() => {
    onCardSelect([]);
  }, [onCardSelect]);

  const handleSelectAll = useCallback(() => {
    onCardSelect([...safeCards]);
  }, [safeCards, onCardSelect]);

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-4">
      <div className="flex justify-between items-center mb-3">
        <div className="text-sm font-medium text-gray-700">
          手牌 ({safeCards.length}张)
          {selectedCards.length > 0 && (
            <span className="ml-2 text-blue-600">
              已选择 {selectedCards.length}张
            </span>
          )}
        </div>
        <div className="flex gap-2">
          {selectedCards.length > 0 && (
            <button
              onClick={handleClearSelection}
              disabled={disabled}
              className="text-xs px-2 py-1 bg-gray-200 text-gray-700 rounded hover:bg-gray-300 disabled:opacity-50"
            >
              清空选择
            </button>
          )}
          <button
            onClick={handleSelectAll}
            disabled={disabled || safeCards.length === 0}
            className="text-xs px-2 py-1 bg-blue-200 text-blue-700 rounded hover:bg-blue-300 disabled:opacity-50"
          >
            全选
          </button>
        </div>
      </div>
      
      <div className="flex flex-wrap gap-x-4 gap-y-2 justify-center max-h-64 overflow-y-auto">
        {sortedRanks.map(rank => (
          <CardGroup
            key={rank}
            rank={rank}
            cards={groupedCards[rank] || []}
            selectedCards={selectedCards}
            onCardSelect={onCardSelect}
            disabled={disabled}
          />
        ))}
      </div>
      
      {safeCards.length === 0 && (
        <div className="text-center text-gray-500 py-8">
          暂无手牌
        </div>
      )}
    </div>
  );
};

export default PlayerHand;