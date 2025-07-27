import React, { useState, useCallback } from 'react';
import { Card } from '../../types';

interface PlayerHandProps {
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  disabled?: boolean;
}

interface CardDisplayProps {
  card: Card;
  isSelected: boolean;
  onClick: () => void;
  disabled?: boolean;
  stackIndex?: number;
}

const CardDisplay: React.FC<CardDisplayProps> = ({ 
  card, 
  isSelected, 
  onClick, 
  disabled = false,
  stackIndex = 0 
}) => {
  const getSuitSymbol = (suit: number) => {
    switch (suit) {
      case 0: return '♠'; // spades
      case 1: return '♥'; // hearts
      case 2: return '♣'; // clubs
      case 3: return '♦'; // diamonds
      default: return '';
    }
  };

  const getSuitColor = (suit: number) => {
    return suit === 1 || suit === 3 ? 'text-red-600' : 'text-black';
  };

  const getRankText = (rank: number) => {
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

  const getCardBackground = () => {
    if (card.is_joker) {
      return card.rank === 16 ? 'bg-red-100' : 'bg-gray-100';
    }
    return 'bg-white';
  };

  return (
    <div
      className={`
        relative w-12 h-16 border border-gray-300 rounded cursor-pointer transition-all duration-200
        ${getCardBackground()}
        ${isSelected ? 'transform -translate-y-2 border-blue-500 shadow-lg' : 'hover:shadow-md'}
        ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
      `}
      style={{ 
        marginLeft: stackIndex > 0 ? '-8px' : '0',
        zIndex: stackIndex 
      }}
      onClick={disabled ? undefined : onClick}
    >
      <div className="absolute inset-0 flex flex-col items-center justify-center text-xs">
        {card.is_joker ? (
          <div className="text-center font-bold">
            {getRankText(card.rank)}
          </div>
        ) : (
          <>
            <div className={`font-bold ${getSuitColor(card.suit)}`}>
              {getRankText(card.rank)}
            </div>
            <div className={`text-lg ${getSuitColor(card.suit)}`}>
              {getSuitSymbol(card.suit)}
            </div>
          </>
        )}
      </div>
    </div>
  );
};

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
  const getRankText = (rank: number) => {
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

  // Sort cards by suit priority: hearts > diamonds > clubs > spades
  const sortedCards = [...cards].sort((a, b) => {
    if (a.is_joker && b.is_joker) return b.rank - a.rank; // Big joker first
    if (a.is_joker) return -1;
    if (b.is_joker) return 1;
    
    const suitPriority = [3, 1, 2, 0]; // hearts, diamonds, clubs, spades
    const aPriority = suitPriority.indexOf(a.suit);
    const bPriority = suitPriority.indexOf(b.suit);
    return aPriority - bPriority;
  });

  const handleCardClick = (clickedCard: Card) => {
    const isSelected = selectedCards.some(c => c.id === clickedCard.id);
    
    if (isSelected) {
      // Remove from selection
      const newSelection = selectedCards.filter(c => c.id !== clickedCard.id);
      onCardSelect(newSelection);
    } else {
      // Add to selection
      const newSelection = [...selectedCards, clickedCard];
      onCardSelect(newSelection);
    }
  };

  return (
    <div className="flex flex-col items-center mb-4">
      <div className="text-xs text-gray-600 mb-1 font-medium">
        {getRankText(rank)} ({cards.length})
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
  // Group cards by rank
  const groupedCards = cards.reduce((groups, card) => {
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
    onCardSelect([...cards]);
  }, [cards, onCardSelect]);

  return (
    <div className="bg-white border border-gray-300 rounded-lg p-4">
      <div className="flex justify-between items-center mb-3">
        <div className="text-sm font-medium text-gray-700">
          手牌 ({cards.length}张)
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
            disabled={disabled || cards.length === 0}
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
            cards={groupedCards[rank]}
            selectedCards={selectedCards}
            onCardSelect={onCardSelect}
            disabled={disabled}
          />
        ))}
      </div>
      
      {cards.length === 0 && (
        <div className="text-center text-gray-500 py-8">
          暂无手牌
        </div>
      )}
    </div>
  );
};

export default PlayerHand;