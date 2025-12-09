import React, { useState, useEffect, useCallback, useRef } from 'react';
import type { Card, PlayAction } from '../../types';
import Countdown from './Countdown';
import {
  fromProtoCard,
  fromCardList,
  FirstPlayRecommender,
  NextPlayRecommender,
} from '@guandan/sdk-ts';

interface GameControlsProps {
  selectedCards: Card[];
  canPlay: boolean;
  isMyTurn: boolean;
  turnDeadlineAtMs: number;
  onPlayCards: (cards: Card[]) => void;
  onPass: () => void;
  onHint: (cards: Card[]) => void;
  handCards: Card[];
  plays: PlayAction[];
  leader?: number;
  playerSeat: number;
  dealLevel: number;
  disabled?: boolean;
}

interface PlayValidationResult {
  isValid: boolean;
  error?: string;
  cardType?: string;
}

const GameControls: React.FC<GameControlsProps> = ({
  selectedCards,
  canPlay,
  isMyTurn,
  turnDeadlineAtMs,
  onPlayCards,
  onPass,
  onHint,
  handCards,
  plays,
  leader,
  playerSeat,
  dealLevel,
  disabled = false
}) => {
  const [validationResult, setValidationResult] = useState<PlayValidationResult>({ isValid: true });
  
  const firstPlayRecommenderRef = useRef<FirstPlayRecommender | null>(null);
  const nextPlayRecommenderRef = useRef<NextPlayRecommender | null>(null);
  const lastHandCardsRef = useRef<Card[]>([]);
  const lastPlaysRef = useRef<PlayAction[]>([]);

  const isFirstPlay = useCallback(() => {
    if (leader === playerSeat && plays.length === 0) return true;
    const nonPassPlays = plays.filter(p => !p.isPass && p.cards.length > 0);
    return nonPassPlays.length === 0;
  }, [leader, playerSeat, plays]);

  const getLastValidPlay = useCallback((): PlayAction | null => {
    for (let i = plays.length - 1; i >= 0; i--) {
      const play = plays[i];
      if (!play.isPass && play.cards.length > 0) {
        return play;
      }
    }
    return null;
  }, [plays]);

  const handleHint = useCallback(() => {
    if (!canPlay || !isMyTurn || handCards.length === 0) return;

    const handCardsChanged = 
      lastHandCardsRef.current.length !== handCards.length ||
      lastHandCardsRef.current.some((c, i) => c.deckIndex !== handCards[i]?.deckIndex);
    const playsChanged =
      lastPlaysRef.current.length !== plays.length;

    if (handCardsChanged || playsChanged) {
      firstPlayRecommenderRef.current = null;
      nextPlayRecommenderRef.current = null;
      lastHandCardsRef.current = handCards;
      lastPlaysRef.current = plays;
    }

    const sdkHandCards = handCards.map(c => fromProtoCard(c, dealLevel));

    if (isFirstPlay()) {
      if (!firstPlayRecommenderRef.current) {
        firstPlayRecommenderRef.current = new FirstPlayRecommender(sdkHandCards);
      }
      const recommended = firstPlayRecommenderRef.current.next();
      if (recommended && recommended.length > 0) {
        const deckIndexes = new Set(recommended.map(c => c.deckIndex));
        const selectedCards = handCards.filter(c => deckIndexes.has(c.deckIndex));
        onHint(selectedCards);
      }
    } else {
      const lastPlay = getLastValidPlay();
      if (!lastPlay) return;

      if (!nextPlayRecommenderRef.current) {
        const lastPlaySdkCards = lastPlay.cards.map(c => fromProtoCard(c, dealLevel));
        const prevComp = fromCardList(lastPlaySdkCards);
        nextPlayRecommenderRef.current = new NextPlayRecommender(sdkHandCards, prevComp);
      }
      
      const recommended = nextPlayRecommenderRef.current.next();
      if (recommended) {
        const recommendedCards = recommended.getCards();
        const deckIndexes = new Set(recommendedCards.map(c => c.deckIndex));
        const selectedCards = handCards.filter(c => deckIndexes.has(c.deckIndex));
        onHint(selectedCards);
      } else {
        onHint([]);
      }
    }
  }, [canPlay, isMyTurn, handCards, plays, dealLevel, isFirstPlay, getLastValidPlay, onHint]);

  const validateCards = useCallback((cards: Card[]): PlayValidationResult => {
    if (cards.length === 0) {
      return { isValid: false, error: '请选择要出的牌' };
    }

    if (cards.length === 1) {
      return { isValid: true, cardType: '单牌' };
    }

    if (cards.length === 2) {
      if (cards[0].rank === cards[1].rank) {
        return { isValid: true, cardType: '对子' };
      }
      return { isValid: false, error: '两张牌必须是对子' };
    }

    if (cards.length === 3) {
      if (cards.every(card => card.rank === cards[0].rank)) {
        return { isValid: true, cardType: '三张' };
      }
      return { isValid: false, error: '三张牌必须是同点数' };
    }

    if (cards.length === 4) {
      if (cards.every(card => card.rank === cards[0].rank)) {
        return { isValid: true, cardType: '炸弹' };
      }
      return { isValid: true, cardType: '四张牌型' };
    }

    if (cards.length >= 5) {
      return { isValid: true, cardType: `${cards.length}张牌型` };
    }

    return { isValid: false, error: '无效的牌型组合' };
  }, []);

  useEffect(() => {
    if (selectedCards.length > 0) {
      const result = validateCards(selectedCards);
      setValidationResult(result);
    } else {
      setValidationResult({ isValid: true });
    }
  }, [selectedCards, validateCards]);

  const handlePlayCards = () => {
    if (selectedCards.length === 0) return;
    
    const validation = validateCards(selectedCards);
    if (!validation.isValid) {
      setValidationResult(validation);
      return;
    }

    onPlayCards(selectedCards);
  };

  const handlePass = () => {
    onPass();
  };

  const isPlayDisabled = disabled || !canPlay || !isMyTurn || selectedCards.length === 0 || !validationResult.isValid;
  const isPassDisabled = disabled || !canPlay || !isMyTurn;

  return (
    <div className="py-2 px-4">
      <div className="flex items-center justify-between gap-3">
        <button
          onClick={handlePass}
          disabled={isPassDisabled}
          className={`
            flex-1 py-2 px-8 rounded-xl font-semibold transition-all duration-200 border min-w-[100px]
            ${isPassDisabled 
              ? 'bg-slate-300/50 text-slate-400 border-slate-200/30' 
              : 'bg-slate-500 text-white border-slate-400/30 shadow-md hover:bg-slate-400 active:bg-slate-600'
            }
          `}
        >
          不出
        </button>

        <div className="flex items-center justify-center">
          <Countdown
            deadlineAtMs={turnDeadlineAtMs}
            isActive={isMyTurn && canPlay && !disabled}
            size="small"
          />
        </div>

        <button
          onClick={handleHint}
          disabled={isPassDisabled}
          className={`
            flex-1 py-2 px-8 rounded-xl font-semibold transition-all duration-200 border min-w-[100px]
            ${isPassDisabled 
              ? 'bg-slate-300/50 text-slate-400 border-slate-200/30' 
              : 'bg-gradient-to-b from-amber-400 to-amber-500 text-white border-amber-300/30 shadow-md hover:from-amber-300 hover:to-amber-400 active:from-amber-500 active:to-amber-600'
            }
          `}
        >
          提示
        </button>

        <button
          onClick={handlePlayCards}
          disabled={isPlayDisabled}
          className={`
            flex-1 py-2 px-8 rounded-xl font-semibold transition-all duration-200 border min-w-[100px]
            ${isPlayDisabled 
              ? 'bg-slate-300/50 text-slate-400 border-slate-200/30' 
              : 'bg-gradient-to-b from-emerald-500 to-emerald-600 text-white border-emerald-400/30 shadow-md hover:from-emerald-400 hover:to-emerald-500 hover:shadow-lg active:from-emerald-600 active:to-emerald-700'
            }
          `}
        >
          出牌
        </button>
      </div>
    </div>
  );
};

export default GameControls;
