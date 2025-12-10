import React, { useCallback, useRef } from 'react';
import type { Card } from '../../types/proto';
import { isJoker } from '../../utils/cardUtils';
import CardDisplay from './CardDisplay';

interface PlayerHandProps {
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  currentLevel?: number;
}

interface CardGroupProps {
  rank: number;
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onCardPointerEnter?: (card: Card) => void;
  currentLevel?: number;
}

const CardGroup: React.FC<CardGroupProps> = ({ 
  rank: _rank, 
  cards, 
  selectedCards, 
  onCardSelect,
  onCardPointerEnter,
  currentLevel
}) => {
  const safeCards = Array.isArray(cards) ? cards : [];

  const sortedCards = [...safeCards].sort((a, b) => {
    if (isJoker(a) && isJoker(b)) return b.rank - a.rank;
    if (isJoker(a)) return -1;
    if (isJoker(b)) return 1;
    
    const suitPriority = [2, 0, 3, 1];
    const aPriority = suitPriority.indexOf(a.suit);
    const bPriority = suitPriority.indexOf(b.suit);
    return aPriority - bPriority;
  });

  const handleCardClick = (clickedCard: Card) => {
    const isSelected = selectedCards.some(c => c.deckIndex === clickedCard.deckIndex);
    
    if (isSelected) {
      const newSelection = selectedCards.filter(c => c.deckIndex !== clickedCard.deckIndex);
      onCardSelect(newSelection);
    } else {
      const newSelection = [...selectedCards, clickedCard];
      onCardSelect(newSelection);
    }
  };

  return (
    <div className="flex flex-col items-center mb-4">
      <div className="flex flex-col items-center">
        {sortedCards.map((card, index) => (
          <CardDisplay
            key={card.deckIndex}
            card={card}
            deckIndex={card.deckIndex}
            isSelected={selectedCards.some(c => c.deckIndex === card.deckIndex)}
            onClick={() => handleCardClick(card)}
            onPointerEnter={() => onCardPointerEnter?.(card)}
            stackIndex={index}
            stackDirection="vertical"
            currentLevel={currentLevel}
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
  currentLevel
}) => {
  const isDraggingRef = useRef(false);
  const touchedIndexesRef = useRef<Set<number>>(new Set());
  const selectedCardsRef = useRef(selectedCards);
  selectedCardsRef.current = selectedCards;

  const safeCards = Array.isArray(cards) ? cards : [];
  
  const groupedCards = safeCards.reduce((groups, card) => {
    const rank = card.rank;
    if (!groups[rank]) {
      groups[rank] = [];
    }
    groups[rank].push(card);
    return groups;
  }, {} as Record<number, Card[]>);

  const sortedRanks = Object.keys(groupedCards)
    .map(Number)
    .sort((a, b) => b - a);

  const handleClearSelection = useCallback(() => {
    onCardSelect([]);
  }, [onCardSelect]);

  const handlePointerDown = useCallback(() => {
    isDraggingRef.current = true;
    touchedIndexesRef.current = new Set();
  }, []);

  const handlePointerUp = useCallback(() => {
    isDraggingRef.current = false;
    touchedIndexesRef.current = new Set();
  }, []);

  const handleCardPointerEnter = useCallback((card: Card) => {
    if (!isDraggingRef.current) return;
    if (touchedIndexesRef.current.has(card.deckIndex)) return;

    touchedIndexesRef.current.add(card.deckIndex);

    const currentSelected = selectedCardsRef.current;
    const isCardSelected = currentSelected.some(c => c.deckIndex === card.deckIndex);

    if (isCardSelected) {
      onCardSelect(currentSelected.filter(c => c.deckIndex !== card.deckIndex));
    } else {
      onCardSelect([...currentSelected, card]);
    }
  }, [onCardSelect]);

  return (
    <div 
      className="relative bg-card-slot backdrop-blur-md rounded-t-3xl shadow-2xl border-t border-white/20 pt-6 px-6 pb-4 select-none"
      onPointerDown={handlePointerDown}
      onPointerUp={handlePointerUp}
    >
      <div className="absolute top-0 left-1/2 -translate-x-1/2 w-2/3 h-1.5 bg-black/20 rounded-b-full blur-[2px]" />
      
      <div className="flex justify-end mb-3">
        <button
          onClick={handleClearSelection}
          className={`text-xs px-3 py-1.5 bg-white/20 text-white/90 rounded-full hover:bg-white/30 transition-colors font-medium ${selectedCards.length > 0 ? '' : 'invisible'}`}
        >
          清空选择
        </button>
      </div>
      
      <div className="flex flex-wrap items-end gap-x-1 gap-y-2 justify-center min-h-[240px] max-h-[280px] overflow-y-auto">
        {sortedRanks.map(rank => (
          <CardGroup
            key={rank}
            rank={rank}
            cards={groupedCards[rank] || []}
            selectedCards={selectedCards}
            onCardSelect={onCardSelect}
            onCardPointerEnter={handleCardPointerEnter}
            currentLevel={currentLevel}
          />
        ))}
      </div>
      
      {safeCards.length === 0 && (
        <div className="text-center text-white/60 py-8 font-medium">
          暂无手牌
        </div>
      )}
    </div>
  );
};

export default PlayerHand;