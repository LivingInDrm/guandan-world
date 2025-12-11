import React, { useState, useEffect, useCallback, useRef, useMemo } from 'react';
import type { Card, PlayAction } from '../../types';
import Countdown from './Countdown';
import { Button } from '@/components/ui-next';
import {
  fromProtoCard,
  fromCardList,
  fromCardListWithType,
  FirstPlayRecommender,
  NextPlayRecommender,
  CompType,
  type CardComp,
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

  const prevComp = useMemo((): CardComp | undefined => {
    // 内联 isFirstPlay 逻辑
    const isFirst = (leader === playerSeat && plays.length === 0) || 
                    plays.filter(p => !p.isPass && p.cards.length > 0).length === 0;
    if (isFirst) return undefined;
    
    // 内联 getLastValidPlay 逻辑
    let lastPlay: PlayAction | null = null;
    for (let i = plays.length - 1; i >= 0; i--) {
      if (!plays[i].isPass && plays[i].cards.length > 0) {
        lastPlay = plays[i];
        break;
      }
    }
    if (!lastPlay) return undefined;
    
    const lastPlaySdkCards = lastPlay.cards.map(c => fromProtoCard(c, dealLevel));
    return fromCardListWithType(lastPlaySdkCards, lastPlay.compType);
  }, [plays, dealLevel, leader, playerSeat]);

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
      } else {
        onHint([]);
      }
    } else {
      if (!prevComp) return;

      if (!nextPlayRecommenderRef.current) {
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
  }, [canPlay, isMyTurn, handCards, plays, dealLevel, isFirstPlay, prevComp, onHint]);

  const validateCards = useCallback((cards: Card[]): PlayValidationResult => {
    if (cards.length === 0) {
      return { isValid: false, error: '请选择要出的牌' };
    }

    const sdkCards = cards.map(c => fromProtoCard(c, dealLevel));
    const comp = fromCardList(sdkCards, prevComp);
    
    if (comp.getType() === CompType.Illegal || !comp.isValid()) {
      return { isValid: false, error: '无效的牌型组合' };
    }

    if (prevComp && !comp.greaterThan(prevComp)) {
      return { isValid: false, error: '打不过上家' };
    }

    return { isValid: true, cardType: comp.toString() };
  }, [dealLevel, prevComp]);

  useEffect(() => {
    if (selectedCards.length === 0) {
      setValidationResult({ isValid: true });
      return;
    }

    // 立即设置为验证中状态，防止按钮闪烁
    setValidationResult({ isValid: false });

    const timer = setTimeout(() => {
      const result = validateCards(selectedCards);
      setValidationResult(result);
    }, 200);

    return () => clearTimeout(timer);
  }, [selectedCards, validateCards]);

  const handlePlayCards = () => {
    if (selectedCards.length === 0) return;
    
    // 同步验证兜底，防止防抖期间点击
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
  const isHintDisabled = disabled || !canPlay || !isMyTurn;
  const isPassDisabled = disabled || !canPlay || !isMyTurn || isFirstPlay();

  return (
    <div className="py-2 px-4">
      <div className="flex items-center justify-between gap-3">
        <Button
          intent="neutral"
          size="sm"
          onClick={handlePass}
          disabled={isPassDisabled}
          className="flex-1"
        >
          不出
        </Button>

        <div className="flex items-center justify-center">
          <Countdown
            deadlineAtMs={turnDeadlineAtMs}
            isActive={isMyTurn && canPlay && !disabled}
            size="small"
          />
        </div>

        <Button
          intent="secondary"
          size="sm"
          onClick={handleHint}
          disabled={isHintDisabled}
          className="flex-1"
        >
          提示
        </Button>

        <Button
          intent="primary"
          size="sm"
          onClick={handlePlayCards}
          disabled={isPlayDisabled}
          className="flex-1"
        >
          出牌
        </Button>
      </div>
    </div>
  );
};

export default GameControls;
