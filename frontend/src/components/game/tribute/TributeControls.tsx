import React, { useRef, useCallback } from 'react';
import type { Card } from '../../../types';
import Countdown from '../Countdown';
import { Button } from '@/components/ui';
import { Lightbulb, Gift } from 'lucide-react';
import { cn } from '@/lib/utils';
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
    <div className={cn(
      "py-2 px-4 mobile-landscape:py-1 mobile-landscape:px-2",
    )}>
      <div className="flex items-center justify-center gap-[var(--control-gap)]">
        <div className="flex items-center justify-center">
          <Countdown
            deadlineAtMs={turnDeadlineAtMs}
            isActive={canReturnTribute && !disabled}
            size="small"
          />
        </div>

        <Button
          intent="secondary"
          size="md"
          onClick={handleHint}
          disabled={isHintDisabled}
          className={cn(
            "min-w-[90px] h-10 px-4 pr-5 text-base",
            "mobile-landscape:min-w-[70px] mobile-landscape:h-8 mobile-landscape:px-3 mobile-landscape:pr-4 mobile-landscape:text-sm",
            "bg-[hsl(158,55%,35%)]/80 hover:bg-[hsl(158,55%,42%)]",
            "text-white",
            "border border-[hsl(158,55%,50%)]/30 hover:border-[hsl(158,55%,50%)]/50",
            "shadow-none hover:shadow-[0_0_12px_hsla(158,55%,42%,0.4)]",
            "disabled:bg-white/5 disabled:text-white/30 disabled:border-white/10",
          )}
        >
          <Lightbulb className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4 mr-1.5" />
          提示
        </Button>

        <Button
          intent="primary"
          size="md"
          onClick={onReturnTribute}
          disabled={isReturnDisabled}
          className={cn(
            "min-w-[90px] h-10 px-4 pr-5 text-base",
            "mobile-landscape:min-w-[70px] mobile-landscape:h-8 mobile-landscape:px-3 mobile-landscape:pr-4 mobile-landscape:text-sm",
            "bg-gradient-to-r from-state-active to-[hsl(42,95%,45%)]",
            "hover:from-[hsl(42,95%,55%)] hover:to-[hsl(42,95%,48%)]",
            "text-primitive-neutral-900 font-bold",
            "border-none",
            "shadow-[0_2px_8px_hsla(42,95%,52%,0.4)]",
            "hover:shadow-[0_4px_16px_hsla(42,95%,52%,0.5)]",
            "disabled:bg-white/10 disabled:from-white/10 disabled:to-white/10",
            "disabled:text-white/30 disabled:shadow-none",
          )}
        >
          <Gift className="w-5 h-5 mobile-landscape:w-4 mobile-landscape:h-4 mr-1.5" />
          还贡
        </Button>
      </div>
    </div>
  );
};

export default TributeControls;
