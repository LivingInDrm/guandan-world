import React, { useRef, useCallback } from 'react';
import type { Card } from '../../../types';
import Countdown from '../Countdown';
import { Button } from '@/components/ui';
import { fromProtoCard, TributeRecommender } from '@guandan/sdk-ts';

interface TributeControlsProps {
  selectedCards: Card[];
  canReturnTribute: boolean;
  turnDeadlineAtMs: number;
  onReturnTribute: () => void;
  onHint: (cards: Card[]) => void;
  handCards: Card[];
  dealLevel: number;
  disabled?: boolean;
}

const TributeControls: React.FC<TributeControlsProps> = ({
  selectedCards,
  canReturnTribute,
  turnDeadlineAtMs,
  onReturnTribute,
  onHint,
  handCards,
  dealLevel,
  disabled = false,
}) => {
  const tributeRecommenderRef = useRef<TributeRecommender | null>(null);
  const lastHandCardsRef = useRef<Card[]>([]);

  const handleHint = useCallback(() => {
    if (!canReturnTribute || handCards.length === 0) return;

    const cardsChanged =
      lastHandCardsRef.current.length !== handCards.length ||
      lastHandCardsRef.current.some((c, i) => c.deckIndex !== handCards[i]?.deckIndex);

    if (cardsChanged) {
      tributeRecommenderRef.current = null;
      lastHandCardsRef.current = handCards;
    }

    if (!tributeRecommenderRef.current) {
      const sdkCards = handCards.map(c => fromProtoCard(c, dealLevel));
      tributeRecommenderRef.current = new TributeRecommender(sdkCards);
    }

    const recommended = tributeRecommenderRef.current.next();
    if (recommended) {
      const selectedCard = handCards.find(c => c.deckIndex === recommended.deckIndex);
      onHint(selectedCard ? [selectedCard] : []);
    } else {
      onHint([]);
    }
  }, [canReturnTribute, handCards, dealLevel, onHint]);

  const isReturnDisabled = disabled || !canReturnTribute || selectedCards.length !== 1;
  const isHintDisabled = disabled || !canReturnTribute;

  return (
    <div className="py-2 px-4">
      <div className="flex items-center justify-center gap-3">
        <div className="flex items-center justify-center">
          <Countdown
            deadlineAtMs={turnDeadlineAtMs}
            isActive={canReturnTribute && !disabled}
            size="small"
          />
        </div>

        <Button
          intent="secondary"
          size="sm"
          onClick={handleHint}
          disabled={isHintDisabled}
        >
          提示
        </Button>

        <Button
          intent="primary"
          size="sm"
          onClick={onReturnTribute}
          disabled={isReturnDisabled}
        >
          还贡
        </Button>
      </div>
    </div>
  );
};

export default TributeControls;
