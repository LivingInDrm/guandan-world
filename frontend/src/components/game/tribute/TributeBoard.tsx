import React, { useRef, useCallback } from 'react';
import type { Player, Card } from '../../../types';
import { TributeType } from '../../../types/generated/event';
import type { TributeStep, FlyingCard } from '../../../store/tributeStore';
import { useTributeStore } from '../../../store/tributeStore';
import TributePool from './TributePool';
import CardFlyAnimation from './CardFlyAnimation';
import PlayerHand from '../PlayerHand';

interface TributeData {
  step: TributeStep;
  tributeStarted: {
    tributeType: TributeType;
    givers: number[];
    receivers: number[];
  } | null;
  tributeExempted: {
    bigJokerHolders: { [key: number]: number };
  } | null;
  submittedCards: { [giverSeat: number]: Card };
  poolCards: (Card | null)[];
  selectedCards: { [receiverSeat: number]: Card };
  returnedCards: Array<{
    fromSeat: number;
    toSeat: number;
    card: Card;
  }>;
  messages: string[];
  currentSelectingSeat: number | null;
  flyingCards: FlyingCard[];
  players: Player[];
  playerSeat: number | null;
}

interface TributeBoardProps {
  tributeData: TributeData;
  players: (Player | null)[];
  currentPlayerSeat: number;
  playerHand: Card[];
  selectedCards: Card[];
  onCardSelect: (cards: Card[]) => void;
  onSelectTribute: (deckIndex: number) => void;
  onReturnTribute: (deckIndex: number) => void;
}

const getTributeTypeName = (type: TributeType): string => {
  switch (type) {
    case TributeType.TRIBUTE_TYPE_DOUBLE_DOWN:
      return '双下';
    case TributeType.TRIBUTE_TYPE_SINGLE_LAST:
      return '单下';
    case TributeType.TRIBUTE_TYPE_PARTNER_LAST:
      return '末游';
    case TributeType.TRIBUTE_TYPE_NONE:
      return '无需进贡';
    default:
      return '未知';
  }
};

interface TributePlayerAreaProps {
  player: Player | null;
  position: 'bottom' | 'left' | 'top' | 'right';
  seat: number;
  isGiver: boolean;
  isReceiver: boolean;
  hasSubmitted: boolean;
  hasSelected: boolean;
  isCurrentSelector: boolean;
  bigJokerCount?: number;
}

const TributePlayerArea: React.FC<TributePlayerAreaProps> = ({
  player,
  position,
  isGiver,
  isReceiver,
  hasSubmitted,
  hasSelected,
  isCurrentSelector,
  bigJokerCount,
}) => {
  const getPositionClasses = () => {
    switch (position) {
      case 'bottom':
        return 'absolute bottom-4 left-1/2 transform -translate-x-1/2';
      case 'left':
        return 'absolute left-4 top-1/2 transform -translate-y-1/2';
      case 'top':
        return 'absolute top-4 left-1/2 transform -translate-x-1/2';
      case 'right':
        return 'absolute right-4 top-1/2 transform -translate-y-1/2';
      default:
        return '';
    }
  };

  const getRoleBadge = () => {
    if (isGiver) {
      return (
        <div className="text-xs px-2 py-1 rounded bg-red-200 text-red-800">
          {hasSubmitted ? '已上贡' : '待上贡'}
        </div>
      );
    }
    if (isReceiver) {
      return (
        <div className={`text-xs px-2 py-1 rounded ${
          isCurrentSelector 
            ? 'bg-yellow-200 text-yellow-800 animate-pulse' 
            : 'bg-green-200 text-green-800'
        }`}>
          {hasSelected ? '已收贡' : isCurrentSelector ? '选牌中' : '待收贡'}
        </div>
      );
    }
    return null;
  };

  if (!player) {
    return (
      <div className={`${getPositionClasses()} w-24 h-16`}>
        <div className="bg-gray-100 border-2 border-dashed border-gray-300 rounded-lg p-2 text-center">
          <div className="text-sm text-gray-400">空座位</div>
        </div>
      </div>
    );
  }

  return (
    <div className={`${getPositionClasses()} w-32`}>
      <div className={`border-2 rounded-lg p-2 text-center bg-white ${
        isCurrentSelector ? 'border-yellow-400 shadow-lg' : 'border-gray-300'
      }`}>
        <div className="text-sm font-medium truncate">{player.username}</div>
        {getRoleBadge()}
        {bigJokerCount !== undefined && bigJokerCount > 0 && (
          <div className="text-xs px-2 py-1 rounded bg-purple-200 text-purple-800 mt-1">
            大王 x{bigJokerCount}
          </div>
        )}
      </div>
    </div>
  );
};

