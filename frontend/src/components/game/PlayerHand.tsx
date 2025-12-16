import React, { useCallback, useEffect, useRef } from 'react';
import type { Card } from '../../types/proto';
import { isJoker } from '../../utils/cardUtils';
import CardDisplay from './CardDisplay';
import { RotateCcw } from 'lucide-react';
import { cn } from '@/lib/utils';

interface PlayerHandProps {
  cards: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  currentLevel?: number;
  finishRank?: number | null;
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
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
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

const getRankText = (rank: number): string => {
  switch (rank) {
    case 1: return '头游';
    case 2: return '二游';
    case 3: return '三游';
    case 4: return '末游';
    default: return '';
  }
};

const PlayerHand: React.FC<PlayerHandProps> = ({ 
  cards, 
  selectedCards, 
  onCardSelect,
  currentLevel,
  finishRank
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

  useEffect(() => {
    const handleGlobalPointerUp = () => {
      isDraggingRef.current = false;
      touchedIndexesRef.current = new Set();
    };
    window.addEventListener('pointerup', handleGlobalPointerUp);
    window.addEventListener('pointercancel', handleGlobalPointerUp);
    return () => {
      window.removeEventListener('pointerup', handleGlobalPointerUp);
      window.removeEventListener('pointercancel', handleGlobalPointerUp);
    };
  }, []);

  return (
    <div 
      className="bg-transparent pt-4 px-4 pb-2 select-none touch-none"
      onPointerDown={handlePointerDown}
    >
      <div className="flex justify-end mb-3">
        <button
          onClick={handleClearSelection}
          className={cn(
            "p-1.5 rounded-md",
            "text-fg-secondary hover:text-fg-primary",
            "hover:bg-surface-elevated",
            "transition-colors duration-fast",
            "opacity-60 hover:opacity-100",
            selectedCards.length > 0 ? '' : 'invisible'
          )}
          aria-label="清空选择"
        >
          <RotateCcw className="w-4 h-4" />
        </button>
      </div>
      
      <div className="flex flex-wrap items-end gap-x-1 gap-y-2 justify-center h-64 overflow-y-auto">
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
        <div className="flex items-start justify-center h-64 pt-4">
          {finishRank ? (
            <div className="text-center">
              <span 
                className={cn(
                  "text-5xl font-black tracking-widest",
                  "bg-gradient-to-b bg-clip-text text-transparent",
                  finishRank === 1 
                    ? "from-state-active via-action-secondary to-action-secondary" 
                    : finishRank === 4 
                      ? "from-fg-secondary via-action-neutral to-action-neutral"
                      : "from-action-primary via-action-primary to-action-primary/70"
                )}
                style={{ 
                  textShadow: finishRank === 1 
                    ? '0 0 30px hsla(var(--primitive-accent-500), 0.6)' 
                    : undefined 
                }}
              >
                {getRankText(finishRank)}
              </span>
            </div>
          ) : (
            <div className="text-center text-fg-secondary">
              暂无手牌
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default PlayerHand;