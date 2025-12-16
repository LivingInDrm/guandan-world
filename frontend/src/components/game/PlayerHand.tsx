import React, { useCallback, useEffect, useRef } from 'react';
import type { Card } from '../../types/proto';
import { isJoker } from '../../utils/cardUtils';
import CardDisplay from './CardDisplay';
import { RotateCcw, Trophy, Medal, Award, Frown } from 'lucide-react';
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
    <div className="flex flex-col items-center mb-4 mobile-landscape:mb-2">
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

const getRankConfig = (rank: number) => {
  switch (rank) {
    case 1:
      return {
        text: '头游',
        icon: Trophy,
        gradient: 'from-state-active via-[hsl(42,95%,60%)] to-[hsl(42,95%,45%)]',
        glow: '0 0 40px hsla(42, 95%, 52%, 0.6), 0 0 20px hsla(42, 95%, 52%, 0.4)',
        iconColor: 'text-state-active',
      };
    case 2:
      return {
        text: '二游',
        icon: Medal,
        gradient: 'from-[hsl(158,55%,42%)] via-[hsl(158,55%,50%)] to-[hsl(158,55%,38%)]',
        glow: '0 0 30px hsla(158, 55%, 42%, 0.5)',
        iconColor: 'text-[hsl(158,55%,42%)]',
      };
    case 3:
      return {
        text: '三游',
        icon: Award,
        gradient: 'from-[hsl(30,50%,50%)] via-[hsl(30,55%,55%)] to-[hsl(30,45%,45%)]',
        glow: '0 0 20px hsla(30, 50%, 50%, 0.4)',
        iconColor: 'text-[hsl(30,50%,50%)]',
      };
    case 4:
      return {
        text: '末游',
        icon: Frown,
        gradient: 'from-white/60 via-white/50 to-white/40',
        glow: 'none',
        iconColor: 'text-white/50',
      };
    default:
      return {
        text: '',
        icon: null,
        gradient: '',
        glow: 'none',
        iconColor: '',
      };
  }
};

const findCardDeckIndex = (x: number, y: number): number | null => {
  const elements = document.elementsFromPoint(x, y);
  for (const element of elements) {
    const cardElement = (element as HTMLElement).closest('[data-deck-index]');
    if (cardElement) {
      const deckIndex = cardElement.getAttribute('data-deck-index');
      return deckIndex ? parseInt(deckIndex, 10) : null;
    }
  }
  return null;
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
  const cardsRef = useRef(cards);
  cardsRef.current = cards;

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

  const handlePointerMove = useCallback((e: React.PointerEvent) => {
    if (!isDraggingRef.current) return;

    const deckIndex = findCardDeckIndex(e.clientX, e.clientY);
    if (deckIndex === null) return;
    if (touchedIndexesRef.current.has(deckIndex)) return;

    touchedIndexesRef.current.add(deckIndex);

    const card = cardsRef.current.find(c => c.deckIndex === deckIndex);
    if (!card) return;

    const currentSelected = selectedCardsRef.current;
    const isCardSelected = currentSelected.some(c => c.deckIndex === deckIndex);

    if (isCardSelected) {
      onCardSelect(currentSelected.filter(c => c.deckIndex !== deckIndex));
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
      className={cn(
        "select-none touch-none",
        "pt-2 px-4 pb-2 mobile-landscape:pt-1 mobile-landscape:px-2 mobile-landscape:pb-1",
      )}
      onPointerDown={handlePointerDown}
      onPointerMove={handlePointerMove}
    >
      {/* 清空选择按钮 */}
      <div className="flex justify-end mb-1">
        <button
          onClick={handleClearSelection}
          className={cn(
            "p-2 mobile-landscape:p-1.5 rounded-lg",
            "text-white/60 hover:text-white",
            "bg-black/20 hover:bg-black/40",
            "border border-white/10 hover:border-white/20",
            "transition-all duration-200",
            "hover:shadow-[0_0_8px_hsla(158,55%,42%,0.3)]",
            selectedCards.length > 0 ? '' : 'invisible'
          )}
          aria-label="清空选择"
        >
          <RotateCcw className="w-4 h-4 mobile-landscape:w-3 mobile-landscape:h-3" />
        </button>
      </div>

      {/* 手牌区域 */}
      <div className="flex flex-wrap items-end gap-x-1 gap-y-2 mobile-landscape:gap-x-0.5 mobile-landscape:gap-y-1 justify-center min-h-[var(--hand-area-height)] pb-2">
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

      {/* 空手牌状态 - 完成排名或等待 */}
      {safeCards.length === 0 && (
        <div className="flex items-start justify-center min-h-[var(--hand-area-height)] pt-4 mobile-landscape:pt-2 pb-2">
          {finishRank ? (
            (() => {
              const config = getRankConfig(finishRank);
              const IconComponent = config.icon;
              return (
                <div className={cn(
                  "flex flex-col items-center gap-3 mobile-landscape:gap-2 p-6 mobile-landscape:p-3 rounded-2xl mobile-landscape:rounded-xl",
                  "bg-black/30 backdrop-blur-sm",
                  "border border-white/10",
                  finishRank === 1 && "animate-pulse-glow",
                )}>
                  {IconComponent && (
                    <IconComponent
                      className={cn("w-12 h-12 mobile-landscape:w-8 mobile-landscape:h-8", config.iconColor)}
                      strokeWidth={1.5}
                    />
                  )}
                  <span
                    className={cn(
                      "text-5xl mobile-landscape:text-3xl font-black tracking-widest font-display",
                      "bg-gradient-to-b bg-clip-text text-transparent",
                      config.gradient,
                    )}
                    style={{
                      textShadow: config.glow !== 'none' ? config.glow : undefined,
                    }}
                  >
                    {config.text}
                  </span>
                </div>
              );
            })()
          ) : (
            <div className={cn(
              "flex flex-col items-center gap-2 p-4 mobile-landscape:p-2 rounded-xl",
              "bg-black/20 backdrop-blur-sm",
              "border border-white/10",
            )}>
              <span className="text-white/50 text-sm mobile-landscape:text-xs">暂无手牌</span>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default PlayerHand;