const TributeBoard: React.FC<TributeBoardProps> = ({
  tributeData,
  players,
  currentPlayerSeat,
  playerHand,
  selectedCards,
  onCardSelect,
  onSelectTribute,
  onReturnTribute,
}) => {
  const boardRef = useRef<HTMLDivElement>(null);
  const poolRef = useRef<HTMLDivElement>(null);
  const poolSlotRefs = useRef<{ [slot: number]: HTMLDivElement | null }>({});

  const { 
    step, 
    tributeStarted, 
    tributeExempted, 
    submittedCards, 
    poolCards, 
    selectedCards: selectedByReceivers, 
    messages,
    currentSelectingSeat,
    flyingCards 
  } = tributeData;

  const removeFlyingCard = useTributeStore((s) => s.removeFlyingCard);

  const isGiver = (seat: number) => tributeStarted?.givers.includes(seat) ?? false;
  const isReceiver = (seat: number) => tributeStarted?.receivers.includes(seat) ?? false;
  const hasSubmitted = (seat: number) => seat in submittedCards;
  const hasSelected = (seat: number) => seat in selectedByReceivers;
  const getBigJokerCount = (seat: number) => tributeExempted?.bigJokerHolders[seat];

  const canSelectFromPool = 
    step === 'selecting' && 
    currentSelectingSeat === currentPlayerSeat;
  
  const canReturnTribute = step === 'returning' && isReceiver(currentPlayerSeat);

  const handlePoolCardSelect = (card: Card) => {
    if (canSelectFromPool) {
      onSelectTribute(card.deckIndex);
    }
  };

  const handleReturnCard = () => {
    if (canReturnTribute && selectedCards.length === 1) {
      onReturnTribute(selectedCards[0].deckIndex);
    }
  };

  const getStepTitle = () => {
    switch (step) {
      case 'started':
        return '上贡阶段开始';
      case 'exempted':
        return '抗贡成功';
      case 'submitting':
        return '等待进贡';
      case 'selecting':
        return currentSelectingSeat === currentPlayerSeat ? '请选择贡牌' : '等待选牌';
      case 'returning':
        return '还贡阶段';
      case 'completed':
        return '上贡完成';
      default:
        return '进贡阶段';
    }
  };

  const getPositionForSeat = useCallback((
    seat: number | 'pool',
    poolSlot?: number
  ): { x: number; y: number } => {
    if (!boardRef.current) return { x: 0, y: 0 };
    const boardRect = boardRef.current.getBoundingClientRect();

    if (seat === 'pool') {
      if (poolSlot !== undefined && poolSlotRefs.current[poolSlot]) {
        const slotRect = poolSlotRefs.current[poolSlot]!.getBoundingClientRect();
        return {
          x: slotRect.left + slotRect.width / 2,
          y: slotRect.top + slotRect.height / 2,
        };
      }
      if (poolRef.current) {
        const poolRect = poolRef.current.getBoundingClientRect();
        return {
          x: poolRect.left + poolRect.width / 2,
          y: poolRect.top + poolRect.height / 2,
        };
      }
      return {
        x: boardRect.left + boardRect.width / 2,
        y: boardRect.top + boardRect.height / 2,
      };
    }

    const relativeSeat = (seat - currentPlayerSeat + 4) % 4;
    switch (relativeSeat) {
      case 0:
        return { x: boardRect.left + boardRect.width / 2, y: boardRect.bottom - 40 };
      case 1:
        return { x: boardRect.left + 40, y: boardRect.top + boardRect.height / 2 };
      case 2:
        return { x: boardRect.left + boardRect.width / 2, y: boardRect.top + 40 };
      case 3:
        return { x: boardRect.right - 40, y: boardRect.top + boardRect.height / 2 };
      default:
        return { x: 0, y: 0 };
    }
  }, [currentPlayerSeat]);

  const handleAnimationComplete = useCallback((id: string) => {
    removeFlyingCard(id);
  }, [removeFlyingCard]);

  const handleSlotRefReady = useCallback((slotIndex: number, element: HTMLDivElement | null) => {
    poolSlotRefs.current[slotIndex] = element;
  }, []);

  return (
    <div className="max-w-6xl mx-auto p-6 space-y-6">
      <div ref={boardRef} className="relative w-full h-96 bg-green-100 border border-gray-300 rounded-lg">
        <div className="absolute top-4 left-4 bg-white border border-gray-300 rounded-lg p-3 shadow-sm z-10">
          <div className="text-sm font-medium mb-2">{getStepTitle()}</div>
          {tributeStarted && (
            <div className="text-xs text-gray-600">
              类型: {getTributeTypeName(tributeStarted.tributeType)}
            </div>
          )}
        </div>

        <div ref={poolRef} className="absolute inset-0 flex items-center justify-center">
          <TributePool
            poolCards={poolCards}
            canSelect={canSelectFromPool}
            onSelectCard={handlePoolCardSelect}
            tributeType={tributeStarted?.tributeType ?? TributeType.TRIBUTE_TYPE_UNSPECIFIED}
            messages={messages}
            onSlotRefReady={handleSlotRefReady}
          />
        </div>

        <TributePlayerArea
          player={players[currentPlayerSeat]}
          position="bottom"
          seat={currentPlayerSeat}
          isGiver={isGiver(currentPlayerSeat)}
          isReceiver={isReceiver(currentPlayerSeat)}
          hasSubmitted={hasSubmitted(currentPlayerSeat)}
          hasSelected={hasSelected(currentPlayerSeat)}
          isCurrentSelector={currentSelectingSeat === currentPlayerSeat}
          bigJokerCount={getBigJokerCount(currentPlayerSeat)}
        />

        <TributePlayerArea
          player={players[(currentPlayerSeat + 1) % 4]}
          position="left"
          seat={(currentPlayerSeat + 1) % 4}
          isGiver={isGiver((currentPlayerSeat + 1) % 4)}
          isReceiver={isReceiver((currentPlayerSeat + 1) % 4)}
          hasSubmitted={hasSubmitted((currentPlayerSeat + 1) % 4)}
          hasSelected={hasSelected((currentPlayerSeat + 1) % 4)}
          isCurrentSelector={currentSelectingSeat === (currentPlayerSeat + 1) % 4}
          bigJokerCount={getBigJokerCount((currentPlayerSeat + 1) % 4)}
        />

        <TributePlayerArea
          player={players[(currentPlayerSeat + 2) % 4]}
          position="top"
          seat={(currentPlayerSeat + 2) % 4}
          isGiver={isGiver((currentPlayerSeat + 2) % 4)}
          isReceiver={isReceiver((currentPlayerSeat + 2) % 4)}
          hasSubmitted={hasSubmitted((currentPlayerSeat + 2) % 4)}
          hasSelected={hasSelected((currentPlayerSeat + 2) % 4)}
          isCurrentSelector={currentSelectingSeat === (currentPlayerSeat + 2) % 4}
          bigJokerCount={getBigJokerCount((currentPlayerSeat + 2) % 4)}
        />

        <TributePlayerArea
          player={players[(currentPlayerSeat + 3) % 4]}
          position="right"
          seat={(currentPlayerSeat + 3) % 4}
          isGiver={isGiver((currentPlayerSeat + 3) % 4)}
          isReceiver={isReceiver((currentPlayerSeat + 3) % 4)}
          hasSubmitted={hasSubmitted((currentPlayerSeat + 3) % 4)}
          hasSelected={hasSelected((currentPlayerSeat + 3) % 4)}
          isCurrentSelector={currentSelectingSeat === (currentPlayerSeat + 3) % 4}
          bigJokerCount={getBigJokerCount((currentPlayerSeat + 3) % 4)}
        />

        {flyingCards.map((fc) => (
          <CardFlyAnimation
            key={fc.id}
            card={fc.card}
            fromPosition={getPositionForSeat(fc.fromSeat, fc.fromPoolSlot)}
            toPosition={getPositionForSeat(fc.toSeat, fc.toPoolSlot)}
            onComplete={() => handleAnimationComplete(fc.id)}
          />
        ))}
      </div>

      {step === 'returning' && canReturnTribute && (
        <>
          <PlayerHand
            cards={playerHand}
            selectedCards={selectedCards}
            onCardSelect={onCardSelect}
          />
          <div className="flex justify-center">
            <button
              onClick={handleReturnCard}
              disabled={selectedCards.length !== 1}
              className={`px-6 py-2 rounded-lg font-medium transition-colors ${
                selectedCards.length === 1
                  ? 'bg-blue-500 text-white hover:bg-blue-600'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              }`}
            >
              确认还贡
            </button>
          </div>
        </>
      )}

      {step !== 'returning' && (
        <div className="bg-white rounded-lg p-4 shadow">
          <h3 className="text-sm font-medium mb-2">消息记录</h3>
          <div className="text-xs text-gray-600 space-y-1">
            {messages.map((msg, i) => (
              <div key={i}>{msg}</div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
};

export default TributeBoard;
